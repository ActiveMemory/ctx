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
//   - **Init**: [Created], [Skipped],
//     [InitSummary]. The `init` subcommand
//     announces each foundation file it
//     materializes (or skipped because it
//     already exists), then summarizes counts.
//   - **List / Preview**: [NoFilesFound],
//     [FileEntry], [FileCount], [NoFilesMatch],
//     [PreviewHeader], [PreviewEntry],
//     [PreviewCount]. Render the available
//     steering files and their inclusion-rule
//     match results against a sample prompt.
//   - **Sync**: [SyncWritten], [SyncSkipped],
//     [SyncError], [SyncSummary], [SyncDirect].
//     Per-tool progress narration during
//     `ctx steering sync`; [SyncDirect] is the
//     polite no-op line for tools that consume
//     steering via ctx agent (claude, codex).
//
// # Concurrency
//
// Pure data → io.Writer. Concurrent calls
// serialize through cobra's output stream.
package steering
