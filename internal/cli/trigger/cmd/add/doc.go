//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package add implements the "ctx trigger add" cobra
// subcommand.
//
// This command creates a new trigger script from a
// shell template and places it in the appropriate
// hook-type subdirectory under .context/hooks/. The
// script is created with executable permissions so it
// is immediately active.
//
// # Usage
//
//	ctx trigger add <hook-type> <name>
//
// # Arguments
//
// Exactly two positional arguments are required:
//
//   - hook-type: the lifecycle event to hook into
//     (e.g. "pre-tool-use", "post-tool-use"). Must
//     be one of the valid trigger types defined in
//     the trigger package.
//   - name: the script name (without .sh extension).
//     A ".sh" suffix is appended automatically.
//
// # Behavior
//
// The command:
//
//   - Validates the hook type against the list of
//     valid trigger types. Returns an error listing
//     valid types if the input is unknown.
//   - Ensures the type subdirectory exists under
//     .context/hooks/, creating it if needed with
//     restricted permissions.
//   - Checks that no script with the same name
//     already exists; returns an error if it does.
//   - Generates the script content from a built-in
//     shell template (tpl.TriggerScript) with the
//     name and hook type interpolated.
//   - Writes the file with executable permissions.
//
// # Output
//
// Prints the path to the newly created script file.
//
// # Delegation
//
// Type validation uses trigger.ValidTypes. Template
// content comes from assets/tpl. File operations use
// internal/io. Output formatting uses write/trigger.
package add
