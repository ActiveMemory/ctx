//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package health holds the **shared helpers** that the
// architecture-map staleness, knowledge-file growth, and
// background-task cleanup hooks all use to evaluate
// "context health" signals — things that are not strictly
// drift but indicate the project is drifting from its own
// process invariants.
//
// The package is the *measurement layer*; the hooks
// (`check_map_staleness`, `check_knowledge`,
// background pruners) decide what to do with the numbers.
//
// # Public Surface
//
//   - **[ReadMapTracking]** — reads the persisted
//     architecture-map last-update tracking record.
//     Used by `check_map_staleness` to decide whether
//     ARCHITECTURE.md has fallen behind code changes.
//   - **[CountModuleCommits](module, since)** — counts
//     git commits touching a module path since a given
//     timestamp. Used to score map staleness.
//   - **[EmitMapStalenessWarning](staleModules)** —
//     produces the formatted nudge sent to the agent
//     via the VERBATIM relay path.
//   - **[UUIDPattern]** — compiled regex for matching
//     session UUIDs in state file names. Used by the
//     auto-pruner.
//   - **[AutoPrune](dir, age)** — removes per-session
//     state files older than `age`. Idempotent and
//     safe to run during a session (skips the active
//     session's marker file).
//
// # Concurrency
//
// All functions are filesystem-bound and stateless.
// Concurrent invocations are safe; the auto-pruner
// uses `os.Remove` which is atomic on POSIX.
//
// # Related Packages
//
//   - [internal/cli/system/cmd/check_map_staleness] —
//     the hook that consumes [ReadMapTracking] and
//     [EmitMapStalenessWarning].
//   - [internal/cli/system/cmd/check_knowledge]   —
//     consumes the entry-count helpers via
//     [internal/drift].
//   - [internal/prune]                            —
//     the `ctx prune` command that calls [AutoPrune]
//     under explicit user invocation.
package health
