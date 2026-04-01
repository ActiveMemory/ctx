# Spec: Memory Publish (`ctx memory publish`)

**Status**: Future — depends on foundation (`specs/memory-bridge.md`)
**Brainstorm**: `ideas/brainstorm-memory-bridge.md`
**Findings**: `ideas/claude-code-project-directory-structure.md`

Push curated context from `.context/` into Claude Code's MEMORY.md
so the agent sees structured project context on session start without
needing hooks.

## Prereqs

- `internal/memory/` package (discover, mirror, state) — from foundation
- `ctx memory sync` working — publish reads current MEMORY.md for merge

## CLI

```bash
ctx memory publish             # Push curated context
ctx memory publish --budget 80 # Line budget (default 80)
ctx memory publish --dry-run   # Show what would be published
```

### Example Output

```
$ ctx memory publish
Publishing .context/ → MEMORY.md...

  Source files: TASKS.md, DECISIONS.md, CONVENTIONS.md, LEARNINGS.md
  Budget: 80 lines (MEMORY.md cap is 200; leaving room for Claude)

  Published block:
    5 pending tasks (from TASKS.md)
    3 recent decisions (from DECISIONS.md, last 7 days)
    5 key conventions (from CONVENTIONS.md, first 5)
    3 recent learnings (from LEARNINGS.md, last 7 days)

  Total: 62 lines (within 80-line budget)

✓ Published to MEMORY.md (markers: <!-- ctx:published --> ... <!-- ctx:end -->)
```

## Content Selection

Priority order for budget allocation:

1. **Pending tasks** (TASKS.md): `[ ]` items, newest first, max 10
2. **Recent decisions** (DECISIONS.md): last 7 days, max 5
3. **Key conventions** (CONVENTIONS.md): first N entries, max 10
4. **Recent learnings** (LEARNINGS.md): last 7 days, max 5

If total exceeds budget, trim from the bottom (learnings first,
then conventions, then decisions). Tasks always fit.

No need to be aggressive about the 200-line cap. If we exceed it,
Claude ignores the overflow. We still have everything in `.context/`.

## Marker-Based Merge

The published block is wrapped in HTML comments:

```markdown
<!-- ctx:published -->
# Project Context (managed by ctx)

## Pending Tasks
- [ ] Implement memory bridge
- [ ] Add hook nudge for memory drift
...

## Recent Decisions
...
<!-- ctx:end -->
```

Rules:
- ctx owns everything between the markers
- Claude owns everything outside the markers
- Import (see `ideas/spec-memory-import.md`) reads only outside the markers
- Publish replaces only inside the markers

## Marker Stripping Recovery

Claude owns MEMORY.md and may strip the markers (rewrite, reorganize).
When `ctx memory publish` detects markers are missing:

1. Warn: "Published markers were removed from MEMORY.md"
2. Treat entire file as Claude-owned content
3. Append new marker block at the end
4. No content is lost — Claude's content is preserved

The warning also surfaces in `ctx memory status` and the hook nudge.

## Wrap-Up Integration

Add a publish step to the `/ctx-wrap-up` skill:

```markdown
### Memory Publish

If auto memory is active (MEMORY.md exists), offer to publish:
"Publish current context to MEMORY.md? (ctx memory publish --dry-run to preview)"
```

This runs after the persist step (decisions, learnings, tasks) and
before the final summary. The user sees what would be published and
can approve or skip.

**Never auto-publish on `ctx commit`.** `git push` is an almost-
irreversible footgun — once content is indexed on the internet, it is
effectively leaked. Give the user freedom to verify. Wrap-up is the
publish trigger, not commit.

## Package Additions

```
internal/
├── memory/
│   ├── publish.go       # Publish(), marker merge, content selection
│   └── publish_test.go
├── cli/
│   └── memory/
│       └── publish.go   # ctx memory publish
```

## Dependencies

- `internal/memory` (foundation: discover, mirror, state)
- `internal/task`: task parsing for pending items
- `internal/index`: entry parsing for decisions/learnings selection

## Testing

| Test | Type | Scope |
|------|------|-------|
| Marker insertion into empty file | Unit | `publish.go` |
| Marker replacement (existing block) | Unit | `publish.go` |
| Marker stripping recovery (append) | Unit | `publish.go` |
| Budget trimming | Unit | `publish.go` |
| Content selection priority | Unit | `publish.go` |
| End-to-end publish cycle | Integration | fixture files |

## Non-Goals

- No auto-publish on commit (wrap-up only)
- No MEMORY.md editing beyond the marker block
- No cross-project publish
