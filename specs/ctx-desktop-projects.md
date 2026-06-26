# ctx Desktop — Multi-Project Dashboard

A workspace-level "Projects" screen that surveys *every* ctx project
across one or more registered workspace roots at once, instead of
one-project-at-a-time. Extends `specs/ctx-desktop.md`; same thin-client
contract (all data via the `ctx` CLI, `.context/` stays the source of
truth).

## Multiple workspaces

The app registers a **list** of workspace roots (persisted as
`ctx.workspaces`, migrated from the legacy single `ctx.workspace`).
All roots are scanned and their projects merged, de-duped by path, and
sorted by name. Roots are added/removed from the Projects screen (chips)
or the sidebar; the top-bar dropdown and the dashboard both draw from
the merged set.

## Problem

Today the GUI views a single active project; the workspace scan only
feeds the switcher dropdown. Someone running several ctx projects has
no way to see, in one glance, which projects need attention — task
backlog, doctor health, drift, context bloat — or to jump into the one
that does. "What's happening across my projects?" requires clicking
through each in turn.

## Approach

A new `Projects` screen (workspace-level: keyed off the chosen
workspace, not the single active `dir`). It fans the existing
per-project `ctx` commands out across every discovered project and
aggregates the results into a card per project.

### Card (per project) — works on stock ctx 0.8.1

- **Task status** — open / done counts, parsed from the `TASKS.md`
  `summary` in `ctx status --json` (e.g. `"237 active, 11 completed"`).
  No new CLI capability required.
- **Health & drift** — `warnings`/`errors` from `ctx doctor --json`;
  drift surfaced from the `drift` check result.
- **Context size** — `total_files` / `total_tokens` from
  `ctx status --json`.
- Affordance to **open** the project (set it active and navigate to
  its Overview).

### Drill-down — "what's happening in each task"

Expanding a card reveals task detail. Graceful degradation by CLI:

- **With `ctx task list --json`** (feat/ctx-artifact-list-json): a row
  per task — status, and provenance (branch, commit, session, added).
  Each task's `session` links to its journal entry via
  `ctx journal source --show <session>` (lazy-loaded on expand).
- **On stock 0.8.1** (no `task list --json`): show the project's
  recent journal sessions (`ctx journal source --limit N`) as the
  activity feed, each expandable to full content via `--show`.

## Wired commands

Reuses `ctx_status`, `ctx_doctor`, `ctx_task_list`, `discover_projects`
per project. Adds one Rust command: `ctx_journal_show(dir, session)`
→ `ctx journal source --show <session>` (raw text, rendered verbatim).

## Constraints

- Local-only; all access through `ctx` (no direct `.context/` reads
  beyond the existing allowlisted viewers).
- Fan-out is concurrency-limited so a large workspace can't spawn a
  process storm.
- Task counts degrade to "—" when a project's `status --json` lacks a
  parseable TASKS.md summary; per-task rows degrade to the journal feed
  when `task list --json` is absent.

## Out of scope (v1)

Cross-project search, bulk actions across projects, editing tasks from
the dashboard, and any background polling — the screen refreshes on
open and on the active project's `ctx-changed` event.
