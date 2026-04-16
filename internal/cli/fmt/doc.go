//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package fmt provides the "ctx fmt" command for formatting
// context files to a consistent line width.
//
// The fmt command rewraps markdown text in .context/ files
// so that lines stay within a target width (default 72
// characters). This produces cleaner diffs and improves
// readability in terminals and editors.
//
// # Wrapping Rules
//
//   - Long lines are broken at word boundaries
//   - Markdown list items use 2-space continuation indent
//     when wrapped
//   - Headings are never wrapped
//   - Tables, YAML frontmatter blocks, and HTML comments
//     are preserved verbatim
//   - Code blocks (fenced and indented) are left untouched
//
// # Subpackages
//
//   - cmd/root: cobra command definition and flag binding
package fmt
