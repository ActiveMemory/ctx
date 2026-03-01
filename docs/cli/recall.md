---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Recall and Journal
icon: lucide/history
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

| Flag                 | Description                                                            |
|----------------------|------------------------------------------------------------------------|
| `--all`              | Export all sessions (only new files by default)                        |
| `--all-projects`     | Export from all projects                                               |
| `--regenerate`       | Re-export existing files (preserves YAML frontmatter by default)       |
| `--keep-frontmatter` | Preserve enriched YAML frontmatter during regeneration (default: true) |
| `--yes`, `-y`        | Skip confirmation prompt                                               |
| `--dry-run`          | Show what would be exported without writing files                      |

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
ctx recall export abc123                 # Export one session
ctx recall export --all                  # Export only new sessions
ctx recall export --all --dry-run        # Preview what would be exported
ctx recall export --all --regenerate     # Re-export existing (prompts)
ctx recall export --all --regenerate -y  # Re-export without prompting
ctx recall export --all --regenerate --keep-frontmatter=false -y  # Discard frontmatter
```

#### `ctx recall lock`

Protect journal entries from being overwritten by `export --regenerate`.

```bash
ctx recall lock <pattern> [flags]
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

| Flag    | Description                |
|---------|----------------------------|
| `--all` | Unlock all journal entries |

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
ctx journal obsidian                        # Generate in .context/journal-obsidian/
ctx journal obsidian --output ~/vaults/ctx  # Custom output directory
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
