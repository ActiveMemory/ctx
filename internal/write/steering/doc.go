//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package steering provides the **terminal-output
// helpers** the `ctx steering` CLI subcommands use to
// narrate their `add`, `init`, `list`, `preview`, and
// `sync` operations.
//
// All exported functions take a `*cobra.Command` so
// they route through cobra's output stream (which
// tests can wire to a buffer for assertion).
//
// # Public Surface
//
// Output families:
//
//   - **Init** — [Created], [Skipped],
//     [InitSummary]. The `init` subcommand
//     announces each foundation file it
//     materializes (or skipped because it
//     already exists), then summarizes counts.
//   - **List / Preview** — [NoFilesFound],
//     [FileEntry], [FileCount], [NoFilesMatch],
//     [PreviewHeader], [PreviewEntry],
//     [PreviewCount]. Render the available
//     steering files and their inclusion-rule
//     match results against a sample prompt.
//   - **Sync** — [SyncWritten], [SyncSkipped],
//     [SyncError], [SyncSummary]. Per-tool
//     progress narration during
//     `ctx steering sync`.
//
// # Concurrency
//
// Pure data → io.Writer. Concurrent calls
// serialize through cobra's output stream.
//
// # Related Packages
//
//   - [internal/cli/steering]    — chief
//     consumer.
//   - [internal/steering]        — the engine
//     this package narrates.
//   - [internal/err/steering]    — typed
//     errors surfaced by [SyncError].
package steering
