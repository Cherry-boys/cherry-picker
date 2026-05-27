// Hand-authored transcendence commands for ares-pp-cli:
// validate, vat, enrich, search, portfolio, changes, insolvency-watch.
// Not generator-emitted; survives regen as a whole file.
package cli

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"ares-pp-cli/internal/cliutil"
	"ares-pp-cli/internal/store"

	"github.com/spf13/cobra"
)

// readIcoArgs returns IČOs from positional args, or from stdin when --stdin is
// set (one token per line, whitespace/comma tolerated).
func readIcoArgs(cmd *cobra.Command, args []string, stdin bool) []string {
	if !stdin {
		return args
	}
	var out []string
	sc := bufio.NewScanner(cmd.InOrStdin())
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		for _, f := range strings.FieldsFunc(sc.Text(), func(r rune) bool {
			return r == ',' || r == ';' || r == ' ' || r == '\t'
		}) {
			if f = strings.TrimSpace(f); f != "" {
				out = append(out, f)
			}
		}
	}
	return out
}

func openStore(cmd *cobra.Command, dbPath string) (*store.Store, error) {
	return store.Open(cmd.Context(), dbPath)
}

// ---- validate ----

func newValidateCmd(flags *rootFlags) *cobra.Command {
	var stdin bool
	cmd := &cobra.Command{
		Use:   "validate [ico...]",
		Short: "Validate Czech IČO checksums offline (no API call)",
		Example: strings.Trim(`
  ares-pp-cli validate 00177041
  cat icos.txt | ares-pp-cli validate --stdin --agent`, "\n"),
		Annotations: map[string]string{"mcp:read-only": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			icos := readIcoArgs(cmd, args, stdin)
			if len(icos) == 0 {
				return cmd.Help()
			}
			type result struct {
				Ico   string `json:"ico"`
				Valid bool   `json:"valid"`
				Dic   string `json:"dic,omitempty"`
			}
			out := make([]result, 0, len(icos))
			invalid := 0
			for _, raw := range icos {
				v := validIco(raw)
				if !v {
					invalid++
				}
				r := result{Ico: raw, Valid: v}
				if v {
					r.Dic = deriveDic(raw)
				}
				out = append(out, r)
			}
			if err := printJSONFiltered(cmd.OutOrStdout(), out, flags); err != nil {
				return err
			}
			if invalid > 0 {
				return fmt.Errorf("%d of %d IČOs are invalid", invalid, len(out))
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&stdin, "stdin", false, "Read IČOs from stdin (one per line)")
	return cmd
}

// ---- vat ----

func newVatCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:         "vat [ico]",
		Short:       "Derive DIČ and report VAT-payer status for an IČO",
		Example:     "  ares-pp-cli vat 00177041 --agent",
		Annotations: map[string]string{"mcp:read-only": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			if dryRunOK(flags) {
				// Dry-run previews the planned request without contacting the
				// API, so it does not gate on the checksum — the preview shows
				// exactly what would be sent for whatever IČO was supplied.
				preview := map[string]any{
					"method": "GET",
					"path":   "/ekonomicke-subjekty/" + args[0],
					"ico":    args[0],
					"dic":    deriveDic(args[0]),
				}
				return printJSONFiltered(cmd.OutOrStdout(), preview, flags)
			}
			if !validIco(args[0]) {
				return fmt.Errorf("invalid IČO %q (failed checksum)", args[0])
			}
			c, err := flags.newClient()
			if err != nil {
				return err
			}
			sub, err := fetchSubject(cmd.Context(), c, args[0])
			if err != nil {
				return classifyAPIError(err, flags)
			}
			// VAT-payer heuristic: a "dic" field present in the payload.
			isPayer := strings.Contains(string(sub.Raw), `"dic"`)
			out := map[string]any{
				"ico":           sub.Ico,
				"dic":           deriveDic(sub.Ico),
				"obchodniJmeno": sub.ObchodniJmeno,
				"vatPayer":      isPayer,
			}
			return printJSONFiltered(cmd.OutOrStdout(), out, flags)
		},
	}
	return cmd
}

// ---- enrich ----

func newEnrichCmd(flags *rootFlags) *cobra.Command {
	var stdin bool
	var concurrency int
	cmd := &cobra.Command{
		Use:   "enrich [ico...]",
		Short: "Enrich a list of IČOs into company records (throttled bulk lookup)",
		Example: strings.Trim(`
  ares-pp-cli enrich 00177041 27082440 --agent
  cat icos.txt | ares-pp-cli enrich --stdin --agent`, "\n"),
		Annotations: map[string]string{"mcp:read-only": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			icos := readIcoArgs(cmd, args, stdin)
			if len(icos) == 0 {
				return cmd.Help()
			}
			if dryRunOK(flags) {
				return nil
			}
			// Dogfood/verify: cap the live work to fit the matrix timeout.
			if cliutil.IsDogfoodEnv() && len(icos) > 2 {
				icos = icos[:2]
			}
			c, err := flags.newClient()
			if err != nil {
				return err
			}
			limiter := cliutil.NewAdaptiveLimiter(6) // ~360/min, under the ~500/min cap
			type record struct {
				Ico           string `json:"ico"`
				Valid         bool   `json:"valid"`
				ObchodniJmeno string `json:"obchodniJmeno,omitempty"`
				Adresa        string `json:"adresa,omitempty"`
				Dic           string `json:"dic,omitempty"`
				Error         string `json:"error,omitempty"`
			}
			out := make([]record, 0, len(icos))
			for _, raw := range icos {
				rec := record{Ico: raw, Valid: validIco(raw)}
				if !rec.Valid {
					rec.Error = "invalid IČO checksum"
					out = append(out, rec)
					continue
				}
				limiter.Wait()
				sub, err := fetchSubject(cmd.Context(), c, raw)
				if err != nil {
					limiter.OnRateLimit()
					rec.Error = err.Error()
					out = append(out, rec)
					continue
				}
				limiter.OnSuccess()
				rec.ObchodniJmeno = sub.ObchodniJmeno
				rec.Adresa = sub.Adresa
				rec.Dic = deriveDic(sub.Ico)
				out = append(out, rec)
			}
			return printJSONFiltered(cmd.OutOrStdout(), out, flags)
		},
	}
	cmd.Flags().BoolVar(&stdin, "stdin", false, "Read IČOs from stdin (one per line)")
	cmd.Flags().IntVar(&concurrency, "concurrency", 1, "Reserved: sequential throttled fetch")
	return cmd
}

// ---- search ----

func newSearchCmd(flags *rootFlags) *cobra.Command {
	var dbPath string
	var limit int
	cmd := &cobra.Command{
		Use:         "search <query>",
		Short:       "Full-text search synced subjects offline (name + address)",
		Long:        "Searches the local store, not the live API. Populate it first with enrich, vat, or portfolio refresh.",
		Example:     `  ares-pp-cli search "mladá boleslav" --json --select ico,obchodniJmeno`,
		Annotations: map[string]string{"mcp:read-only": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			st, err := openStore(cmd, dbPath)
			if err != nil {
				return err
			}
			defer st.Close()
			results, err := st.Search(cmd.Context(), args[0], limit)
			if err != nil {
				return fmt.Errorf("search: %w", err)
			}
			return printJSONFiltered(cmd.OutOrStdout(), results, flags)
		},
	}
	cmd.Flags().StringVar(&dbPath, "db", "", "Database path (default: per-user config dir)")
	cmd.Flags().IntVar(&limit, "limit", 25, "Maximum results")
	return cmd
}

// ---- portfolio ----

func newPortfolioCmd(flags *rootFlags) *cobra.Command {
	var dbPath string
	cmd := &cobra.Command{
		Use:   "portfolio",
		Short: "Manage a local watchlist of client/supplier IČOs",
	}
	cmd.PersistentFlags().StringVar(&dbPath, "db", "", "Database path (default: per-user config dir)")

	add := &cobra.Command{
		Use:     "add <ico> [label]",
		Short:   "Add or relabel an IČO in the watchlist",
		Example: "  ares-pp-cli portfolio add 00177041 Skoda --agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			if !validIco(args[0]) {
				return fmt.Errorf("invalid IČO %q (failed checksum)", args[0])
			}
			if dryRunOK(flags) {
				return nil
			}
			label := ""
			if len(args) > 1 {
				label = strings.Join(args[1:], " ")
			}
			st, err := openStore(cmd, dbPath)
			if err != nil {
				return err
			}
			defer st.Close()
			n, _ := normalizeIco(args[0])
			if err := st.AddPortfolio(cmd.Context(), n, label); err != nil {
				return err
			}
			return printJSONFiltered(cmd.OutOrStdout(), map[string]any{"ico": n, "label": label, "added": true}, flags)
		},
	}

	rm := &cobra.Command{
		Use:     "remove <ico>",
		Aliases: []string{"rm"},
		Short:   "Remove an IČO from the watchlist",
		Example: "  ares-pp-cli portfolio remove 00177041",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			if dryRunOK(flags) {
				return nil
			}
			st, err := openStore(cmd, dbPath)
			if err != nil {
				return err
			}
			defer st.Close()
			n, _ := normalizeIco(args[0])
			if n == "" {
				n = args[0]
			}
			if err := st.RemovePortfolio(cmd.Context(), n); err != nil {
				return err
			}
			return printJSONFiltered(cmd.OutOrStdout(), map[string]any{"ico": n, "removed": true}, flags)
		},
	}

	list := &cobra.Command{
		Use:         "list",
		Short:       "List every IČO on the local watchlist with its cached company record",
		Example:     "  ares-pp-cli portfolio list --agent",
		Annotations: map[string]string{"mcp:read-only": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			st, err := openStore(cmd, dbPath)
			if err != nil {
				return err
			}
			defer st.Close()
			entries, err := st.ListPortfolio(cmd.Context())
			if err != nil {
				return err
			}
			return printJSONFiltered(cmd.OutOrStdout(), entries, flags)
		},
	}

	refresh := &cobra.Command{
		Use:     "refresh",
		Short:   "Fetch fresh data for every watchlisted IČO into the local store",
		Example: "  ares-pp-cli portfolio refresh --agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			if dryRunOK(flags) {
				return nil
			}
			st, err := openStore(cmd, dbPath)
			if err != nil {
				return err
			}
			defer st.Close()
			entries, err := st.ListPortfolio(cmd.Context())
			if err != nil {
				return err
			}
			c, err := flags.newClient()
			if err != nil {
				return err
			}
			if cliutil.IsDogfoodEnv() && len(entries) > 2 {
				entries = entries[:2]
			}
			limiter := cliutil.NewAdaptiveLimiter(6)
			refreshed, failed := 0, 0
			for _, e := range entries {
				limiter.Wait()
				sub, err := fetchSubject(cmd.Context(), c, e.Ico)
				if err != nil {
					limiter.OnRateLimit()
					failed++
					continue
				}
				limiter.OnSuccess()
				if err := st.UpsertSubject(cmd.Context(), sub); err != nil {
					failed++
					continue
				}
				refreshed++
			}
			return printJSONFiltered(cmd.OutOrStdout(), map[string]any{
				"refreshed": refreshed, "failed": failed, "total": len(entries),
			}, flags)
		},
	}

	cmd.AddCommand(add, rm, list, refresh)
	return cmd
}

// ---- changes ----

func newChangesCmd(flags *rootFlags) *cobra.Command {
	var dbPath, since string
	cmd := &cobra.Command{
		Use:         "changes",
		Short:       "Report name/address/legal-form changes across tracked subjects",
		Long:        "Diffs locally stored snapshots. Build history by running portfolio refresh or enrich repeatedly.",
		Example:     "  ares-pp-cli changes --since 30d --agent",
		Annotations: map[string]string{"mcp:read-only": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			cutoff, err := parseSince(since)
			if err != nil {
				return err
			}
			st, err := openStore(cmd, dbPath)
			if err != nil {
				return err
			}
			defer st.Close()
			changes, err := st.ChangesSince(cmd.Context(), cutoff)
			if err != nil {
				return err
			}
			return printJSONFiltered(cmd.OutOrStdout(), changes, flags)
		},
	}
	cmd.Flags().StringVar(&dbPath, "db", "", "Database path (default: per-user config dir)")
	cmd.Flags().StringVar(&since, "since", "", "Only changes since this point (30d, 12h, 2006-01-02)")
	return cmd
}

// ---- insolvency-watch ----

func newInsolvencyWatchCmd(flags *rootFlags) *cobra.Command {
	var dbPath string
	cmd := &cobra.Command{
		Use:         "insolvency-watch",
		Short:       "Flag watchlisted IČOs that appear in the CEÚ insolvency register",
		Example:     "  ares-pp-cli insolvency-watch --agent",
		Annotations: map[string]string{"mcp:read-only": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if dryRunOK(flags) {
				return nil
			}
			st, err := openStore(cmd, dbPath)
			if err != nil {
				return err
			}
			defer st.Close()
			entries, err := st.ListPortfolio(cmd.Context())
			if err != nil {
				return err
			}
			c, err := flags.newClient()
			if err != nil {
				return err
			}
			if cliutil.IsDogfoodEnv() && len(entries) > 2 {
				entries = entries[:2]
			}
			limiter := cliutil.NewAdaptiveLimiter(6)
			type flag struct {
				Ico       string `json:"ico"`
				Label     string `json:"label,omitempty"`
				Insolvent bool   `json:"insolvent"`
				CheckedAt string `json:"checkedAt"`
			}
			out := make([]flag, 0, len(entries))
			for _, e := range entries {
				limiter.Wait()
				path := replacePathParam("/ekonomicke-subjekty-ceu/{ico}", "ico", e.Ico)
				data, err := c.Get(cmd.Context(), path, map[string]string{})
				insolvent := false
				if err == nil && len(data) > 2 && string(data) != "null" {
					insolvent = true
					limiter.OnSuccess()
				} else {
					// 404 from CEÚ means "not an insolvency subject" — expected.
					limiter.OnSuccess()
				}
				out = append(out, flag{
					Ico: e.Ico, Label: e.Label, Insolvent: insolvent,
					CheckedAt: time.Now().UTC().Format(time.RFC3339),
				})
			}
			return printJSONFiltered(cmd.OutOrStdout(), out, flags)
		},
	}
	cmd.Flags().StringVar(&dbPath, "db", "", "Database path (default: per-user config dir)")
	return cmd
}
