// Copyright 2026 patrik-zita. Licensed under Apache-2.0. See LICENSE.

package cli

import (
	"encoding/json"
	"fmt"
	"strings"
)

// coerceStringArray normalizes a flag value into a slice of strings, as required
// by ARES array[string] filter fields (czNace, pravniFormaRos). It accepts a
// JSON array (`["461","469"]`), a comma-separated list (`461,469`), or a bare
// value (`6210`), always producing []string. Without this, a bare numeric value
// like `6210` JSON-parses to a number and ARES rejects it (HTTP 500) because the
// field must be an array of strings.
func coerceStringArray(raw string) ([]string, error) {
	if arr := []string(nil); json.Unmarshal([]byte(raw), &arr) == nil {
		return arr, nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("empty value")
	}
	return out, nil
}
