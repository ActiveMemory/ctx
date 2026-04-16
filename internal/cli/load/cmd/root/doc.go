//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the "ctx load" command.
//
// # Overview
//
// The load command reads context files from the .context/
// directory and outputs them in the recommended read
// order, suitable for piping into an AI assistant or
// reviewing manually.
//
// Two output modes are available: assembled (default)
// and raw. Assembled mode applies token-budget-aware
// truncation and adds section headers. Raw mode outputs
// file contents verbatim without headers or assembly.
//
// # Flags
//
//	--budget <n>   Token budget for assembled output
//	               (default 8000, or the value from
//	               the project config if set).
//	--raw          Output raw file contents without
//	               headers or assembly formatting.
//
// # Behavior
//
// [Cmd] builds the cobra.Command and registers the two
// flags. If --budget is not explicitly set on the
// command line, the configured project budget from rc
// is used instead.
//
// [Run] loads the context via context/load.Do, sorts
// files by read order, and dispatches to either
// writeLoad.Raw or writeLoad.Assembled depending on
// the --raw flag.
//
// If the .context/ directory does not exist, the command
// returns a "not initialized" error prompting the user
// to run "ctx init" first.
//
// # Output
//
// In assembled mode, outputs prioritized context with
// section headers and a token summary. In raw mode,
// outputs each file's contents separated by blank lines.
package root
