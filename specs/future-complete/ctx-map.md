# Architecture Mapping (`/ctx-map`)

## Overview

A skill that builds and maintains architecture documentation
incrementally. Each run reads existing documents, identifies stale or
uncovered areas, analyzes code, and updates the documents. Coverage
state is tracked so subsequent runs extend the map rather than
re-analyzing everything.

Two output documents:
- **ARCHITECTURE.md** — succinct project map (~4000 tokens max).
  Loaded at session start via FileReadOrder.
- **DETAILED_DESIGN.md** — deep per-module reference. NOT in
  FileReadOrder; consulted on-demand when working on specific modules.

## Behavior

### First Run (no tracking file)

1. Survey: directory structure, key files, imports, go.mod
2. Create/update ARCHITECTURE.md: overview, package map, dependency
   graph, data flow diagrams, key patterns, file layout
3. Create DETAILED_DESIGN.md: per-module sections
4. Create map-tracking.json with initial coverage

### Subsequent Runs (tracking file exists)

1. Read map-tracking.json for coverage state
2. Check staleness per module:
   ```bash
   git log --oneline --since="<last_analyzed>" -- <module>/
   ```
3. Identify frontier: uncovered modules, stale modules (commits after
   last_analyzed), low-confidence modules (<0.7)
4. Analyze frontier areas (read code, trace data flow)
5. Update ARCHITECTURE.md only if boundaries/flow/dependencies changed
6. Update DETAILED_DESIGN.md with deeper module analysis
7. Update map-tracking.json with new coverage data
8. Report: what was analyzed, what changed, overall coverage

### Opt-Out

- User says "never" → agent sets `opted_out: true` + timestamp
- Agent confirms: "Noted — won't ask again. Delete
  `.context/map-tracking.json` to re-enable."
- If `opted_out: true`, skill exits immediately with brief message

### Nudge (Proactive Suggestion)

The agent MAY suggest running `/ctx-map` during session start when:
- No map-tracking.json exists → "This project doesn't have an
  architecture map yet. Want me to run `/ctx-map`?"
- Tracking exists, `last_run` > 30 days, and git shows changes →
  "The architecture map hasn't been updated since <date> and there
  are commits touching <N> modules. Want me to refresh?"
- `opted_out: true` → say nothing

The nudge is a VERBATIM suggestion, not automatic execution.

## Tracking File

**Location:** `.context/map-tracking.json` (committed to git)

```json
{
  "version": 1,
  "opted_out": false,
  "opted_out_at": null,
  "last_run": "2026-02-23T14:30:00Z",
  "coverage": {
    "internal/cli/pad": {
      "last_analyzed": "2026-02-23T14:30:00Z",
      "confidence": 0.85,
      "files_seen": ["pad.go", "store.go", "merge.go", "pad_test.go"],
      "notes": "Merge logic well understood; encryption flow covered"
    },
    "internal/notify": {
      "last_analyzed": "2026-02-20T10:00:00Z",
      "confidence": 0.6,
      "files_seen": ["notify.go"],
      "notes": "Basic webhook dispatch; retry/error paths unexplored"
    }
  }
}
```

### Confidence Rubric

| Level     | Meaning |
|-----------|---------|
| 0.0 – 0.3 | Stubbed — directory listed but contents not examined |
| 0.4 – 0.6 | Shallow — purpose understood, key exports known, internal flow unclear |
| 0.7 – 0.8 | Solid — can explain exports, data flow, and main code paths |
| 0.9 – 1.0 | Deep — can explain edge cases, error handling, design rationale |

The skill instructions define these levels so the agent can self-assess
honestly. A 0.9 means "I could explain every exported function's purpose
and the data flow through this module." Not "I read the file."

### Staleness Detection

For each module in `coverage`:

```bash
git log --oneline --since="<last_analyzed>" -- <module_path>/
```

If commits exist after `last_analyzed`, the module is stale and goes
back on the frontier regardless of confidence. The skill should
re-analyze and update confidence.

## ARCHITECTURE.md Constraints

- **Size target**: under 4000 tokens (~16KB of markdown)
- **Sections**: Overview, Package Dependency Graph (mermaid),
  Component Map (tables), Data Flow (mermaid sequence diagrams),
  Key Patterns, File Layout (ASCII tree)
- **Update threshold**: only when boundaries, dependencies, data flow,
  or key patterns change. Internal implementation changes do NOT
  trigger ARCHITECTURE.md updates.

## DETAILED_DESIGN.md Structure

Per-module sections with consistent format:

```markdown
## internal/cli/pad

**Purpose**: Encrypted scratchpad CRUD operations.

**Key types**: `Store` (read/write), `Entry` (single pad item)

**Exported API**:
- `Cmd()` — Cobra command registration
- `Store.Read()` / `Store.Write()` — encrypted storage
- `Merge()` — dedup merge of two scratchpads

**Data flow**: User → CLI flags → Store.Read() → decrypt → modify →
encrypt → Store.Write()

**Edge cases**:
- Missing encryption key → falls back to plaintext check
- Duplicate entries during merge → content-hash dedup

**Dependencies**: `internal/config`, `internal/crypto`
```

This file can grow large. Each section is self-contained and can be
read independently. The agent should read specific sections when
working on a module, not the entire file.

NOT in FileReadOrder — never loaded at session start. Consulted
on-demand via Read tool.

## Files

### New
- `internal/assets/claude/skills/ctx-map/SKILL.md` — Skill definition
- `specs/ctx-map.md` — This spec

### Modified
- `internal/config/file.go` — Add constants + permission

### Created at Runtime (by the skill, not by Go code)
- `.context/map-tracking.json` — Coverage tracking
- `.context/DETAILED_DESIGN.md` — Deep module reference

## Non-Goals

- No Go code for tracking.json parsing — agent reads/writes directly
- No CLI `ctx map` command — skill-only for now
- No automatic scheduling — user-invoked or agent-suggested
- No `ctx drift` integration in this iteration
- No `ctx init` template for DETAILED_DESIGN.md or map-tracking.json
- No modifications to FileReadOrder (DETAILED_DESIGN.md stays out)
