---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Architecture Deep Dive
icon: lucide/layers
---
![ctx](../images/ctx-banner.png)


## The Problem

Understanding a codebase at the surface level is easy. Understanding
where it will break under real-world conditions takes three passes:
mapping what exists, quantifying how it connects, and hunting for
where it silently fails. Most teams stop at the first pass.

## TL;DR

```bash
# Pass 1: Map the system
/ctx-architecture

# Pass 2: Enrich with code intelligence
/ctx-architecture-enrich

# Pass 3: Hunt for failure modes
/ctx-architecture-failure-analysis
```

Each pass builds on the previous one. Run them in order. The
output accumulates in `.context/` — each pass reads the prior
artifacts and extends them.

## Commands and Skills Used

| Tool                                    | Type  | Purpose                                          |
|-----------------------------------------|-------|--------------------------------------------------|
| `/ctx-architecture`                     | Skill | Map modules, dependencies, data flow, patterns   |
| `/ctx-architecture-enrich`              | Skill | Verify blast radius and flows with code intel     |
| `/ctx-architecture-failure-analysis`    | Skill | Generate falsifiable incident hypotheses          |
| `ctx drift`                             | CLI   | Detect stale paths and broken references          |
| `ctx status`                            | CLI   | Quick structural overview                         |

## The Workflow

### Pass 1: Map what exists

```text
/ctx-architecture
```

Produces:

- **ARCHITECTURE.md** — succinct project map (< 4000 tokens),
  loaded at every session start
- **DETAILED_DESIGN*.md** — deep per-module reference with
  exported API, data flow, danger zones, extension points
- **CHEAT-SHEETS.md** — lifecycle flow diagrams
- **map-tracking.json** — coverage state with confidence scores

This pass forces deep code reading. No shortcuts, no code
intelligence tools — the agent reads every module it analyzes.
That forced reading is what makes the subsequent passes useful.

**When to run**: First time on a codebase, or after significant
structural changes (new packages, moved files, changed
dependencies).

**Principal mode**: Add `principal` to get strategic analysis
(ARCHITECTURE-PRINCIPAL.md, DANGER-ZONES.md from P4):

```text
/ctx-architecture principal
```

### Pass 2: Enrich with code intelligence

```text
/ctx-architecture-enrich
```

Takes the Pass 1 artifacts as baseline and layers on verified,
graph-backed data from GitNexus:

- Blast radius numbers for key functions
- Execution flow traces through hot paths
- Domain clustering validation
- Registration site discovery

This pass does not replace reading — it quantifies what reading
found. If Pass 1 says "module X depends on module Y," Pass 2
says "module X has 47 callers in module Y, and changing function
Z would affect 12 downstream consumers."

**When to run**: After Pass 1, when you need quantified
confidence for refactoring decisions or risk assessment.

**Requires**: GitNexus MCP server connected.

### Pass 3: Hunt for failure modes

```text
/ctx-architecture-failure-analysis
```

The adversarial pass. Reads all prior artifacts, then
systematically hunts for correctness bugs across 9 failure
categories:

1. Concurrency (races, deadlocks, goroutine leaks)
2. Ordering assumptions (init, registration, shutdown)
3. Cache staleness (TTL-less, read-your-writes, cross-process)
4. Fan-out amplification (N+1, retry storms)
5. Ownership and lifecycle (orphans, double-close)
6. Error handling (silent swallowing, partial failure)
7. Scaling cliffs (quadratic, unbounded, global locks)
8. Idempotency failures (duplicate processing, retry mutations)
9. State machine drift (illegal states, unvalidated transitions)

Every finding must meet an evidence standard: code path, trigger,
failure path, silence reason, and code evidence. A mandatory
challenge phase attempts to disprove each finding before it is
accepted. Findings carry a confidence level (High/Medium/Low) and
explicit risk score.

Produces **DANGER-ZONES.md** — a ranked inventory of findings
split into Critical and Elevated tiers.

**When to run**: Before releases, after major refactors, when
investigating incident categories, or when onboarding.

## What You Get

After all three passes, `.context/` contains:

| File | From | Purpose |
|------|------|---------|
| `ARCHITECTURE.md` | Pass 1 | System map (session-start context) |
| `DETAILED_DESIGN*.md` | Pass 1 | Module-level deep reference |
| `CHEAT-SHEETS.md` | Pass 1 | Lifecycle flow diagrams |
| `map-tracking.json` | Pass 1 | Coverage and confidence data |
| `CONVERGENCE-REPORT.md` | Pass 1 | What's covered, what's not |
| `DANGER-ZONES.md` | Pass 3 | Ranked failure hypotheses |

Pass 2 enriches Pass 1 artifacts in-place rather than creating
new files.

## Tips

- **Run Pass 1 with focus areas** if the codebase is large.
  The skill asks what to go deep on — name the modules you're
  about to change.
- **You don't need all three passes every time.** Pass 1 is
  the foundation. Pass 2 and 3 are for when you need
  quantified confidence or adversarial rigor.
- **Re-run Pass 1 incrementally.** It tracks coverage in
  `map-tracking.json` and only re-analyzes stale modules.
- **Pass 3 is most valuable before releases.** The ranked
  DANGER-ZONES.md is a pre-release checklist.
- **The trilogy maps to a question progression**: How does it
  work? How well does it connect? Where will it break?

## See Also

*See also: [Detecting and Fixing Context Drift](context-health.md)
— keep architecture artifacts fresh between deep-dive sessions.*

*See also: [Detecting and Fixing Context Drift](context-health.md)
— structural checks that complement architecture analysis.*
