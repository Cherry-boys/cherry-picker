# Cherrypick

Internal business system for **Cherry-boys** — find leads, manage prospects, drive sales, and create offers from one place.

## What it does

- **Lead finder** — discover and qualify potential customers
- **Prospects** — track and manage the prospect pipeline
- **Sales** — move deals from first contact to close
- **Offers** — generate and send offers/quotes

## Stack

Turborepo monorepo (Bun) with LangChain / LangGraph agent tooling.

- `apps/agent` — AI agent backend (LangGraph)
- `apps/web` — web frontend
- `packages/*` — shared config (eslint, typescript)
- `tools/ares` — Go CLI for the Czech ARES business registry (see below)

## Getting started

```bash
bun install
bun run dev
```

## `tools/ares` — ARES business-registry CLI

A Go CLI that wraps the Czech **ARES** company registry: lookups by IČO across
VR/RES/RŽP/ROS, plus a local SQLite store, offline full-text search, bulk IČO
enrichment, portfolio change-tracking, and insolvency watch. Handy for qualifying
and enriching Czech prospects/suppliers.

This is a vendored copy of a [printing-press](https://github.com/mvanhorn/printing-press-library)–generated
CLI. It is **not** part of the Bun/Turbo workspace, so `bun run build` never touches
Go — you build it on demand.

**Prerequisite:** Go 1.26+.

```bash
bun run ares:build      # -> tools/ares/bin/ares-pp-cli
tools/ares/bin/ares-pp-cli --help
tools/ares/bin/ares-pp-cli doctor   # check auth + connectivity

bun run ares:install    # go install onto your $PATH
bun run ares:test       # go test ./...
```

Optional MCP server (for Claude Desktop): `make -C tools/ares build-mcp` → `tools/ares/bin/ares-pp-mcp`.

> Upstream lives at `~/printing-press/library/ares`. To pull in upstream changes,
> re-vendor the source (cmd/, internal/, go.mod/sum, spec.json, Makefile, etc.).
