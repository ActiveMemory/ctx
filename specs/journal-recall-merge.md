# Merge recall into journal

## Problem

`ctx recall` and `ctx journal` are two commands operating on the same
domain — session history. `recall` handles ingest (list sources, show,
import, lock/unlock/sync) while `journal` handles publish (site,
obsidian). Users must know which command owns which stage. The name
"recall" also collides semantically with the `ctx-remember` skill,
which uses "recall" in its own help text.

## Approach

Absorb all `recall` subcommands into `journal`, introducing a `source`
subcommand for the pre-import inspection commands. Delete `ctx recall`
as a top-level command.

### Before → After

| Before                          | After                                   |
|---------------------------------|-----------------------------------------|
| `ctx recall list`               | `ctx journal source --list`             |
| `ctx recall show <id>`          | `ctx journal source --show <id>`        |
| `ctx recall import [id]`        | `ctx journal import [id]`               |
| `ctx recall lock <pattern>`     | `ctx journal lock <pattern>`            |
| `ctx recall unlock <pattern>`   | `ctx journal unlock <pattern>`          |
| `ctx recall sync`               | `ctx journal sync`                      |
| `ctx journal site`              | `ctx journal site` (unchanged)          |
| `ctx journal obsidian`          | `ctx journal obsidian` (unchanged)      |

### Why `source` with flags, not subcommands

`source` is a noun describing what you're looking at (raw session
files). `--list` and `--show` are view modes on that noun, not
independent entities. This keeps nesting to two levels max:
`ctx journal source --list` rather than `ctx journal source list`.

Default behavior (bare `ctx journal source`): equivalent to `--list`.

### MCP tool rename

The MCP tool `ctx_recall` becomes `ctx_journal_source`. This is a
breaking change for MCP clients, but the tool is young and has no
external consumers beyond the project's own skills.

### Skill rename

| Before           | After                    |
|------------------|--------------------------|
| `ctx-recall`     | `ctx-journal-browse`     |

The `ctx-remember` skill stays as-is. Its description can drop
the word "recall" now that there's no command name collision.

### Parser package stays

`internal/recall/parser` stays at its current path. It's a domain
parser for session files, not tied to the CLI command name. Renaming
it would churn 3+ consumers for no user-visible benefit.

## Behavior

### Happy Path

1. User runs `ctx journal source` → sees list of available sessions
   (same output as old `ctx recall list`)
2. User runs `ctx journal source --show <slug>` → sees session detail
3. User runs `ctx journal import --all` → imports to `.context/journal/`
4. User runs `ctx journal lock --all` → protects entries
5. User runs `ctx journal site` → generates static site (unchanged)

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| User runs old `ctx recall list` | Command not found (clean break, no alias) |
| `ctx journal source` with no flags | Defaults to `--list` behavior |
| `ctx journal source --show` without ID | Error: session ID required |
| `ctx journal source --list --show <id>` | Error: mutually exclusive flags |
| MCP client calls `ctx_recall` | Tool not found; `ctx_journal_source` is the replacement |

### No backwards compatibility shim

`recall` is removed cleanly. No hidden alias, no deprecation warning.
The command is internal-facing (AI agents, not end users with muscle
memory). Skills and hooks are updated in the same commit.

## Interface

### CLI

```
ctx journal source [flags]
ctx journal import [session-id] [flags]
ctx journal lock <pattern> [flags]
ctx journal unlock <pattern> [flags]
ctx journal sync
ctx journal site [flags]
ctx journal obsidian [flags]
```

#### `ctx journal source` flags

| Flag             | Short | Default | Description                              |
|------------------|-------|---------|------------------------------------------|
| `--list`         | `-l`  | true    | List available sessions (default)        |
| `--show`         | `-s`  | ""      | Show details of a specific session       |
| `--limit`        | `-n`  | 20      | Max sessions to list                     |
| `--project`      | `-p`  | ""      | Filter by project name                   |
| `--tool`         | `-t`  | ""      | Filter by tool                           |
| `--all-projects` |       | false   | Include all projects                     |
| `--full`         |       | false   | Full conversation (with --show)          |
| `--latest`       |       | false   | Show most recent session (with --show)   |
| `--since`        |       | ""      | Filter sessions since date               |
| `--until`        |       | ""      | Filter sessions until date               |

## Implementation

### Phase 1: Move recall subcommands into journal

| Layer | Files | Change |
|-------|-------|--------|
| CLI parent | `internal/cli/journal/journal.go` | Add import, lock, unlock, sync subcommands |
| New source cmd | `internal/cli/journal/cmd/source/` | New subcommand combining list+show with flag dispatch |
| Bootstrap | `internal/bootstrap/bootstrap.go` | Remove `recall.Cmd` registration |
| Config constants | `internal/config/embed/cmd/recall.go` | Rename to journal additions or delete; add `UseJournalSource` etc. |
| Config constants | `internal/config/embed/cmd/journal.go` | Add new Use/DescKey constants |
| Config constants | `internal/config/embed/cmd/base.go` | Remove `UseRecall` |
| Flag constants | `internal/config/embed/flag/recall.go` | Rename keys to `journal.source.*` / `journal.import.*` etc. |
| Text constants | `internal/config/embed/text/recall.go` | Rename keys to `journal.*` domain |
| Text constants | `internal/config/embed/text/err_recall.go` | Rename to `err_journal_source.go` or merge into existing |
| Text constants | `internal/config/embed/text/import.go` | Update `write.recall-import-*` → `write.journal-import-*` |
| Text constants | `internal/config/embed/text/mcp_recall.go` | Rename to `mcp_journal.go`, update keys |
| Text constants | `internal/config/embed/text/mcp_tool.go` | Update `DescKeyMCPToolRecallDesc` |
| YAML | `internal/assets/commands/commands.yaml` | Rename `recall.*` → `journal.source` / `journal.import` etc. |
| YAML | `internal/assets/commands/flags.yaml` | Rename `recall.*` → `journal.*` flag keys |
| YAML | `internal/assets/commands/text/write.yaml` | Rename `write.recall-*` → `write.journal-*` |
| YAML | `internal/assets/commands/text/ui.yaml` | Rename `recall.*` → `journal.*` |
| Write package | `internal/write/recall/` | Rename to `internal/write/journal/` (or merge if one exists) |
| Error package | `internal/err/recall/` | Rename to `internal/err/journal/` (or merge) |
| MCP tool | `internal/mcp/server/route/tool/dispatch.go` | Route `ctx_journal_source` instead of `ctx_recall` |
| MCP tool | `internal/mcp/server/route/tool/tool.go` | Rename `recall()` func |
| MCP handler | `internal/mcp/handler/tool.go` | Update method name if needed |
| MCP test | `internal/mcp/server/server_test.go` | Update tool name in tests |
| MCP text | `internal/config/embed/text/mcp_tool.go` | Rename desc key |

### Phase 2: Move core/ packages

| From | To |
|------|-----|
| `internal/cli/recall/core/extract/` | `internal/cli/journal/core/extract/` (if no collision) |
| `internal/cli/recall/core/format/` | Merge or namespace under journal |
| `internal/cli/recall/core/frontmatter/` | Check for collision with journal's frontmatter |
| `internal/cli/recall/core/index/` | `internal/cli/journal/core/index/` (if no collision) |
| `internal/cli/recall/core/lock/` | `internal/cli/journal/core/lock/` |
| `internal/cli/recall/core/plan/` | `internal/cli/journal/core/plan/` |
| `internal/cli/recall/core/query/` | `internal/cli/journal/core/query/` |
| `internal/cli/recall/core/slug/` | `internal/cli/journal/core/slug/` |
| `internal/cli/recall/core/validate/` | `internal/cli/journal/core/validate/` |
| `internal/cli/recall/core/confirm/` | `internal/cli/journal/core/confirm/` |
| `internal/cli/recall/core/execute/` | `internal/cli/journal/core/execute/` |

**Collision check needed**: `recall/core/format/` vs `journal/core/format/`,
and `recall/core/frontmatter/` vs `journal/core/frontmatter/`. If both exist,
merge or namespace (e.g., `core/importformat/`).

### Phase 3: Delete old recall CLI package

Remove `internal/cli/recall/` entirely after all references are moved.

### Phase 4: Update skills

| File | Change |
|------|--------|
| `internal/assets/claude/skills/ctx-recall/SKILL.md` | Rename dir to `ctx-journal-browse/`, update all commands |
| `internal/assets/claude/skills/ctx-remember/SKILL.md` | Remove `ctx recall list` reference, use `ctx journal source` |
| `internal/assets/claude/skills/ctx-journal-enrich/SKILL.md` | Check for recall references |
| `internal/assets/claude/skills/ctx-journal-enrich-all/SKILL.md` | Check for recall references |
| `internal/assets/claude/skills/ctx-journal-normalize/SKILL.md` | Check for recall references |
| `.claude/skills/generated/recall/SKILL.md` | Delete (regenerated on next init) |

### Phase 5: Update hooks and CLAUDE.md

| File | Change |
|------|--------|
| `CLAUDE.md` | `ctx recall list` → `ctx journal source` |
| `internal/assets/claude/CLAUDE.md` | Same |

### Phase 6: Update documentation

| File | Change |
|------|--------|
| `docs/cli/recall.md` | Rewrite as unified journal reference; rename file to `journal.md` |
| `docs/recipes/session-archaeology.md` | Update all `ctx recall` → `ctx journal` commands |
| `docs/recipes/session-lifecycle.md` | Same |
| `docs/recipes/session-ceremonies.md` | Same |
| `docs/recipes/publishing.md` | Same |
| `docs/recipes/autonomous-loops.md` | Same |
| `docs/recipes/multi-tool-setup.md` | Same |
| `docs/recipes/claude-code-permissions.md` | Same |
| `docs/recipes/parallel-worktrees.md` | Same |
| `docs/recipes/multilingual-sessions.md` | Same |
| `docs/recipes/hook-output-patterns.md` | Same |
| `docs/recipes/system-hooks-audit.md` | Same |
| `docs/recipes/index.md` | Same |
| `docs/cli/index.md` | Update nav if recall.md → journal.md |

### Phase 7: Update context files

| File | Change |
|------|--------|
| `.context/ARCHITECTURE.md` | Update recall references |
| `.context/AGENT_PLAYBOOK.md` | Update `ctx recall` → `ctx journal source` in tables |
| `.context/architecture-dia-*.md` | Update data flow diagrams |
| `specs/future-complete/recall-sync.md` | Update spec references |
| `specs/future-complete/recall-export-safety.md` | Update spec references |

### Phase 8: Rebuild site

```bash
make docs  # or equivalent site rebuild
```

## Configuration

No new `.ctxrc` keys. Existing `session_prefixes` config stays
in the parser package (unchanged).

## Testing

- **Unit**: Existing recall tests move with the code; `source`
  subcommand gets new tests for flag dispatch (`--list` default,
  `--show` requires ID, mutual exclusion)
- **Integration**: `recall_test.go` / `run_test.go` at recall
  package level → adapt for journal package
- **Smoke**: `ctx journal source`, `ctx journal source --show --latest`,
  `ctx journal import --dry-run`

## Non-Goals

- Renaming `internal/recall/parser` — it's a domain package, not a CLI surface
- Adding new functionality — this is a pure restructure
- Backwards-compatible aliases for `ctx recall`
- Changing the journal pipeline stages (import → normalize → enrich → lock)
- Blog post updates — historical references to `ctx recall` stay as-is

## Resolved Questions

1. **`internal/write/journal/` collision**: Both exist. No function name
   collisions — recall has 33 functions (import, list, show, lock output),
   journal has 4 functions (site/orphan output). Merge recall functions
   into `write/journal/` in a separate file (e.g., `source.go`, `import.go`).
   Delete `write/recall/`.

2. **`internal/err/journal/` collision**: Both exist. No function name
   collisions — recall has 6 error constructors, journal has 13. Merge
   recall errors into `err/journal/` in a separate file (e.g., `source.go`).
   Delete `err/recall/`.

3. **core/ package collisions**: `format/` and `frontmatter/` exist in
   both recall and journal core. Functions are completely different — no
   name collisions. Recall's `format/` has session formatting (Duration,
   Tokens, JournalFilename, etc.). Journal's `format/` has site formatting
   (FormatSize, KeyFileSlug, FormatSessionLink). Recall's `frontmatter/`
   has YAML writing helpers. Journal's has frontmatter transformation.
   **Strategy**: rename recall's packages during move to avoid ambiguity:
   `recall/core/format/` → `journal/core/sourceformat/`,
   `recall/core/frontmatter/` → `journal/core/sourcefm/`.
   Alternative: merge if the functions are complementary enough.
