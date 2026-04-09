---
name: ctx-consolidate
description: "Consolidate redundant entries in LEARNINGS.md or DECISIONS.md. Use when context files grow large with overlapping entries."
tools: [bash, read, write, edit]
---

Find and merge duplicate or overlapping entries in context files.

## When to Use

- Context files have grown large (50+ entries)
- Multiple entries cover the same topic from different sessions
- After a long project phase where many similar learnings accumulated
- When `ctx status` shows high token counts for context files

## When NOT to Use

- Files are small and manageable
- Entries are all distinct
- Just after a fresh `ctx init`

## Process

### 1. Read the target file

Read the full content of the file to consolidate
(LEARNINGS.md or DECISIONS.md).

### 2. Identify clusters

Group entries by topic. Look for:
- Same subject with different wording
- Entries that build on each other chronologically
- Contradictory entries (later one supersedes)

### 3. Propose merges

For each cluster, propose a consolidated entry:

> **Cluster: [topic]** (N entries → 1)
> - Entry A: "..."
> - Entry B: "..."
> - **Merged**: "..."
>
> Approve?

### 4. Apply approved merges

Replace the cluster entries with the merged version.
Archive originals if requested.

## Quality Checklist

- [ ] No information lost in merges
- [ ] Contradictions resolved (latest wins)
- [ ] User approved each merge
- [ ] File is valid markdown after edits
