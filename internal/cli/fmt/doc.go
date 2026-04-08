//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package fmt provides the "ctx fmt" command for formatting
// context files to a consistent line width.
//
// It wraps long lines at word boundaries, using 2-space
// continuation indent for markdown list items. Headings,
// tables, frontmatter, and HTML comments are preserved.
package fmt
