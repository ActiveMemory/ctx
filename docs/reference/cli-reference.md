---
#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: CLI
icon: lucide/terminal
---

![ctx](../images/ctx-banner.png)

## `ctx` CLI

This is a complete reference for all `ctx` commands.

## Global Options

All commands support these flags:

| Flag                   | Description                                              |
|------------------------|----------------------------------------------------------|
| `--help`               | Show command help                                        |
| `--version`            | Show version                                             |
| `--context-dir <path>` | Override context directory (default: `.context/`)        |
| `--no-color`           | Disable colored output                                   |
| `--allow-outside-cwd`  | Allow context directory outside current working directory |

> The `NO_COLOR=1` environment variable also disables colored output.

## Commands

| Command                           | Description                                               |
|-----------------------------------|-----------------------------------------------------------|
| [`ctx init`](#ctx-init)           | Initialize `.context/` directory with templates           |
| [`ctx status`](#ctx-status)       | Show context summary (files, tokens, drift)               |
| [`ctx agent`](#ctx-agent)         | Print token-budgeted context packet for AI consumption    |
| [`ctx load`](#ctx-load)           | Output assembled context in read order                    |
| [`ctx add`](#ctx-add)             | Add a task, decision, learning, or convention             |
| [`ctx complete`](#ctx-complete)   | Mark a task as done                                       |
| [`ctx drift`](#ctx-drift)         | Detect stale paths, secrets, missing files                |
| [`ctx sync`](#ctx-sync)           | Reconcile context with codebase state                     |
| [`ctx compact`](#ctx-compact)     | Archive completed tasks, clean up files                   |
| [`ctx tasks`](#ctx-tasks)         | Task archival and snapshots                               |
| [`ctx permissions`](#ctx-permissions) | Permission snapshots (golden image)                   |
| [`ctx decisions`](#ctx-decisions) | Manage `DECISIONS.md` (reindex)                           |
| [`ctx learnings`](#ctx-learnings) | Manage `LEARNINGS.md` (reindex)                           |
| [`ctx recall`](#ctx-recall)       | Browse and export AI session history                      |
| [`ctx journal`](#ctx-journal)     | Generate static site from journal entries                 |
| [`ctx serve`](#ctx-serve)         | Serve static site locally                                 |
| [`ctx watch`](#ctx-watch)         | Auto-apply context updates from AI output                 |
| [`ctx hook`](#ctx-hook)           | Generate AI tool integration configs                      |
| [`ctx loop`](#ctx-loop)           | Generate autonomous loop script                           |
| [`ctx notify`](#ctx-notify)       | Send webhook notifications                                |
| [`ctx pad`](#ctx-pad)             | Encrypted scratchpad for sensitive one-liners             |
| [`ctx system`](#ctx-system)       | System diagnostics and hook commands                     |

---

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
| `--merge`   |       | Auto-merge ctx content into existing CLAUDE.md and PROMPT.md          |
| `--ralph`   |       | Agent works autonomously without asking questions                     |

**Creates**:

- `.context/` directory with all template files
- `.claude/settings.local.json` with pre-approved ctx permissions
- `PROMPT.md` with session prompt (autonomous mode with `--ralph`)
- `IMPLEMENTATION_PLAN.md` with high-level project direction
- `CLAUDE.md` with bootstrap instructions (or merges into existing)

Claude Code hooks and skills are provided by the **ctx plugin**
(see [Integrations](../operations/integrations.md#claude-code-full-integration)).

**Example**:

```bash
# Collaborative mode (agent asks questions when unclear)
ctx init

# Autonomous mode (agent works independently)
ctx init --ralph

# Minimal setup (just core files)
ctx init --minimal

# Force overwrite existing
ctx init --force

# Merge into existing files
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

| Flag         | Default | Description                                                     |
|--------------|---------|-----------------------------------------------------------------|
| `--budget`   | 8000    | Token budget — controls content selection and prioritization    |
| `--format`   | md      | Output format: `md` or `json`                                   |
| `--cooldown` | 10m     | Suppress repeated output within this duration (requires `--session`) |
| `--session`  | (none)  | Session ID for cooldown isolation (e.g., `$PPID`)               |

**How budget works**:

The budget controls how much context is included. Entries are selected
in priority tiers:

1. **Constitution** — always included in full (inviolable rules)
2. **Tasks** — all active tasks, up to 40% of budget
3. **Conventions** — all conventions, up to 20% of budget
4. **Decisions** — scored by recency and relevance to active tasks
5. **Learnings** — scored by recency and relevance to active tasks

Decisions and learnings are ranked by a combined score (how recent + how
relevant to your current tasks). High-scoring entries are included with
their full body. Entries that don't fit get title-only summaries in an
"Also Noted" section. Superseded entries are excluded.

**Output sections**:

| Section              | Source           | Selection                          |
|----------------------|------------------|------------------------------------|
| Read These Files     | all `.context/`  | Non-empty files in priority order  |
| Constitution         | CONSTITUTION.md  | All rules (never truncated)        |
| Current Tasks        | TASKS.md         | All unchecked tasks (budget-capped)|
| Key Conventions      | CONVENTIONS.md   | All items (budget-capped)          |
| Recent Decisions     | DECISIONS.md     | Full body, scored by relevance     |
| Key Learnings        | LEARNINGS.md     | Full body, scored by relevance     |
| Also Noted           | overflow         | Title-only summaries               |

**Example**:

```bash
# Default (8000 tokens, markdown)
ctx agent

# Smaller packet for tight context windows
ctx agent --budget 4000

# JSON format for programmatic use
ctx agent --format json

# Pipe to file
ctx agent --budget 4000 > context.md

# With cooldown (hooks/automation — requires --session)
ctx agent --session $PPID
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
- Entry count — warns when LEARNINGS.md or DECISIONS.md exceed configurable
  thresholds (default: 30 learnings, 20 decisions), or when CONVENTIONS.md
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

### `ctx completion`

Generate shell autocompletion scripts.

```bash
ctx completion <shell>
```

#### Subcommands

| Shell        | Command                     |
|--------------|-----------------------------|
| `bash`       | `ctx completion bash`       |
| `zsh`        | `ctx completion zsh`        |
| `fish`       | `ctx completion fish`       |
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

Manage the DECISIONS.md file.

```bash
ctx decisions <subcommand>
```

#### `ctx decisions reindex`

Regenerate the quick-reference index at the top of DECISIONS.md.

```bash
ctx decisions reindex
```

The index is a compact table showing the date and title for each decision,
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

The index is a compact table showing the date and title for each learning,
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

| Flag             | Description                                               |
|------------------|-----------------------------------------------------------|
| `--all`          | Export all sessions (only new files by default)            |
| `--all-projects` | Export from all projects                                   |
| `--regenerate`        | Re-export existing files (preserves YAML frontmatter by default) |
| `--keep-frontmatter` | Preserve enriched YAML frontmatter during regeneration (default: true) |
| `--yes`, `-y`         | Skip confirmation prompt                                   |
| `--dry-run`           | Show what would be exported without writing files           |

**Safe by default**: `--all` only exports new sessions. Existing files are
skipped. Use `--regenerate` to re-export existing files (conversation content
is regenerated, YAML frontmatter from enrichment is preserved by default).
Use `--keep-frontmatter=false` to discard enriched frontmatter during
regeneration.

Locked entries (via `ctx recall lock`) are always skipped, regardless of flags.

Single-session export (`ctx recall export <id>`) always writes without
prompting, since you are explicitly targeting one session.

The `journal/` directory should be gitignored (like `sessions/`) since it
contains raw conversation data.

**Example**:

```bash
ctx recall export abc123                              # Export one session
ctx recall export --all                               # Export only new sessions
ctx recall export --all --dry-run                     # Preview what would be exported
ctx recall export --all --regenerate                  # Re-export existing (prompts)
ctx recall export --all --regenerate -y               # Re-export without prompting
ctx recall export --all --regenerate --keep-frontmatter=false -y  # Discard frontmatter
```

#### `ctx recall lock`

Protect journal entries from being overwritten by `export --regenerate`.

```bash
ctx recall lock <pattern> [flags]
```

**Flags**:

| Flag    | Description                      |
|---------|----------------------------------|
| `--all` | Lock all journal entries          |

The pattern matches filenames by slug, date, or short ID. Locking a
multi-part entry locks all parts. The lock is recorded in
`.context/journal/.state.json` and a `locked: true` line is added to the
file's YAML frontmatter for visibility.

**Example**:

```bash
ctx recall lock abc12345
ctx recall lock 2026-01-21-session-abc12345.md
ctx recall lock --all
```

#### `ctx recall unlock`

Remove lock protection from journal entries.

```bash
ctx recall unlock <pattern> [flags]
```

**Flags**:

| Flag    | Description                        |
|---------|------------------------------------|
| `--all` | Unlock all journal entries          |

**Example**:

```bash
ctx recall unlock abc12345
ctx recall unlock --all
```

#### `ctx recall sync`

Sync lock state from journal frontmatter to `.state.json`.

```bash
ctx recall sync
```

Scans all journal markdowns and updates `.state.json` to match each file's
frontmatter. Files with `locked: true` in frontmatter are marked locked in
state; files without a `locked:` line have their lock cleared.

This is the inverse of `ctx recall lock`: instead of state driving
frontmatter, frontmatter drives state. Useful after batch enrichment where
you add `locked: true` to frontmatter manually.

**Example**:

```bash
# After enriching entries and adding locked: true to frontmatter
ctx recall sync
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
pipx install zensical
```

**Example**:

```bash
ctx journal site                    # Generate in .context/journal-site/
ctx journal site --output ~/public  # Custom output directory
ctx journal site --build            # Generate and build HTML
ctx journal site --serve            # Generate and serve locally
```

#### `ctx journal obsidian`

Generate an Obsidian vault from journal entries in `.context/journal/`.

```bash
ctx journal obsidian [flags]
```

**Flags**:

| Flag       | Short | Description                                             |
|------------|-------|---------------------------------------------------------|
| `--output` | `-o`  | Output directory (default: .context/journal-obsidian)   |

Creates an Obsidian-compatible vault with:

- **Wikilinks** (`[[target|display]]`) for all internal navigation
- **MOC pages** (Map of Content) for topics, key files, and session types
- **Related sessions footer** linking entries that share topics
- **Transformed frontmatter** (`topics` → `tags` for Obsidian integration)
- **Minimal `.obsidian/`** config enforcing wikilink mode

No external dependencies are required:
Open the output directory as an Obsidian  vault directly.

**Example**:

```bash
ctx journal obsidian                          # Generate in .context/journal-obsidian/
ctx journal obsidian --output ~/vaults/ctx    # Custom output directory
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
pipx install zensical
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
ctx hook <tool> [flags]
```

**Flags**:

| Flag      | Short | Description                                                            |
|-----------|-------|------------------------------------------------------------------------|
| `--write` | `-w`  | Write the generated config to disk (e.g. `.github/copilot-instructions.md`) |

**Supported tools**:

| Tool          | Description                                  |
|---------------|----------------------------------------------|
| `claude-code` | Redirects to plugin install instructions     |
| `cursor`      | Cursor IDE                                   |
| `aider`       | Aider CLI                                    |
| `copilot`     | GitHub Copilot                               |
| `windsurf`    | Windsurf IDE                                 |

!!! note "Claude Code uses the plugin system"
    Claude Code integration is now provided via the ctx plugin.
    Running `ctx hook claude-code` prints plugin install instructions.

**Example**:

```bash
# Print hook instructions to stdout
ctx hook cursor
ctx hook aider

# Generate and write .github/copilot-instructions.md
ctx hook copilot --write
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

See [Autonomous Loops](../operations/autonomous-loop.md) for detailed workflow documentation.

---

### `ctx notify`

Send fire-and-forget webhook notifications from skills, loops, and hooks.

```bash
ctx notify --event <name> [--session-id <id>] "message"
```

**Flags**:

| Flag           | Short | Description                    |
|----------------|-------|--------------------------------|
| `--event`      | `-e`  | Event name (required)          |
| `--session-id` | `-s`  | Session ID (optional)          |

**Behavior**:

- No webhook configured: silent noop (exit 0)
- Webhook set but event not in `events` list: silent noop (exit 0)
- Webhook set and event matches: fire-and-forget HTTP POST
- HTTP errors silently ignored (no retry)

**Example**:

```bash
ctx notify --event loop "Loop completed after 5 iterations"
ctx notify -e nudge -s session-abc "Context checkpoint at prompt #20"
```

#### `ctx notify setup`

Configure the webhook URL interactively. The URL is encrypted with AES-256-GCM
using the scratchpad key and stored in `.context/.notify.enc`.

```bash
ctx notify setup
```

The encrypted file is safe to commit. The key (`.context/.scratchpad.key`) is
gitignored and never committed.

#### `ctx notify test`

Send a test notification and report the HTTP response status.

```bash
ctx notify test
```

**Payload format** (JSON POST):

```json
{
  "event": "loop",
  "message": "Loop completed after 5 iterations",
  "session_id": "abc123-...",
  "timestamp": "2026-02-22T14:30:00Z",
  "project": "ctx"
}
```

| Field        | Type   | Description                           |
|--------------|--------|---------------------------------------|
| `event`      | string | Event name from `--event` flag        |
| `message`    | string | Notification message                  |
| `session_id` | string | Session ID (omitted if empty)         |
| `timestamp`  | string | UTC RFC3339 timestamp                 |
| `project`    | string | Project directory name                |

---

### `ctx pad`

Encrypted scratchpad for sensitive one-liners that travel with the project.

When invoked without a subcommand, lists all entries.

```bash
ctx pad
ctx pad <subcommand>
```

#### `ctx pad add`

Append a new entry to the scratchpad.

```bash
ctx pad add <text>
ctx pad add <label> --file <path>
```

**Flags**:

| Flag     | Short | Description                                |
|----------|-------|--------------------------------------------|
| `--file` | `-f`  | Ingest a file as a blob entry (max 64 KB)  |

**Examples**:

```bash
ctx pad add "DATABASE_URL=postgres://user:pass@host/db"
ctx pad add "deploy config" --file ./deploy.yaml
```

#### `ctx pad show`

Output the raw text of an entry by number. For blob entries, prints
decoded file content (or writes to disk with `--out`).

```bash
ctx pad show <n>
ctx pad show <n> --out <path>
```

**Arguments**:

- `n`: 1-based entry number

**Flags**:

| Flag    | Description                                      |
|---------|--------------------------------------------------|
| `--out` | Write decoded blob content to a file (blobs only)|

**Examples**:

```bash
ctx pad show 3
ctx pad show 2 --out ./recovered.yaml
```

#### `ctx pad rm`

Remove an entry by number.

```bash
ctx pad rm <n>
```

**Arguments**:

- `n`: 1-based entry number

#### `ctx pad edit`

Replace, append to, or prepend to an entry.

```bash
ctx pad edit <n> [text]
```

**Arguments**:

- `n`: 1-based entry number
- `text`: Replacement text (mutually exclusive with `--append`/`--prepend`)

**Flags**:

| Flag        | Description                                      |
|-------------|--------------------------------------------------|
| `--append`  | Append text to the end of the entry              |
| `--prepend` | Prepend text to the beginning of entry           |
| `--file`    | Replace blob file content (preserves label)      |
| `--label`   | Replace blob label (preserves content)           |

**Examples**:

```bash
ctx pad edit 2 "new text"
ctx pad edit 2 --append " suffix"
ctx pad edit 2 --prepend "prefix "
ctx pad edit 1 --file ./v2.yaml
ctx pad edit 1 --label "new name"
```

#### `ctx pad mv`

Move an entry from one position to another.

```bash
ctx pad mv <from> <to>
```

**Arguments**:

- `from`: Source position (1-based)
- `to`: Destination position (1-based)

#### `ctx pad resolve`

Show both sides of a merge conflict in the encrypted scratchpad.

```bash
ctx pad resolve
```

#### `ctx pad import`

Bulk-import lines from a file into the scratchpad. Each non-empty line
becomes a separate entry. All entries are written in a single encrypt/write
cycle.

```bash
ctx pad import <file>
ctx pad import -          # read from stdin
```

**Arguments**:

- `file`: Path to a text file, or `-` for stdin

**Examples**:

```bash
ctx pad import notes.txt
grep TODO *.go | ctx pad import -
```

#### `ctx pad export`

Export all blob entries from the scratchpad to a directory as files.
Each blob's label becomes the filename. Non-blob entries are skipped.

```bash
ctx pad export [dir]
```

**Arguments**:

- `dir`: Target directory (default: current directory)

**Flags**:

| Flag        | Short | Description                                         |
|-------------|-------|-----------------------------------------------------|
| `--force`   | `-f`  | Overwrite existing files instead of timestamping     |
| `--dry-run` |       | Print what would be exported without writing         |

When a file already exists, a unix timestamp is prepended to avoid
collisions (e.g., `1739836200-label`). Use `--force` to overwrite instead.

**Examples**:

```bash
ctx pad export ./ideas
ctx pad export --dry-run
ctx pad export --force ./backup
```

#### `ctx pad merge`

Merge entries from one or more scratchpad files into the current pad.
Each input file is auto-detected as encrypted or plaintext. Entries are
deduplicated by exact content.

```bash
ctx pad merge FILE...
```

**Arguments**:

- `FILE...`: One or more scratchpad files to merge (encrypted or plaintext)

**Flags**:

| Flag        | Short | Description                                         |
|-------------|-------|-----------------------------------------------------|
| `--key`     | `-k`  | Path to key file for decrypting input files          |
| `--dry-run` |       | Print what would be merged without writing           |

**Examples**:

```bash
ctx pad merge worktree/.context/scratchpad.enc
ctx pad merge notes.md backup.enc
ctx pad merge --key /other/.scratchpad.key foreign.enc
ctx pad merge --dry-run pad-a.enc pad-b.md
```

---

### `ctx system`

System diagnostics and hook commands.

```bash
ctx system <subcommand>
```

The parent command shows available subcommands. Hidden plumbing subcommands
(`ctx system mark-journal`) are used by skills and automation. Hidden hook
subcommands (`ctx system check-*`) are used by the Claude Code plugin — see
[AI Tool Integrations](../operations/integrations.md#plugin-hooks) for details.

#### `ctx system resources`

Show system resource usage with threshold-based alerts.

```bash
ctx system resources [flags]
```

Displays memory, swap, disk, and CPU load with two severity tiers:

| Resource | WARNING | DANGER |
|----------|---------|--------|
| Memory | >= 80% used | >= 90% used |
| Swap | >= 50% used | >= 75% used |
| Disk (cwd) | >= 85% full | >= 95% full |
| Load (1m) | >= 0.8x CPUs | >= 1.5x CPUs |

**Flags**:

| Flag     | Description           |
|----------|-----------------------|
| `--json` | Output in JSON format |

**Examples**:

```bash
ctx system resources               # Text output with status indicators
ctx system resources --json        # Machine-readable JSON
ctx system resources --json | jq '.alerts'   # Extract alerts only
```

**Text output**:

```
System Resources
====================

Memory:    4.2 / 16.0 GB (26%)                     ✓ ok
Swap:      0.0 /  8.0 GB (0%)                      ✓ ok
Disk:    180.2 / 500.0 GB (36%)                     ✓ ok
Load:     0.52 / 0.41 / 0.38  (8 CPUs, ratio 0.07) ✓ ok

All clear — no resource warnings.
```

When resources breach thresholds, alerts are listed below the summary:

```
Alerts:
  ✖ Memory 92% used (14.7 / 16.0 GB)
  ✖ Swap 78% used (6.2 / 8.0 GB)
  ✖ Load 1.56x CPU count
```

**Platform support**: Full metrics on Linux and macOS. Windows shows
disk only; memory and load report as unsupported.

#### `ctx system mark-journal`

Update processing state for a journal entry. Records the current date
in `.context/journal/.state.json`. Used by journal skills to record
pipeline progress.

```bash
ctx system mark-journal <filename> <stage>
```

**Stages**: `exported`, `enriched`, `normalized`, `fences_verified`

| Flag      | Description                           |
|-----------|---------------------------------------|
| `--check` | Check if stage is set (exit 1 if not) |

**Example**:

```bash
ctx system mark-journal 2026-01-21-session-abc12345.md enriched
ctx system mark-journal 2026-01-21-session-abc12345.md normalized
ctx system mark-journal --check 2026-01-21-session-abc12345.md fences_verified
```

---

## Exit Codes

| Code | Meaning                                |
|------|----------------------------------------|
| 0    | Success                                |
| 1    | General error / warnings (e.g. drift)  |
| 2    | Context not found                      |
| 3    | Violations found (e.g. drift)          |
| 4    | File operation error                   |

## Environment Variables

| Variable           | Description                             |
|--------------------|-----------------------------------------|
| `CTX_DIR`          | Override default context directory path |
| `CTX_TOKEN_BUDGET` | Override default token budget           |
| `NO_COLOR`         | Disable colored output when set         |

## Configuration File

Optional `.ctxrc` (*YAML format*) at project root:

```yaml
# .ctxrc
context_dir: .context                # Context directory name
token_budget: 8000                   # Default token budget
priority_order:                      # File loading priority
  - TASKS.md
  - DECISIONS.md
  - CONVENTIONS.md
auto_archive: true                   # Auto-archive old items
archive_after_days: 7                # Days before archiving tasks
scratchpad_encrypt: true             # Encrypt scratchpad (default: true)
allow_outside_cwd: false             # Skip boundary check (default: false)
entry_count_learnings: 30            # Drift warning threshold (0 = disable)
entry_count_decisions: 20            # Drift warning threshold (0 = disable)
convention_line_count: 200           # Line count warning for CONVENTIONS.md (0 = disable)
notify:                              # Webhook notification settings
  events:                            # Required — only listed events fire
    - loop
    - nudge
    - relay
  key_rotation_days: 90              # Days before key rotation nudge
```

| Field                           | Type       | Default      | Description                                          |
|---------------------------------|------------|--------------|------------------------------------------------------|
| `context_dir`                   | `string`   | `.context`   | Context directory name (relative to project root)    |
| `token_budget`                  | `int`      | `8000`       | Default token budget for `ctx agent`                 |
| `priority_order`                | `[]string` | *(all files)* | File loading priority for context packets           |
| `auto_archive`                  | `bool`     | `false`      | Auto-archive completed tasks                         |
| `archive_after_days`            | `int`      | `7`          | Days before completed tasks are archived             |
| `scratchpad_encrypt`            | `bool`     | `true`       | Encrypt scratchpad with AES-256-GCM                  |
| `allow_outside_cwd`             | `bool`     | `false`      | Skip boundary check for external context dirs        |
| `entry_count_learnings`         | `int`      | `30`         | Drift warning when LEARNINGS.md exceeds this count   |
| `entry_count_decisions`         | `int`      | `20`         | Drift warning when DECISIONS.md exceeds this count   |
| `convention_line_count`         | `int`      | `200`        | Line count warning for CONVENTIONS.md                |
| `notify.events`                 | `[]string` | *(all)*      | Event filter for webhook notifications (empty = all) |
| `notify.key_rotation_days`      | `int`      | `90`         | Days before encryption key rotation nudge            |

**Priority order:** CLI flags > Environment variables > `.ctxrc` > Defaults

All settings are optional. Missing values use defaults.
