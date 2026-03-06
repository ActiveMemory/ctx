# Spec: Memory Import (`ctx memory import`)

**Status**: Future — depends on foundation (`specs/memory-bridge.md`)
**Brainstorm**: `ideas/brainstorm-memory-bridge.md`
**Findings**: `ideas/claude-code-project-directory-structure.md`

Import entries from Claude Code's MEMORY.md into structured `.context/`
files using heuristic classification.

## Prereqs

- `internal/memory/` package (discover, mirror, state) — from foundation
- `ctx memory sync` working — import diffs mirror against current MEMORY.md

## CLI

```bash
ctx memory import              # Heuristic classification (default)
ctx memory import --dry-run    # Show what would be imported
ctx memory import --interactive # Ask for each entry (future)
```

### Example Output

```
$ ctx memory import
Scanning MEMORY.md for new entries...
  Found 3 new entries (15 lines) since last import

  → "always use bun for this project"
    Classified: CONVENTIONS.md
    ✓ Added to CONVENTIONS.md

  → "decided to use marker-based merge for bidirectional sync"
    Classified: DECISIONS.md
    ✓ Added to DECISIONS.md

  → "golangci-lint v2 ignores inline nolint for some linters"
    Classified: LEARNINGS.md
    ✓ Added to LEARNINGS.md

  → "Session 2026-03-05: Memory Bridge Design"
    Classified: skip (session notes)

Imported: 3 entries (1 convention, 1 decision, 1 learning)
Skipped: 1 entry (session notes)
```

## Classification Heuristics

Entries are classified by keyword matching on the entry text:

| Pattern | Target | Examples |
|---------|--------|----------|
| `always use`, `prefer`, `convention`, `never use`, `standard` | CONVENTIONS.md | "always use bun", "prefer filepath.Join" |
| `decided`, `chose`, `trade-off`, `approach`, `over`, `instead of` | DECISIONS.md | "decided to use SQLite over Postgres" |
| `gotcha`, `learned`, `watch out`, `bug`, `caveat`, `careful` | LEARNINGS.md | "learned that nolint is ignored in v2" |
| `todo`, `need to`, `follow up`, `should`, `task` | TASKS.md | "need to add tests for import" |
| Everything else | Skip | Session summaries, notes |

Matching is case-insensitive. Multiple matches resolve to the first
match in priority order (conventions > decisions > learnings > tasks).

Skipped entries stay in mirror.md — the user can cherry-pick from
the raw file, git history, or timestamped archives at any time.

## Entry Boundary Detection

MEMORY.md entries are delimited by:
- Markdown headers (`##`, `###`)
- Blank-line-separated paragraphs
- List items (`-`, `*`)

Each distinct block is classified independently.

## Deduplication

Before promoting, check the entry hash (first 64 chars + length) against
`memory-import.json`. This prevents re-importing identical entries without
full-text comparison on every run.

```json
{
  "imported_hashes": [
    "a1b2c3d4:decision:2026-03-05",
    "e5f6a7b8:learning:2026-03-05"
  ]
}
```

## Promotion

Classified entries are appended to the appropriate `.context/` file
using the same format as `ctx add`:

- Decisions → DECISIONS.md with `## [timestamp] Title` header
- Learnings → LEARNINGS.md with `## [timestamp] Title` header
- Conventions → CONVENTIONS.md appended to relevant section
- Tasks → TASKS.md as `- [ ] Description`

Each promoted entry gets a source annotation:

```markdown
## [2026-03-04-150000] Use pnpm for package management

Source: auto-memory import

Chose pnpm over npm because...
```

## Package Additions

```
internal/
├── memory/
│   ├── classify.go      # Classify(), heuristic rules
│   └── classify_test.go
├── cli/
│   └── memory/
│       └── import.go    # ctx memory import
```

## Dependencies

- `internal/memory` (foundation: discover, mirror, state)
- `internal/index`: entry header formatting for promoted entries
- Reuse `ctx add` internals for promotion to context files

## Testing

| Test | Type | Scope |
|------|------|-------|
| Classification heuristics (table-driven) | Unit | `classify.go` |
| Entry boundary detection | Unit | `classify.go` |
| Deduplication via hash | Unit | `state.go` |
| Import with mixed entry types | Integration | fixture MEMORY.md |

## Non-Goals

- No LLM-based classification (heuristics only for v1)
- No interactive mode (future)
- No cross-project import (shelved)
