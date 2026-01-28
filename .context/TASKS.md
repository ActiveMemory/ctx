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

## Phase 1.a: Cleanup and Release

- [ ] T1.2.0 feat: ctx add learning requires --context, --lesson, --application flags (matching decision's ADR pattern) #priority:high #added:2026-01-28-053941
- [ ] T1.2.1 fix: update context-update XML tag format to include required fields (context, lesson, application for learnings; context, rationale, consequences for decisions) #priority:medium #added:2026-01-28-054914
- [ ] T1.2.2 chore: add tests to verify docs match implementation (caught drift in context-update format) #priority:low #added:2026-01-28-054915
- [ ] T1.2.3 refactor: ctx watch should use shared validation with ctx add (currently bypasses CLI, writes directly to files) #priority:medium #added:2026-01-28-055110
- [ ] T1.2.4 feat: /ctx-audit-docs slash command for semantic doc drift detection - reads docs and implementation, reports inconsistencies (AI-assisted, not deterministic tests) #priority:low #added:2026-01-28-055151
- [ ] T1.2.5: upstream CI is broken (again)
- [ ] T1.2.6: Human code review
- [ ] T1.2.7: cut a release (version number is already bumped)

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

- [ ] feat: ctx journal - LLM-powered session analysis and synthesis

Parent command for working with exported sessions (.context/journal/):

Subcommands to explore:
- ctx journal enrich: Add frontmatter/tags (topics, type, outcome, key files)
- ctx journal cluster: Group related sessions, build continuation chains
- ctx journal summarize: Generate timeline summaries, feature narratives
- ctx journal analyze: Find patterns (recurring mistakes, revisited decisions, coupling)
- ctx journal brief <topic>: Generate compressed context packet for a topic
- ctx journal site: Generate static site via zensical (browse, search, timeline)

Depends on: ctx recall export (Phase 2)
#priority:low #phase:future #added:2026-01-28-071638
