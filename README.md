# CiS Agent

AI assistant for **CiS systems s.r.o.** that automates the creation of production work procedures (*pracovní postupy*) from customer technical drawings.

## What is this?

CiS manufactures cable harnesses and electromechanical components for international industrial customers. Technologists currently convert customer PDF drawings (often in German) into Word work procedure documents by hand — reading the drawing, assembling a bill of materials in the proAlpha ERP, and writing the ordered operation sequence manually.

This project builds an AI agent that handles that pipeline end-to-end.

## Current state (PoC)

The agent is a [LangGraph.js](https://langchain-ai.github.io/langgraphjs/) graph with two nodes running in sequence:

```
pdf_extractor → bom_lookup
```

**`pdf_extractor`** — sends the PDF to Claude claude-sonnet-4-6 via the Anthropic API and extracts the drawing number from the technical drawing.

**`bom_lookup`** — calls the proAlpha REST API with the drawing number and retrieves the bill of materials (BOM) for that product.

The graph state carries:

| Field | Type | Description |
|---|---|---|
| `pdf_path` | `string` | Path to the input customer PDF drawing |
| `drawing_number` | `string \| null` | Extracted from the PDF |
| `bom` | `BomEntry[] \| null` | Materials returned from proAlpha |

## Monorepo structure

```
apps/
  agent/           # LangGraph.js agent (the main AI pipeline)
  proalpha-mock/   # Mock proAlpha REST API backed by SQLite (for local dev)
  web/             # Next.js app (placeholder)
packages/
  eslint-config/
  typescript-config/
assets/
  input_drawing.pdf          # Sample customer drawing for testing
  output-working-procedure.docx  # Example output document (target format)
docs/
  project-overview.md        # Domain context and full goal description
  agent-architecture.md      # Agent design notes
  plans/poc-agent-design.md  # PoC design decisions
```

## What's planned but not built yet

- Operation sequencing (cut → crimp → assemble → inspect → pack)
- Time estimation based on historically similar procedures
- Word document generation conforming to `Formblatt CiS-0038`
- Consistency check (drawing ↔ BOM ↔ procedure alignment)
- Revision management (diff against previous procedure when drawing index changes)
- Frontend UI for technologists

## Running locally

```sh
# Install dependencies
bun install

# Start the mock proAlpha API (port 3001)
cd apps/proalpha-mock
bun run dev

# Start the LangGraph dev server (agent)
cd apps/agent
bun run agent
```

The agent requires `ANTHROPIC_API_KEY` set in `apps/agent/.env` (see `.env.example`).

## Tech stack

- **LangGraph.js** — agent orchestration
- **Anthropic Claude** — PDF parsing and AI reasoning
- **Hono + SQLite** — mock proAlpha ERP REST API
- **Bun** — runtime and package manager
- **Turborepo** — monorepo build orchestration
- **TypeScript** throughout
