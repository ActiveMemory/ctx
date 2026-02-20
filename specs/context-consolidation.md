# Context Consolidation: Entry Merging Skill

Ref: https://github.com/ActiveMemory/ctx/issues/19 (Phase 3)

## Problem

LEARNINGS.md and DECISIONS.md accumulate entries that overlap, repeat, or
cover the same topic from different sessions. Five learnings about hook edge
cases should be one consolidated entry covering all the gotchas.

Archival (Phase 5, `specs/knowledge-scaling.md`) moves *old* entries out.
Consolidation replaces *redundant* entries with denser ones. They're
complementary: archival reduces by age, consolidation reduces by content.

## Solution

A `/ctx-consolidate` skill that analyzes entries, proposes groupings,
and — with user approval — merges groups into denser combined entries.
Originals are archived, not deleted.

### Key Distinction

**Consolidation != archival.** Archival moves entries to
`.context/archive/`. Consolidation *replaces* verbose entries with tighter
ones — the file stays useful, just denser. The originals move to archive
as a paper trail.

### Workflow

```
User: /ctx-consolidate

Agent:
1. Parse all entries from LEARNINGS.md (and/or DECISIONS.md)
2. Group entries by topic similarity (keyword overlap in title + body)
3. Present consolidation candidates:

   Group 1: "Hook behavior" (5 entries)
   - [2026-01-15] Hook scripts can lose execute permission
   - [2026-01-20] Two-tier hook output is sufficient
   - [2026-02-03] Claude Code Hook Key Names
   - [2026-02-09] Agent ignores repeated hook output
   - [2026-02-16] Security docs vulnerable after migrations
   → Proposed: merge into 1 consolidated entry

   Group 2: "Path handling" (3 entries)
   - [2026-01-10] Path construction uses stdlib
   - [2026-02-05] G304 gosec false positives
   - [2026-02-16] gosec G301/G306 permissions
   → Proposed: merge into 1 consolidated entry

   Ungrouped: 12 entries (no consolidation needed)

4. User approves/modifies/rejects each group
5. For approved groups:
   a. Generate a consolidated entry (AI-written, covers all points)
   b. Add consolidated entry to the file
   c. Move originals to .context/archive/learnings-consolidated-YYYY-MM-DD.md
   d. Rebuild index
```

### Grouping Strategy

Entries are grouped by keyword overlap using the same keyword extraction
from Phase 1 (`extractTaskKeywords` logic):

1. Extract keywords from each entry's title + body
2. Build a keyword → entries map
3. Entries sharing 2+ non-trivial keywords are candidates for the same group
4. Minimum group size: 2 entries (nothing to consolidate with 1)
5. Maximum group size: 8 entries (larger groups suggest the topic needs
   splitting, not merging)

This is intentionally simple — no embeddings, no LLM calls for grouping.
The LLM call happens only when generating the consolidated entry text,
which the skill handles naturally since it runs inside an AI session.

### Consolidated Entry Format

```markdown
## [YYYY-MM-DD-HHMMSS] Hook behavior (consolidated)

**Consolidated from**: 5 entries (2026-01-15 to 2026-02-16)

- Hook scripts can lose execute permission without warning; always
  restore +x after sync operations
- Two-tier output (stdout for AI context, stderr+exit for blocks)
  is sufficient; don't over-engineer severity levels
- Claude Code hook key names are case-sensitive: PreToolUse, not
  pre_tool_use
- Agents develop repetition fatigue: vary hook output phrasing
  across invocations
- After infrastructure migrations, audit security docs first —
  stale paths in security guidance give false confidence
```

The consolidated entry:
- Uses today's timestamp
- Appends "(consolidated)" to the title
- Lists the date range of originals
- Distills each original into 1-2 lines
- Preserves all unique information (nothing is lost)

### Archive Format

Originals go to `.context/archive/learnings-consolidated-YYYY-MM-DD.md`
(or `decisions-consolidated-YYYY-MM-DD.md`). Same format as age-based
archival, distinguished by the `-consolidated-` infix.

## Skill Design

This is a **skill** (`/ctx-consolidate`), not a CLI command, because:

1. It requires AI judgment to generate consolidated entry text
2. It needs interactive approval (present candidates, get confirmation)
3. It's a low-frequency maintenance operation, not a daily workflow

### Skill File Structure

```
.claude/skills/ctx-consolidate/SKILL.md
```

### Skill Behavior

The skill should:

1. **Parse entries** using `ctx` CLI or direct file reading
2. **Group by similarity** using keyword overlap
3. **Present candidates** with clear before/after preview
4. **Wait for approval** before modifying any files
5. **Execute approved merges**:
   - Write consolidated entry at the top of the file (newest first)
   - Append originals to archive file
   - Remove originals from source file
   - Run `ctx reindex` (decisions or learnings) to rebuild index
6. **Report results**: "Consolidated 5 entries into 1. Originals archived."

### What the Skill Does NOT Do

- Automatic consolidation (always requires user approval)
- Cross-file consolidation (learnings stay learnings, decisions stay decisions)
- Semantic understanding (uses keyword matching, not embeddings)
- Delete entries (always archives originals)

## Non-Goals

- **LLM-based grouping**: Keyword overlap is sufficient for structured
  Markdown entries with consistent formatting. Adding LLM calls for
  grouping adds latency and cost for marginal improvement.
- **Automatic triggers**: Consolidation is judgment-heavy. The drift
  nudge (Phase 2) tells users *when* to consolidate; this skill handles
  *how*.
- **Cross-file merging**: A learning about "hook edge cases" stays in
  LEARNINGS.md even if there's a related decision. Files have distinct
  purposes.
- **Undo**: Originals are archived, which serves as the undo mechanism.
  Restoring means copying entries back from archive.

## Dependencies

- **Phase 1 (Smart Retrieval)**: Provides `extractTaskKeywords` for
  keyword extraction. Already implemented.
- **Phase 2 (Drift Nudges)**: Provides the nudge that suggests running
  `/ctx-consolidate`. Should be implemented first so the workflow is
  complete.
- **Phase 5 (Knowledge Scaling)**: Provides `ctx learnings archive` and
  `ctx decisions archive` for the archival step. The skill can use these
  commands directly or replicate the archive logic.

## Testing

The skill itself is a Markdown prompt file — no Go tests needed. But the
grouping logic should be tested if extracted into Go code:

- Entries with 2+ shared keywords are grouped together
- Entries with 0-1 shared keywords are not grouped
- Groups larger than 8 are split
- Single-entry groups are listed as "ungrouped"
- Empty files produce no candidates
- Entries with no meaningful keywords (all stop words) are ungrouped

## Implementation Order

1. Write the skill file (`SKILL.md`) with the full prompt
2. Test manually on a project with 20+ learnings
3. Iterate on grouping quality and consolidated entry format
4. Document in `docs/skills.md` and `docs/cli-reference.md`
