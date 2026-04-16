//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package edit implements the "ctx hook message edit"
// command.
//
// # Overview
//
// The edit command creates a local override file for a
// hook message template, allowing the user to customize
// what ctx injects into AI tool hooks. The override
// copy is written to the .context/messages/ directory
// where ctx reads it instead of the embedded default.
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
//  1. Looks up the hook/variant pair in the message
//     registry. Returns an error if unknown.
//  2. Checks whether an override file already exists;
//     returns an error if it does.
//  3. For ctx-specific templates, prints a warning
//     about scope limitations.
//  4. Reads the embedded default template content.
//  5. Creates the override directory and writes the
//     template file with restricted permissions.
//  6. Prints the override path, an edit hint, and
//     the available template variables.
//
// # Output
//
// Prints the path to the newly created override file,
// a hint to edit it, and a list of template variables
// that can be used in the template content.
package edit
