//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package show implements the "ctx hook message show"
// command.
//
// # Overview
//
// The show command displays the content of a hook
// message template. If a local override exists in
// .context/messages/, that version is shown; otherwise
// the embedded default template is displayed.
//
// This lets the user inspect what content ctx will
// inject into a hook without needing to locate the
// override file or browse embedded assets.
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
//  1. Validates the hook/variant pair against the
//     message registry. Returns an error if unknown.
//  2. Checks for a local override file. If found,
//     prints a "source: override" header with the
//     file path, template variables, and content.
//  3. If no override exists, reads the embedded
//     default template and prints a "source: default"
//     header with template variables and content.
//
// # Output
//
// Prints the source label (override path or "default"),
// available template variables, and the full template
// content in a formatted block.
package show
