//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package dismiss implements the "ctx remind dismiss"
// subcommand for removing reminders by ID.
//
// # Behavior
//
// The command removes one or more reminders from the
// JSON store. Arguments can be individual IDs or
// ranges (e.g., "3 5-7"). When the --all flag is
// set, all reminders are dismissed regardless of
// any positional arguments.
//
// If no IDs and no --all flag are provided, the
// command returns an error requesting at least one
// ID. Invalid or missing IDs also produce errors
// without modifying the store.
//
// The command is aliased so it can be invoked with
// a shorter name as well.
//
// # Flags
//
//	--all    Dismiss every reminder in the store.
//	         When set, positional ID arguments are
//	         ignored.
//
// # Output
//
// Each dismissed reminder prints a confirmation line.
// When --all is used, a summary of the total count
// is printed. On failure, returns an error for
// missing IDs or write problems.
//
// # Delegation
//
// Bulk dismissal is handled by [coreDismiss.All].
// Selective dismissal uses [coreDismiss.Many]. ID
// parsing and range expansion reuse [parse.IDs]
// from the pad core package.
package dismiss
