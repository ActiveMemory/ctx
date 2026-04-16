//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package steering implements the "ctx steering" command
// group for managing steering files that define persistent
// behavioral rules for AI coding assistants.
//
// Steering files are markdown documents in .context/
// steering/ that contain glob-matched instructions. When
// an AI tool accesses files matching a steering file's
// globs, the corresponding rules are injected into the
// session context. This enables path-scoped behavioral
// customization without modifying the AI tool's config.
//
// # Subcommands
//
//   - add: create a new steering file with glob patterns
//     and rule content
//   - list: display all steering files with their globs
//     and rule summaries
//   - preview: given a prompt or file path, show which
//     steering rules would match
//   - initcmd: generate a set of foundation steering files
//     for common project patterns
//   - synccmd: export steering rules into the AI tool's
//     native format (e.g. .cursorrules)
//
// # Subpackages
//
//	cmd/add, cmd/list, cmd/preview: CRUD operations
//	cmd/initcmd: foundation file generation
//	cmd/synccmd: tool-native export
//	core: glob matching, rule parsing, and formatting
package steering
