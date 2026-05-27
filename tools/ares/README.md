# ARES CLI

**Every ARES lookup, plus a local store, offline search, bulk IČO enrichment, and portfolio change-tracking no other ARES tool has.**

ares wraps the Czech business registry's public REST API across all of its source registers (VR, RES, RŽP, ROS, CEÚ insolvency, RÚIAN addresses), then adds a local SQLite layer the official API has no notion of: bulk-enrich a column of IČOs for invoicing, keep a watchlist of clients and suppliers, and diff their name/address/insolvency status over time.

Printed by [@patrikzita](https://github.com/patrikzita) (Patrik Zita).

## Install

The recommended path installs both the `ares-pp-cli` binary and the `pp-ares` agent skill (Claude Code, Codex, Cursor, Gemini CLI, GitHub Copilot, and other agents supported by the upstream [`skills`](https://github.com/vercel-labs/skills) CLI) in one shot:

```bash
npx -y @mvanhorn/printing-press-library install ares
```

For CLI only (no skill):

```bash
npx -y @mvanhorn/printing-press-library install ares --cli-only
```

For skill only — installs the skill into the same agents as the default command above, but skips the CLI binary (use this to update or reinstall just the skill):

```bash
npx -y @mvanhorn/printing-press-library install ares --skill-only
```

To constrain the skill install to one or more specific agents (repeatable — agent names match the [`skills`](https://github.com/vercel-labs/skills) CLI):

```bash
npx -y @mvanhorn/printing-press-library install ares --agent claude-code
npx -y @mvanhorn/printing-press-library install ares --agent claude-code --agent codex
```

### Without Node (Go fallback)

If `npx` isn't available (no Node, offline), install the CLI directly via Go (requires Go 1.26.3 or newer):

```bash
go install github.com/mvanhorn/printing-press-library/library/sales-and-crm/ares/cmd/ares-pp-cli@latest
```

This installs the CLI only — no skill.

### Pre-built binary

Download a pre-built binary for your platform from the [latest release](https://github.com/mvanhorn/printing-press-library/releases/tag/ares-current). On macOS, clear the Gatekeeper quarantine: `xattr -d com.apple.quarantine <binary>`. On Unix, mark it executable: `chmod +x <binary>`.

<!-- pp-hermes-install-anchor -->
## Install for Hermes

From the Hermes CLI:

```bash
hermes skills install mvanhorn/printing-press-library/cli-skills/pp-ares --force
```

Inside a Hermes chat session:

```bash
/skills install mvanhorn/printing-press-library/cli-skills/pp-ares --force
```

## Install for OpenClaw

Tell your OpenClaw agent (copy this):

```
Install the pp-ares skill from https://github.com/mvanhorn/printing-press-library/tree/main/cli-skills/pp-ares. The skill defines how its required CLI can be installed.
```

## Use with Claude Desktop

This CLI ships an [MCPB](https://github.com/modelcontextprotocol/mcpb) bundle — Claude Desktop's standard format for one-click MCP extension installs (no JSON config required).

To install:

1. Download the `.mcpb` for your platform from the [latest release](https://github.com/mvanhorn/printing-press-library/releases/tag/ares-current).
2. Double-click the `.mcpb` file. Claude Desktop opens and walks you through the install.

Requires Claude Desktop 1.0.0 or later. Pre-built bundles ship for macOS Apple Silicon (`darwin-arm64`) and Windows (`amd64`, `arm64`); for other platforms, use the manual config below.

<details>
<summary>Manual JSON config (advanced)</summary>

If you can't use the MCPB bundle (older Claude Desktop, unsupported platform), install the MCP binary and configure it manually.


```bash
go install github.com/mvanhorn/printing-press-library/library/sales-and-crm/ares/cmd/ares-pp-mcp@latest
```

Add to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "ares": {
      "command": "ares-pp-mcp"
    }
  }
}
```

</details>

## Quick Start

```bash
# Look up one company by IČO — no auth needed, ARES is a public registry.
ares-pp-cli ekonomicke-subjekty get 00177041


# Check an IČO checksum offline before spending an API call.
ares-pp-cli validate 00177041


# Search by company name.
ares-pp-cli ekonomicke-subjekty vyhledej --obchodni-jmeno Alza --pocet 5


# Turn a column of IČOs into enriched company records.
ares-pp-cli enrich 00177041 27082440 --agent

```

## Unique Features

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

## Usage

Run `ares-pp-cli --help` for the full command reference and flag list.

## Commands

### ciselniky-nazevniky

ciselniky-nazevniky operations

- **`ares-pp-cli ciselniky-nazevniky`** - Vyhledání číselníků používaných v IS ARES podle komplexního filtru

### ekonomicke-subjekty

ekonomicke-subjekty operations

- **`ares-pp-cli ekonomicke-subjekty vrat-ekonomicky-subjekt`** - Vyhledání ekonomického subjektu ARES podle zadaného iča
- **`ares-pp-cli ekonomicke-subjekty vyhledej`** - Vyhledání seznamu ekonomických subjektů ARES podle komplexního filtru

### ekonomicke-subjekty-ceu

ekonomicke-subjekty-ceu operations

- **`ares-pp-cli ekonomicke-subjekty-ceu vrat-ekonomicky-subjekt-ceu`** - Vyhledání konkrétního úpadce ze zdroje CEÚ
- **`ares-pp-cli ekonomicke-subjekty-ceu vyhledej-seznam-ekonomickych-subjektu-ceu`** - Vyhledání seznamu úpadců ze zdroje CEÚ

### ekonomicke-subjekty-notifikace

ekonomicke-subjekty-notifikace operations

- **`ares-pp-cli ekonomicke-subjekty-notifikace vrat-notifikacni-davku`** - Vyhledání  konkrétní notifikační dávky zdroje ARES podle vstupcách parametrů (zdroj, číslo notifikační dávky)
- **`ares-pp-cli ekonomicke-subjekty-notifikace vyhledej-seznam-notifikacnich-davek`** - Vyhledání seznamu notifikačních dávek ekonomických subjektů ARES podle zvoleného filtru

### ekonomicke-subjekty-nrpzs

ekonomicke-subjekty-nrpzs operations

- **`ares-pp-cli ekonomicke-subjekty-nrpzs vrat-ekonomicky-subjekt-nrpzs`** - Vyhledání konkrétního ekonomického subjektu ze zdroje NRPZS
- **`ares-pp-cli ekonomicke-subjekty-nrpzs vyhledej-seznam-ekonomickych-subjektu-nrpzs`** - Vyhledání seznamu ekonomických subjektu ze zdroje NRPZS

### ekonomicke-subjekty-rcns

ekonomicke-subjekty-rcns operations

- **`ares-pp-cli ekonomicke-subjekty-rcns vrat-ekonomicky-subjekt-rcns`** - Vyhledání konkrétního ekonomického subjektu ze zdroje RCNS
- **`ares-pp-cli ekonomicke-subjekty-rcns vyhledej-seznam-ekonomickych-subjektu-rcns`** - Vyhledání seznamu ekonomických subjektu ze zdroje RCNS

### ekonomicke-subjekty-res

ekonomicke-subjekty-res operations

- **`ares-pp-cli ekonomicke-subjekty-res vrat-ekonomicky-subjekt-res`** - Vyhledání konkrétního ekonomického subjektu ze zdroje RES
- **`ares-pp-cli ekonomicke-subjekty-res vyhledej-seznam-ekonomickych-subjektu-res`** - Vyhledání seznamu ekonomických subjektu ze zdroje RES

### ekonomicke-subjekty-ros

ekonomicke-subjekty-ros operations

- **`ares-pp-cli ekonomicke-subjekty-ros vrat-ekonomicky-subjekt-ros`** - Vyhledání konkrétního ekonomického subjektu ze zdroje ROS
- **`ares-pp-cli ekonomicke-subjekty-ros vyhledej-seznam-ekonomickych-subjektu-ros`** - Vyhledání seznamu ekonomických subjektu ze zdroje ROS

### ekonomicke-subjekty-rpsh

ekonomicke-subjekty-rpsh operations

- **`ares-pp-cli ekonomicke-subjekty-rpsh vrat-ekonomicky-subjekt-rpsh`** - Vyhledání konkrétního ekonomického subjektu ze zdroje RPSH
- **`ares-pp-cli ekonomicke-subjekty-rpsh vyhledej-seznam-ekonomickych-subjektu-rpsh`** - Vyhledání seznamu ekonomických subjektu ze zdroje RPSH

### ekonomicke-subjekty-rs

ekonomicke-subjekty-rs operations

- **`ares-pp-cli ekonomicke-subjekty-rs vrat-ekonomicky-subjekt-rs`** - Vyhledání konkrétního ekonomického subjektu ze zdroje RŠ
- **`ares-pp-cli ekonomicke-subjekty-rs vyhledej-seznam-ekonomickych-subjektu-rs`** - Vyhledání seznamu ekonomických subjektu ze zdroje RŠ

### ekonomicke-subjekty-rzp

ekonomicke-subjekty-rzp operations

- **`ares-pp-cli ekonomicke-subjekty-rzp vrat-ekonomicky-subjekt-rzp`** - Vyhledání konkrétního ekonomického subjektu ze zdroje RŽP
- **`ares-pp-cli ekonomicke-subjekty-rzp vyhledej-seznam-ekonomickych-subjektu-rzp`** - Vyhledání seznamu ekonomických subjektu ze zdroje RŽP

### ekonomicke-subjekty-szr

ekonomicke-subjekty-szr operations

- **`ares-pp-cli ekonomicke-subjekty-szr vrat-ekonomicky-subjekt-szr`** - Vyhledání konkrétního ekonomického subjektu ze zdroje SZR - subregistr EZP
- **`ares-pp-cli ekonomicke-subjekty-szr vyhledej-seznam-ekonomickych-subjektu-szr`** - Vyhledání seznamu ekonomických subjektu ze zdroje SZR - subregistr EZP

### ekonomicke-subjekty-vr

ekonomicke-subjekty-vr operations

- **`ares-pp-cli ekonomicke-subjekty-vr vrat-ekonomicky-subjekt-vr`** - Vyhledání konkrétního ekonomického subjektu ze zdroje VR
- **`ares-pp-cli ekonomicke-subjekty-vr vyhledej-seznam-ekonomickych-subjektu-vr`** - Vyhledání seznamu ekonomických subjektu ze zdroje VR

### standardizovane-adresy

standardizovane-adresy operations

- **`ares-pp-cli standardizovane-adresy`** - Vyhledání seznamu standardizovaných adres RÚIAN podle komplexního filtru


## Output Formats

```bash
# Human-readable table (default in terminal, JSON when piped)
ares-pp-cli ciselniky-nazevniky

# JSON for scripting and agents
ares-pp-cli ciselniky-nazevniky --json

# Filter to specific fields
ares-pp-cli ciselniky-nazevniky --json --select id,name,status

# Dry run — show the request without sending
ares-pp-cli ciselniky-nazevniky --dry-run

# Agent mode — JSON + compact + no prompts in one flag
ares-pp-cli ciselniky-nazevniky --agent
```

## Agent Usage

This CLI is designed for AI agent consumption:

- **Non-interactive** - never prompts, every input is a flag
- **Pipeable** - `--json` output to stdout, errors to stderr
- **Filterable** - `--select id,name` returns only fields you need
- **Previewable** - `--dry-run` shows the request without sending
- **Explicit retries** - add `--idempotent` to create retries when a no-op success is acceptable
- **Confirmable** - `--yes` for explicit confirmation of destructive actions
- **Piped input** - write commands can accept structured input when their help lists `--stdin`
- **Agent-safe by default** - no colors or formatting unless `--human-friendly` is set

Exit codes: `0` success, `2` usage error, `3` not found, `5` API error, `7` rate limited, `10` config error.

## Health Check

```bash
ares-pp-cli doctor
```

Verifies configuration and connectivity to the API.

## Configuration

Config file: `~/.config/ares-verejne-pp-cli/config.toml`

Static request headers can be configured under `headers`; per-command header overrides take precedence.

## Troubleshooting
**Not found errors (exit code 3)**
- Check the resource ID is correct
- Run the `list` command to see available items

### API-specific

- **HTTP 429 / throttled during bulk enrich** — ARES allows ~500 req/min; enrich and refresh throttle automatically, lower concurrency with --concurrency if you still hit limits.
- **search returns nothing** — Run ares-pp-cli sync first — search reads the local store, not the live API.

---

## Sources & Inspiration

This CLI was built by studying these projects and resources:

- [**ares-mcp-server**](https://github.com/vzeman/ares-mcp-server) — Python
- [**czech-company-registry-api**](https://github.com/XrayHunter/czech-company-registry-api) — TypeScript
- [**ares-cz**](https://github.com/borisgrigorov/ares-cz) — JavaScript
- [**ares-api-tool**](https://github.com/mnamnau/ares-api-tool) — Python
- [**ares_util**](https://github.com/illagrenan/ares_util) — Python
- [**dfridrich/Ares**](https://github.com/dfridrich/Ares) — PHP

Generated by [CLI Printing Press](https://github.com/mvanhorn/cli-printing-press)
