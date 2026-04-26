---
summary: PoC design plan for the draft-generation agent (Phase 1) — decisions from design session 2026-04-26
read_when: Before implementing any part of Phase 1 agent, proAlpha mock, or RAG pipeline
---

# PoC Agent Design Plan

> **Status:** 🟡 Design agreed — not yet implemented
> **Date:** 2026-04-26
> **Scope:** Phase 1 (deterministic draft generation) + supporting infrastructure for PoC

---

## What We're Building

An agent that takes two inputs:

1. **Customer PDF drawing** (e.g. Belimed AG technical drawing, in German)
2. **Technologist's free-text prompt** (Czech, optional context not visible on the drawing)

…and produces a **first-draft working procedure** (`Formblatt CiS-0038`) ready for Phase 2 conversational refinement.

---

## Phase 1 Pipeline (LangGraph, deterministic)

```
upload PDF + technologist prompt
        │
   pdf_extractor          ← Claude Vision → structured drawing_info
        │
   bom_lookup             ← proAlpha mock API → resolved CiS material codes
        │
   rag_retriever          ← vector search → similar historical procedures
        │
   drafter                ← LLM (Claude) + drawing_info + bom + similar_procedures → draft_operations
        │
   ⏸ INTERRUPT            ← technologist reviews the draft
```

Phase 2 (conversational refinement + Word export) is out of scope for this PoC.

---

## Input: Technologist's Prompt

The technologist writes free text in Czech alongside the drawing. This captures information
**not present on the drawing** that improves the first draft:

| What technologist writes | Why it helps |
|---|---|
| Revision context: "Index 03→04, length 2000→2500mm, rest unchanged" | Agent only needs to update the diff, not re-derive everything |
| Special operations: "Add heat shrink 8mm after crimping on both ends" | Adds an operation the drawing doesn't show |
| Packaging: "Pack 5 pcs per bag" | Fills operation 080 (Balení) parameters |
| Production context: "Prototype run, 10 pcs, no soak test" | Removes or marks operations as optional |
| Customer requirement: "Belimed requires 100% pull test on all contacts" | Adds a QC step not on the drawing |

**Key insight:** ~90% of the draft comes from the drawing + RAG. The technologist's prompt covers
the ~10% that requires human knowledge.

---

## RAG Strategy (Option A — whole-procedure embeddings)

- **What we embed:** each historical working procedure as a single document
- **Similarity search on:** product type, connector family, wire count, cable cross-section
- **Result:** 2–3 most similar past procedures returned to the `drafter` node as context
- **Rationale:** Simpler to implement for PoC; good enough to demonstrate the pattern
- **Future upgrade (Option B):** embed individual operations for more granular retrieval

**PoC seed data:**
- 1 real procedure (`assets/output-working-procedure.docx`)
- 2–3 synthetic procedures covering different cable types

---

## proAlpha Mock — Separate App

### Why a separate app

The real proAlpha is an external service. Building a mock as a standalone HTTP server means:
- The agent calls it via HTTP tool — same as it would call real proAlpha
- Swapping mock → real is just a base URL + credentials change; agent code doesn't change

### Location

```
apps/proalpha-mock/
```

### Tech stack

- **Runtime:** Bun
- **HTTP framework:** Hono
- **Database:** bun:sqlite (SQLite, built-in)

### API pattern (mirroring real proAlpha REST API)

The real proAlpha REST API uses an **async job pattern**:

```
# Step 1 — submit a job
POST /api/targets/{target_id}/jobs/
Authorization: Bearer {token}
{ "api_name": "GetMaterial", "params": { "belimed_nr": "18301" } }
→ { "id": "job-abc", "status": "pending" }

# Step 2 — poll result
GET /api/targets/{target_id}/jobs/{job_id}
→ { "status": "completed", "result": { "cis_code": "2601086", "name": "...", "unit": "ks" } }
```

For PoC the mock can resolve jobs synchronously (status always `completed` immediately),
but keeping the same URL structure means the agent tool works against real proAlpha without changes.

### Named API operations (minimum for PoC)

| api_name | Input | Output |
|---|---|---|
| `GetMaterial` | `belimed_nr` or `tyco_nr` | CiS material code, name, unit, description |
| `GetBOM` | `cis_product_code` | flat list of `{cis_code, name, qty, unit}` |

### SQLite schema (proposed)

```sql
-- Component catalog: maps customer part numbers → CiS internal codes
CREATE TABLE materials (
  id INTEGER PRIMARY KEY,
  cis_code TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  unit TEXT NOT NULL,               -- ks, m, bal, ...
  description TEXT,
  belimed_nr TEXT,                  -- customer-facing number (Belimed)
  tyco_nr TEXT                      -- supplier number (Tyco/TE Connectivity)
);

-- Bill of materials: which materials go into which CiS product
CREATE TABLE bom_entries (
  id INTEGER PRIMARY KEY,
  cis_product_code TEXT NOT NULL,
  material_cis_code TEXT NOT NULL,
  qty REAL NOT NULL,
  unit TEXT NOT NULL,
  FOREIGN KEY (material_cis_code) REFERENCES materials(cis_code)
);
```

### Seed data (from drawing 74678)

| belimed_nr | tyco_nr | cis_code | name |
|---|---|---|---|
| 18301 | 172157-1 | TBD | Aufnahmegehäuse 6.3mm 2-pol transparent |
| 54770 | 170362-1 | 2601086 | MATE-N-LOK Buchse 0.30-0.89mm² |
| 18286 | 966156-1 | TBD | Aderendhülse Ø0.5mm×14mm DIN 46228 |
| — | — | 5100560 | Kabel 0.5mm² (schwarz/braun) |
| — | — | 9100008 | Sáček (packaging bag) |

CIS codes marked TBD need to be filled in by the technologist.

---

## PoC Assumptions

These are simplifications explicitly chosen for the first run.
Each is a known gap to be addressed in later iterations.

| Assumption | Reality | Plan |
|---|---|---|
| proAlpha mock returns static data | Real proAlpha is live ERP | Replace mock with real adapter |
| 1 real + synthetic historical procedures | Up to 1800 real procedures | Bulk import from CiS file storage |
| Whole-procedure embeddings (Option A) | Per-operation embeddings (Option B) may be more accurate | Evaluate after first run |
| Output is structured JSON (not Word) | Final output must be `.docx` (Formblatt CiS-0038) | Word generation is Phase 2 |
| Technologist prompt is optional free text | May need structured fields for better consistency | Evaluate with real technologists |
| CiS has proAlpha REST API module enabled | Unknown — not confirmed | Verify with CiS before production integration |

---

## Open Questions

See [docs/internal/open-questions.md](../internal/open-questions.md) for the full list.

Critical blockers before implementation:

1. **Exact fields of Formblatt CiS-0038** — do operations have fields beyond `{code, name, side, time_s, steps[], materials[]}`?
2. **proAlpha API availability** — does CiS actually have the REST API module purchased?
3. **Material lookup key** — does the technologist search by Belimed-Nr, Tyco-Nr, or free text?

---

## What's Not in This Plan

- Phase 2 (conversational refinement)
- Word document generation (`.docx` export)
- Frontend UI
- Authentication / multi-user
- Revision management (diff between drawing index versions)
