//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package list implements the "ctx trigger list" cobra
// subcommand.
//
// This command discovers all trigger scripts under
// .context/hooks/ and displays them grouped by hook
// type. Each entry shows the script name, whether it
// is enabled or disabled, and the file path.
//
// # Usage
//
//	ctx trigger list
//
// # Arguments
//
// None. The command takes no positional arguments.
//
// # Behavior
//
// The command:
//
//   - Calls trigger.Discover to scan all hook-type
//     subdirectories under .context/hooks/.
//   - Iterates over valid trigger types in a stable
//     order.
//   - For each type that has scripts, prints a type
//     header followed by one line per script showing
//     name, status (enabled/disabled), and path.
//   - Prints a total count at the end.
//   - If no scripts exist, prints a "no hooks found"
//     message.
//
// # Output
//
// A grouped listing of trigger scripts with type
// headers, per-script status lines, and a summary
// count. Example:
//
//	pre-tool-use:
//	  my-hook  enabled  .context/hooks/pre-tool-use/my-hook.sh
//
// # Delegation
//
// Script discovery uses trigger.Discover. Output
// formatting uses write/trigger.
package list
