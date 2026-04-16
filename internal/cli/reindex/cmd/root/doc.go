//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the "ctx reindex" command
// for regenerating quick-reference indices in context
// files.
//
// # Behavior
//
// The command rebuilds the index section at the top
// of both DECISIONS.md and LEARNINGS.md in a single
// invocation. The index provides a compact table of
// contents that lets agents and humans quickly scan
// available entries without reading the full file.
//
// Each file is processed independently: if the first
// succeeds but the second fails, the first index is
// still written. The command reads each file from the
// context directory, parses its entries, generates
// the index block, and writes the updated file back.
//
// # Flags
//
// None. This command takes no flags.
//
// # Output
//
// For each file, prints a confirmation line to stdout
// indicating the file was reindexed. On failure,
// returns an error identifying which file could not
// be read or written.
//
// # Delegation
//
// Index generation is handled by [index.Reindex],
// which accepts an update function and entry parser.
// The context directory path comes from [rc.ContextDir].
// File names are defined in the [ctx] config package.
package root
