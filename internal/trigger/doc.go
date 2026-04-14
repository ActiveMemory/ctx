//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package trigger implements ctx's **lifecycle automation**
// layer: project-owned shell scripts that run when an AI
// session crosses a defined boundary (start, end, before a
// tool call, after a tool call, file save, context-add).
//
// Triggers are the **imperative** counterpart to
// [internal/steering]. Steering says "when the user asks
// about Y, prepend these rules"; triggers say "when X
// happens, run this code." Together they form the
// behavior-customization plane the CLI exposes via
// `ctx trigger ...` and the docs cover in
// `docs/home/triggers.md`.
//
// # Trigger Types
//
// Six lifecycle events are supported (returned by
// [ValidTypes], see [internal/config/trigger]):
//
//   - **session-start** — a new AI session begins. Common
//     uses: inject standup notes, rotate context, log a
//     start marker.
//   - **session-end** — an AI session ends. Common uses:
//     persist summaries, ship notifications, capture
//     transcripts.
//   - **pre-tool-use** — before a tool call executes. Can
//     **block** the call (cancel:true) — used for safety
//     gates, audit logging, and policy enforcement.
//   - **post-tool-use** — after a tool call completes. Used
//     for reactions, lint-on-save, and post-processing.
//   - **file-save** — a file is saved. Lint, regenerate
//     indices, update derived files.
//   - **context-add** — a new entry was added to
//     `.context/`. Cross-link, notify, enrich.
//
// Each script receives a JSON [HookInput] on stdin and is
// expected to emit a JSON [HookOutput] on stdout (`cancel`,
// `context`, `message`).
//
// # Discovery — Layout under `.context/hooks/`
//
// [Discover] scans `.context/hooks/<type>/*.sh` (and any
// other extension; the executable bit is what counts, not
// the suffix) and returns one [HookInfo] per discovered
// script. A script's `Enabled` flag is **the executable
// permission bit**: an enabled hook fires, a non-executable
// one is reported by `ctx trigger list` but skipped at
// run-time.
//
// [FindByName] is the lookup helper used by `ctx trigger
// enable/disable/test` to address a single script by stem.
//
// # Security — The Disabled-by-Default Contract
//
// **A trigger is a shell script that runs with the same
// privileges as your AI tool.** A buggy or malicious one can
// block tool calls, corrupt context files, or exfiltrate
// data. The package therefore enforces a strict
// security-first workflow:
//
//  1. `ctx trigger add` creates new scripts **without** the
//     executable bit. They are inert until the user opts in.
//
//  2. The user reviews the script and runs `ctx trigger
//     enable <name>`, which sets the executable bit
//     **after** [ValidatePath] has run.
//
//  3. [ValidatePath] enforces three rules at every
//     execution:
//
//     - **No symlinks** — `os.Lstat` is used; symlinks
//     under `.context/hooks/` are rejected outright.
//     - **Boundary check** — the resolved absolute path
//     must lie within the absolute hooks directory; a
//     path that escapes via `..` is rejected.
//     - **Executable bit must be set** — the same bit
//     that gates discovery.
//
//     A failure here returns a typed error from
//     [internal/err/trigger] — the runner refuses to
//     `exec`.
//
// `docs/home/triggers.md` makes this contract explicit to
// users; the package enforces it.
//
// # Execution — How `RunAll` Behaves
//
// [RunAll] runs every enabled hook for the given type **in
// alphabetical order**, marshals [HookInput] to JSON on
// stdin, reads [HookOutput] from stdout, and aggregates the
// result into an [AggregatedOutput]. Per-hook semantics:
//
//   - `cancel:true` in stdout → **halt the chain**, set
//     `Cancelled` and `Message` on the aggregate, return
//     immediately.
//   - non-empty `context` field → append to
//     `AggregatedOutput.Context` (concatenated across
//     hooks).
//   - non-zero exit code → log via [ctxLog], record in
//     `Errors`, **continue** with the next hook.
//   - invalid JSON on stdout → log warning, record in
//     `Errors`, continue.
//   - timeout exceeded → kill the process group, log
//     warning, continue. Default is [DefaultTimeout]; can
//     be overridden by the caller.
//
// "One bad hook does not abort the chain" is intentional:
// security gates fire-and-forget, automation hooks fail
// loud-but-non-fatal, and only an explicit
// `cancel:true` short-circuits the rest.
//
// # The Three Hook-Like Layers
//
// The user-facing docs (`docs/home/triggers.md`) call out
// that ctx has **three** distinct hook concepts; only this
// package owns the first:
//
//   - **`ctx trigger`** (this package) — project-authored
//     scripts under `.context/hooks/`, fire on lifecycle
//     events, work with any AI tool.
//   - **`ctx system` hooks** — built-in nudges shipped by
//     ctx itself (see `internal/cli/system/cmd/check_*`).
//     Wired into tool configs at `ctx init` time.
//   - **Claude Code hooks** — Claude-Code-specific entries
//     in `.claude/settings.local.json`. Tool-native, not
//     portable.
//
// # Concurrency
//
// The package holds no mutable global state. [RunAll] runs
// hooks **sequentially** within a single invocation —
// alphabetical order is part of the contract — but
// concurrent invocations from different goroutines are
// safe.
//
// # Related Packages
//
//   - [github.com/ActiveMemory/ctx/internal/cli/trigger] —
//     the `ctx trigger` CLI: add, list, test, enable,
//     disable.
//   - [github.com/ActiveMemory/ctx/internal/config/trigger]
//     — the [TriggerType] enum and lifecycle constants.
//   - [github.com/ActiveMemory/ctx/internal/entity] —
//     [TriggerSession], [TriggerInput] — the input payload
//     types.
//   - [github.com/ActiveMemory/ctx/internal/err/trigger] —
//     typed error constructors used by [ValidatePath] and
//     the runner.
//   - [github.com/ActiveMemory/ctx/internal/drift] —
//     `checkHookPerms` flags any trigger script lacking the
//     executable bit so the user can re-enable it.
package trigger
