---
name: create-doc
description: Creates a new markdown documentation file in the docs/ folder with proper AI-readable frontmatter (summary and read_when). Use when the user explicitly asks to document something — a workflow, deployment process, project structure, architecture decision, integration, or any knowledge an AI coding agent should be able to look up.
disable-model-invocation: true
---

You are creating a new documentation file for the `docs/` folder. These docs are consumed by an AI coding agent that runs `bun scripts/docs-list.ts` before working on a task, scanning summaries and `read_when` hints to decide what to read. Write for that agent, not for a human tutorial reader.

## Step 1: Survey what already exists

```bash
bun scripts/docs-list.ts
```

Review the output to avoid duplicating existing coverage. If the user's topic is partially covered, focus on the gap.

## Step 2: Clarify scope if needed

If the user hasn't specified what to document, suggest 3–4 concrete options based on what's missing. Good candidates for this project:

- **Workflows** — end-to-end process flows (drawing extraction → BOM lookup → procedure generation)
- **Project structure** — where code lives, what each directory and entry point does
- **Deployment** — how to run, configure, and deploy the agent locally or in production
- **Integrations** — external systems (proAlpha ERP, Word template engine, PDF parser)
- **Architecture** — how AI components connect and hand off work to each other
- **Data formats** — shapes of key inputs/outputs (drawing schema, procedure schema, BOM structure)
- **Conventions** — naming rules, patterns, decisions that are non-obvious from the code

If the intent is clear, skip straight to Step 3.

## Step 3: Explore the codebase

Don't write from memory or assumptions — read the actual code for the topic. Use `find`, `ls`, and file reads to understand:
- Entry points and main modules involved
- Key interfaces, types, or schemas
- Configuration and environment variables
- How data flows through the relevant components

## Step 4: Write the file

Place the file at `docs/<kebab-case-name>.md`. For focused subtopics, use a subdirectory: `docs/workflows/drawing-extraction.md`.

Do not place files inside `archive/`, `research/`, `plans/`, or `internal/` — those are excluded from the agent's doc listing.

Every doc must open with this frontmatter block:

```yaml
---
summary: "One precise sentence describing what knowledge the agent gains from reading this."
read_when:
  - specific scenario where this doc is relevant
  - another specific scenario
  - ...
---
```

### Writing a good `summary`

The summary is the only thing the agent sees before deciding whether to open the file. Make it count:

- Say what knowledge the doc *contains*, not what the topic is
- One sentence, under ~20 words
- No filler ("comprehensive guide to...", "overview of...")

| Bad | Good |
|-----|------|
| `"Overview of the deployment process"` | `"Steps to run the agent locally with Docker and configure proAlpha and LLM credentials"` |
| `"About the project structure"` | `"Directory layout and the role of each top-level package: agent, api, worker, scripts"` |

### Writing good `read_when` hints

Each hint is a concrete situation the agent might be in. It should be specific enough for the agent to match it to its current task.

| Bad | Good |
|-----|------|
| `"when working on deployment"` | `"setting up a local dev environment or configuring environment variables"` |
| `"when touching the agent"` | `"adding a new step to the drawing extraction pipeline or modifying agent tools"` |

Aim for 3–6 hints. Cover the main scenarios where an agent would genuinely need this doc, including non-obvious ones.

### Body content

Write for an AI reader. Prefer:
- **Tables** for structured data (file → purpose, env var → default → meaning)
- **Code blocks** for commands, file paths, schemas, example values
- **Short headed sections** so the agent can scan and skip to the relevant part
- **Concrete examples** over abstract descriptions

Avoid long narrative prose. State facts directly.

## Step 5: Verify

Run `bun scripts/docs-list.ts` again and confirm the new file appears with its summary and `read_when` hints rendered correctly. Fix any frontmatter formatting issues if it shows an error.
