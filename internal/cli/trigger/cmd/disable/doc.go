//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package disable implements the "ctx trigger disable"
// cobra subcommand.
//
// This command disables a trigger script by removing
// its executable permission bit. The script file is
// preserved on disk but will not run during lifecycle
// events because the trigger runner skips non-
// executable files.
//
// # Usage
//
//	ctx trigger disable <name>
//
// # Arguments
//
// Exactly one positional argument is required:
//
//   - name: the trigger script name to disable.
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
//   - Removes the executable permission bits (user,
//     group, and other) from the script file via
//     chmod.
//   - Prints a confirmation with the trigger name
//     and file path.
//
// # Output
//
// A confirmation line showing the disabled trigger
// name and its file path.
//
// # Delegation
//
// Trigger discovery uses trigger.FindByName. Output
// formatting uses write/trigger.
package disable
