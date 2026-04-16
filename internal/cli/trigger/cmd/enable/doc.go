//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package enable implements the "ctx trigger enable"
// cobra subcommand.
//
// This command enables a previously disabled trigger
// script by adding executable permission bits. Once
// enabled, the script will run during the matching
// lifecycle event.
//
// # Usage
//
//	ctx trigger enable <name>
//
// # Arguments
//
// Exactly one positional argument is required:
//
//   - name: the trigger script name to enable.
//     The script is located by searching all hook-type
//     subdirectories under .context/hooks/.
//
// # Behavior
//
// The command:
//
//   - Searches for the named trigger across all
//     hook-type directories using trigger.FindByName.
//   - Returns an error if no matching script is found.
//   - Adds executable permission bits (user, group,
//     and other) to the script file via chmod.
//   - Prints a confirmation with the trigger name
//     and file path.
//
// # Output
//
// A confirmation line showing the enabled trigger
// name and its file path.
//
// # Delegation
//
// Trigger discovery uses trigger.FindByName. Output
// formatting uses write/trigger.
package enable
