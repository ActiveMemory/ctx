//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package write centralizes user-facing terminal output
// for every CLI command in ctx.
//
// All formatted messages, error output, progress lines,
// and informational text that CLI commands print to the
// user are routed through subpackages of write. This
// ensures consistent prefixes, templates, and output
// routing (stdout vs. stderr) across the entire CLI
// surface.
//
// # Organization
//
// Each subpackage corresponds to one command or feature
// area. For example, write/archive handles task archival
// output, write/bootstrap handles the bootstrap command,
// and write/err provides shared error formatting.
//
// # Conventions
//
// Functions accept a *cobra.Command to write to the
// correct output stream. Nil commands are treated as
// no-ops, making it safe to call from code paths where
// a command may not be available. Functions accept
// primitive types (strings, ints) rather than domain
// types to avoid coupling write packages to business
// logic.
//
// Message text is never hardcoded in write packages.
// All strings are loaded from the embedded descriptor
// system via internal/assets/read/desc, keyed by
// constants defined in internal/config/embed/text.
//
// # Message Categories
//
// Output falls into three categories:
//
//   - Info: progress confirmations, status lines, and
//     success messages printed to stdout via cmd.Println.
//   - Errors: prefixed error messages printed to stderr
//     via cmd.PrintErrln (see write/err).
//   - Warnings: non-fatal file or config warnings printed
//     to stderr with a warning prefix.
package write
