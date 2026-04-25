---
summary: "AI assistant for CiS systems technologists — automates creation of work procedures (pracovní postupy) from customer PDF drawings"
read_when:
  - starting work on any feature or agent capability
  - unclear what the end-to-end user workflow looks like
  - adding integration with proAlpha ERP or Word output
  - working on drawing extraction, BOM matching, or operation sequencing
---

# Project Overview

**CiS systems s.r.o.** (Nové Město pod Smrkem) manufactures cable harnesses and electromechanical components for international industrial customers — up to 250 000 units/week across 1 800 product variants.

## Problem

The **technology department** (technologists) manually converts customer drawings into production work procedures. The full process is:

1. Receive customer PDF drawing (typically in German)
2. Read and interpret it (dimensions, components, pin mapping, tolerances)
3. Look up / assemble the bill of materials (BOM) in **proAlpha** ERP
4. Open a Word template (`Formblatt CiS-0038`) and manually write the ordered operation sequence with timing
5. Hand the document to production

This is entirely manual, error-prone, slow, and dependent on senior technologist tribal knowledge.

## Goal

Build an AI agent that assists technologists in producing a validated work procedure from a customer drawing — covering:

- **Drawing extraction** — parse German PDF drawings (components, dimensions, pin table, notes)
- **BOM lookup** — match drawing components to proAlpha material numbers
- **Operation sequencing** — derive the standard sequence (cut → crimp → assemble → inspect → pack)
- **Time estimation** — based on historically similar procedures
- **Word generation** — produce output conforming to `Formblatt CiS-0038`
- **Consistency check** — verify drawing ↔ BOM ↔ procedure alignment
- **Revision management** — diff against previous procedure version when drawing index changes

## Key Actors

| Actor | Role |
|---|---|
| Technologist | Primary user; drives the procedure creation workflow |
| proAlpha (ERP/PPS) | Source of BOM / material numbers |
| Customer | Provides the PDF drawing (e.g. Belimed AG) |
| Production | Consumes the finished work procedure |

## Output Format

Word document per template `Formblatt CiS-0038` containing:
- Header (harness type, revision, date, article, customer)
- Ordered operations table with codes (`010`, `020`, …), names, side, time (seconds), step descriptions, and proAlpha material references
