---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Context Management
icon: lucide/layers
---

![ctx](../images/ctx-banner.png)

### Adding entries

Each context-artifact noun (`task`, `decision`, `learning`,
`convention`) owns its own `add` subcommand under the
noun-first command tree:

```bash
ctx task add <content> [flags]
ctx decision add <content> [flags]
ctx learning add <content> [flags]
ctx convention add <content> [flags]
```

**Target files**:

| Subcommand              | Target File      |
|-------------------------|------------------|
| `ctx task add`          | `TASKS.md`       |
| `ctx decision add`      | `DECISIONS.md`   |
| `ctx learning add`      | `LEARNINGS.md`   |
| `ctx convention add`    | `CONVENTIONS.md` |

**Flags** (shared by every `add` subcommand; per-noun
required-flag rules surface as command errors):

| Flag                      | Short | Description                                                 |
|---------------------------|-------|-------------------------------------------------------------|
| `--priority <level>`      | `-p`  | Priority for tasks: `high`, `medium`, `low`                 |
| `--section <name>`        | `-s`  | Target section within file                                  |
| `--context`               | `-c`  | Context (required for decisions and learnings)              |
| `--rationale`             | `-r`  | Rationale for decisions (required for decisions)            |
| `--consequence`           |       | Consequence for decisions (required for decisions)          |
| `--lesson`                | `-l`  | Key insight (required for learnings)                        |
| `--application`           | `-a`  | How to apply going forward (required for learnings)         |
| `--file`                  | `-f`  | Read content from file instead of argument                  |

**Examples**:

```bash
# Add a task
ctx task add "Implement user authentication" \
  --session-id abc12345 --branch main --commit 68fbc00a
ctx task add "Fix login bug" --priority high \
  --session-id abc12345 --branch main --commit 68fbc00a

# Record a decision (requires all ADR (Architectural Decision Record) fields)
ctx decision add "Use PostgreSQL for primary database" \
  --context "Need a reliable database for production" \
  --rationale "PostgreSQL offers ACID compliance and JSON support" \
  --consequence "Team needs PostgreSQL training" \
  --session-id abc12345 --branch main --commit 68fbc00a

# Note a learning (requires context, lesson, and application)
ctx learning add "Vitest mocks must be hoisted" \
  --context "Tests failed with undefined mock errors" \
  --lesson "Vitest hoists vi.mock() calls to top of file" \
  --application "Always place vi.mock() before imports in test files" \
  --session-id abc12345 --branch main --commit 68fbc00a

# Add to specific section
ctx convention add "Use kebab-case for filenames" --section "Naming"
```

---

### `ctx drift`

Detect stale or invalid context.

```bash
ctx drift [flags]
```

**Flags**:

| Flag     | Description                  |
|----------|------------------------------|
| `--json` | Output machine-readable JSON |
| `--fix`  | Auto-fix simple issues       |

**Checks**:

* Path references in `ARCHITECTURE.md` and `CONVENTIONS.md` exist
* Task references are valid
* Constitution rules aren't violated (*heuristic*)
* Staleness indicators (*old files, many completed tasks*)
* Missing packages: warns when `internal/` directories exist on disk but are
  not referenced in `ARCHITECTURE.md` (*suggests running `/ctx-architecture`*)
* Entry count: warns when `LEARNINGS.md` or `DECISIONS.md` exceed configurable
  thresholds (*default: 30 learnings, 20 decisions*), or when `CONVENTIONS.md`
  exceeds a line count threshold (default: 200). Configure via `.ctxrc`:
  ```yaml
  entry_count_learnings: 30      # warn above this (0 = disable)
  entry_count_decisions: 20      # warn above this (0 = disable)
  convention_line_count: 200     # warn above this (0 = disable)
  ```

**Example**:

```bash
ctx drift
ctx drift --json
ctx drift --fix
```

**Exit codes**:

| Code | Meaning           |
|------|-------------------|
| 0    | All checks passed |
| 1    | Warnings found    |
| 3    | Violations found  |

---

### `ctx sync`

Reconcile context with the current codebase state.

```bash
ctx sync [flags]
```

**Flags**:

| Flag        | Description                              |
|-------------|------------------------------------------|
| `--dry-run` | Show what would change without modifying |

**What it does:**

* Scans codebase for structural changes
* Compares with ARCHITECTURE.md
* Suggests documenting dependencies if package files exist
* Identifies stale or outdated context

**Example**:

```bash
ctx sync
ctx sync --dry-run
```

---

### `ctx compact`

Consolidate and clean up context files.

* Moves completed tasks older than 7 days to the archive
* Removes empty sections

```bash
ctx compact [flags]
```

**Flags**:

| Flag             | Description                                |
|------------------|--------------------------------------------|
| `--archive`      | Create `.context/archive/` for old content |

**Example**:

```bash
ctx compact
ctx compact --archive
```

---

### `ctx fmt`

Format context files to a consistent line width.

Wraps long lines in `TASKS.md`, `DECISIONS.md`, `LEARNINGS.md`, and
`CONVENTIONS.md` at word boundaries. Markdown list items get 2-space
continuation indent. Headings, tables, frontmatter, and HTML comments
are preserved as-is.

Idempotent: running twice produces the same output.

```bash
ctx fmt [flags]
```

**Flags**:

| Flag        | Type  | Default | Description                                |
|-------------|-------|---------|--------------------------------------------|
| `--width`   | `int` | `80`    | Target line width                          |
| `--check`   | `bool`| `false` | Check only, exit 1 if files would change   |

**Examples**:

```bash
ctx fmt              # format all context files
ctx fmt --check      # CI mode: check without modifying
ctx fmt --width 100  # custom width
```

Also available as a Makefile target:

```bash
make fmt-context
```

---

### `ctx task`

Manage task completion, archival, and snapshots.

```bash
ctx task <subcommand>
```

#### `ctx task complete`

Mark a task as completed.

```bash
ctx task complete <task-id-or-text>
```

**Arguments**:

* `task-id-or-text`: Task number or partial text match

**Examples**:

```bash
# By text (partial match)
ctx task complete "user auth"

# By task number
ctx task complete 3
```

#### `ctx task archive`

Move completed tasks from `TASKS.md` to a timestamped archive file.

```bash
ctx task archive [flags]
```

**Flags**:

| Flag        | Description                              |
|-------------|------------------------------------------|
| `--dry-run` | Preview changes without modifying files  |

Archive files are stored in `.context/archive/` with timestamped names
(`tasks-YYYY-MM-DD.md`). Completed tasks (marked with `[x]`) are moved;
pending tasks (`[ ]`) remain in `TASKS.md`.

**Example**:

```bash
ctx task archive
ctx task archive --dry-run
```

#### `ctx task snapshot`

Create a point-in-time snapshot of `TASKS.md` without modifying the original.

```bash
ctx task snapshot [name]
```

**Arguments**:

- `name`: Optional name for the snapshot (defaults to "snapshot")

Snapshots are stored in `.context/archive/` with timestamped names
(`tasks-<name>-YYYY-MM-DD-HHMM.md`).

**Example**:

```bash
ctx task snapshot
ctx task snapshot "before-refactor"
```

---

### `ctx permission`

Manage Claude Code permission snapshots.

```bash
ctx permission <subcommand>
```

#### `ctx permission snapshot`

Save `.claude/settings.local.json` as the golden image.

```bash
ctx permission snapshot
```

Creates `.claude/settings.golden.json` as a byte-for-byte copy of the
current settings. Overwrites if the golden file already exists.

The golden file is meant to be committed to version control and shared
with the team.

**Example**:

```bash
ctx permission snapshot
# Saved golden image: .claude/settings.golden.json
```

#### `ctx permission restore`

Replace `settings.local.json` with the golden image.

```bash
ctx permission restore
```

Prints a diff of dropped (session-accumulated) and restored permissions.
No-op if the files already match.

**Example**:

```bash
ctx permission restore
# Dropped 3 session permission(s):
#   - Bash(cat /tmp/debug.log:*)
#   - Bash(rm /tmp/test-*:*)
#   - Bash(curl https://example.com:*)
# Restored from golden image.
```

---

### `ctx reindex`

Regenerate the quick-reference index for both `DECISIONS.md` and `LEARNINGS.md`
in a single invocation.

```bash
ctx reindex
```

This is a convenience wrapper around `ctx decision reindex` and
`ctx learning reindex`. Both files grow at similar rates and users
typically want to reindex both after manual edits.

The index is a compact table of date and title for each entry, allowing
AI tools to scan entries without reading the full file.

**Example**:

```bash
ctx reindex
# ✓ Index regenerated with 12 entries
# ✓ Index regenerated with 8 entries
```

---

### `ctx decision`

Manage the `DECISIONS.md` file.

```bash
ctx decision <subcommand>
```

#### `ctx decision reindex`

Regenerate the quick-reference index at the top of `DECISIONS.md`.

```bash
ctx decision reindex
```

The index is a compact table showing the date and title for each decision,
allowing AI tools to quickly scan entries without reading the full file.

Use this after manual edits to `DECISIONS.md` or when migrating existing
files to use the index format.

**Example**:

```bash
ctx decision reindex
# ✓ Index regenerated with 12 entries
```

---

### `ctx learning`

Manage the `LEARNINGS.md` file.

```bash
ctx learning <subcommand>
```

#### `ctx learning reindex`

Regenerate the quick-reference index at the top of `LEARNINGS.md`.

```bash
ctx learning reindex
```

The index is a compact table showing the date and title for each learning,
allowing AI tools to quickly scan entries without reading the full file.

Use this after manual edits to `LEARNINGS.md` or when migrating existing
files to use the index format.

**Example**:

```bash
ctx learning reindex
# ✓ Index regenerated with 8 entries
```
