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

Check if **Gemini Search** MCP is available by attempting a
simple query. Gemini is used for upstream documentation, design
rationale, KEPs, peer-project patterns - anything outside the
local codebase that helps understand *why* the code is shaped
the way it is.

**If available**: note it silently. Use Gemini throughout the
analysis for upstream lookups. Prefer it over built-in web search.

**If not available**: ask the user once:

```
Gemini Search is not connected. It helps me look up upstream
design docs, KEPs, and peer-project patterns during analysis.

Want to set it up now, or proceed without it?
```

Respect the answer and continue either way.

**Important**: Gemini is for *upstream* and *external* context
only. Do not use it to understand the local codebase - read the
code directly. The depth of analysis comes from forced reading,
not from search shortcuts.

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
## Convergence Report

### By Module

| Module | Confidence | Status | Blocker |
|--------|------------|--------|---------|
| pkg/foo | 0.9 | ✅ Converged | - |
| pkg/bar | 0.6 | 🔶 Shallow | Internal flow unclear |
| pkg/baz | 0.2 | 🔴 Stubbed | Not analyzed |

### By Domain (if natural groupings exist)

Group related modules and show aggregate coverage:
  e.g. "Auth layer: 2/3 modules converged (avg 0.72)"

### Overall

- Total modules: N
- Converged (≥ 0.9): N  ✅
- Solid (0.7-0.89): N   🟡
- Shallow (0.4-0.69): N 🔶
- Stubbed (< 0.4): N    🔴

### What Would Help Next

For each non-converged module, print a specific suggestion:

🔶 pkg/bar (0.6) - Shallow
  → Read the test files to understand expected behavior under
    edge cases: `pkg/bar/*_test.go`
  → Trace the internal flow through <specific function identified>
  → Ask: "walk me through what happens when X"

🔴 pkg/baz (0.2) - Not analyzed
  → Run /ctx-architecture with focus area: pkg/baz
  → Or: open pkg/baz/README.md if present

### Convergence Verdict

One of:
- ✅ CONVERGED - all modules ≥ 0.9, frontier empty. Further runs
  without code changes won't improve coverage.
- 🟡 MOSTLY CONVERGED - core modules ≥ 0.9, peripheral modules
  shallow. Diminishing returns on full re-run; use focus areas.
- 🔶 PARTIAL - significant modules below 0.7. Re-run with focus
  areas or read tests.
- 🔴 INCOMPLETE - substantial portions unanalyzed. Run again.
```

**Convergence thresholds:**
- Module is **converged** at confidence ≥ 0.9
- Project is **converged** when all non-peripheral modules ≥ 0.9
- Peripheral = no other modules depend on it AND it has no
  exported API surface (pure internal helpers, generated code,
  vendor)

**Blocker vocabulary** (use these consistently in the table):
- `Internal flow unclear` - exports known, internals not traced
- `Not analyzed` - directory listed only
- `Tests not read` - implementation known, behavior under edge
  cases unknown
- `Design rationale unknown` - code understood, "why" is unclear
- `Converged` - nothing left to learn from static reading

---

After printing the convergence verdict, append a **Search Prompts**
section. The skill has just read the codebase and knows its jargon -
this is the most useful thing it can hand back to someone who is
not blocked by intelligence but by not knowing the right words.

**Format:**

```
## Search Prompts

The right keyword changes everything. Based on what I found in
the codebase, here are targeted searches worth running - in your
internal docs, Confluence, Notion, Slack, or publicly:

### Fill the gaps (ranked by how much they'd help)

For modules/areas still below 0.9:

🔶 pkg/bar - Internal flow unclear
  Try searching:
  - "<SpecificTypeName> design" or "<SpecificTypeName> internals"
  - "<pattern observed, e.g. 'leader election'>  <project name>"
  - "why does <ProjectName> use <pattern>" (ADR or design doc)

🔴 pkg/baz - Not analyzed
  Try searching:
  - "<package name> <project name> explained"
  - "<key interface or type found> behavior"

### Concepts worth understanding deeply

List 3-5 technical concepts the codebase clearly depends on but
that can't be learned from the code alone. Give the exact search
phrase, not a topic:

- "<ExactConceptName> explained" - e.g. "etcd watch semantics
  explained", "CRDT merge strategies", "OIDC token refresh flow"
- "<pattern name> tradeoffs" - e.g. "saga pattern vs 2PC tradeoffs"

### Architecture decision records (if relevant)

If the code shows signs of a deliberate non-obvious choice
(e.g. custom retry logic instead of a library, unusual data
structure), suggest:
  - "<ProjectName> <decision> ADR"
  - "<ProjectName> <decision> RFC"
  - "why <ProjectName> doesn't use <obvious alternative>"

---
Note: I won't run these searches for you - you may have internal
docs where these are more useful than public results, and you know
which sources to trust. Pick the phrases that match what's blocking
you.
```

**Rules for this section:**
- Always generate search prompts, even for converged modules -
  there's always design rationale that code can't express
- Phrases must be concrete and use actual names/types from the
  codebase - no generic "learn more about X" fluff
- Rank by usefulness: gaps in shallow modules first, concepts
  second, ADRs third
- Maximum ~10 phrases total; fewer sharp ones beat many vague ones
- Default: do NOT run the searches yourself
- Exception: if Gemini Search is available, you MAY run upstream
  searches for KEPs, design docs, peer-project patterns, and ADRs
  - but only for concepts the codebase shows clear dependency on.
  Note what you searched and what you found. This applies in any
  mode, not just principal mode.
- If Gemini is not available and the user requested principal-mode
  depth, you may fall back to built-in web search for the same
  purpose

---

## Principal Mode (Phases 0-5 + P1-P3)

Run all default mode phases first (0-5), then continue below.
Principal mode is for strategic thinking - beyond "what is" to
"what could be" and "what should concern us."

### Phase P1: Extended Context Gathering

In addition to the default phase sources, read:

- `.context/TASKS.md` - outstanding work, future plans
- `CHANGELOG.md` or `docs/changelog.md` - trajectory of decisions
- `docs/` - any design rationale in user-facing docs
- Recent git log: `git log --oneline -30`

### Phase P2: Gather Strategic Context

Two-tier behavior - do not stall:

**If answers are available** (user provided them in the prompt,
or they exist in `.context/TASKS.md` / `DECISIONS.md`): use them.
Do not ask for what you already have.

**If answers are not available**: do NOT stop. Generate a
provisional principal analysis with assumptions explicitly labeled
(see Principal Mode Fallback below). Include a "Questions That
Would Sharpen This" section at the end of ARCHITECTURE-PRINCIPAL.md.

When asking the user, present all questions at once as a numbered
list - do not ask one-at-a-time:

```
Before I write the principal analysis, a few questions - skip
or say "unsure" on anything you don't know:

0. **Focus areas** (if not already set in Phase 0.5)

1. **Vision**: What is this project trying to become in 12-24 months?

2. **Future direction**: Any architectural pivots being considered?
   (plugin system, multi-tenant, cloud sync, daemon model, etc.)

3. **Known bottlenecks**: Where does the current design hurt you?

4. **Implementation alternatives**: Any decisions you'd do
   differently starting fresh?

5. **Gaps**: What's missing that you expect to need?

6. **Areas of improvement**: Known tech debt or structural awkwardness?
```

### Phase P3: Write Principal Analysis

After collecting answers, write `.context/ARCHITECTURE-PRINCIPAL.md`
(separate from `ARCHITECTURE.md` - speculation must not pollute
the authoritative doc).

```markdown
# Architecture - Principal Analysis
_Generated <date>. Strategic analysis only; see ARCHITECTURE.md
for the authoritative architecture reference._

## Current State Summary
[Condensed narrative of the current architecture - ~1 page max]

## Vision Alignment
[How does the current architecture support or constrain the stated
vision? What structural changes would enable it?]

## Future Direction
[Architectural implications of planned pivots or new capabilities.
What would need to change if [feature X] were added?]

## Known Bottlenecks
[Analysis of performance, scalability, or dev-experience pain
points identified in the codebase or raised by the user]

## Implementation Alternatives
[For 2-3 key design decisions: current approach, alternatives,
tradeoffs]

## Gaps
[Missing capabilities or abstractions the architecture doesn't
handle yet but probably will need to]

## Areas of Improvement
Ranked by impact/effort:
- **High impact, low effort** (do first)
- **High impact, high effort** (plan for)
- **Low impact** (defer or skip)

## Risks
[Architectural risks as the system scales, team grows, or
requirements evolve]

## Intervention Points
Top 5 highest-leverage places to implement new features or
improvements, ranked by impact/effort:
1. `<symbol or subsystem>` - what kind of change fits here and why
2. ...

(These are concrete locations - package paths, interface names,
function boundaries - not vague subsystem labels.)

## Upstream Proposals
2-3 changes worth proposing to the project upstream (KEP / RFC /
issue style thinking). For each:
- **What**: one-sentence description of the change
- **Why**: what problem it solves that the current design can't
- **Where**: which abstraction boundary it touches
- **Risk**: what it breaks or complicates

Each proposal must cross an abstraction boundary - it must affect
how modules interact, not just refactor internals. If it doesn't
change an interface, a contract, or an ownership boundary, it's
not upstream-worthy; it's a local improvement (put it in
Improvement Ideas instead).

## Productization Gaps
What would need to change for this to work at enterprise scale?
- Multi-cluster / multi-tenant gaps
- Observability and debuggability holes
- Operational hardening missing from current design
- What a large customer would hit first

## Failure-First Analysis
[Hidden assumptions baked into the architecture. What breaks
silently vs. loudly? What would cause a cascade? What does the
system assume about its environment that may not hold?]

## Onboarding Friction
[Practical, not theoretical - this is what a new engineer actually
hits in week one:]
- What makes this system hard to understand quickly?
- Which modules require tribal knowledge to use safely?
- Where would a new engineer get stuck first, and why?
- What isn't written down anywhere?
```

**Boundary hygiene** - ARCHITECTURE-PRINCIPAL.md is for synthesis,
leverage, risk, direction, and judgment. Do NOT restate module
details that already exist in DETAILED_DESIGN.md. Reference module
paths only where needed to ground an argument. If you find yourself
summarizing what a module does, stop - link to it instead.

**Principal mode fallback** - if Phase P2 answers were not provided,
label speculative sections clearly and add at the end:

```markdown
## Questions That Would Sharpen This Analysis

Answering any of these would move speculative sections to grounded ones:

1. **Vision** - What is this project trying to become in 12-24 months?
2. **Future direction** - Any architectural pivots being considered?
3. **Known bottlenecks** - Where does the current design hurt?
4. **Assumptions marked** - These sections are labeled [inferred]:
   [list them]
```

**Autonomous inferences** - principal mode must also answer the
following from the codebase alone, without waiting for user input.
These are things the code is silently deciding. Surface them:

- Where are abstraction boundaries likely to calcify under growth?
- Which current APIs are accidentally becoming public contracts?
- What will become expensive when team size or data volume doubles?
- Where is the architecture optimized for current workflow rather
  than long-term extensibility?
- Which parts are structurally elegant but strategically wrong for
  the likely future?

These go in a dedicated "Silent Choices" section in
ARCHITECTURE-PRINCIPAL.md. The code is making bets - name them.

**Opinion floor** - ARCHITECTURE-PRINCIPAL.md must contain at minimum:
- 3 risks (specific, not "this could be slow")
- 3 improvement ideas (concrete, not "add more tests")
- 2 upstream opportunities (actionable, not "contribute more")

Generate opinions, not just descriptions. If you find yourself
writing neutral summaries, push harder.

When in doubt, prefer a strong, falsifiable opinion over a safe,
generic one. Weak opinions are noise; strong opinions can be
corrected.

**Cross-project comparison** (include when the codebase shows
non-obvious design choices or when focus areas have well-known
peers):

For any module where a comparable exists in another project, add:
```markdown
### Compared to <PeerProject>/<Component>

- What <ThisProject> does differently
- What <PeerProject> does better
- What could be unified or learned from
```

Examples worth comparing when relevant:
- Velero vs Stash (backup)
- controller-runtime reconciler vs custom loops
- Gatekeeper vs Kyverno (policy)
- Any CNCF project vs its closest peer

Skip if no meaningful peer exists. Do not force comparisons.

Be direct. This document is for engineering judgment, not external
audiences.

### Phase P4: Write DANGER-ZONES.md

Extract danger zones from all DETAILED_DESIGN.md module sections
and compile them into a standalone `.context/DANGER-ZONES.md`.
This is the consolidated view - one document a reviewer or new
engineer can read to know where the dragons live.

```markdown
# Danger Zones

_Generated <date> from DETAILED_DESIGN.md danger zone sections.
Run `/ctx-architecture-enrich` to add verified blast radius data._

## Summary

| Module | Zone | Risk | Why |
|--------|------|------|-----|
| <path> | <symbol/area> | HIGH/MEDIUM/LOW | one-line reason |

## By Module

### <module_path>

1. **<symbol or area>** - <why it's dangerous>
   - Hidden coupling / ordering assumption / shared mutable state
   - Modification advice: <what to check before changing>

2. ...
```

**Rules:**
- Only include danger zones from modules actually analyzed
  (confidence ≥ 0.4)
- Risk level is the skill's judgment based on code reading:
  HIGH (will break things), MEDIUM (likely to cause subtle bugs),
  LOW (worth knowing but manageable)
- `/ctx-architecture-enrich` can later add verified blast radius
  numbers - leave room for that (don't claim precision you don't
  have from reading alone)
- If no danger zones were identified, skip the file entirely
  rather than writing an empty one

---

## Confidence Rubric

Score by **decision usefulness**, not descriptive completeness.
Ask: "What could an engineer safely do with this understanding?"

| Level      | Decision usefulness                                                          |
|------------|------------------------------------------------------------------------------|
| 0.0 - 0.3  | Stubbed: not safe to make any decisions; directory listed only               |
| 0.4 - 0.6  | Shallow: can describe purpose; not safe to modify without more reading       |
| 0.7 - 0.79 | Safe to make localized changes with care; can review simple PRs              |
| 0.8 - 0.89 | Can reason about design tradeoffs; safe to design changes in this module     |
| 0.9 - 1.0  | Can predict likely breakage from non-trivial changes; safe to own the module |

Inflate scores and you lie to the next agent that reads the tracking
file. Under-score and the convergence report will never clear.
Score the decision-usefulness honestly.

## Opt-Out Handling

If the user says "never", "don't ask again", or similar:

1. Set `opted_out: true` and `opted_out_at: "<timestamp>"` in
   map-tracking.json
2. Confirm: "Noted: won't ask again. Delete
   `.context/map-tracking.json` to re-enable."
3. On future invocations, exit immediately with brief message

## Nudge Behavior

The agent MAY suggest `/ctx-architecture` during session start when:

- **No tracking file**: "This project doesn't have an architecture
  map yet. Want me to run `/ctx-architecture`?"
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
  Data flow, Edge cases, Performance considerations, Control
  loop & ownership (if applicable), Danger zones, Extension
  points, Improvement ideas, Dependencies
- [ ] ASCII sequence diagram included when 3+ actors or
  non-obvious ordering
- [ ] ASCII state diagram included when module manages lifecycle
  or status transitions
- [ ] No mermaid in DETAILED_DESIGN.md (ASCII only)
- [ ] If DETAILED_DESIGN.md > ~600 lines or 3+ domains: split
  into domain files with shallow index
- [ ] map-tracking.json records domain_file for each module
  when split
- [ ] map-tracking.json is valid JSON with version, coverage entries
- [ ] Confidence levels are honest (not inflated)
- [ ] Stale modules were re-analyzed, not just marked current
- [ ] ARCHITECTURE.md was only updated for boundary/flow/dependency
  changes, not internal implementation details
- [ ] Convergence report printed with per-module table
- [ ] Domain groupings shown if natural groupings exist
- [ ] Each non-converged module has a specific "what would help"
  suggestion (not generic advice)
- [ ] Overall convergence verdict stated (CONVERGED / MOSTLY /
  PARTIAL / INCOMPLETE)
- [ ] Blocker column uses consistent vocabulary
- [ ] Search Prompts section printed after convergence verdict
- [ ] Search phrases use actual type/function/pattern names from
  the codebase (not generic topics)
- [ ] Phrases ranked: shallow-module gaps first, concepts second,
  ADRs third
- [ ] No more than ~10 phrases total
- [ ] Skill did NOT run local-code searches itself (upstream
  searches via Gemini are allowed)
- [ ] CONVERGENCE-REPORT.md written to .context/ (not just printed)
- [ ] Phase 0.25 Gemini check completed (available or user declined)
- [ ] Phase 0.5 structure scan was run before any deep analysis
- [ ] Focus areas question was asked with actual package names (not
  open-ended)
- [ ] If focus areas given: deep analysis concentrated there; other
  packages stubbed at 0.2 unless direct dependencies
- [ ] Principal mode: P2 answers used if available; if not,
  provisional analysis written with [inferred] labels
- [ ] Principal mode: "Questions That Would Sharpen This" section
  present if P2 answers were not provided
- [ ] Principal mode: output written to `ARCHITECTURE-PRINCIPAL.md`,
  not overwriting `ARCHITECTURE.md`
- [ ] Principal mode: "Silent Choices" section present (autonomous
  inferences from code - abstraction calcification, accidental
  contracts, scale costs, strategic bets)
- [ ] Principal mode: ARCHITECTURE-PRINCIPAL.md does not restate
  DETAILED_DESIGN.md content - links to module paths instead
- [ ] CHEAT-SHEETS.md written with at least one lifecycle flow
- [ ] Each cheat sheet fits ~one screen; long flows are split
- [ ] Danger zones section present in each DETAILED_DESIGN module
  (top 3, with reasoning - not just "this is complex")
- [ ] Extension points section present in each module
- [ ] Principal mode: Failure-First Analysis section written
- [ ] Principal mode: Onboarding Friction section present (practical,
  week-one concerns - not generic "hard to understand")
- [ ] Principal mode: Upstream Proposals cross abstraction boundaries
  (not internal refactors)
- [ ] Principal mode: Intervention Points section present (concrete
  locations, not vague labels)
- [ ] Principal mode: Upstream Proposals section present (2-3 items
  with what/why/where/risk)
- [ ] Principal mode: Productization Gaps section present
- [ ] Principal mode: opinion floor met (≥3 risks, ≥3 improvements,
  ≥2 upstream opportunities - specific, not generic)
- [ ] Principal mode: cross-project comparisons included where
  meaningful peers exist (not forced)
- [ ] Principal mode: DANGER-ZONES.md written with consolidated
  danger zones from all analyzed modules (skip if none found)
- [ ] Principal mode: DANGER-ZONES.md includes summary table and
  per-module breakdown with risk levels and modification advice
- [ ] GLOSSARY.md: new terms added alphabetically (max 10, project-
  specific only, skipped if file doesn't exist)
- [ ] Convergence report includes "Glossary additions" line if
  terms were added
