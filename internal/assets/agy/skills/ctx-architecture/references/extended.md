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
- Exception: if a web-search-with-citations MCP is available
  (Gemini Search is the canonical example; equivalents include
  Firecrawl, Exa, Tavily), you MAY run upstream searches for
  KEPs, design docs, peer-project patterns, and ADRs — but only
  for concepts the codebase shows clear dependency on. Note
  what you searched and what you found. This applies in any
  mode, not just principal mode.
- If no such MCP is available and the user requested
  principal-mode depth, fall back to built-in web search for
  the same purpose

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
