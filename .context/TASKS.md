# Tasks — Context CLI

# Tasks

## Phase 0: Cleanup from the previous version

(All tasks archived)

### Phase 1: Parser (DONE)

- [x] T1.1.0: Create CLI command (`ctx recall list`, `ctx recall show`) and
    slash command (`/ctx-recall`) for browsing AI session history.
- [x] T1.1.1: Define data structures in `internal/recall/parser/types.go`
- [x] T1.1.2: Implement line parser in `internal/recall/parser/claude.go`
- [x] T1.1.3: Implement session grouper (ParseFile with streaming)
- [x] T1.1.4: Implement directory scanner (ScanDirectory, FindSessions)

### Phase 2: Export & Search

- [ ] feat: `ctx recall export` - export sessions to editable journal files
  - `ctx recall export <session-id>` - export one session
  - `ctx recall export --all` - export all sessions
  - Skip existing files (user may have edited), `--force` to overwrite
  - Output to `.context/journal/YYYY-MM-DD-slug-shortid.md`
  #added:2026-01-28

- [ ] feat: `ctx recall search <query>` - CLI-based search across sessions
  - Simple text search, no server needed
  - IDE grep is alternative, this is convenience
  #priority:low

- [ ] explore: `ctx recall stats` - analytics/statistics
  - Token usage over time, tool patterns, session durations
  - Explore when we have a clear use case
  #priority:deferred

## Backlog

- [ ] Bug: `ctx tasks archive` doesn't archive nested content under completed tasks.
  When a parent `[x]` item has indented child lines (without checkboxes), only the
  parent line is archived, leaving orphaned content behind. The archive logic should
  include all indented lines that belong to a completed task.
  #added:2026-01-27 #priority:medium

- [ ] feat: ctx add learning requires --context, --lesson, --application flags (matching decision's ADR pattern) #priority:high #added:2026-01-28-053941
