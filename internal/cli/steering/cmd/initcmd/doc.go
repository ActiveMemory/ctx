//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package initcmd implements the **`ctx steering init`**
// subcommand, which scaffolds the foundation steering
// files for a new project.
//
// # What It Does
//
// The command generates a set of foundation steering
// files (product, tech, structure, workflow) inside the
// .context/steering/ directory. Each file is created
// with YAML frontmatter set to:
//
//   - inclusion: always (injected on every prompt)
//   - priority: 10 (high priority)
//
// The body of each file contains a starter template
// drawn from [steering.FoundationFiles]. Files that
// already exist are silently skipped to prevent
// accidental overwrites.
//
// # Arguments
//
// None. The command accepts no positional arguments
// (cobra.NoArgs).
//
// # Output
//
// For each file, prints whether it was created or
// skipped. At the end, prints a summary line showing
// the count of created and skipped files.
//
// # Delegation
//
// [Cmd] builds the cobra command. [Run] checks that
// .context/ exists, ensures the steering subdirectory
// is present, iterates over [steering.FoundationFiles],
// serializes each via [steering.Print], and writes it
// to disk. Output is emitted through the
// [write/steering] formatters.
package initcmd
