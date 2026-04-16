//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package preview extracts short content previews from
// context files for status display.
//
// # Content Extraction
//
// [Content] returns the first n non-empty, meaningful
// lines from file content. The extraction algorithm:
//
//  1. Splits content on newline boundaries.
//  2. Skips empty lines.
//  3. Tracks YAML frontmatter delimiters (---) and
//     skips all lines within frontmatter blocks.
//  4. Skips HTML comment lines that start with <!--.
//  5. Truncates lines longer than 60 characters,
//     appending an ellipsis.
//  6. Collects up to n qualifying lines.
//
// This produces a compact preview suitable for both
// JSON output and terminal text display, filtering
// out metadata and boilerplate that would not help
// the user understand file contents at a glance.
package preview
