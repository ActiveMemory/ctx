# Source Map Schema

Shape for `.context/kb/source-map.md`. One row per source
cited from `evidence-index.md`. Records what a source *is*
and whether it was admitted against the kb's scope.

Distinct from `source-coverage.md`: this file records WHAT a
source is (identity, locator, admission); the coverage ledger
records HOW COMPLETE the kb's work against it is (state
machine over discovered → comprehensive).

## Fields

| Column                 | Description                                                  |
|------------------------|--------------------------------------------------------------|
| `short-name`           | Stable identifier used in `evidence-index.md` cites.         |
| `kind`                 | One of `transcript`, `code`, `doc`, `url`, `mcp`, `kb`.      |
| `locator`              | URL, repo path, MCP resource ID, or kb pointer.              |
| `admission-status`     | `admitted`, `rejected`, or `pending`.                        |
| `admission-rationale`  | One-sentence reason the source was admitted or rejected.     |
| `dated`                | Optional ISO date the source itself is dated.                |

## Example

| short-name   | kind | locator                                | admission-status | admission-rationale                              | dated      |
|--------------|------|----------------------------------------|------------------|--------------------------------------------------|------------|
| WIDGET-SPEC  | url  | `https://example.org/widget/v2/spec`   | admitted         | Canonical upstream spec; in scope per index.md.  | 2026-05-16 |

## Rules

- Short names are stable; rename via alias-add, not by
  rewriting in place. Existing `EV-###` rows pin to the
  original short name.
- Admission is against the kb's declared `## Scope` in
  `.context/kb/index.md`; rejected sources stay in the map
  with their rationale as audit trail.
- `kind: kb` marks federation into another kb (per the
  KB-of-KBs organizing principle); the locator points at
  another `.context/kb/` directory.
- A source can appear here at `admission-status: pending`
  before it is admitted; the coverage ledger picks up only
  admitted sources.
- Do not conflate this file with `source-coverage.md`; the
  identity-vs-progress split is load-bearing for the
  state-machine ledger.
