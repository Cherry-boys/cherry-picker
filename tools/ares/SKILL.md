---
name: pp-ares
description: "Every ARES lookup, plus a local store, offline search, bulk IČO enrichment Trigger phrases: `look up IČO`, `find Czech company`, `validate IČO`, `enrich these IČOs`, `check insolvency ARES`, `use ares`, `run ares`."
author: "Patrik Zita"
license: "Apache-2.0"
argument-hint: "<command> [args] | install cli|mcp"
allowed-tools: "Read Bash"
metadata:
  openclaw:
    requires:
      bins:
        - ares-pp-cli
    install:
      - kind: go
        bins: [ares-pp-cli]
        module: github.com/mvanhorn/printing-press-library/library/sales-and-crm/ares/cmd/ares-pp-cli
---

# ARES — Printing Press CLI

## Prerequisites: Install the CLI

This skill drives the `ares-pp-cli` binary. **You must verify the CLI is installed before invoking any command from this skill.** If it is missing, install it first:

1. Install via the Printing Press installer:
   ```bash
   npx -y @mvanhorn/printing-press-library install ares --cli-only
   ```
2. Verify: `ares-pp-cli --version`
3. Ensure `$GOPATH/bin` (or `$HOME/go/bin`) is on `$PATH`.

If the `npx` install fails (no Node, offline, etc.), fall back to a direct Go install (requires Go 1.26.3 or newer):

```bash
go install github.com/mvanhorn/printing-press-library/library/sales-and-crm/ares/cmd/ares-pp-cli@latest
```

If `--version` reports "command not found" after install, the install step did not put the binary on `$PATH`. Do not proceed with skill commands until verification succeeds.

ares wraps the Czech business registry's public REST API across all of its source registers (VR, RES, RŽP, ROS, CEÚ insolvency, RÚIAN addresses), then adds a local SQLite layer the official API has no notion of: bulk-enrich a column of IČOs for invoicing, keep a watchlist of clients and suppliers, and diff their name/address/insolvency status over time.

## When to Use This CLI

Use ares when an agent or script needs Czech company data by IČO or name: prefilling invoices, validating counterparties, enriching a CRM, or watching a set of clients and suppliers for address or insolvency changes. It is the right choice over a raw API call when you need offline search, bulk processing, or change detection across runs.

## Unique Capabilities

These capabilities aren't available in any other tool for this API.

### Bulk + local state
- **`enrich`** — Pipe a list of IČOs and get one enriched company record per line for invoicing or CRM import.

  _Reach for this to turn a column of IČOs into company profiles in one pass instead of N manual lookups._

  ```bash
  ares-pp-cli enrich 00177041 27082440 --agent
  ```
- **`portfolio`** — Keep a local watchlist of client and supplier IČOs and bulk-refresh them into the store.

  _Use when an agent must track a fixed set of counterparties across runs._

  ```bash
  ares-pp-cli portfolio add 00177041 --agent
  ```
- **`changes`** — Report name, address, legal-form, and status changes across tracked subjects since a point in time.

  _Reach for this to detect when a counterparty moved, renamed, or changed legal form._

  ```bash
  ares-pp-cli changes --since 30d --agent
  ```
- **`search`** — Full-text search over synced company names and addresses with no network call.

  _Reach for this when you need fast, composable lookups without burning API quota._

  ```bash
  ares-pp-cli search "mladá boleslav" --json --select ico,obchodniJmeno
  ```

### Risk
- **`insolvency-watch`** — Flag portfolio members that newly appear among CEÚ insolvency subjects.

  _Use before invoicing or extending credit to a counterparty._

  ```bash
  ares-pp-cli insolvency-watch --agent
  ```

### Validation
- **`validate`** — Validate Czech IČO checksums in bulk with zero API calls.

  _Use to filter a list down to structurally valid IČOs before hitting the API._

  ```bash
  ares-pp-cli validate 00177041 --agent
  ```

## Command Reference

**ciselniky-nazevniky** — ciselniky-nazevniky operations

- `ares-pp-cli ciselniky-nazevniky` — Vyhledání číselníků používaných v IS ARES podle komplexního filtru

**ekonomicke-subjekty** — ekonomicke-subjekty operations

- `ares-pp-cli ekonomicke-subjekty vrat-ekonomicky-subjekt` — Vyhledání ekonomického subjektu ARES podle zadaného iča
- `ares-pp-cli ekonomicke-subjekty vyhledej` — Vyhledání seznamu ekonomických subjektů ARES podle komplexního filtru

**ekonomicke-subjekty-ceu** — ekonomicke-subjekty-ceu operations

- `ares-pp-cli ekonomicke-subjekty-ceu vrat-ekonomicky-subjekt-ceu` — Vyhledání konkrétního úpadce ze zdroje CEÚ
- `ares-pp-cli ekonomicke-subjekty-ceu vyhledej-seznam-ekonomickych-subjektu-ceu` — Vyhledání seznamu úpadců ze zdroje CEÚ

**ekonomicke-subjekty-notifikace** — ekonomicke-subjekty-notifikace operations

- `ares-pp-cli ekonomicke-subjekty-notifikace vrat-notifikacni-davku` — Vyhledání konkrétní notifikační dávky zdroje ARES podle vstupcách parametrů (zdroj, číslo notifikační dávky)
- `ares-pp-cli ekonomicke-subjekty-notifikace vyhledej-seznam-notifikacnich-davek` — Vyhledání seznamu notifikačních dávek ekonomických subjektů ARES podle zvoleného filtru

**ekonomicke-subjekty-nrpzs** — ekonomicke-subjekty-nrpzs operations

- `ares-pp-cli ekonomicke-subjekty-nrpzs vrat-ekonomicky-subjekt-nrpzs` — Vyhledání konkrétního ekonomického subjektu ze zdroje NRPZS
- `ares-pp-cli ekonomicke-subjekty-nrpzs vyhledej-seznam-ekonomickych-subjektu-nrpzs` — Vyhledání seznamu ekonomických subjektu ze zdroje NRPZS

**ekonomicke-subjekty-rcns** — ekonomicke-subjekty-rcns operations

- `ares-pp-cli ekonomicke-subjekty-rcns vrat-ekonomicky-subjekt-rcns` — Vyhledání konkrétního ekonomického subjektu ze zdroje RCNS
- `ares-pp-cli ekonomicke-subjekty-rcns vyhledej-seznam-ekonomickych-subjektu-rcns` — Vyhledání seznamu ekonomických subjektu ze zdroje RCNS

**ekonomicke-subjekty-res** — ekonomicke-subjekty-res operations

- `ares-pp-cli ekonomicke-subjekty-res vrat-ekonomicky-subjekt-res` — Vyhledání konkrétního ekonomického subjektu ze zdroje RES
- `ares-pp-cli ekonomicke-subjekty-res vyhledej-seznam-ekonomickych-subjektu-res` — Vyhledání seznamu ekonomických subjektu ze zdroje RES

**ekonomicke-subjekty-ros** — ekonomicke-subjekty-ros operations

- `ares-pp-cli ekonomicke-subjekty-ros vrat-ekonomicky-subjekt-ros` — Vyhledání konkrétního ekonomického subjektu ze zdroje ROS
- `ares-pp-cli ekonomicke-subjekty-ros vyhledej-seznam-ekonomickych-subjektu-ros` — Vyhledání seznamu ekonomických subjektu ze zdroje ROS

**ekonomicke-subjekty-rpsh** — ekonomicke-subjekty-rpsh operations

- `ares-pp-cli ekonomicke-subjekty-rpsh vrat-ekonomicky-subjekt-rpsh` — Vyhledání konkrétního ekonomického subjektu ze zdroje RPSH
- `ares-pp-cli ekonomicke-subjekty-rpsh vyhledej-seznam-ekonomickych-subjektu-rpsh` — Vyhledání seznamu ekonomických subjektu ze zdroje RPSH

**ekonomicke-subjekty-rs** — ekonomicke-subjekty-rs operations

- `ares-pp-cli ekonomicke-subjekty-rs vrat-ekonomicky-subjekt-rs` — Vyhledání konkrétního ekonomického subjektu ze zdroje RŠ
- `ares-pp-cli ekonomicke-subjekty-rs vyhledej-seznam-ekonomickych-subjektu-rs` — Vyhledání seznamu ekonomických subjektu ze zdroje RŠ

**ekonomicke-subjekty-rzp** — ekonomicke-subjekty-rzp operations

- `ares-pp-cli ekonomicke-subjekty-rzp vrat-ekonomicky-subjekt-rzp` — Vyhledání konkrétního ekonomického subjektu ze zdroje RŽP
- `ares-pp-cli ekonomicke-subjekty-rzp vyhledej-seznam-ekonomickych-subjektu-rzp` — Vyhledání seznamu ekonomických subjektu ze zdroje RŽP

**ekonomicke-subjekty-szr** — ekonomicke-subjekty-szr operations

- `ares-pp-cli ekonomicke-subjekty-szr vrat-ekonomicky-subjekt-szr` — Vyhledání konkrétního ekonomického subjektu ze zdroje SZR - subregistr EZP
- `ares-pp-cli ekonomicke-subjekty-szr vyhledej-seznam-ekonomickych-subjektu-szr` — Vyhledání seznamu ekonomických subjektu ze zdroje SZR - subregistr EZP

**ekonomicke-subjekty-vr** — ekonomicke-subjekty-vr operations

- `ares-pp-cli ekonomicke-subjekty-vr vrat-ekonomicky-subjekt-vr` — Vyhledání konkrétního ekonomického subjektu ze zdroje VR
- `ares-pp-cli ekonomicke-subjekty-vr vyhledej-seznam-ekonomickych-subjektu-vr` — Vyhledání seznamu ekonomických subjektu ze zdroje VR

**standardizovane-adresy** — standardizovane-adresy operations

- `ares-pp-cli standardizovane-adresy` — Vyhledání seznamu standardizovaných adres RÚIAN podle komplexního filtru


### Finding the right command

When you know what you want to do but not which command does it, ask the CLI directly:

```bash
ares-pp-cli which "<capability in your own words>"
```

`which` resolves a natural-language capability query to the best matching command from this CLI's curated feature index. Exit code `0` means at least one match; exit code `2` means no confident match — fall back to `--help` or use a narrower query.

## Recipes


### Enrich an invoicing list

```bash
ares-pp-cli enrich 00177041 27082440 --agent
```

Reads IČOs from stdin and emits one enriched company record per line.

### Validate before lookup

```bash
ares-pp-cli validate 00177041 27082440 --agent
```

Filters a list to structurally valid IČOs with no API calls.

### Narrow a verbose subject payload

```bash
ares-pp-cli ekonomicke-subjekty get 00177041 --agent --select ico,obchodniJmeno,sidlo.textovaAdresa,pravniForma
```

ARES subject records are deeply nested; --select pulls only the fields you need.

### Watch counterparties for insolvency

```bash
ares-pp-cli insolvency-watch --agent
```

Flags portfolio members newly appearing in the CEÚ insolvency register.

## Auth Setup

No authentication required.

Run `ares-pp-cli doctor` to verify setup.

## Agent Mode

Add `--agent` to any command. Expands to: `--json --compact --no-input --no-color --yes`.

- **Pipeable** — JSON on stdout, errors on stderr
- **Filterable** — `--select` keeps a subset of fields. Dotted paths descend into nested structures; arrays traverse element-wise. Critical for keeping context small on verbose APIs:

  ```bash
  ares-pp-cli ciselniky-nazevniky --agent --select id,name,status
  ```
- **Previewable** — `--dry-run` shows the request without sending
- **Non-interactive** — never prompts, every input is a flag
- **Explicit retries** — use `--idempotent` only when an already-existing create should count as success

## Agent Feedback

When you (or the agent) notice something off about this CLI, record it:

```
ares-pp-cli feedback "the --since flag is inclusive but docs say exclusive"
ares-pp-cli feedback --stdin < notes.txt
ares-pp-cli feedback list --json --limit 10
```

Entries are stored locally at `~/.ares-pp-cli/feedback.jsonl`. They are never POSTed unless `ARES_FEEDBACK_ENDPOINT` is set AND either `--send` is passed or `ARES_FEEDBACK_AUTO_SEND=true`. Default behavior is local-only.

Write what *surprised* you, not a bug report. Short, specific, one line: that is the part that compounds.

## Output Delivery

Every command accepts `--deliver <sink>`. The output goes to the named sink in addition to (or instead of) stdout, so agents can route command results without hand-piping. Three sinks are supported:

| Sink | Effect |
|------|--------|
| `stdout` | Default; write to stdout only |
| `file:<path>` | Atomically write output to `<path>` (tmp + rename) |
| `webhook:<url>` | POST the output body to the URL (`application/json` or `application/x-ndjson` when `--compact`) |

Unknown schemes are refused with a structured error naming the supported set. Webhook failures return non-zero and log the URL + HTTP status on stderr.

## Named Profiles

A profile is a saved set of flag values, reused across invocations. Use it when a scheduled agent calls the same command every run with the same configuration - HeyGen's "Beacon" pattern.

```
ares-pp-cli profile save briefing --json
ares-pp-cli --profile briefing ciselniky-nazevniky
ares-pp-cli profile list --json
ares-pp-cli profile show briefing
ares-pp-cli profile delete briefing --yes
```

Explicit flags always win over profile values; profile values win over defaults. `agent-context` lists all available profiles under `available_profiles` so introspecting agents discover them at runtime.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 2 | Usage error (wrong arguments) |
| 3 | Resource not found |
| 5 | API error (upstream issue) |
| 7 | Rate limited (wait and retry) |
| 10 | Config error |

## Argument Parsing

Parse `$ARGUMENTS`:

1. **Empty, `help`, or `--help`** → show `ares-pp-cli --help` output
2. **Starts with `install`** → ends with `mcp` → MCP installation; otherwise → see Prerequisites above
3. **Anything else** → Direct Use (execute as CLI command with `--agent`)

## MCP Server Installation

1. Install the MCP server:
   ```bash
   go install github.com/mvanhorn/printing-press-library/library/sales-and-crm/ares/cmd/ares-pp-mcp@latest
   ```
2. Register with Claude Code:
   ```bash
   claude mcp add ares-pp-mcp -- ares-pp-mcp
   ```
3. Verify: `claude mcp list`

## Direct Use

1. Check if installed: `which ares-pp-cli`
   If not found, offer to install (see Prerequisites at the top of this skill).
2. Match the user query to the best command from the Unique Capabilities and Command Reference above.
3. Execute with the `--agent` flag:
   ```bash
   ares-pp-cli <command> [subcommand] [args] --agent
   ```
4. If ambiguous, drill into subcommand help: `ares-pp-cli <command> --help`.
