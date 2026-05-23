# Evidence Index Schema

Shape for `.context/kb/evidence-index.md`. Each row is one
atomic claim, minted by `/ctx-kb-ingest` and never renumbered.
The evidence index is the spine of the kb: glossary entries,
contradictions, domain-decisions, timeline events, and topic
pages all cite back to specific `EV-###` IDs here.

Rows are append-only. To retire a stale claim, demote its
`confidence` band per the demotion policy in `KB-RULES.md`.
Never delete the row.

## Fields

| Column       | Description                                                |
|--------------|------------------------------------------------------------|
| `id`         | Stable `EV-###` identifier; zero-padded; never renumbered. |
| `claim`      | One-sentence atomic claim; declarative, source-backed.     |
| `source`     | Short name from `source-map.md` (e.g. `CURSOR-HOOKS`).     |
| `locator`    | Line range, timestamp, anchor, or symbol within source.    |
| `sha`        | Optional short SHA; required only for in-repo citations.   |
| `confidence` | One of `high`, `medium`, `low`, `speculative`.             |
| `tags`       | Comma-separated; `evidence-only` is additive (mode tag).   |
| `occurred`   | Optional ISO date the underlying event occurred.           |
| `extracted`  | ISO date the row was extracted into the kb.                |

## Example

| id     | claim                                                                                | source       | locator   | sha     | confidence | tags                  | occurred   | extracted  |
|--------|--------------------------------------------------------------------------------------|--------------|-----------|---------|------------|-----------------------|------------|------------|
| EV-042 | The widget contract pins its schema version in the first frontmatter line of `SKILL.md`. | WIDGET-SPEC  | §Schema   |         | high       | widget, contract      |            | 2026-05-16 |

## Rules

- Append-only. Renumbering or deleting an `EV-###` row is a
  hard anti-pattern; demote in place per the demotion policy.
- Rows tagged `evidence-only` come from `evidence-only`-mode
  passes and have no topic-page prose backing them yet;
  citing one requires re-reading the underlying source.
- In-repo citations include `sha:` set to the short SHA at
  extraction time. Out-of-repo citations omit `sha:`.
- A claim whose only backing is a transcript or meeting note
  is `confidence: low` until a written source reinforces it.
- Never invent an `EV-###` to make a topic page look complete;
  every cited ID must resolve to a row in this file.
