---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Journal
icon: lucide/history
---

### `ctx journal`

Browse and search AI session history from Claude Code and other tools.

```bash
ctx journal <subcommand>
```

#### `ctx journal source`

List all parsed sessions.

```bash
ctx journal source [flags]
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
ctx journal source
ctx journal source --limit 5
ctx journal source --project ctx
ctx journal source --tool claude-code
```

#### `ctx journal source --show`

Show details of a specific session.

```bash
ctx journal source --show [session-id] [flags]
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
ctx journal source --show abc123
ctx journal source --show gleaming-wobbling-sutherland
ctx journal source --show --latest
ctx journal source --show --latest --full
```

#### `ctx journal import`

Import sessions to editable journal files in `.context/journal/`.

```bash
ctx journal import [session-id] [flags]
```

**Flags**:

| Flag                 | Description                                                             |
|----------------------|-------------------------------------------------------------------------|
| `--all`              | Import all sessions (only new files by default)                         |
| `--all-projects`     | Import from all projects                                                |
| `--regenerate`       | Re-import existing files (preserves YAML frontmatter by default)        |
| `--keep-frontmatter` | Preserve enriched YAML frontmatter during regeneration (default: true)  |
| `--yes`, `-y`        | Skip confirmation prompt                                                |
| `--dry-run`          | Show what would be imported without writing files                       |

**Safe by default**: `--all` only imports new sessions. Existing files are
skipped. Use `--regenerate` to re-import existing files (conversation content
is regenerated, YAML frontmatter from enrichment is preserved by default).
Use `--keep-frontmatter=false` to discard enriched frontmatter during
regeneration.

Locked entries (via `ctx journal lock`) are always skipped, regardless of flags.

Single-session import (`ctx journal import <id>`) always writes without
prompting, since you are explicitly targeting one session.

The `journal/` directory should be gitignored (like `sessions/`) since it
contains raw conversation data.

**Example**:

```bash
ctx journal import abc123                 # Import one session
ctx journal import --all                  # Import only new sessions
ctx journal import --all --dry-run        # Preview what would be imported
ctx journal import --all --regenerate     # Re-import existing (prompts)
ctx journal import --all --regenerate -y  # Re-import without prompting
ctx journal import --all --regenerate --keep-frontmatter=false -y  # Discard frontmatter
```

#### `ctx journal lock`

Protect journal entries from being overwritten by `import --regenerate` or
modified by enrichment skills (`/ctx-journal-enrich`, `/ctx-journal-enrich-all`).

```bash
ctx journal lock <pattern> [flags]
```

**Flags**:

| Flag    | Description              |
|---------|--------------------------|
| `--all` | Lock all journal entries |

The pattern matches filenames by slug, date, or short ID. Locking a
multi-part entry locks all parts. The lock is recorded in
`.context/journal/.state.json` and a `locked: true` line is added to the
file's YAML frontmatter for visibility.

**Example**:

```bash
ctx journal lock abc12345
ctx journal lock 2026-01-21-session-abc12345.md
ctx journal lock --all
```

#### `ctx journal unlock`

Remove lock protection from journal entries.

```bash
ctx journal unlock <pattern> [flags]
```

**Flags**:

| Flag    | Description                |
|---------|----------------------------|
| `--all` | Unlock all journal entries |

**Example**:

```bash
ctx journal unlock abc12345
ctx journal unlock --all
```

#### `ctx journal sync`

Sync lock state from journal frontmatter to `.state.json`.

```bash
ctx journal sync
```

Scans all journal markdowns and updates `.state.json` to match each file's
frontmatter. Files with `locked: true` in frontmatter are marked locked in
state; files without a `locked:` line have their lock cleared.

This is the inverse of `ctx journal lock`: instead of state driving
frontmatter, frontmatter drives state. Useful after batch enrichment where
you add `locked: true` to frontmatter manually.

**Example**:

```bash
# After enriching entries and adding locked: true to frontmatter
ctx journal sync
```

---

### `ctx journal`

Analyze and synthesize imported session files.

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
ctx journal obsidian                        # Generate in .context/journal-obsidian/
ctx journal obsidian --output ~/vaults/ctx  # Custom output directory
```

#### `ctx journal schema check`

Validate JSONL session files against the embedded schema and report drift.

```bash
ctx journal schema check [flags]
```

**Flags**:

| Flag              | Short | Description                                  |
|-------------------|-------|----------------------------------------------|
| `--dir`           |       | Directory to scan for JSONL files             |
| `--all-projects`  |       | Scan all Claude Code project directories      |
| `--quiet`         | `-q`  | Exit code only (0 = clean, 1 = drift)         |

Scans JSONL files for unknown fields, missing required fields, unknown record
types, and unknown content block types. When drift is found, writes a Markdown
report to `.context/reports/schema-drift.md`. When drift resolves, the report
is automatically deleted.

Designed for interactive use, CI pipelines, and nightly cron jobs.

**Example**:

```bash
ctx journal schema check                    # Current project
ctx journal schema check --all-projects     # All projects
ctx journal schema check --quiet            # Exit code only
ctx journal schema check --dir /path/to     # Custom directory
```

#### `ctx journal schema dump`

Print the embedded JSONL schema definition.

```bash
ctx journal schema dump
```

Shows all known record types with their required and optional fields, and all
recognized content block types with their parse status. Useful for inspecting
what the schema validator expects.

**Example**:

```bash
ctx journal schema dump
```

---

### `ctx serve`

Serve any zensical directory locally. This is a **serve-only** command: It
does not generate or regenerate site content.

```bash
ctx serve [directory]
```

If no directory is specified, defaults to the journal site (`.context/journal-site`).

Requires `zensical` to be installed:

```bash
pipx install zensical
```

!!! tip "`ctx serve` vs. `ctx journal site --serve`"
    `ctx journal site --serve` **generates** the journal site *then* serves
    it: an all-in-one command. `ctx serve` only **serves** an existing
    directory, and works with any zensical site (journal, docs, etc.).

**Example**:

```bash
ctx serve                        # Serve journal site (no regeneration)
ctx serve .context/journal-site  # Same, explicit path
ctx serve ./site                 # Serve the docs site
```
