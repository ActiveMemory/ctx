//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package synccmd implements the **`ctx steering sync`**
// subcommand, which converts steering files into
// tool-native instruction formats.
//
// # What It Does
//
// Steering files are tool-agnostic Markdown documents.
// This command reads them from the steering directory
// and writes tool-specific output (for example,
// CLAUDE.md for Claude Code, .github/copilot-
// instructions.md for Copilot). The target tool is
// resolved from the --tool flag or from .ctxrc config.
//
// When the --all flag is provided, the command syncs
// to every supported tool in a single pass.
//
// # Flags
//
//   - **--all**: sync to all supported tools instead
//     of a single resolved tool. Overrides --tool and
//     .ctxrc tool settings.
//
// # Output
//
// Prints a sync report showing which files were written
// and to which tool targets. The report is emitted via
// the [cli/steering/core/sync.PrintReport] formatter.
//
// # Delegation
//
// [Cmd] builds the cobra command and binds the --all
// flag. [Run] calls [steering.SyncAll] when --all is
// set, or resolves the target tool via [resolve.Tool]
// and calls [steering.SyncTool] for a single tool.
// Both paths delegate to [cli/steering/core/sync] for
// report formatting.
package synccmd
