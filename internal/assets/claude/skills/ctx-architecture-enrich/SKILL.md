---
name: ctx-architecture-enrich
description: "Enrich architecture artifacts with code intelligence data. Takes existing /ctx-architecture output as baseline, verifies and quantifies with GitNexus MCP (blast radius, execution flows, domain clustering, registration sites). Run after /ctx-architecture, not instead of it."
allowed-tools: Bash(ctx:*), Bash(git:*), Read, Write, Edit, Glob, Grep, mcp__gitnexus__*, mcp__gemini-search__*
---

Enrich existing architecture artifacts with verified data from
code intelligence tools. This skill reads the output of
`/ctx-architecture` (which forces deep code reading) and layers
on quantified, graph-backed data that reading alone cannot
efficiently provide.

## Design Principle

**Reading first, tools second.** `/ctx-architecture` produces
deep artifacts through forced code reading - no code intelligence
tools, no shortcuts. This skill runs AFTER that pass, using the
deep artifacts as a baseline. It verifies, quantifies, and extends
- it never substitutes for reading.

The separation exists because agents take shortcuts when code
intelligence tools are available during analysis. A structural
query returns an answer without opening the file - so the agent
never discovers the operational details (defaults, timeouts, scale
math, edge cases) that only emerge from line-by-line reading. The
tool answers the question asked but prevents discovery of answers
to questions never asked.

## When to Use

- After `/ctx-architecture` or `/ctx-architecture principal` has
  produced artifacts
- After `npx gitnexus analyze --embeddings` completes on a project
  that already has architecture artifacts
- When the user says "enrich the architecture", "run enrichment
  pass", "add graph data", "quantify the danger zones"
- When DANGER-ZONES.md exists but lacks blast radius numbers
- When CONVERGENCE-REPORT.md shows shallow modules that could
  benefit from semantic search

## When NOT to Use

- As a substitute for `/ctx-architecture` - if no architecture
  artifacts exist, run `/ctx-architecture` first
- When GitNexus MCP is not connected or the index is stale -
  preflight will catch this
- Immediately after `/ctx-architecture` in the same session without
  user request - let the user review the base artifacts first

---

## Inputs (Required)

The skill refuses to run if these are missing:

- `.context/ARCHITECTURE.md` - the authoritative architecture map
- `.context/DETAILED_DESIGN.md` (or domain split files) - per-module
  deep reference
- `.context/map-tracking.json` - coverage state and confidence scores

The skill checks and warns if these are missing but proceeds
without them:

- `.context/DANGER-ZONES.md` - consolidated danger zones (if absent,
  extracts from DETAILED_DESIGN.md danger zone sections)
- `.context/CONVERGENCE-REPORT.md` - convergence state
- `.context/ARCHITECTURE-PRINCIPAL.md` - principal analysis
- `.context/CHEAT-SHEETS.md` - lifecycle flow cheat sheets

---

## Phase 1: Preflight

### 1.1 Verify Architecture Artifacts

Read the required files listed above. If any required file is
missing, stop and say:

```
Architecture artifacts not found. Run `/ctx-architecture` first
to generate the baseline, then run this skill to enrich it.
```

### 1.2 Verify Code Intelligence Tools

Check each tool silently:

**GitNexus** (required):
- Call `mcp__gitnexus__list_repos` to verify connection
- Check that the current project is indexed
- Compare the index timestamp against the latest git commit to
  detect staleness

If GitNexus is not connected or the project is not indexed:

```
GitNexus is not connected (or this project is not indexed).
This skill requires a code intelligence graph.

To set up: add the GitNexus MCP server, then run:
  npx gitnexus analyze --embeddings
```

If the index is stale (commits after last index):

- **≤ 5 commits behind**: warn and proceed.
  ```
  GitNexus index is slightly stale (last indexed: <date>,
  <N> commits since). Proceeding - results may be incomplete
  for recently changed code.
  ```
- **> 5 commits behind**: hard stop.
  ```
  GitNexus index is stale (last indexed: <date>, <N> commits
  since). Results would be unreliable.

  Run `npx gitnexus analyze` to update, then re-run this skill.
  ```

**Gemini Search** (optional):
- Attempt a simple query to verify connection
- If available: note silently, use for upstream pattern lookups
- If not available: note silently, proceed without it

### 1.3 Read Baseline Artifacts

Read all architecture artifacts into context. Pay attention to:
- Module list and confidence scores from `map-tracking.json`
- Danger zone entries from DANGER-ZONES.md or DETAILED_DESIGN.md
- Extension points from DETAILED_DESIGN.md module sections
- Shallow modules (confidence < 0.75) from `map-tracking.json`
- Convergence state from CONVERGENCE-REPORT.md

---

## Phase 2: Danger Zone Enrichment

For each danger zone identified in DANGER-ZONES.md (or extracted
from DETAILED_DESIGN.md if no standalone file exists):

1. **Run impact analysis** on the named symbol:
   ```
   mcp__gitnexus__impact({target: "<symbol>", direction: "upstream"})
   ```

2. **Record blast radius**:
   - d=1 count (direct callers - WILL BREAK)
   - d=2 count (indirect dependents - LIKELY AFFECTED)
   - d=3 count (transitive - MAY NEED TESTING)
   - Affected process/execution flow count

3. **Assign verified risk level**:

   Risk thresholds should consider repository scale. In a small
   repo (<1000 files), d=1=4 might be critical. In a large repo
   (>10k files), d=1=15 might be routine. When unsure, bias toward
   HIGH over MEDIUM.

   Guidelines (adjust for scale):
   - CRITICAL: d=1 > 10 or crosses 3+ domains
   - HIGH: d=1 > 5 or crosses 2 domains
   - MEDIUM: d=1 2-5, single domain
   - LOW: d=1 ≤ 1, localized

   Graph data is a lower bound, not ground truth. Dynamic dispatch,
   reflection, config-driven wiring, and runtime registration can
   make graphs incomplete. If blast radius seems suspiciously low
   for a symbol you know is critical from reading the code, flag:

   ```
   Risk: HIGH (enriched <date> via GitNexus)
   ⚠ Possible undercount - dynamic or indirect usage suspected
   ```

4. **Update DANGER-ZONES.md** with enrichment data:
   ```markdown
   1. **<symbol>** - <original description>
      - Blast radius: d=1: N, d=2: N, d=3: N
      - Affected flows: <list of process names>
      - Risk: HIGH (enriched 2026-03-25 via GitNexus)
      - Modification advice: <updated based on blast radius>
   ```

Update the summary table with verified risk levels.

---

## Phase 3: Extension Point Enrichment

For each extension point identified in DETAILED_DESIGN.md module
sections:

1. **Query the call graph** for registration patterns:
   ```
   mcp__gitnexus__context({name: "<registration function>"})
   ```

2. **Build a registration inventory**:
   - All call sites with file:line references
   - Count of registrations per pattern
   - Any unregistered implementations (defined but never wired)

3. **Write or update `.context/EXTENSION-POINTS.md`**:
   ```markdown
   # Extension Points

   _Generated <date>. Enriched via GitNexus call graph analysis._

   ## Summary

   | Pattern | Registration Function | Count | Files |
   |---------|----------------------|-------|-------|
   | <name>  | <func>               | N     | N     |

   ## By Pattern

   ### <Pattern Name>

   Registration function: `<func>` in `<file>`

   Registered implementations:
   1. `<impl>` - `<file>:<line>`
   2. ...

   Unregistered (defined but not wired):
   - `<impl>` - `<file>:<line>` (potential dead code or
     conditional registration)
   ```

---

## Phase 4: Execution Flow Enrichment

1. **Read all processes** from GitNexus:
   ```
   READ gitnexus://repo/<name>/processes
   ```

2. **Select the most significant flows** (10-15). Prefer flows
   that:
   - Share symbols with other flows (high centrality)
   - Originate from public APIs or entry points
   - Cross multiple domains
   Step count alone is not a good signal - a 50-step internal
   flow matters less than a 10-step cross-domain API flow.

3. **Identify multi-flow hotspots** - symbols that appear
   in 3+ execution flows are integration points worth knowing

4. **Update CHEAT-SHEETS.md** with an execution flow index:
   ```markdown
   ## Execution Flow Index (via GitNexus)

   _Enriched <date>. These flows are auto-detected from the call
   graph and complement the manually written cheat sheets above._

   | Flow | Steps | Entry Point | Key Symbols |
   |------|-------|-------------|-------------|
   | <name> | N | <entry> | <hotspot symbols> |

   ### Multi-Flow Hotspots

   Symbols participating in 3+ flows (high-impact modification
   points):

   | Symbol | Flows | Location |
   |--------|-------|----------|
   | <name> | N     | <file>:<line> |
   ```

---

## Phase 5: Domain Clustering Comparison

1. **Read auto-detected clusters** from GitNexus:
   ```
   READ gitnexus://repo/<name>/clusters
   ```

2. **Read manual domain splits** from DETAILED_DESIGN.md (the
   domain file index if split, or section headers if monolithic)

3. **Compare** - surface mismatches:
   - Modules that GitNexus groups together but the manual split
     separates → potential hidden coupling
   - Modules that GitNexus separates but the manual split groups
     → potential artificial grouping

4. **Write comparison** to ARCHITECTURE-PRINCIPAL.md (if it exists)
   or print to terminal:
   ```markdown
   ## Domain Clustering Comparison (enriched <date>)

   | Manual Domain | GitNexus Cluster | Match | Notes |
   |---------------|-----------------|-------|-------|
   | <domain>      | <cluster>       | yes/partial/no | <mismatch detail> |

   ### Hidden Coupling (GitNexus groups, manual splits)
   - <module A> ↔ <module B>: N calls between them, but in
     separate manual domains

   ### Artificial Grouping (GitNexus splits, manual groups)
   - <module A> and <module B>: 0 calls between them, but in
     same manual domain

   ### Boundary Violations
   If hidden coupling exceeds 10 calls between manually separated
   domains, flag as:

   ⚠ Architectural boundary violation: <domain A> ↔ <domain B>
     N cross-boundary calls - consider refactoring or redefining
     the boundary.
   ```

---

## Phase 6: Shallow Module Deep-Dive

For each module at confidence < 0.75 in `map-tracking.json`:

1. **Run semantic queries** for the module's key concepts:
   ```
   mcp__gitnexus__query({query: "<module purpose or key concept>"})
   ```

2. **Get context on key symbols**:
   ```
   mcp__gitnexus__context({name: "<key exported symbol>"})
   ```

3. **Read process participation** - which execution flows does
   this module participate in?

4. **Update the module's DETAILED_DESIGN.md section** with
   findings - callers, execution flow participation, cross-module
   relationships that reading alone may have missed

5. **Update confidence score** in `map-tracking.json` with
   justification:
   ```json
   "<module>": {
     "confidence": 0.7,
     "notes": "Enriched via GitNexus: 12 callers found, participates in 3 flows. Bumped from 0.5.",
     "enriched_at": "<ISO-8601>"
   }
   ```

   Only bump confidence if the enrichment genuinely improved
   understanding. Adding caller counts without comprehension
   does not justify a bump.

   **Litmus test**: if the module's purpose cannot be restated
   in 1-2 sentences after enrichment, do NOT increase confidence.
   Numbers without narrative are noise.

---

## Phase 7: Update Artifacts

### 7.1 Update CONVERGENCE-REPORT.md

Re-read `map-tracking.json` (source of truth) and regenerate
`CONVERGENCE-REPORT.md` with updated confidence scores and an
enrichment summary section:

```markdown
## Enrichment Summary

_Last enrichment: <date> via GitNexus (index: <commit hash>)_

| Phase | Items Processed | Key Findings |
|-------|----------------|--------------|
| Danger zones | N entries | N upgraded to CRITICAL/HIGH |
| Extension points | N patterns | N registrations found |
| Execution flows | N flows indexed | N multi-flow hotspots |
| Clustering | N domains compared | N mismatches found |
| Shallow modules | N modules enriched | N confidence bumps |
```

### 7.2 Timestamp Annotations

All enrichment edits include a timestamp annotation:
```
(enriched <date> via GitNexus)
```

This allows `/ctx-architecture` on subsequent runs to distinguish
manual analysis from enrichment data.

### 7.3 Print Summary

Print a concise summary to the terminal:

```
Enrichment complete:
- Danger zones: N entries, N with blast radius data
- Extension points: N patterns, N total registrations
- Execution flows: N indexed, N multi-flow hotspots
- Clustering: N domains compared, N mismatches
- Shallow modules: N enriched, N confidence bumps
- Artifacts updated: DANGER-ZONES.md, EXTENSION-POINTS.md,
  CHEAT-SHEETS.md, CONVERGENCE-REPORT.md, map-tracking.json
```

---

## Deliverables

All changes are in-place edits to existing `.context/` files
plus these standalone files:

| File                        | Created by                      | Updated by enrichment              |
|-----------------------------|---------------------------------|------------------------------------|
| `DANGER-ZONES.md`           | `/ctx-architecture` (principal) | Blast radius, risk levels          |
| `EXTENSION-POINTS.md`       | This skill                      | Registration inventory             |
| `CONVERGENCE-REPORT.md`     | `/ctx-architecture`             | Updated scores, enrichment summary |
| `CHEAT-SHEETS.md`           | `/ctx-architecture`             | Execution flow index               |
| `ARCHITECTURE-PRINCIPAL.md` | `/ctx-architecture` (principal) | Clustering comparison              |
| `DETAILED_DESIGN.md`        | `/ctx-architecture`             | Shallow module updates             |
| `map-tracking.json`         | `/ctx-architecture`             | Confidence bumps, enriched_at      |

No new files are created beyond EXTENSION-POINTS.md (which only
this skill produces).

---

## Design Constraints

- **Sequential phases**: each phase MUST complete before the next
  begins. Do not interleave phases - complete one, write its
  results, then move to the next.
- **Idempotent**: running enrichment twice updates existing
  annotations, does not duplicate them. Before adding enrichment
  data, remove or update prior `(enriched <date>)` entries for
  the same symbol. Timestamp annotations make previous enrichment
  data identifiable.
- **Incremental**: each phase is independent. If one phase fails
  (e.g., no danger zones exist), skip it and continue.
- **Composable**: can chain with `/ctx-architecture --principal`
  for a full "analyze + enrich" pipeline across sessions.
- **Tool-aware**: auto-detects available tools. GitNexus is
  required; Gemini is optional (used for upstream pattern lookups
  when comparing clustering or researching extension patterns).

---

## Quality Checklist

After running, verify:
- [ ] Preflight confirmed GitNexus connected and index fresh
- [ ] Required architecture artifacts were read before enrichment
- [ ] Danger zones enriched with blast radius (d=1/d=2/d=3 counts)
- [ ] Risk levels assigned using verified criteria (not guessed)
- [ ] Extension points have file:line references (not just names)
- [ ] Execution flow index added to CHEAT-SHEETS.md
- [ ] Cross-community hotspots identified (symbols in 3+ flows)
- [ ] Domain clustering compared (mismatches surfaced)
- [ ] Shallow modules enriched only if understanding improved
- [ ] Confidence bumps justified (not inflated by raw counts)
- [ ] All edits include timestamp annotations
- [ ] CONVERGENCE-REPORT.md regenerated from map-tracking.json
- [ ] Enrichment summary printed to terminal
- [ ] DANGER-ZONES.md summary table updated with verified risk
- [ ] EXTENSION-POINTS.md written (if extension points exist)
- [ ] map-tracking.json updated with enriched_at timestamps
