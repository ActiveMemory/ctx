---
name: ctx-journal-enrich-all
description: "Full journal pipeline: import unimported sessions, then batch-enrich all unenriched entries. Use when the user says 'process the journal' or to catch up on the backlog."
allowed-tools: Bash(ctx:*), Read, Glob, Grep, Edit, Write, Task
---

Full journal pipeline: import if needed, then batch-enrich.

## When to Use

- When the user says "enrich everything" or "process the journal"
- When there is a backlog of unenriched or unimported sessions
- Periodically to catch up on recent sessions
- After the `check-journal` hook reports unimported or unenriched entries

## When NOT to Use

- For a single specific session (use `/ctx-journal-enrich` instead)

## Process

### Step 0: Import If Needed

Before enriching, check whether there are unimported sessions. If
the journal directory has no `.md` files at all, or if there are
`.jsonl` session files newer than the newest journal entry, import
them first.

```bash
CTX_DIR=$(ctx system bootstrap -q)
JOURNAL_DIR="$CTX_DIR/journal"

# Check if any .md files exist
md_count=$(ls "$JOURNAL_DIR"/*.md 2>/dev/null | wc -l)

if [ "$md_count" -eq 0 ]; then
  echo "No journal entries found: importing all sessions."
  ctx recall import --all --yes
else
  # Compare newest .md mtime against .jsonl files
  newest_md=$(stat -c %Y $(ls -t "$JOURNAL_DIR"/*.md | head -1))
  unimported=$(find ~/.claude/projects -name "*.jsonl" -newermt @${newest_md} 2>/dev/null | wc -l)
  if [ "$unimported" -gt 0 ]; then
    echo "$unimported unimported session(s) found: importing first."
    ctx recall import --all --yes
  fi
fi
```

Report how many sessions were imported (or "none needed") before
moving to enrichment.

### Step 1: Find Unenriched Entries

List all journal entries that lack enrichment using the state file:

```bash
# List .md files in journal dir and check state
CTX_DIR=$(ctx system bootstrap -q)
for f in "$CTX_DIR/journal/"*.md; do
  name=$(basename "$f")
  ctx system mark-journal --check "$name" enriched || echo "$f"
done
```

Or read `.state.json` in the journal directory directly and list
entries without an `enriched` date set.

### Fallback: Detect Enrichment from Frontmatter

If `mark-journal --check` is unavailable (no state file, command
fails), fall back to frontmatter inspection. An entry is considered
**already enriched** if its YAML frontmatter contains **both** `type`
and `outcome` fields: these are set exclusively by enrichment, never
by import.

Do NOT use `title` or `date` to detect enrichment: those are always
present from import. The enrichment-only fields are:

| Field          | Set by         |
|----------------|----------------|
| `title`        | Import         |
| `date`         | Import         |
| `time`         | Import         |
| `model`        | Import         |
| `tokens_in`    | Import         |
| `tokens_out`   | Import         |
| `session_id`   | Import         |
| `project`      | Import         |
| `type`         | **Enrichment** |
| `outcome`      | **Enrichment** |
| `topics`       | **Enrichment** |
| `technologies` | **Enrichment** |
| `summary`      | **Enrichment** |

If all entries already have enrichment recorded, report that and stop.

### Step 2: Filter Out Noise

Skip entries that are not worth enriching:

- **Locked entries**: a file is locked if `.state.json` has a
  `locked` date OR the frontmatter contains `locked: true`. Never
  modify locked files: neither metadata nor body. Check via:
  `ctx system mark-journal --check <filename> locked`
  or look for `locked: true` in the YAML frontmatter.
- **Suggestion sessions**: files under ~20 lines or containing
  only auto-complete fragments. Check with:
  ```bash
  wc -l <file>
  ```
- **Multi-part continuations**: files ending in `-p2.md`, `-p3.md`
  etc. Enrich only the first part; continuation parts inherit
  the frontmatter topic.

Report how many entries will be processed and how many were
filtered out.

### Step 3: Process Each Entry

For each entry, read the conversation and extract:

1. **Title**: a short descriptive title for the session
2. **Type**: feature, bugfix, refactor, exploration, debugging,
   or documentation
3. **Outcome**: completed, partial, abandoned, or blocked
4. **Topics**: 2-5 topic tags
5. **Technologies**: languages, frameworks, tools used
6. **Summary**: 2-3 sentences describing what was accomplished

Apply YAML frontmatter to each file:

```yaml
---
title: "Session title"
date: 2026-01-27
type: feature
outcome: completed
topics:
  - authentication
  - caching
technologies:
  - go
  - redis
---
```

### Step 4: Mark Enriched

After writing frontmatter to each file, update the state file:

```bash
ctx system mark-journal <filename> enriched
```

### Step 5: Report

After processing, report:

- How many sessions were imported (or "none needed")
- How many entries were enriched
- How many were skipped (already enriched, too short, etc.)
- Remind the user to rebuild: `ctx journal site --build`

## Confirmation Mode

**Interactive** (default when user is present): show a summary
of proposed enrichments before applying. Group by type/outcome
so the user can scan quickly rather than reviewing one by one.

**Unattended** (when running in a loop or explicitly told
"just do it"): apply enrichments directly and report results.

## Large Backlogs (20+ entries)

For large backlogs, use the heuristic enrichment script bundled
in `references/enrich-heuristic.py`. This script infers type,
outcome, topics, and technologies from the title and filename
patterns, then inserts frontmatter and marks state automatically.

### How to use

1. Build a file list of eligible entries (non-multipart, 20+ lines,
   missing `type:` and `outcome:` fields):
   ```bash
   CTX_DIR=$(ctx system bootstrap -q)
   for f in "$CTX_DIR"/journal/*.md; do
     [ -f "$f" ] || continue
     has_type=$(head -30 "$f" | grep -c '^type:' || true)
     has_outcome=$(head -30 "$f" | grep -c '^outcome:' || true)
     if [ "$has_type" -eq 0 ] || [ "$has_outcome" -eq 0 ]; then
       name=$(basename "$f")
       case "$name" in *-p[0-9].md|*-p[0-9][0-9].md) continue ;; esac
       lines=$(wc -l < "$f")
       [ "$lines" -ge 20 ] && echo "$f"
     fi
   done > /tmp/enrich-list.txt
   ```

2. Run the heuristic enrichment script. The script path is relative
   to this skill's directory: copy it to /tmp or reference it via
   the full embedded path:
   ```bash
   python3 references/enrich-heuristic.py /tmp/enrich-list.txt
   ```

3. The script handles everything: reads files, inserts frontmatter,
   runs `ctx system mark-journal` for each, and reports counts.

### When to use heuristic vs. per-file enrichment

| Backlog size | Approach                                          |
|--------------|---------------------------------------------------|
| 1-5 entries  | Read each file, enrich manually with full context |
| 6-20 entries | Sequential processing in the main conversation    |
| 20+ entries  | Use `enrich-heuristic.py` for bulk processing     |

The heuristic script produces good-enough enrichment from titles
and filenames. For higher quality, follow up with manual review
of entries where the type or topics look wrong.

Subagent parallelization is an alternative for 20+ entries, but
requires that subagents have Edit and Bash permissions granted.
If permissions are restricted, the heuristic script is faster
and more reliable.

## Quality Checklist

- [ ] Unimported sessions detected and imported before enrichment
- [ ] Suggestion sessions and multi-part continuations filtered
- [ ] Each enriched entry has all required frontmatter fields
- [ ] Summary is specific to the session, not generic
- [ ] User was shown a summary before applying (unless unattended)
- [ ] State file updated for each enriched entry
