//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package add implements the **`ctx steering add`**
// subcommand, which creates a new steering file with
// default frontmatter in the project steering directory.
//
// # What It Does
//
// Given a single positional argument (the file name
// without the .md extension), the command creates a
// new Markdown file under .context/steering/ with
// pre-populated YAML frontmatter:
//
//   - inclusion: manual (must be explicitly requested)
//   - priority: 50 (the default midpoint)
//
// The file body is left empty for the user to fill in.
// If a file with the same name already exists, the
// command returns an error instead of overwriting.
//
// # Arguments
//
//   - NAME (required): the steering file name. The
//     .md extension is appended automatically.
//
// # Output
//
// On success, prints the path of the created file and
// a hint about switching the inclusion mode. On failure
// (missing .context/, duplicate name) prints the
// corresponding error message.
//
// # Delegation
//
// [Cmd] builds the cobra command and enforces
// ExactArgs(1). [Run] checks that .context/ exists,
// ensures the steering subdirectory is present,
// serializes a [steering.SteeringFile] struct via
// [steering.Print], and writes the result to disk.
// Output messages are emitted through the
// [write/steering] formatters.
package add
