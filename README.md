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

## Getting started

```bash
bun install
bun run dev
```
