---
summary: "The north star of this repo: it is Cherry-boys' internal operating brain; current sole focus is the sales/networking engine, with other layers explicitly parked or out of scope."
read_when:
  - starting any new feature, app, package, or tool and unsure whether it belongs in this repo
  - deciding between building infrastructure vs. shipping something that drives sales
  - tempted to flesh out apps/web, packages/db, or apps/agent
  - figuring out which modality to use (.md wiki vs. Go CLI vs. LangChain/LangGraph agent)
  - unsure whether client automation work belongs here or elsewhere
  - writing or reviewing docs and need the single source of truth for "what is this repo for"
---

# Repository Vision — Cherry-boys' Internal Brain

This repo is the **internal operating system ("brain") of Cherry-boys** — a small
founding team (currently 1–2 devs) building a company around **AI automation for
mid-sized companies**.

Two things to hold at once:

- **The company thesis is not fixed.** "AI automation for mid-sized companies" is
  the current bet, not a contract. The repo must stay flexible enough to follow a
  pivot, not lock us into one shape.
- **Broad north star, narrow current focus.** The long-term vision below is wide.
  What we actually build *right now* is deliberately narrow. Vision ≠ priority.

> Founder-trap warning, on purpose: our real bottleneck is **networking and sales**,
> not tooling. Devs find it easy to build infrastructure and hard to sell. Every
> addition to this repo should be checked against the question: *does this help us
> talk to potential customers, or is it a comfortable way to avoid it?*

---

## Current focus (the only active priority)

**The sales / networking engine.** Find and reach mid-sized Czech companies, enrich
them, and track outreach.

- Capture leads (companies/people worth contacting).
- Enrich Czech companies via `tools/ares` (the ARES business-registry CLI).
- Track outreach state (who, when, next step) as lightweight structured files.
- Use agents to assist: lead discovery, qualification, drafting outreach.

**Measured by:** number of companies reached and meetings booked — not lines of code
or features shipped.

Until the sales engine is working, the other north-star layers below are **not built**.

---

## North star (what this repo may eventually hold)

These are in-scope long-term, in roughly this priority order. Do **not** start them
before the current focus is delivering.

| Layer | What it is | Status |
|---|---|---|
| **Sales / networking engine** | Lead capture, enrichment, outreach tracking | **Active** (current focus) |
| **Internal knowledge / wiki** | SOPs, decisions, notes, vision — prose in `.md` (`docs/` = repo-meta, `wiki/` = business/domain knowledge) | Active as needed (this file is part of it) |
| **Reusable building blocks** | General tools & agent modules recyclable across uses (e.g. `tools/ares`, ERP-connector / PDF-extraction patterns) | Emerging |
| **Structured CRM / pipeline** | Deals, stages, contacts, follow-ups with real persistence + UI | **Parked** (future) |

---

## Out of scope (explicitly not goals of this repo)

Naming these prevents quiet scope creep. If you find yourself building one of these,
stop and re-read this section.

| Not here | Why / where instead |
|---|---|
| **Client production deliverables / portfolio** | Real client deployments (e.g. CiS-style automations) live in their **own repos**, not here. This repo is our internal brain, not a client codebase. |
| **Task / project management** | Use an existing tool (Linear/Jira/…). We don't build our own. |
| **Finance / ops back-office** | Invoicing, expenses, time tracking — out. |

---

## How we build

### Modalities
The repo deliberately mixes three ways of holding capability — pick the lightest one
that fits:

| Modality | Use for | Example |
|---|---|---|
| **`.md` knowledge** | Prose. `docs/` = repo-meta (vision, architecture); `wiki/` = business/domain knowledge (ICP, sales playbook, SOPs) — boundary test: *"would it matter even without this repo?"* yes → `wiki/`, no → `docs/` | `docs/`, `wiki/` |
| **`.md` operational records** | Structured per-entity files driven by CLIs/agents — one folder per company, state in frontmatter | `companies/` |
| **Go CLIs** | Reusable, scriptable tooling | `tools/ares` |
| **LangChain / LangGraph agents** | AI automation & assistance | `apps/agent` (template only — see below) |

### Medium for now: markdown + CLI + agents
Structured data (e.g. the sales pipeline) lives as **files in git** (markdown / simple
structured files), driven by **CLIs and agents**. No web app, no database yet — the
team is 1–2 devs and everything is dev-facing and git-native.

A real database + UI become worth it only when (a) volume demands it, or (b) a
non-dev needs access. Until then they stay parked.

---

## Status of the repo right now

What physically exists today, and how to treat it:

| Path | Role | Treat as |
|---|---|---|
| `tools/ares` | Czech ARES business-registry CLI — lead enrichment | **Active**, core to the sales engine |
| `docs/` | Repo-meta knowledge — vision, architecture | **Active** |
| `wiki/` | Business/domain knowledge — ICP, sales playbook, SOPs (flat, grows as needed) | **Active** |
| `companies/` | Operational sales-engine records — one folder per company (research, people, meeting notes, outreach state in frontmatter) | **Active**, core to the sales engine |
| `apps/agent` | LangGraph scaffold (carries leftover CiS POC content) | **Template only** — ignore its contents; it's a starting shape, not a deliverable |
| `apps/web` | Next.js scaffold | **Parked / template** — part of the future CRM-UI layer |
| `packages/db` | Postgres + Prisma scaffold | **Parked / template** — part of the future structured-pipeline layer |
| `packages/*` (eslint, ts config) | Shared monorepo config | Active plumbing |

> Note on `apps/`: the current contents are **template scaffolding**, not meaningful
> product. Don't mine them, don't maintain their business logic, don't treat the
> CiS-drawing workflow as a goal of this repo.

---

## Principles

1. **Sales first.** If it doesn't help us reach customers, it waits.
2. **Start simple, record the vision.** Build the minimum; write down where it could go (this doc).
3. **Broad north star, narrow focus.** Aspire wide, ship narrow.
4. **The thesis can change.** Keep the repo flexible; don't over-fit to today's bet.
5. **Don't let docs lie.** If reality diverges from this file, update the file.
