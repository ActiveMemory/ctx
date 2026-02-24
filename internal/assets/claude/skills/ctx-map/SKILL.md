---
name: ctx-map
description: "Build and maintain architecture maps. Use to create or refresh ARCHITECTURE.md and DETAILED_DESIGN.md."
allowed-tools: Bash(ctx:*), Bash(git:*), Read, Write, Edit, Glob, Grep
---

Build and maintain two architecture documents incrementally:
**ARCHITECTURE.md** (succinct project map, loaded at session start)
and **DETAILED_DESIGN.md** (deep per-module reference, consulted
on-demand). Coverage is tracked in `map-tracking.json` so each run
extends the map rather than re-analyzing everything.

## When to Use

- First time setting up architecture documentation for a project
- Periodically to refresh stale module coverage after significant
  changes
- After major refactors, new package additions, or dependency changes
- When the agent nudges that the map is stale (>30 days, commits
  detected)
- When you need deep understanding of a module before working on it

## When NOT to Use

- For minor code changes that don't affect module boundaries or
  data flow
- When ARCHITECTURE.md just needs a quick path fix (use `/ctx-drift`
  instead)
- Repeatedly in the same session without intervening code changes
- When the user has opted out (`opted_out: true` in
  map-tracking.json)

## Execution

### Phase 0: Check Opt-Out

Read `.context/map-tracking.json`. If it exists and
`opted_out: true`, say:

> Architecture mapping is opted out for this project. Delete
> `.context/map-tracking.json` to re-enable.

Then stop.

### Phase 1: Assess Current State

Determine if this is a **first run** or **subsequent run**:

- **First run**: no `.context/map-tracking.json` exists
- **Subsequent run**: tracking file exists with coverage data

For subsequent runs, identify the **frontier** — modules that need
analysis:

1. Read `map-tracking.json` for coverage state
2. For each covered module, check staleness:
   ```bash
   git log --oneline --since="<last_analyzed>" -- <module_path>/
   ```
3. Frontier = uncovered modules + stale modules (commits after
   `last_analyzed`) + low-confidence modules (confidence < 0.7)

### Phase 2: Survey (First Run) or Analyze Frontier (Subsequent Run)

**First run — full survey:**

1. Read `go.mod` (or equivalent) for project identity and deps
2. Explore directory structure:
   ```bash
   ctx status
   ```
3. Read key files in each package: exported types, functions,
   imports
4. Trace data flow through main entry points
5. Identify architectural patterns (dependency injection,
   interfaces, registries)

**Subsequent run — targeted analysis:**

1. For each frontier module, read its source files
2. Trace data flow and dependencies
3. Note changes since last analysis
4. Update confidence based on depth of understanding

### Phase 3: Update Documents

**ARCHITECTURE.md** — update ONLY if module boundaries, dependency
graph, data flow, or key patterns changed. Internal implementation
changes do NOT warrant updates. Target: under 4000 tokens (~16KB).

Required sections:
- Overview (design philosophy, key concepts)
- Package Dependency Graph (mermaid `graph TD`)
- Component Map (tables: package, purpose, depends on)
- Data Flow (mermaid sequence diagrams for key operations)
- Key Architectural Patterns
- File Layout (ASCII tree)

**DETAILED_DESIGN.md** — update per-module sections using this
format:

```markdown
## <module_path>

**Purpose**: One-line description.

**Key types**: List main structs/interfaces.

**Exported API**:
- `FuncName()` — what it does
- `Type.Method()` — what it does

**Data flow**: Entry → Processing → Output

**Edge cases**:
- Condition → behavior

**Dependencies**: list of internal packages used
```

Each section is self-contained. The agent reads specific sections
when working on a module, not the entire file.

### Phase 4: Update Tracking

Write `.context/map-tracking.json` with:

```json
{
  "version": 1,
  "opted_out": false,
  "opted_out_at": null,
  "last_run": "<ISO-8601 timestamp>",
  "coverage": {
    "<module_path>": {
      "last_analyzed": "<ISO-8601 timestamp>",
      "confidence": <0.0-1.0>,
      "files_seen": ["file1.go", "file2.go"],
      "notes": "Brief summary of understanding"
    }
  }
}
```

### Phase 5: Report

Summarize what was done:

1. **Modules analyzed**: list with old → new confidence
2. **Documents updated**: which sections changed in each doc
3. **Overall coverage**: fraction of modules at confidence ≥ 0.7
4. **Remaining frontier**: modules still below 0.7 or unanalyzed

## Confidence Rubric

Use these levels for honest self-assessment:

| Level     | Meaning |
|-----------|---------|
| 0.0 – 0.3 | Stubbed — directory listed but contents not examined |
| 0.4 – 0.6 | Shallow — purpose understood, key exports known, internal flow unclear |
| 0.7 – 0.8 | Solid — can explain exports, data flow, and main code paths |
| 0.9 – 1.0 | Deep — can explain edge cases, error handling, design rationale |

A confidence of 0.9 means "I could explain every exported function's
purpose and the data flow through this module." Not "I read the file."

## Opt-Out Handling

If the user says "never", "don't ask again", or similar:

1. Set `opted_out: true` and `opted_out_at: "<timestamp>"` in
   map-tracking.json
2. Confirm: "Noted — won't ask again. Delete
   `.context/map-tracking.json` to re-enable."
3. On future invocations, exit immediately with brief message

## Nudge Behavior

The agent MAY suggest `/ctx-map` during session start when:

- **No tracking file**: "This project doesn't have an architecture
  map yet. Want me to run `/ctx-map`?"
- **Stale (>30 days)**: "The architecture map hasn't been updated
  since <date> and there are commits touching <N> modules. Want me
  to refresh?"
- **Opted out**: say nothing

The nudge is a suggestion, not automatic execution.

## Quality Checklist

After running, verify:
- [ ] ARCHITECTURE.md is under 4000 tokens (~16KB)
- [ ] ARCHITECTURE.md has all required sections (Overview, Dependency
      Graph, Component Map, Data Flow, Key Patterns, File Layout)
- [ ] DETAILED_DESIGN.md uses consistent per-module format
- [ ] Each module section has Purpose, Key types, Exported API,
      Data flow, Edge cases, Dependencies
- [ ] map-tracking.json is valid JSON with version, coverage entries
- [ ] Confidence levels are honest (not inflated)
- [ ] Stale modules were re-analyzed, not just marked current
- [ ] ARCHITECTURE.md was only updated for boundary/flow/dependency
      changes, not internal implementation details
- [ ] Report was provided summarizing what changed
