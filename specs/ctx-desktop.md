# ctx Desktop — Tauri GUI Client (Sub-project B)

A cross-platform desktop client that is a thin lens over the
local `ctx` CLI. The `.context/` files remain the single source
of truth; every read and write shells out to `ctx`. This spec
covers the runnable scaffold (this commit) and the P0 roadmap.

Depends on Sub-project A (`specs/ctx-artifact-list-json.md`):
the GUI's counts and list views consume `ctx task|decision|
learning list --json`.

## Problem

`ctx`'s value is that context is durable, inspectable state, but
the human half of the loop is terminal-only. Reviewing stored
tasks/decisions/learnings, authoring a decision record, and
previewing the `ctx agent` packet within a token budget are all
clumsy in a shell. A GUI is the window into that state.

## Approach

Tauri 2 + React + TypeScript + Tailwind v4, in `ctx-desktop/`
(a subdirectory of this repo, versioned alongside the CLI).

- **Rust adapter** (`src-tauri/src/ctx_adapter.rs`): the single
  module that runs `ctx` via `std::process::Command` and returns
  results. No `tauri-plugin-shell` — direct process spawning, so
  no shell capability wiring. Prepends `/usr/local/bin` and
  `/opt/homebrew/bin` to PATH so a bundled app (minimal launchd
  PATH on macOS) can still find a user-installed `ctx`. Each call
  sets `current_dir` to the selected project root, because the
  CLI resolves its context from `$PWD/.context`.
- **TS adapter** (`src/adapter/ctx.ts`): typed `invoke` wrappers
  mirroring the Rust commands and the CLI JSON schemas. One file,
  so a CLI/output change is a one-file fix.
- **Screens** (`src/screens/`): React views. The scaffold ships
  Overview; P0 adds the rest.

### Scaffold (this commit)

A runnable shell: `npm install && npm run tauri dev` opens a
native window. The Overview screen detects the `ctx` binary
(version pill), takes a project path, and shows live counts —
open/done tasks, decisions, learnings (via A's `list --json`),
plus context-file and token totals from `status --json`.
Commands that fail (e.g. an old `ctx` without `list`) degrade
gracefully with an inline note rather than crashing the view.

Wired Rust commands: `ctx_info`, `ctx_status`, `ctx_doctor`,
`ctx_task_list`, `ctx_decision_list`, `ctx_learning_list`.

## P0 roadmap (not in the scaffold)

- Project switcher (recent projects, open-folder dialog,
  `ctx init`) replacing the hardcoded default path.
- `ctx doctor` health pill on the Overview.
- Tasks: list + inline add (`ctx task add`) + complete + filter.
- Decisions: browse/search + three-field authoring form
  (`ctx decision add --context --rationale --consequence`),
  synthesizing required provenance (`--session-id`, and
  `--branch`/`--commit` from git).
- Learnings: list + quick add.
- Journal timeline (`ctx journal source` — text adapter).
- Context Packet: budget slider re-running `ctx agent
  --format json --budget N`, with copy-packet / copy-command.
- fs-watch on the active `.context/` so external CLI/agent
  writes refresh the UI.

### Shipped surface (as of the ctx-desktop PR)

The client outgrew the P0 list above. The nav now ships, in
order: Projects, Overview, Search, Tasks, Reminders, Decisions,
Learnings, Conventions, Constitution (the last two via a shared
read-only `CanonicalDoc` viewer), Context Packet, Knowledge Base,
Scratchpad, Journal, Drift, Health, and Hub. The
Tasks/Decisions/Learnings/Search screens are gated behind a
capability probe for `ctx <artifact> list --json` (Sub-project A)
and show a "needs a newer ctx" notice when it is absent rather
than failing. Multi-workspace discovery (depth-bounded Rust scan)
and a top-bar `ctx doctor` health pill round out the shell.

## Constraints

- Local-only; no network except an explicit, user-triggered
  update check. No telemetry.
- Never read/write `.context/` directly — all access via `ctx`,
  preserving its audit trail and invariants.
- Detect (don't bundle) `ctx`; link to releases if missing.

## Out of scope (v1)

Hub/team sync, mobile, raw `.ctxrc` editing, an embedded LLM
assistant. The GUI manages the memory an external assistant
reads; it does not run model calls.
