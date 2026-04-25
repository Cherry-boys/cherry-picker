---
summary: "Two-phase agent design: deterministic draft generator (Phase 1) and conversational refinement agent (Phase 2) that together produce a finalized work procedure."
read_when:
  - designing or modifying the draft-generation pipeline
  - adding tools to the conversational agent (Phase 2)
  - unclear how Phase 1 hands off state to Phase 2
  - implementing the Word document export step
  - deciding what the agent should do automatically vs. ask the technologist about
---

# Agent Architecture — Two-Phase Design

The assistant is a two-phase system. Phase 1 produces a structured draft deterministically from the uploaded drawing. Phase 2 lets the technologist refine that draft through natural-language conversation before exporting to Word.

---

## Phase 1 — Deterministic Draft Generator

**Trigger:** Technologist uploads a PDF drawing.  
**Goal:** Produce an ~80% complete work procedure draft with no user interaction.  
**Character:** Linear, deterministic — same input should yield the same draft.

### Steps (in order)

| Step | What happens |
|------|-------------|
| 1. Drawing extraction | Parse the PDF visually; extract wire list, lengths, cross-sections, colors, connectors, pin maps, certifications, notes. Handles German text. |
| 2. Material matching (proAlpha) | For each component in the drawing, look up the matching material in proAlpha. Unambiguous matches are auto-linked; ambiguous ones are flagged `[TO VERIFY]`. |
| 3. Operation sequencing | Query the historical procedure database for similar products. Propose an operation sequence (cutting → crimping → assembly → QC → packaging) with estimated times derived from historical analogues. |

### Output

A structured draft containing:
- Proposed operation list with times
- Material references from proAlpha (with unresolved items flagged)
- Source drawing context preserved for Phase 2 lookups

The draft is surfaced to the technologist as-is — no approval gate before Phase 2 begins.

---

## Phase 2 — Conversational Refinement Agent

**Trigger:** Technologist starts chatting after reviewing the Phase 1 draft.  
**Goal:** Iteratively update the draft based on natural-language instructions until the technologist says it's done.  
**Character:** Flexible, stateful — maintains the current draft across the conversation.

### What the technologist can say

| Intent | Example |
|--------|---------|
| Correct operation time | *"Operace 030 má špatný čas, dej 15 sekund"* |
| Add an operation | *"Přidej operaci pro smršťovačku před balení"* |
| Change a material | *"Ten konektor je špatně, použij raději číslo 2601099"* |
| Ask about history | *"Jak vypadal postup u podobného výrobku z loňska?"* |
| Verify consistency | *"Zkontroluj jestli sedí všechny part numbers s výkresem"* |

### What the agent does

- Applies each change to the live draft
- Re-queries proAlpha or the historical procedure DB when needed
- Flags inconsistencies (drawing vs. BOM vs. procedure)
- Shows the current draft state after each update

### End condition

Technologist says they're satisfied and requests document generation.

---

## Document Export

After Phase 2 sign-off the system generates a Word document:
- Template: `Formblatt CiS-0038`
- Contents: filled header, operations table, proAlpha material references
- Technologist downloads and hands to production

---

## End-to-End Flow

```
1. Upload PDF drawing
        ↓
2. Phase 1 runs (~1 min) → structured draft delivered
        ↓
3. Phase 2: chat-based refinement (iterate until satisfied)
        ↓
4. "Generate" command → Word document downloaded
        ↓
5. Hand to production
```

---

## Automation Boundary

| Task | Phase 1 (auto) | Phase 2 (human-in-loop) |
|------|---------------|------------------------|
| Read and interpret drawing (DE) | Yes | — |
| Match materials in proAlpha | Yes (flags ambiguous) | Technologist resolves flags |
| Propose operation sequence | Yes | Technologist corrects/extends |
| Estimate operation times | Yes | Technologist overrides |
| Write Word document | Yes (on demand) | — |
| Consistency check | Yes | Technologist reviews warnings |

---
