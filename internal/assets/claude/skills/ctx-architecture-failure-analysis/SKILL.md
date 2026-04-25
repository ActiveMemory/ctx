---
name: ctx-architecture-failure-analysis
description: "Adversarial failure analysis for codebases. Generates falsifiable incident hypotheses: race conditions, ordering assumptions, cache staleness, error swallowing, ownership gaps, idempotency failures, and scaling cliffs. Produces DANGER-ZONES.md with evidence-backed, ranked findings."
allowed-tools: Bash(ctx:*), Bash(git:*), Bash(go:*), Read, Write, Edit, Glob, Grep, mcp__gitnexus__*, mcp__gemini-search__*
---

Adversarial analysis that identifies where a codebase will
silently betray you. Think like a correctness chaos monkey:
someone who knows the code intimately and exploits logical
bugs, not security holes.

## Design Principle

**Generate falsifiable incident hypotheses.** `/ctx-architecture`
maps what exists. `/ctx-architecture-enrich` improves map fidelity.
This skill generates concrete, disprovable claims about where the
map will break under real-world conditions. Every finding is a
hypothesis with evidence, not a suspicion or a vibe.

The goal is to find failure modes that code review misses: the
ones that ship, pass tests, and break in production at 3am.

This skill requires `/ctx-architecture` artifacts as input.
If they don't exist, stop and tell the user to run
`/ctx-architecture` first.

## When to Use

- After `/ctx-architecture` has run and artifacts exist
- Before a release or major deployment
- After a significant refactor that changed data flow
- When investigating a production incident category
- When onboarding to an unfamiliar codebase and need to
  know where the dragons are

## When NOT to Use

- Without architecture artifacts (run `/ctx-architecture` first)
- For security analysis (that's `ctx-threat-model`, a separate
  concern: auth bypass, injection, privilege escalation)
- On trivial or small codebases where the analysis cost exceeds
  the risk
- Mid-refactor when the code is intentionally in flux

## Inputs

**Required** (must exist before running):
- `.context/ARCHITECTURE.md`: system map
- `.context/DETAILED_DESIGN*.md`: module-level detail
- `.context/map-tracking.json`: coverage data

**Optional** (enhances analysis):
- `.context/DANGER-ZONES.md`: existing danger zones from
  `/ctx-architecture` principal mode (used as starting points,
  not as the final word)
- GitNexus MCP: blast radius estimation, shared-state detection
- Gemini Search: cross-reference against known failure patterns

## Process

### Phase 0: Validate Prerequisites

1. Check that architecture artifacts exist. If missing:
   > Architecture artifacts not found. Run `/ctx-architecture`
   > first; this skill analyzes existing maps, it doesn't
   > create them.
2. Load `map-tracking.json` to identify which modules have
   sufficient coverage (confidence >= 0.7). Low-confidence
   modules get flagged as "unanalyzed, risk unknown" rather
   than skipped.
3. If `.context/DANGER-ZONES.md` exists, load it as seed
   findings to extend, not as the complete picture.

### Phase 1: Build the Attack Surface

Read architecture artifacts and identify **mutation points**:
places where state changes, data transforms, or side effects
occur. These are the attack surface for correctness failures.

For each module with confidence >= 0.7:

1. Read the DETAILED_DESIGN entry for the module
2. Identify:
   - Shared mutable state (package-level vars, singletons,
     caches, registries)
   - Concurrent access points (goroutines, channels, locks)
   - External I/O boundaries (file, network, database)
   - Error handling chains (where errors are caught, wrapped,
     or swallowed)
   - Implicit ordering dependencies (init order, registration
     order, shutdown sequence)
   - State machines and transition points

3. For each mutation point, read the actual source code,
   DETAILED_DESIGN summaries are not enough for failure
   analysis. You need to see the actual lock scope, the actual
   error check, the actual nil guard.

### Phase 2: Adversarial Analysis

Apply each failure category systematically to the mutation
points identified in Phase 1. For each category, ask:
"How would a correctness chaos monkey make this fail silently?"

Every candidate finding must meet the **evidence standard**
before it can be recorded:

1. **Code path observed**: Exact file, function, line range
2. **Triggering precondition**: What must be true for the
   failure to occur
3. **Failure path**: Step-by-step sequence from trigger to
   observable effect
4. **Why it is silent**: Why tests, logs, or monitoring miss it
5. **Code evidence**: The specific code pattern that supports
   the claim (the missing lock, the unchecked error, the
   unbounded loop)

If you cannot provide all five, the finding is a hypothesis,
not a confirmed hazard. Label it accordingly (see Confidence
below).

#### Category 1: Concurrency

- **Races**: Shared state accessed from multiple goroutines
  without synchronization. Check: is every field of every
  shared struct protected? Are map reads concurrent-safe?
- **Deadlocks**: Lock ordering violations, channel operations
  that can block forever
- **Goroutine leaks**: Goroutines started but never joined,
  context cancellation not propagated

#### Category 2: Ordering Assumptions

- **Init order**: Code that assumes packages initialize in a
  specific order. `init()` functions with side effects that
  depend on other packages' `init()`
- **Registration order**: Slices or maps where order matters
  but insertion order isn't guaranteed
- **Shutdown sequence**: Resources closed in wrong order,
  goroutines still running during cleanup

#### Category 3: Cache Staleness

- **TTL-less caches**: In-memory caches with no expiration
  or invalidation
- **Read-your-writes**: Code that writes then reads expecting
  the write to be visible, but a cache serves stale data
- **Cross-process**: Caches that don't account for multiple
  instances or processes

#### Category 4: Fan-out Amplification

- **N+1 patterns**: Loops that make one call per item instead
  of batching
- **Recursive expansion**: Tree traversals that expand
  exponentially
- **Retry storms**: Retries without backoff that amplify under
  load

#### Category 5: Ownership and Lifecycle

- **Orphaned resources**: Objects created but never cleaned up
  (file handles, goroutines, temp files)
- **Double-close**: Resources closed by multiple owners
- **Use-after-close**: References held to closed resources
- **Force-delete orphans**: Deletion that doesn't cascade to
  dependent resources

#### Category 6: Error Handling

- **Silent swallowing**: `_ = someFunc()` or empty `if err`
  blocks
- **Error shadowing**: Inner errors that mask outer context
- **Partial failure**: Operations that succeed partially and
  leave inconsistent state
- **Panic recovery**: `recover()` that catches too broadly and
  masks real bugs

#### Category 7: Scaling Cliffs

- **O(n^2) hidden in loops**: Quadratic behavior that's fine
  for 10 items but kills at 10,000
- **Unbounded growth**: Data structures that grow without limit
  (in-memory lists, log files, caches)
- **Single-threaded bottlenecks**: Sequential operations that
  can't be parallelized
- **Global locks**: Locks that serialize all operations across
  unrelated requests

#### Category 8: Idempotency Failures

- **Duplicate processing**: At-least-once semantics without
  deduplication, causing double writes or double side effects
- **Retry-induced mutations**: Retries that re-execute
  non-idempotent operations (increment counters, send emails,
  append to lists)
- **Missing idempotency keys**: Operations that should be
  keyed but aren't, making replay indistinguishable from
  first execution

#### Category 9: State Machine and Invariant Drift

- **Illegal intermediate states**: Objects that pass through
  states the system never validates (half-initialized structs,
  partially-migrated records)
- **Unvalidated transitions**: State changes that skip
  validation (e.g., moving from "pending" to "complete"
  without passing through "running")
- **Cross-store inconsistency**: State split across multiple
  stores (file + memory, database + cache) that can diverge
  after partial updates

### Phase 3: Challenge Each Finding

Before a candidate finding is accepted, attempt to disprove it.
For each candidate, explicitly check:

- Is there existing locking or serialization that invalidates
  the race condition concern?
- Is the cache intentionally immutable or write-once?
- Is the shutdown ordering already coordinated by context
  cancellation or a shutdown hook?
- Is the apparent N+1 actually bounded by a small constant?
- Is the "missing" error check handled by a defer or a
  higher-level wrapper?
- Does the test suite already cover this failure path?

If the challenge succeeds (the concern is invalid), drop the
finding. If the challenge partially succeeds (the concern is
mitigated but not eliminated), note the mitigation and adjust
confidence downward.

This phase is mandatory. Skipping it produces smart fiction.

### Phase 4: Quantify and Cross-Reference

For each surviving finding:

1. **Blast radius**: If this fails, what breaks?
   - If GitNexus is available, use `impact` to get caller
     chains and dependency graphs
   - Otherwise, estimate from the architecture dependency graph

2. **Detection gap**: Would existing tests catch this?
   - Check test coverage for the affected code path
   - Check if the failure mode is tested (not just the happy
     path)

3. **Likelihood**: How likely is this to trigger?
   - Is the precondition common (every request) or rare
     (only under load)?
   - Has this pattern caused issues before? (check git log
     for related fixes)

4. **Risk score**: Assign explicit scores on a 1-3 scale:
   - **Likelihood**: 1 (rare) | 2 (uncommon) | 3 (common)
   - **Blast radius**: 1 (localized) | 2 (module) | 3 (system)
   - **Detection gap**: 1 (tested) | 2 (partially tested) |
     3 (untested/silent)
   - **Total**: Sum of all three (range 3-9)
   - **Critical**: Total >= 7 AND failure is silent or cascading
   - **Elevated**: All others

5. **Cross-reference**: If Gemini Search is available, search
   for the pattern in known incident databases, blog posts,
   and similar project post-mortems. This grounds findings in
   real-world evidence.

### Phase 5: Write DANGER-ZONES.md

Write `.context/DANGER-ZONES.md` with findings ranked by
risk score (highest first within each tier).

```markdown
# Danger Zones

_Generated YYYY-MM-DD by /ctx-architecture-failure-analysis._
_Run after /ctx-architecture for full coverage._

## Critical (risk score >= 7, silent or cascading)

### DZ-1: [Location]: [Failure Mode]

**Category**: Concurrency | Ordering | Cache | Amplification
  | Ownership | Error Handling | Scaling | Idempotency
  | State Machine
**Confidence**: High | Medium | Low
**Risk score**: L:N + B:N + D:N = N
**Location**: `package/file.go:function`
**Failure mode**: What goes wrong and why
**Triggering precondition**: What must be true for this to fire
**Failure path**: Step-by-step from trigger to effect
**Blast radius**: What breaks when this fails
**Detection gap**: Why tests don't catch it
**Evidence**: The specific code pattern (with file:line)
**Suggested fix**: Concrete code change, not vague advice

### DZ-2: ...

## Elevated (risk score < 7 or non-silent)

### DZ-3: ...

## Unanalyzed Modules

Modules with coverage < 0.7 in map-tracking.json:
- `module/path` (confidence: 0.3): risk unknown
```

**Confidence levels:**
- **High**: All five evidence fields present, challenge phase
  did not weaken the finding
- **Medium**: Evidence is strong but triggering precondition
  is hard to verify statically (may need runtime proof)
- **Low**: Plausible hypothesis based on code patterns, but
  could not be fully confirmed or disproved from static reading

**DZ numbering**: Sequential across the file. Stable across
re-runs (findings for the same location keep their number).

### Phase 6: Summary Report

Print a summary to the terminal:

```
## Failure Analysis Report

Modules analyzed: N (of M total)
Findings: N critical, M elevated
  High confidence: N
  Medium confidence: N
  Low confidence: N
Unanalyzed modules: K (coverage < 0.7)
Findings challenged and dropped: N

### Critical Findings
- DZ-1: [one-line summary] (L:N B:N D:N = N)
- DZ-2: [one-line summary] (L:N B:N D:N = N)

### Top Recommendation
[Single most impactful fix across all findings]
```

## Relationship to Other Skills

| Skill | Mode |
|-------|------|
| `/ctx-architecture` | Map what exists |
| `/ctx-architecture-enrich` | Improve map fidelity |
| `/ctx-architecture-failure-analysis` | Generate falsifiable incident hypotheses |
| `ctx-threat-model` (future) | Security-focused analysis |
| `/ctx-architecture` P4 | Surface danger zones noticed during mapping |

The key distinction: P4 extracts danger zones that were
*noticed during mapping*. This skill *generates hypotheses*
and *tests them against the code*. P4 findings are
observations; this skill's findings are tested claims.

## Quality Checklist

Before writing DANGER-ZONES.md, verify:
- [ ] Architecture artifacts exist and were loaded
- [ ] All 9 failure categories were applied (not skipped)
- [ ] Each finding meets the evidence standard (code path,
  trigger, failure path, silence reason, code evidence)
- [ ] Each finding includes a confidence level (High/Med/Low)
- [ ] Each finding has an explicit risk score (L+B+D)
- [ ] Each finding was challenged in Phase 3 with a "why this
  might be false" pass
- [ ] Findings that failed the challenge were dropped
- [ ] Findings distinguish confirmed hazards from plausible
  hypotheses via the confidence field
- [ ] Findings are ranked by risk score, not discovery order
- [ ] Source code was read for each finding (not just
  DETAILED_DESIGN summaries)
- [ ] Unanalyzed modules are listed with their coverage level
- [ ] Existing DANGER-ZONES.md findings were incorporated
  (not duplicated or lost)
- [ ] At least one concrete fix is suggested per finding
- [ ] Each finding names the triggering precondition explicitly
- [ ] Summary report includes confidence breakdown and
  challenge drop count
- [ ] If GitNexus available: blast radius was verified with
  impact analysis
- [ ] If Gemini available: findings were cross-referenced
  against known patterns
