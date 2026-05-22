# Source Coverage Schema

Shape for `.context/kb/source-coverage.md`. Per-source
completeness ledger maintained by `/ctx-kb-ingest`. Each row
tracks **coverage state**, not just existence: a source can be
known to the kb at any of several stages between discovery and
comprehensive understanding.

This file is the canonical answer to *"what is incomplete?"*
Distinct from `source-map.md` (which records WHAT a source
is); the ledger records HOW COMPLETE the work against it is.

## Fields

| Column        | Description                                                            |
|---------------|------------------------------------------------------------------------|
| `Source`      | Short name from `source-map.md`.                                       |
| `Topic`       | Kb topic slug, or `n/a` if no topic page applies (e.g. `triage` runs). |
| `State`       | One of the state-machine values listed in the transitions table.       |
| `EV coverage` | Range like `EV-018..EV-034`, comma list, or `none`.                    |
| `Residue`     | One-line free text describing what was deliberately left out.          |
| `Next action` | Exact resumption invocation (e.g. `/ctx-kb-ingest <slug> (resume ...)`). |
| `Updated`     | ISO date the row was last written.                                     |

## Allowed Transitions

| state                   | next states                                                                                  |
|-------------------------|----------------------------------------------------------------------------------------------|
| `discovered`            | `admitted`, `skipped`                                                                        |
| `admitted`              | `highlights-extracted`, `partially-ingested`, `topic-page-drafted`, `comprehensive`          |
| `highlights-extracted`  | `partially-ingested`, `topic-page-drafted`, `comprehensive`                                  |
| `partially-ingested`    | `topic-page-drafted`, `comprehensive`                                                        |
| `topic-page-drafted`    | `comprehensive`                                                                              |
| `comprehensive`         | terminal until source updates                                                                |
| `superseded`            | terminal                                                                                     |
| `skipped`               | terminal until scope changes                                                                 |

Backward steps (e.g. `comprehensive → highlights-extracted`)
require an explicit `superseded` step first; otherwise the
transition is illegal and `AdvanceLedger` refuses it.

## Example

| Source       | Topic   | State                  | EV coverage     | Residue                            | Next action                              | Updated    |
|--------------|---------|------------------------|-----------------|------------------------------------|------------------------------------------|------------|
| WIDGET-SPEC  | widgets | highlights-extracted   | EV-042..EV-051  | composition examples, error cases  | /ctx-kb-ingest widgets (resume topic-page) | 2026-05-16 |

## Rules

- Lying to the ledger is a hard anti-pattern; record the
  state honestly even when it means writing
  `topic-page-drafted` instead of `comprehensive`.
- Illegal transitions are refused at write time; backing out
  of `comprehensive` requires an explicit `superseded` step
  first, naming the superseder.
- `comprehensive` requires cited bands ≥ `medium`, no
  `TBD-cite` markers, and a passing cold-reader rubric;
  topic-page-drafted is the right state when any of those
  three conditions fail.
- `Next action` is the exact resumption invocation a future
  agent should run; free prose like "finish later" defeats
  the ledger's purpose.
- Every pass that touches a source updates its row *before*
  writing the closeout; doctor advisory flags rows whose
  `Updated` predates the file's last-modified time.
