//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package reset implements the "ctx hook message reset"
// command.
//
// # Overview
//
// The reset command removes a local override file for a
// hook message template, restoring the embedded default.
// After reset, ctx will use the built-in template for
// that hook/variant pair.
//
// # Arguments
//
// Requires exactly two positional arguments:
//
//  1. hook     The hook name (e.g. "PreToolUse").
//  2. variant  The template variant (e.g. "default").
//
// # Flags
//
// This command accepts no flags.
//
// # Behavior
//
// [Cmd] builds a cobra.Command requiring exactly two
// positional arguments. [Run] performs these steps:
//
//  1. Validates the hook/variant pair against the
//     message registry. Returns an error if unknown.
//  2. Removes the override file at the computed path.
//  3. If the file does not exist, prints a "no
//     override" message and returns nil.
//  4. Attempts to clean up the now-empty parent
//     directories (hook dir and messages dir).
//  5. Prints a confirmation that the override was
//     removed.
//
// Directory cleanup failures are logged as warnings
// but do not cause the command to fail, since the
// directories may contain other override files.
//
// # Output
//
// Prints either a "no override found" notice or an
// "override removed" confirmation with the hook and
// variant names.
package reset
