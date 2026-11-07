---
name: ctx-architecture
description: "Build and maintain architecture maps. Use to create or refresh ARCHITECTURE.md and DETAILED_DESIGN.md. Supports principal mode for deeper analysis: vision, future direction, bottlenecks, implementation alternatives, gaps, upstream proposals, and intervention points."
---

Build and maintain two architecture documents incrementally:
**ARCHITECTURE.md** (succinct project map, loaded at session start)
and **DETAILED_DESIGN.md** (deep per-module reference, consulted
on-demand). Coverage is tracked in `map-tracking.json` so each run
extends the map rather than re-analyzing everything.

## Execution Priority

When time or context budget runs short, execute in this order.
Never skip a tier to do a lower one:

1. **Authoritative truth first**: ARCHITECTURE.md + DETAILED_DESIGN.md
   must be accurate and honest. Incomplete is fine; wrong is not.
2. **Surface uncertainty honestly**: partial coverage with correct
   confidence scores beats inflated scores. Mark what you don't know.
3. **Offer judgment only where grounded**: danger zones, extension
   points, improvement ideas only for modules you actually analyzed.
4. **Prefer fewer sharp insights over many shallow sections**: a
   CHEAT-SHEETS.md with one excellent cheat sheet beats five thin ones.
   An ARCHITECTURE-PRINCIPAL.md with three concrete risks beats ten
   vague ones.

## Mode Detection

Read the invocation for a mode keyword:

- **No keyword** (or `default`) → run **Default mode** (Phases 0-5 below)
- `principal` → run **Principal mode** (Phases 0-5 + Principal phases P1-P3)

Examples:
```text
/ctx-architecture
/ctx-architecture principal
/ctx-architecture (principal)
```

---

## When to Use

- First time setting up architecture documentation for a project
- Periodically to refresh stale module coverage after significant
  changes
- After major refactors, new package additions, or dependency changes
- When the agent nudges that the map is stale (>30 days, commits
  detected)
- When you need deep understanding of a module before working on it
- When you want strategic analysis of the architecture (principal mode)

## When NOT to Use

- For minor code changes that don't affect module boundaries or
  data flow
- When ARCHITECTURE.md just needs a quick path fix (use `/ctx-drift`
  instead)
- Repeatedly in the same session without intervening code changes
- When the user has opted out (`opted_out: true` in
  map-tracking.json)

---

## Default Mode (Phases 0-5)

### Phase 0: Check Opt-Out

Read `.context/map-tracking.json`. If it exists and
`opted_out: true`, say:

> Architecture mapping is opted out for this project. Delete
> `.context/map-tracking.json` to re-enable.

Then stop.

### Phase 0.25: Companion Tool Check

Check if a **web-search-with-citations MCP** is available by
attempting a simple query. The canonical implementation is
Gemini Search (`mcp__gemini-search__search_with_grounding`);
if your toolchain provides the same capability via a different
server (Firecrawl, Exa, Tavily, etc.), use whatever is
connected. This capability is for upstream documentation,
design rationale, KEPs, peer-project patterns — anything
outside the local codebase that helps understand *why* the
code is shaped the way it is.

**If available**: note it silently. Use it throughout the
analysis for upstream lookups. Prefer it over built-in web
search.

**If not available**: silently fall back to built-in web search
for upstream lookups. Do not prompt the user to install
anything — ctx does not vouch for companion-tool install
paths (see DECISIONS.md, 2026-05-23).

**Important**: this capability is for *upstream* and *external*
context only. Do not use it to understand the local codebase —
read the code directly. The depth of analysis comes from forced
reading, not from search shortcuts.

### Phase 0.5: Quick Structure Scan + Focus Areas

Before any deep analysis, do a lightweight structural survey to
discover what the project actually contains. This takes seconds
and makes the focus-area question concrete instead of open-ended.

**Scan steps** (no file reads - structure only):

```bash
# Detect ecosystem
ls go.mod package.json Cargo.toml pyproject.toml 2>/dev/null

# List top-level source directories / packages
# Go:
go list ./... 2>/dev/null | sed 's|.*/||' | sort -u | head -40
# or: ls internal/ cmd/ pkg/ 2>/dev/null

# Node/other: ls src/ lib/ packages/ 2>/dev/null

# Large monorepo guard: if >100 packages, limit to top 2 levels only
find . -mindepth 1 -maxdepth 2 -type d \
  ! -path './.git/*' ! -path './vendor/*' ! -path './node_modules/*' \
  | sort | head -60
```

**Then ask** (present the discovered package/module names):

```
I found these top-level packages/modules:
  [list from scan]

Any specific areas you'd like me to go deep on? You can name
packages from the list above, describe subsystems (e.g. "the
reconciler loop", "auth handling"), or say "all" for a uniform
pass.

Skip or press enter to do a standard uniform pass.
```

**If focus areas are given**, carry them forward:
- Phase 2 goes deep on focus packages (target confidence ≥ 0.8)
- Direct dependencies of focus packages get a solid pass (≥ 0.7)
- All other packages are stubbed (0.2) unless they appear as
  transitive dependencies
- DETAILED_DESIGN.md sections for focus packages are written first
  and in full detail
- Principal mode Phase P2 strategic questions reference the focus
  areas explicitly

**If "all" or no answer**, proceed with standard uniform analysis.

### Phase 1: Assess Current State

Determine if this is a **first run** or **subsequent run**:

- **First run**: no `.context/map-tracking.json` exists
- **Subsequent run**: tracking file exists with coverage data

For subsequent runs, identify the **frontier**: modules that need
analysis:

1. Read `map-tracking.json` for coverage state
2. For each covered module, check staleness:

```bash
git log --oneline --since="<last_analyzed>" \
-- <module_path>/
```

3. Frontier = uncovered modules + stale modules (commits after
   `last_analyzed`) + low-confidence modules (confidence < 0.7)

### Phase 2: Survey (First Run) or Analyze Frontier (Subsequent Run)

**First run: full survey:**

0. Run `ctx deps` to bootstrap the dependency graph:
   ```bash
   ctx deps
   ```
   Auto-detects the ecosystem (Go, Node.js, Python, Rust) from
   manifest files. Use this as the starting point for "Package
   Dependency Graph": verify and enrich with semantic context.

1. Read the project manifest for project identity (name, version,
   description): `ctx deps` covers the dependency tree
2. Explore directory structure:
   ```bash
   ctx status
   ```
3. Read key files in each package: exported types, functions,
   imports
4. Trace data flow through main entry points
5. Identify architectural patterns (dependency injection,
   interfaces, registries)

**Subsequent run: targeted analysis:**

1. For each frontier module, read its source files
2. Trace data flow and dependencies
3. Note changes since last analysis
4. Update confidence based on depth of understanding

### Phase 3: Update Documents

**ARCHITECTURE.md**: update ONLY if module boundaries, dependency
graph, data flow, or key patterns changed. Internal implementation
changes do NOT warrant updates. Target: under 4000 tokens (~16KB)
so ARCHITECTURE.md loads within the session-start context budget.

Required sections:
- Overview (design philosophy, key concepts)
- Package Dependency Graph (mermaid `graph TD`)
- Component Map (tables: package, purpose, depends on)
- Data Flow (mermaid sequence diagrams for key operations)
- Key Architectural Patterns
- File Layout (ASCII tree)

**DETAILED_DESIGN.md**: update per-module sections using this
format:

```markdown
## <module_path>

**Purpose**: One-line description.

**Key types**: List main structs/interfaces.

**Exported API**:
- `FuncName()`: what it does
- `Type.Method()`: what it does

**Data flow**: Entry → Processing → Output

Include an ASCII sequence diagram when there are 3+ actors or
non-obvious ordering:

```
Caller          Scheduler       Worker
|--schedule()-->|               |
|               |--dispatch()-->|
|               |<--result------|
|<--done--------|               |
```

Include an ASCII state diagram when the module manages lifecycle
or status transitions:

```
[Init] --configure()--> [Ready] --start()--> [Running]
|    |
error()---------|    |--stop()-->[ Stopped]
|                              [Stopped] --reset()--> [Ready]
[Failed]
```

Use plain ASCII (not mermaid) for DETAILED_DESIGN.md - it renders
in any terminal, editor, or raw file view without a renderer.
Reserve mermaid for ARCHITECTURE.md only.

**Edge cases**:
- Condition → behavior

**Performance considerations**:
- Known or likely bottlenecks (hot paths, allocation pressure,
  lock contention, I/O bound operations)
- Scale assumptions baked into the design (e.g. "assumes <1000
  items", "single-threaded reconcile loop")
- What breaks first under load

**Danger zones** (top 3 riskiest modification points):
1. `<symbol or area>` - why it's dangerous (hidden coupling,
   ordering assumption, shared mutable state, etc.)
2. ...
3. ...

**Control loop & ownership** (if the module participates in
reconciliation or state management):
- What owns the reconciliation for this module's resources?
- What is source of truth vs. derived/cached state?
- What triggers re-reconciliation?

**Extension points** (where features would naturally attach):
- `<symbol or pattern>` - what kind of extension fits here

**Improvement ideas** (1-3 concrete suggestions, not generic):
- `<specific change>` - what it fixes and why it's feasible

**Dependencies**: list of internal packages used
```

**Splitting DETAILED_DESIGN.md when it grows large:**

When DETAILED_DESIGN.md exceeds ~600 lines or covers 3+ natural
domains, split into domain files and keep a shallow index:

- `DETAILED_DESIGN.md` - index only (domain name, file pointer,
  module list, one-line domain purpose)
- `DETAILED_DESIGN-<domain>.md` - full module sections for that
  domain

Domains are natural groupings, not arbitrary splits. Examples:
- storage, auth, api, reconciler, cli, observability
- If no natural grouping exists, split by: core vs. peripheral

Index format:
```markdown
# Detailed Design Index

| Domain  | File                       | Modules              | Summary           |
|---------|----------------------------|----------------------|-------------------|
| storage | DETAILED_DESIGN-storage.md | pkg/store, pkg/cache | Persistence layer |
| auth    | DETAILED_DESIGN-auth.md    | pkg/authn, pkg/authz | Identity + policy |

> See individual files for module-level detail.
```

Update `map-tracking.json` to record which domain file each module
lives in:
```json
"pkg/store": {
  "domain_file": "DETAILED_DESIGN-storage.md",
  ...
}
```

Each section is self-contained. The agent reads specific sections
when working on a module, not the entire file.

**CHEAT-SHEETS.md**: write (or update) short mental models for
key lifecycle flows. One cheat sheet per major lifecycle or flow
identified in the codebase. Format:

```markdown
## <Lifecycle or Flow Name>

Steps:
1. <event or trigger>
2. <what happens next>
3. ...

Key invariants:
- <thing that must always be true>

Common failure modes:
- <condition> → <outcome>

Flow (ASCII - include when sequence or state is non-obvious):

  [Trigger] --> [Step A] --> [Step B] --> [Done]
                               |
                            [Error] --> [Retry] --> [Dead Letter]
```

Aim for cheat sheets that fit on one screen. If a flow needs more
than ~15 steps, split it. Write cheat sheets for at minimum:
- The main entry-point lifecycle (e.g. controller reconcile loop,
  request handler, CLI command dispatch)
- Any policy or rule evaluation flow
- Any significant async or background job lifecycle

Skip if the project has no meaningful lifecycles (e.g. a pure
library with no runtime behavior).

**GLOSSARY.md**: append project-specific terms discovered during
analysis. This captures the vocabulary that makes the codebase
searchable: type names, internal concepts, abbreviations, and
domain jargon that a new reader wouldn't know to search for.

Rules:
- Skip entirely if `.context/GLOSSARY.md` does not exist (the
  project hasn't opted into a glossary)
- Additive only: never modify or remove existing entries
- Maximum 10 new terms per run to avoid flooding
- Project-specific terms only: skip generic programming concepts
  (e.g. "mutex", "goroutine") and well-known patterns (e.g.
  "singleton"). Include terms that are unique to this codebase or
  used in a project-specific way
- Insert alphabetically into the existing list
- Format: `**Term**: one-line definition`
- Print added terms in the convergence report under a
  "Glossary additions" line

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

### Phase 5: Convergence Report + Search Prompts

Print a structured convergence report AND write it to
`.context/CONVERGENCE-REPORT.md`. The printed version is the
primary output the user reads. The file version is the artifact
that `/ctx-architecture-enrich` and future sessions consume.

The source of truth for confidence scores is `map-tracking.json`.
`CONVERGENCE-REPORT.md` is a human-readable view of that data -
if they ever conflict, `map-tracking.json` wins.

**Format:**

```

## Extended Reference

Read `references/extended.md` in this skill directory for the remaining specification and detailed steps.
