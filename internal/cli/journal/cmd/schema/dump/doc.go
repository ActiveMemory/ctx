//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package dump implements the "ctx journal schema dump"
// command.
//
// # Overview
//
// The dump command prints the embedded JSONL schema
// definition to stdout in a human-readable format. It
// shows the schema version, supported Claude Code version
// range, all known record types with their required and
// optional fields, and all recognized content block types
// with their parse status.
//
// # Output Format
//
// The output is structured as follows:
//
//  1. Schema version and CC version range header.
//  2. Record types section listing each type with its
//     required and optional field names. Metadata-only
//     types are shown without field lists.
//  3. Block types section listing each content block
//     type and whether it is "known" or "parsed".
//
// # Behavior
//
// [Cmd] builds a simple cobra.Command with no flags.
// [Run] loads the default schema, sorts record types
// and block types alphabetically, and writes each
// section to the command output stream. The command
// always returns nil.
//
// This is useful for understanding what the schema
// validator expects before running a check, or for
// documenting the current schema in external tooling.
package dump
