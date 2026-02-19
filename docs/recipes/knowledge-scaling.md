---
title: "Managing Knowledge at Scale"
icon: lucide/archive
---

![ctx](../images/ctx-banner.png)

## Problem

`DECISIONS.md` and `LEARNINGS.md` grow monotonically. Every session can add
entries, but nothing removes them. After months of active development, these
files consume significant token budget, and AI tools spend attention on
decisions that were superseded weeks ago or learnings about bugs that have
long been fixed.

**How do you keep knowledge files lean without losing history?**

## Commands and Skills Used

| Tool                       | Type    | Purpose                                     |
|----------------------------|---------|---------------------------------------------|
| `ctx decisions archive`    | Command | Archive old or superseded decisions         |
| `ctx learnings archive`    | Command | Archive old or superseded learnings         |
| `ctx compact --archive`    | Command | Archive tasks, decisions, and learnings     |
| `ctx decisions reindex`    | Command | Rebuild the quick-reference index           |
| `ctx learnings reindex`    | Command | Rebuild the quick-reference index           |
| `/ctx-archive`             | Skill   | AI-guided archival                          |
| `/ctx-drift`               | Skill   | Detect staleness and bloat                  |

## The Workflow

### Step 1: Check File Size

Start by understanding how large your knowledge files have grown:

```bash
ctx status --verbose
```

Or ask your agent:

```text
How many decisions and learnings do we have? Are any stale?
```

### Step 2: Preview with Dry-Run

Always preview before archiving:

```bash
ctx decisions archive --dry-run
ctx learnings archive --dry-run
```

This shows which entries would be archived and why (old or superseded),
without modifying any files.

### Step 3: Archive

Archive entries older than the threshold (default 90 days):

```bash
ctx decisions archive
ctx learnings archive
```

Or archive everything in one pass:

```bash
ctx compact --archive
```

Archived entries are written to `.context/archive/decisions-YYYY-MM-DD.md`
and `.context/archive/learnings-YYYY-MM-DD.md`. The source files are
cleaned up and reindexed automatically.

### Step 4: Configure Thresholds

Adjust the defaults in `.contextrc`:

```yaml
archive_knowledge_after_days: 90   # Days before archiving (default: 90)
archive_keep_recent: 5             # Recent entries to always keep (default: 5)
```

For fast-moving projects, lower the threshold:

```yaml
archive_knowledge_after_days: 30
archive_keep_recent: 3
```

For stable projects with long-lived decisions, raise it:

```yaml
archive_knowledge_after_days: 180
archive_keep_recent: 10
```

### Step 5: Supersede Decisions

When a decision is replaced by a newer one, mark the old entry as
superseded by adding a line to its body:

```markdown
~~Superseded by [2026-02-15-120000] Use Redis instead of Memcached~~
```

Superseded entries are archived regardless of age, even if they are
within the `--days` threshold.

### Step 6: Auto-Archive via Compact

If `auto_archive: true` is set in `.contextrc`, running `ctx compact`
automatically archives old decisions and learnings alongside tasks:

```bash
ctx compact
```

This is the easiest way to keep all context files lean as part of
regular maintenance.

## The Conversational Approach

```text
You: "Our DECISIONS.md is getting long. Archive the old ones."

Agent: "I'll preview first... 8 decisions are older than 90 days and
       2 are marked as superseded. The 5 most recent will be kept.
       Want me to proceed?"

You: "Yes, go ahead."

Agent: "Done. Archived 10 entries to .context/archive/decisions-2026-02-19.md.
       DECISIONS.md now has 5 entries. Index regenerated."
```

## CLI Reference

```bash
# Preview what would be archived
ctx decisions archive --dry-run
ctx learnings archive --dry-run

# Archive with defaults (90 days, keep 5)
ctx decisions archive
ctx learnings archive

# Custom thresholds
ctx decisions archive --days 30 --keep 3
ctx learnings archive --days 60 --keep 10

# Archive all except most recent
ctx decisions archive --all --keep 3
ctx learnings archive --all --keep 5

# One-pass archival via compact
ctx compact --archive
```

## Tips

* **Preview first.** Always use `--dry-run` before a large archive operation.
  It costs nothing and prevents surprises.
* **Supersede, don't delete.** When a decision is replaced, mark it as
  superseded. The archive command picks it up automatically.
* **Keep recent entries generous.** The `--keep` flag protects the most recent
  entries. When in doubt, keep more than you think you need.
* **Archive files are searchable.** Everything moves to `.context/archive/`,
  not oblivion. Use `grep` or your editor to search archived entries.
* **Reindex happens automatically.** After archival, the quick-reference index
  in the source file is regenerated. No need to run `reindex` separately.
* **Compact is the simplest path.** If you just want everything clean,
  `ctx compact --archive` handles tasks, decisions, and learnings in one pass.

## Next Up

**[Detecting and Fixing Drift](context-health.md)**:
Keep context files accurate as your codebase evolves.

## See Also

* [Persisting Decisions, Learnings, and Conventions](knowledge-capture.md):
  how to create the entries that this recipe archives
* [Tracking Work Across Sessions](task-management.md): task archival with
  `ctx tasks archive`
* [CLI Reference](../cli-reference.md): full flag documentation for
  `ctx decisions archive`, `ctx learnings archive`, and `ctx compact`
