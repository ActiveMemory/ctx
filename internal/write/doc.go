//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package write centralizes user-facing output for CLI commands.
//
// All formatted messages, error output, and informational lines that CLI
// commands to print to the user are routed through this package. This ensures
// consistent prefixes, templates, and output routing (stdout vs. stderr)
// across the entire CLI surface.
//
// Functions accept a *cobra.Command to write to the correct output stream.
// Nil commands are treated as no-ops, making it safe to call from code
// paths where a command may not be available.
package write
