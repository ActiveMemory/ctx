---
name: ctx-consolidate
description: "Consolidate redundant entries in LEARNINGS.md or DECISIONS.md. Use when ctx drift reports high entry counts or entries overlap."
allowed-tools: Bash(ctx:*), Read, Edit, Write
---

Analyze entries in LEARNINGS.md and/or DECISIONS.md, group overlapping
entries by topic, and — with user approval — merge groups into denser
consolidated entries. Originals are archived, not deleted.

## Key Distinction

**Consolidation != archival.** Archival moves old entries to
the archive directory. Consolidation *replaces* verbose entries with
tighter ones — the file stays useful, just denser. The originals move
to archive as a paper trail.

## When to Use

- When `ctx drift` reports entry counts above threshold
  (default: 30 learnings, 20 decisions)
- When you notice 3+ entries about the same topic
- When the user asks "clean up learnings", "consolidate context",
  "reduce noise in decisions"
- Before a release, to keep context lean

## When NOT to Use

- When there are fewer than 10 entries (nothing meaningful to group)
- When the user wants to *delete* entries (offer archival instead)
- Automatically — always require user approval before modifying files
- Mid-task when the user is focused on shipping

## Execution

### Step 1: Parse Entries

Read the target file(s):

```bash
# Check entry counts first
ctx drift --json
```

Then read the files directly:
- LEARNINGS.md (in the context directory)
- DECISIONS.md (in the context directory)

Parse entries by their `## [YYYY-MM-DD-HHMMSS] Title` headers. Each
entry extends from its header to the line before the next header or
end of file.

### Step 2: Extract Keywords and Group

For each entry, extract keywords from its title and body:

1. Split text on whitespace and punctuation
2. Lowercase everything
3. Filter out stop words (the, and, for, with, from, are, was, etc.)
   and words shorter than 3 characters
4. Deduplicate

Build a keyword-to-entries map. Entries sharing **2 or more
non-trivial keywords** are candidates for the same group.

**Grouping rules:**
- Minimum group size: 2 entries (nothing to consolidate with 1)
- Maximum group size: 8 entries (larger groups suggest the topic
  needs splitting, not merging)
- An entry can only belong to one group (assign to the best match)

### Step 3: Present Candidates

Show the user what you found. Format:

```
Consolidation candidates for LEARNINGS.md:

Group 1: "Hook behavior" (5 entries)
  - [2026-01-15] Hook scripts can lose execute permission
  - [2026-01-20] Two-tier hook output is sufficient
  - [2026-02-03] Claude Code Hook Key Names
  - [2026-02-09] Agent ignores repeated hook output
  - [2026-02-16] Security docs vulnerable after migrations
  -> Proposed: merge into 1 consolidated entry

Group 2: "Path handling" (3 entries)
  - [2026-01-10] Path construction uses stdlib
  - [2026-02-05] G304 gosec false positives
  - [2026-02-16] gosec G301/G306 permissions
  -> Proposed: merge into 1 consolidated entry

Ungrouped: 12 entries (no consolidation needed)
```

**Wait for the user to approve, modify, or reject each group.**
Do NOT proceed without explicit confirmation.

### Step 4: Generate Consolidated Entries

For each approved group, write a consolidated entry that:

- Uses today's timestamp in `YYYY-MM-DD-HHMMSS` format
- Appends "(consolidated)" to the title
- Lists the date range of originals in a `**Consolidated from**` line
- Distills each original into 1-2 lines
- **Preserves all unique information** (nothing is lost)

**Format:**

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

### Step 5: Execute Approved Merges

For each approved group:

1. **Add the consolidated entry** at the top of the file (below
   the `# Learnings` or `# Decisions` header)
2. **Remove the original entries** from the source file
3. **Append originals to archive** at
   `archive/learnings-consolidated-YYYY-MM-DD.md` in the context
   directory (or `decisions-consolidated-YYYY-MM-DD.md`)
4. **Rebuild the index**:

```bash
ctx learnings reindex
# or
ctx decisions reindex
```

### Step 6: Report Results

```
Consolidated LEARNINGS.md:
  - Group "Hook behavior": 5 entries -> 1 (originals archived)
  - Group "Path handling": 3 entries -> 1 (originals archived)
  Total: 8 entries consolidated into 2. File reduced from 47 to 41 entries.
  Archive: archive/learnings-consolidated-2026-02-19.md (in context dir)
```

## Archive Format

The archive file uses the same Markdown format as the source file.
Each archived entry keeps its original timestamp and content,
preceded by a header noting which consolidated entry replaced it:

```markdown
# Archived Learnings (consolidated 2026-02-19)

Originals replaced by consolidated entries in LEARNINGS.md.

## Group: Hook behavior

## [2026-01-15-120000] Hook scripts can lose execute permission
(original content preserved verbatim)

## [2026-01-20-093000] Two-tier hook output is sufficient
(original content preserved verbatim)
```

## What This Skill Does NOT Do

- **Automatic consolidation**: always requires user approval
- **Cross-file consolidation**: learnings stay in LEARNINGS.md,
  decisions stay in DECISIONS.md
- **Delete entries**: always archives originals as a paper trail
- **Semantic understanding via embeddings**: uses keyword matching,
  which is sufficient for structured entries with consistent formatting
- **Consolidate TASKS.md or CONVENTIONS.md**: use `ctx tasks archive`
  for tasks; conventions rarely need consolidation

## Quality Checklist

Before reporting results:
- [ ] Presented all candidate groups before making changes
- [ ] Waited for explicit user approval per group
- [ ] Each consolidated entry preserves all unique information
- [ ] Original entries are archived, not deleted
- [ ] Ran `ctx reindex` after modifications
- [ ] Reported what changed and where archives were written
