//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package tag implements the "ctx pad tag" subcommand
// for listing tags found across scratchpad entries.
//
// # Behavior
//
// The command scans all scratchpad entries and
// extracts inline tags (e.g., #topic). It collects
// unique tags, counts how many entries each tag
// appears in, and displays them sorted alphabetically.
//
// When no tags are found, the command prints a notice
// and exits. When entries contain tags, each tag is
// printed with its occurrence count.
//
// # Flags
//
//	--json    Output the tag list as a JSON array
//	          of objects with "tag" and "count"
//	          fields instead of the default
//	          human-readable format.
//
// # Output
//
// Default mode prints one line per tag with the tag
// name and occurrence count. JSON mode prints a
// single JSON array to stdout, suitable for piping
// into jq or other tools.
//
// When no tags exist, prints a "no tags" notice
// regardless of output mode.
//
// # Delegation
//
// Tag extraction from individual entries is handled
// by [tag.Extract] in the core/tag package. Entry
// reading goes through [store.ReadEntries]. Output
// formatting is routed through [writePad].
package tag
