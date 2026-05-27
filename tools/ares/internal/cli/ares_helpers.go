// Hand-authored ARES domain helpers for the transcendence commands.
// Not generator-emitted; survives regen as a whole file.
package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"ares-pp-cli/internal/client"
	"ares-pp-cli/internal/store"
)

// normalizeIco strips spaces and left-pads a numeric IČO to 8 digits.
// Returns ok=false if the input is not 1-8 digits.
func normalizeIco(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if s == "" || len(s) > 8 {
		return "", false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return "", false
		}
	}
	return fmt.Sprintf("%08s", s), true
}

// validIco reports whether s is a structurally valid Czech IČO using the
// mod-11 checksum over the first 7 digits (weights 8..2).
func validIco(s string) bool {
	n, ok := normalizeIco(s)
	if !ok {
		return false
	}
	sum := 0
	for i := 0; i < 7; i++ {
		sum += int(n[i]-'0') * (8 - i)
	}
	mod := sum % 11
	var check int
	switch mod {
	case 0:
		check = 1
	case 1:
		check = 0
	default:
		check = 11 - mod
	}
	return check == int(n[7]-'0')
}

// deriveDic returns the conventional Czech DIČ (CZ + IČO) for a valid IČO.
// This is the structural form; VAT-payer status is read from registry data.
func deriveDic(ico string) string {
	n, ok := normalizeIco(ico)
	if !ok {
		return ""
	}
	return "CZ" + strings.TrimLeft(n, "0")
}

// fetchSubject calls GET /ekonomicke-subjekty/{ico} and normalizes the
// high-gravity fields into a store.Subject. The full payload is kept in Raw.
func fetchSubject(ctx context.Context, c *client.Client, ico string) (store.Subject, error) {
	n, ok := normalizeIco(ico)
	if !ok {
		return store.Subject{}, fmt.Errorf("invalid IČO %q: must be 1-8 digits", ico)
	}
	path := replacePathParam("/ekonomicke-subjekty/{ico}", "ico", n)
	data, err := c.Get(ctx, path, map[string]string{})
	if err != nil {
		return store.Subject{}, err
	}
	var doc struct {
		Ico           string `json:"ico"`
		ObchodniJmeno string `json:"obchodniJmeno"`
		PravniForma   string `json:"pravniForma"`
		Sidlo         struct {
			TextovaAdresa string `json:"textovaAdresa"`
		} `json:"sidlo"`
	}
	_ = json.Unmarshal(data, &doc)
	sub := store.Subject{
		Ico:           n,
		ObchodniJmeno: doc.ObchodniJmeno,
		Adresa:        doc.Sidlo.TextovaAdresa,
		PravniForma:   doc.PravniForma,
		Raw:           data,
		FetchedAt:     time.Now().UTC(),
	}
	if sub.Ico == "" {
		sub.Ico = n
	}
	return sub, nil
}

// parseSince converts "30d", "12h", "2026-01-01" (or RFC3339) to a cutoff time.
func parseSince(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, nil
	}
	if strings.HasSuffix(s, "d") {
		var days int
		if _, err := fmt.Sscanf(s, "%dd", &days); err == nil {
			return time.Now().UTC().Add(-time.Duration(days) * 24 * time.Hour), nil
		}
	}
	if d, err := time.ParseDuration(s); err == nil {
		return time.Now().UTC().Add(-d), nil
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse --since %q (use 30d, 12h, or 2006-01-02)", s)
}
