---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Context Management
icon: lucide/layers
---

### `ctx add`

Add a new item to a context file.

```bash
ctx add <type> <content> [flags]
```

**Types**:

| Type         | Target File      |
|--------------|------------------|
| `task`       | `TASKS.md`       |
| `decision`   | `DECISIONS.md`   |
| `learning`   | `LEARNINGS.md`   |
| `convention` | `CONVENTIONS.md` |

**Flags**:

| Flag                      | Short | Description                                                 |
|---------------------------|-------|-------------------------------------------------------------|
| `--priority <level>`      | `-p`  | Priority for tasks: `high`, `medium`, `low`                 |
| `--section <name>`        | `-s`  | Target section within file                                  |
| `--context`               | `-c`  | Context (required for decisions and learnings)              |
| `--rationale`             | `-r`  | Rationale for decisions (required for decisions)            |
| `--consequences`          |       | Consequences for decisions (required for decisions)         |
| `--lesson`                | `-l`  | Key insight (required for learnings)                        |
| `--application`           | `-a`  | How to apply going forward (required for learnings)         |
| `--file`                  | `-f`  | Read content from file instead of argument                  |

**Examples**:

```bash
# Add a task
ctx add task "Implement user authentication"
ctx add task "Fix login bug" --priority high

# Record a decision (requires all ADR (Architectural Decision Record) fields)
ctx add decision "Use PostgreSQL for primary database" \
  --context "Need a reliable database for production" \
  --rationale "PostgreSQL offers ACID compliance and JSON support" \
  --consequences "Team needs PostgreSQL training"

# Note a learning (requires context, lesson, and application)
ctx add learning "Vitest mocks must be hoisted" \
  --context "Tests failed with undefined mock errors" \
  --lesson "Vitest hoists vi.mock() calls to top of file" \
  --application "Always place vi.mock() before imports in test files"

# Add to specific section
ctx add convention "Use kebab-case for filenames" --section "Naming"
```

---

### `ctx complete`

Mark a task as completed.

```bash
ctx complete <task-id-or-text>
```

**Arguments**:

* `task-id-or-text`: Task number or partial text match

**Examples**:

```bash
# By text (partial match)
ctx complete "user auth"

# By task number
ctx complete 3
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
  not referenced in `ARCHITECTURE.md` (*suggests running `/ctx-map`*)
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

### `ctx tasks`

Manage task archival and snapshots.

```bash
ctx tasks <subcommand>
```

#### `ctx tasks archive`

Move completed tasks from `TASKS.md` to a timestamped archive file.

```bash
ctx tasks archive [flags]
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
ctx tasks archive
ctx tasks archive --dry-run
```

#### `ctx tasks snapshot`

Create a point-in-time snapshot of `TASKS.md` without modifying the original.

```bash
ctx tasks snapshot [name]
```

**Arguments**:

- `name`: Optional name for the snapshot (defaults to "snapshot")

Snapshots are stored in `.context/archive/` with timestamped names
(`tasks-<name>-YYYY-MM-DD-HHMM.md`).

**Example**:

```bash
ctx tasks snapshot
ctx tasks snapshot "before-refactor"
```

---

### `ctx permissions`

Manage Claude Code permission snapshots.

```bash
ctx permissions <subcommand>
```

#### `ctx permissions snapshot`

Save `.claude/settings.local.json` as the golden image.

```bash
ctx permissions snapshot
```

Creates `.claude/settings.golden.json` as a byte-for-byte copy of the
current settings. Overwrites if the golden file already exists.

The golden file is meant to be committed to version control and shared
with the team.

**Example**:

```bash
ctx permissions snapshot
# Saved golden image: .claude/settings.golden.json
```

#### `ctx permissions restore`

Replace `settings.local.json` with the golden image.

```bash
ctx permissions restore
```

Prints a diff of dropped (session-accumulated) and restored permissions.
No-op if the files already match.

**Example**:

```bash
ctx permissions restore
# Dropped 3 session permission(s):
#   - Bash(cat /tmp/debug.log:*)
#   - Bash(rm /tmp/test-*:*)
#   - Bash(curl https://example.com:*)
# Restored from golden image.
```

---

### `ctx decisions`

Manage the `DECISIONS.md` file.

```bash
ctx decisions <subcommand>
```

#### `ctx decisions reindex`

Regenerate the quick-reference index at the top of `DECISIONS.md`.

```bash
ctx decisions reindex
```

The index is a compact table showing the date and title for each decision,
allowing AI tools to quickly scan entries without reading the full file.

Use this after manual edits to `DECISIONS.md` or when migrating existing
files to use the index format.

**Example**:

```bash
ctx decisions reindex
# ✓ Index regenerated with 12 entries
```

---

### `ctx learnings`

Manage the `LEARNINGS.md` file.

```bash
ctx learnings <subcommand>
```

#### `ctx learnings reindex`

Regenerate the quick-reference index at the top of `LEARNINGS.md`.

```bash
ctx learnings reindex
```

The index is a compact table showing the date and title for each learning,
allowing AI tools to quickly scan entries without reading the full file.

Use this after manual edits to `LEARNINGS.md` or when migrating existing
files to use the index format.

**Example**:

```bash
ctx learnings reindex
# ✓ Index regenerated with 8 entries
```
