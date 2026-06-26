# ctx-desktop API Hardening

## Problem

A line-by-line review of the desktop IPC surface (the Tauri
command layer wrapping `ctx` CLI spawns — there is no HTTP
server) found 13 issues: 3 high, 5 medium, 5 low. The common
threads: the IPC layer inherited CLI assumptions that don't
hold in a long-lived GUI (blocking spawns on the UI thread,
positional text parsed as flags), and the watcher/event design
carried no project identity, forcing full-grid refreshes.

## Design

One hardening pass, no behavior additions:

1. **Async command layer** (H1): every `#[tauri::command]` is
   async; spawns and fs walks run via `spawn_blocking`. The UI
   thread never parks in `recv_timeout`.
2. **Exact-text task completion** (H2): the UI completes tasks
   by exact text, not by a cached pending-number. Rationale:
   `.context/` has no file locking and agents write concurrently;
   number resolution at file-read time silently completes the
   wrong task. The CLI's ambiguity error surfaces in the UI.
3. **Conditional provenance flags** (H3): `provenance_args()`
   omits `--branch`/`--commit` when the project has no usable
   git state, since the CLI rejects empty values by default.
4. **Capability probe** (M1): `ctx_capabilities` runs
   `ctx task list --json` once per project; screens that need
   the list contract gate on it instead of hard-failing against
   a ctx built without the sibling list-json branch. No `--help`
   probing (cobra reports success for unknown subcommands).
5. **Child reaping** (M2): `run_bin` spawns, polls `try_wait`,
   kills and reaps on the 30s deadline. No leaked threads or
   orphaned ctx processes.
6. **Scan supersession** (M3): `scanAll` carries a request id;
   stale scans cannot overwrite newer project sets or re-point
   the watcher.
7. **Flag/positional separation** (M4): every write path emits
   flags, then `--`, then user text — user input can never be
   parsed as a flag.
8. **Attributed change events** (M5/L1): watcher events carry
   the project root; only the affected card refreshes; the
   summary grid never blanks out.
9. **Defense-in-depth** (L2-L5): kb reads canonicalized and
   prefix-checked, topic recursion depth-bounded, drill-down
   cache invalidated on change events, custom ctx binary path
   validated (`--version` mentions ctx) before persisting.

## Out of Scope

- File locking for `.context/` (repo-wide gap, predates this
  branch; tracked separately)
- Merging the `feat/ctx-artifact-list-json` contract (the
  capability probe makes its absence survivable)

## Verification

`cargo check`, `cargo clippy --all-targets` (0 warnings),
`cargo test`, and `tsc && vite build` all pass; full repo
`make lint` + `make test` green.
