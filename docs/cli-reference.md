---
#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: CLI Reference
icon: lucide/terminal
---

![ctx](images/ctx-banner.png)

## `ctx` CLI

This is a complete reference for all `ctx` commands.

## Global Options

All commands support these flags:

| Flag                   | Description                                       |
|------------------------|---------------------------------------------------|
| `--help`               | Show command help                                 |
| `--version`            | Show version                                      |
| `--context-dir <path>` | Override context directory (default: `.context/`) |
| `--no-color`           | Disable colored output                            |

> The `NO_COLOR=1` environment variable also disables colored output.

## Commands

### `ctx init`

Initialize a new `.context/` directory with template files.

```bash
ctx init [flags]
```

**Flags**:

| Flag        | Short | Description                                                           |
|-------------|-------|-----------------------------------------------------------------------|
| `--force`   | `-f`  | Overwrite existing context files                                      |
| `--minimal` | `-m`  | Only create essential files (TASKS.md, DECISIONS.md, CONSTITUTION.md) |
| `--merge`   |       | Auto-merge ctx content into existing CLAUDE.md                        |

**Creates**:

- `.context/` directory with all template files
- `.claude/hooks/` with auto-save and enforcement scripts (for Claude Code)
- `.claude/commands/` with ctx slash command definitions
- `.claude/settings.local.json` with hook configuration and pre-approved ctx permissions
- `CLAUDE.md` with bootstrap instructions (or merges into existing)

**Example**:

```bash
# Standard initialization
ctx init

# Minimal setup (just core files)
ctx init --minimal

# Force overwrite existing
ctx init --force

# Merge into existing CLAUDE.md
ctx init --merge
```

---

### `ctx status`

Show the current context summary.

```bash
ctx status [flags]
```

**Flags**:

| Flag        | Short | Description                   |
|-------------|-------|-------------------------------|
| `--json`    |       | Output as JSON                |
| `--verbose` | `-v`  | Include file contents summary |

**Output**:

- Context directory path
- Total files and token estimate
- Status of each file (*loaded, empty, missing*)
- Recent activity (*modification times*)
- Drift warnings if any

**Example**:

```bash
ctx status
ctx status --json
ctx status --verbose
```

---

### `ctx agent`

Print an AI-ready context packet optimized for LLM consumption.

```bash
ctx agent [flags]
```

**Flags**:

| Flag                | Description                  |
|---------------------|------------------------------|
| `--budget <tokens>` | Token budget (default: 8000) |
| `--format md\|json` | Output format (default: md)  |

**Output**:

- Read order for context files
- Constitution rules (never truncated)
- Current tasks
- Key conventions
- Recent decisions

**Example**:

```bash
# Default (8000 tokens, markdown)
ctx agent

# Custom budget
ctx agent --budget 4000

# JSON format
ctx agent --format json
```

**Use case**: Copy-paste into AI chat, pipe to system prompt, or use in hooks.

---

### `ctx load`

Load and display assembled context as AI would see it.

```bash
ctx load [flags]
```

**Flags**:

| Flag                | Description                               |
|---------------------|-------------------------------------------|
| `--budget <tokens>` | Token budget for assembly (default: 8000) |
| `--raw`             | Output raw file contents without assembly |

**Example**:

```bash
ctx load
ctx load --budget 16000
ctx load --raw
```

---

### `ctx add`

Add a new item to a context file.

```bash
ctx add <type> <content> [flags]
```

**Types**:

| Type         | Target File    |
|--------------|----------------|
| `task`       | TASKS.md       |
| `decision`   | DECISIONS.md   |
| `learning`   | LEARNINGS.md   |
| `convention` | CONVENTIONS.md |

**Flags**:

| Flag                      | Short | Description                                                 |
|---------------------------|-------|-------------------------------------------------------------|
| `--priority <level>`      |       | Priority for tasks: `high`, `medium`, `low`                 |
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

# Record a decision (requires all ADR—Architectural Decision Record—fields)
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

- `task-id-or-text`: Task number or partial text match

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

- Path references in ARCHITECTURE.md and CONVENTIONS.md exist
- Task references are valid
- Constitution rules aren't violated (*heuristic*)
- Staleness indicators (*old files, many completed tasks*)

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
* Deduplicates the "*learning*"s with similar content
* Removes empty sections

```bash
ctx compact [flags]
```

**Flags**:

| Flag             | Description                                |
|------------------|--------------------------------------------|
| `--archive`      | Create `.context/archive/` for old content |
| `--no-auto-save` | Skip auto-saving session before compact    |

**Example**:

```bash
ctx compact
ctx compact --archive
ctx compact --no-auto-save
```

---

### `ctx completion`

Generate shell autocompletion scripts.

```bash
ctx completion <shell>
```

#### Subcommands

| Shell        | Command                  |
|--------------|--------------------------|
| `bash`       | `ctx completion bash`    |
| `zsh`        | `ctx completion zsh`     |
| `fish`       | `ctx completion fish`    |
| `powershell` | `ctx completion powershell` |

#### Installation

=== "Bash"

    ```bash
    # Add to ~/.bashrc
    source <(ctx completion bash)
    ```

=== "Zsh"

    ```bash
    # Add to ~/.zshrc
    source <(ctx completion zsh)
    ```

=== "Fish"

    ```bash
    ctx completion fish | source
    # Or save to completions directory
    ctx completion fish > ~/.config/fish/completions/ctx.fish
    ```

---

### `ctx tasks`

Manage task archival and snapshots.

```bash
ctx tasks <subcommand>
```

#### `ctx tasks archive`

Move completed tasks from TASKS.md to a timestamped archive file.

```bash
ctx tasks archive [flags]
```

**Flags**:

| Flag        | Description                              |
|-------------|------------------------------------------|
| `--dry-run` | Preview changes without modifying files  |

Archive files are stored in `.context/archive/` with timestamped names
(`tasks-YYYY-MM-DD.md`). Completed tasks (marked with `[x]`) are moved;
pending tasks (`[ ]`) remain in TASKS.md.

**Example**:

```bash
ctx tasks archive
ctx tasks archive --dry-run
```

#### `ctx tasks snapshot`

Create a point-in-time snapshot of TASKS.md without modifying the original.

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

### `ctx decisions`

Manage the DECISIONS.md file.

```bash
ctx decisions <subcommand>
```

#### `ctx decisions reindex`

Regenerate the quick-reference index at the top of DECISIONS.md.

```bash
ctx decisions reindex
```

The index is a compact table showing date and title for each decision,
allowing AI tools to quickly scan entries without reading the full file.

Use this after manual edits to DECISIONS.md or when migrating existing
files to use the index format.

**Example**:

```bash
ctx decisions reindex
# ✓ Index regenerated with 12 entries
```

---

### `ctx learnings`

Manage the LEARNINGS.md file.

```bash
ctx learnings <subcommand>
```

#### `ctx learnings reindex`

Regenerate the quick-reference index at the top of LEARNINGS.md.

```bash
ctx learnings reindex
```

The index is a compact table showing date and title for each learning,
allowing AI tools to quickly scan entries without reading the full file.

Use this after manual edits to LEARNINGS.md or when migrating existing
files to use the index format.

**Example**:

```bash
ctx learnings reindex
# ✓ Index regenerated with 8 entries
```

---

### `ctx recall`

Browse and search AI session history from Claude Code and other tools.

```bash
ctx recall <subcommand>
```

#### `ctx recall list`

List all parsed sessions.

```bash
ctx recall list [flags]
```

**Flags**:

| Flag             | Short | Description                               |
|------------------|-------|-------------------------------------------|
| `--limit`        | `-n`  | Maximum sessions to display (default: 20) |
| `--project`      | `-p`  | Filter by project name                    |
| `--tool`         | `-t`  | Filter by tool (e.g., `claude-code`)      |
| `--all-projects` |       | Include sessions from all projects        |

Sessions are sorted by date (newest first) and display slug, project,
start time, duration, turn count, and token usage.

**Example**:

```bash
ctx recall list
ctx recall list --limit 5
ctx recall list --project ctx
ctx recall list --tool claude-code
```

#### `ctx recall show`

Show details of a specific session.

```bash
ctx recall show [session-id] [flags]
```

**Flags**:

| Flag             | Description                        |
|------------------|------------------------------------|
| `--latest`       | Show the most recent session       |
| `--full`         | Show full message content          |
| `--all-projects` | Search across all projects         |

The session ID can be a full UUID, partial match, or session slug name.

**Example**:

```bash
ctx recall show abc123
ctx recall show gleaming-wobbling-sutherland
ctx recall show --latest
ctx recall show --latest --full
```

#### `ctx recall export`

Export sessions to editable journal files in `.context/journal/`.

```bash
ctx recall export [session-id] [flags]
```

**Flags**:

| Flag             | Description                              |
|------------------|------------------------------------------|
| `--all`          | Export all sessions                      |
| `--force`        | Overwrite existing files                 |
| `--all-projects` | Export from all projects                 |

Exported files include session metadata, tool usage summary, and the full
conversation. Existing files are skipped by default to preserve your edits.

The `journal/` directory should be gitignored (like `sessions/`) since it
contains raw conversation data.

**Example**:

```bash
ctx recall export abc123          # Export one session
ctx recall export --all           # Export all sessions
ctx recall export --all --force   # Overwrite existing exports
```

---

### `ctx journal`

Analyze and synthesize exported session files.

```bash
ctx journal <subcommand>
```

#### `ctx journal site`

Generate a static site from journal entries in `.context/journal/`.

```bash
ctx journal site [flags]
```

**Flags**:

| Flag       | Short | Description                                       |
|------------|-------|---------------------------------------------------|
| `--output` | `-o`  | Output directory (default: .context/journal-site) |
| `--build`  |       | Run zensical build after generating               |
| `--serve`  |       | Run zensical serve after generating               |

Creates a `zensical`-compatible site structure with an index page listing
all sessions by date, and individual pages for each journal entry.

Requires `zensical` to be installed for `--build` or `--serve`:

```bash
pip install zensical
```

**Example**:

```bash
ctx journal site                    # Generate in .context/journal-site/
ctx journal site --output ~/public  # Custom output directory
ctx journal site --build            # Generate and build HTML
ctx journal site --serve            # Generate and serve locally
```

---

### `ctx serve`

Serve a static site locally via `zensical`.

```bash
ctx serve [directory]
```

If no directory is specified, serves the journal site (`.context/journal-site`).

Requires `zensical` to be installed:

```bash
pip install zensical
```

**Example**:

```bash
ctx serve                           # Serve journal site
ctx serve .context/journal-site     # Serve specific directory
ctx serve ./docs                    # Serve docs folder
```

---

### `ctx watch`

Watch for AI output and auto-apply context updates.

Parses `<context-update>` XML commands from AI output and applies
them to context files.

```bash
ctx watch [flags]
```

**Flags**:

| Flag           | Description                         |
|----------------|-------------------------------------|
| `--log <file>` | Log file to watch (default: stdin)  |
| `--dry-run`    | Preview updates without applying    |
| `--auto-save`  | Periodically save session snapshots |

**Example**:

```bash
# Watch stdin
ai-tool | ctx watch

# Watch a log file
ctx watch --log /path/to/ai-output.log

# Preview without applying
ctx watch --dry-run
```

---

### `ctx hook`

Generate AI tool integration configuration.

```bash
ctx hook <tool>
```

**Supported tools**:

| Tool          | Description     |
|---------------|-----------------|
| `claude-code` | Claude Code CLI |
| `cursor`      | Cursor IDE      |
| `aider`       | Aider CLI       |
| `copilot`     | GitHub Copilot  |
| `windsurf`    | Windsurf IDE    |

**Example**:

```bash
ctx hook claude-code
ctx hook cursor
ctx hook aider
```

---

### `ctx session`

Manage session snapshots.

#### ctx session save

Save the current context snapshot.

```bash
ctx session save [topic] [flags]
```

**Flags**:

| Flag            | Short | Description                                              |
|-----------------|-------|----------------------------------------------------------|
| `--type <type>` | `-t`  | Session type: `feature`, `bugfix`, `refactor`, `session` |

**Example**:

```bash
ctx session save
ctx session save "feature-auth"
ctx session save "bugfix" --type bugfix
```

#### `ctx session list`

List saved sessions.

```bash
ctx session list [flags]
```

**Flags**:

| Flag      | Short | Description                                  |
|-----------|-------|----------------------------------------------|
| `--limit` | `-n`  | Maximum sessions to display (default: 10)    |

**Output**: Table of sessions with index, date, topic, and type.

**Example**:

```bash
ctx session list
ctx session list --limit 5
```

#### `ctx session load`

Load and display a previous session.

```bash
ctx session load <index|date|topic>
```

**Arguments**:

* `index`: Numeric index from `session list`
* `date`: Date pattern (e.g., `2026-01-21`)
* `topic`: Topic keyword match

**Example:**

```bash
ctx session load 1           # by index
ctx session load 2026-01-21  # by date
ctx session load auth        # by topic
```

#### `ctx session parse`

Parse JSONL transcript to readable markdown.

```bash
ctx session parse <file> [flags]
```

**Flags:**

| Flag         | Short | Description                                     |
|--------------|-------|-------------------------------------------------|
| `--output`   | `-o`  | Output file (default: stdout)                   |
| `--extract`  |       | Extract decisions and learnings from transcript |

**Example**:

```bash
ctx session parse ~/.claude/projects/.../transcript.jsonl
ctx session parse transcript.jsonl --extract
ctx session parse transcript.jsonl -o conversation.md
```

---

### `ctx loop`

Generate a shell script for running an autonomous loop.

An autonomous loop continuously runs an AI assistant with the same prompt until
a completion signal is detected, enabling iterative development where the
AI builds on its previous work.

```bash
ctx loop [flags]
```

**Flags**:

| Flag                     | Short | Description                              | Default            |
|--------------------------|-------|------------------------------------------|--------------------|
| `--tool <tool>`          | `-t`  | AI tool: `claude`, `aider`, or `generic` | `claude`           |
| `--prompt <file>`        | `-p`  | Prompt file to use                       | `PROMPT.md`        |
| `--max-iterations <n>`   | `-n`  | Maximum iterations (0 = unlimited)       | `0`                |
| `--completion <signal>`  | `-c`  | Completion signal to detect              | `SYSTEM_CONVERGED` |
| `--output <file>`        | `-o`  | Output script filename                   | `loop.sh`          |

**Example**:

```bash
# Generate loop.sh for Claude Code
ctx loop

# Generate for Aider with custom prompt
ctx loop --tool aider --prompt TASKS.md

# Limit to 10 iterations
ctx loop --max-iterations 10

# Output to custom file
ctx loop -o my-loop.sh
```

**Usage**:

```bash
# Generate and run the loop
ctx loop
chmod +x loop.sh
./loop.sh
```

See [Autonomous Loops](autonomous-loop.md) for detailed workflow documentation.

---

## Exit Codes

| Code | Meaning              |
|------|----------------------|
| 0    | Success              |
| 1    | General error        |
| 2    | Context not found    |
| 3    | Invalid arguments    |
| 4    | File operation error |

## Environment Variables

| Variable           | Description                             |
|--------------------|-----------------------------------------|
| `CTX_DIR`          | Override default context directory path |
| `CTX_TOKEN_BUDGET` | Override default token budget           |
| `NO_COLOR`         | Disable colored output when set         |

## Configuration File

Optional `.contextrc` (YAML format) at project root:

```yaml
# .contextrc
context_dir: .context # Context directory name
token_budget: 8000    # Default token budget
priority_order:       # File loading priority
  - TASKS.md
  - DECISIONS.md
  - CONVENTIONS.md
auto_archive: true    # Auto-archive old items
archive_after_days: 7 # Days before archiving
```

**Priority order:** CLI flags > Environment variables > `.contextrc` > Defaults

All settings are optional. Missing values use defaults.
