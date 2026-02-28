---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Tools and Utilities
icon: lucide/wrench
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
    Claude Code integration is now provided via the `ctx` plugin.
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
using the encryption key and stored in `.context/.notify.enc`.

```bash
ctx notify setup
```

The encrypted file is safe to commit. The key (`.context/.context.key`) is
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
ctx pad merge --key /other/.context.key foreign.enc
ctx pad merge --dry-run pad-a.enc pad-b.md
```

---

### `ctx remind`

Session-scoped reminders that surface at session start. Reminders are
stored verbatim and relayed verbatim — no summarization, no categories.

When invoked with a text argument and no subcommand, adds a reminder.

```bash
ctx remind "text"
ctx remind <subcommand>
```

#### `ctx remind add`

Add a reminder. This is the default action — `ctx remind "text"` and
`ctx remind add "text"` are equivalent.

```bash
ctx remind "refactor the swagger definitions"
ctx remind add "check CI after the deploy" --after 2026-02-25
```

**Arguments**:

- `text`: The reminder message (verbatim)

**Flags**:

| Flag      | Short | Description                                |
|-----------|-------|--------------------------------------------|
| `--after` | `-a`  | Don't surface until this date (YYYY-MM-DD) |

**Examples**:

```bash
ctx remind "refactor the swagger definitions"
ctx remind "check CI after the deploy" --after 2026-02-25
```

#### `ctx remind list`

List all pending reminders. Date-gated reminders that aren't yet due
are annotated with `(after DATE, not yet due)`.

```bash
ctx remind list
```

**Aliases**: `ls`

#### `ctx remind dismiss`

Remove a reminder by ID, or remove all reminders with `--all`.

```bash
ctx remind dismiss <id>
ctx remind dismiss --all
```

**Arguments**:

- `id`: Reminder ID (shown in `list` output)

**Flags**:

| Flag    | Description              |
|---------|--------------------------|
| `--all` | Dismiss all reminders    |

**Aliases**: `rm`

**Examples**:

```bash
ctx remind dismiss 3
ctx remind dismiss --all
```

---

### `ctx pause`

Pause all context nudge and reminder hooks for the current session.
Security hooks (dangerous command blocking) and housekeeping hooks still fire.

```bash
ctx pause [flags]
```

**Flags**:

| Flag             | Description                         |
|------------------|-------------------------------------|
| `--session-id`   | Session ID (overrides stdin)        |

**Example**:

```bash
# Pause hooks for a quick investigation
ctx pause

# Resume when ready
ctx resume
```

**See also**: [Pausing Context Hooks](../recipes/session-pause.md)

---

### `ctx resume`

Resume context hooks after a pause. Silent no-op if not paused.

```bash
ctx resume [flags]
```

**Flags**:

| Flag             | Description                         |
|------------------|-------------------------------------|
| `--session-id`   | Session ID (overrides stdin)        |

**Example**:

```bash
ctx resume
```

**See also**: [Pausing Context Hooks](../recipes/session-pause.md)

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

### `ctx site`

Site management commands for the ctx.ist static site.

```bash
ctx site <subcommand>
```

#### `ctx site feed`

Generate an Atom 1.0 feed from finalized blog posts in `docs/blog/`.

```bash
ctx site feed [flags]
```

Scans `docs/blog/` for files matching `YYYY-MM-DD-*.md`, parses YAML
frontmatter, and generates a valid Atom feed. Only posts with
`reviewed_and_finalized: true` are included. Summaries are extracted
from the first paragraph after the heading.

**Flags**:

| Flag         | Short | Type   | Default          | Description              |
|--------------|-------|--------|------------------|--------------------------|
| `--out`      | `-o`  | string | `site/feed.xml`  | Output path              |
| `--base-url` |       | string | `https://ctx.ist`| Base URL for entry links |

**Output**:

```
Generated site/feed.xml (21 entries)

Skipped:
  2026-02-25-the-homework-problem.md — not finalized

Warnings:
  2026-02-09-defense-in-depth.md — no summary paragraph found
```

Three buckets: **included** (count), **skipped** (with reason),
**warnings** (included but degraded). Exit 0 always — warnings
inform but do not block.

**Frontmatter requirements**:

| Field                    | Required | Feed mapping         |
|--------------------------|----------|----------------------|
| `title`                  | Yes      | `<title>`            |
| `date`                   | Yes      | `<updated>`          |
| `reviewed_and_finalized` | Yes      | Draft gate (must be `true`) |
| `author`                 | No       | `<author><name>`     |
| `topics`                 | No       | `<category term="">` |

**Examples**:

```bash
ctx site feed                                # Generate site/feed.xml
ctx site feed --out /tmp/feed.xml            # Custom output path
ctx site feed --base-url https://example.com # Custom base URL
make site-feed                               # Makefile shortcut
make site                                    # Builds site + feed
```

---

### `ctx why`

Read ctx's philosophy documents directly in the terminal, stripped of
MkDocs artifacts so they display cleanly without a browser.

```bash
ctx why [DOCUMENT]
```

**Documents**:

| Name         | Description                                      |
|--------------|--------------------------------------------------|
| `manifesto`  | The ctx Manifesto — creation, not code           |
| `about`      | About ctx — what it is and why it exists         |
| `invariants` | Design invariants — properties that must hold    |

**Usage**:

```bash
# Interactive numbered menu
ctx why

# Show a specific document
ctx why manifesto
ctx why about
ctx why invariants

# Pipe to a pager
ctx why manifesto | less
```
