//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package zensical

// MkDocs stripping constants (used by "ctx why" to clean embedded docs).
const (
	// MkDocsAdmonitionPrefix is the prefix for admonition lines in MkDocs.
	MkDocsAdmonitionPrefix = "!!!"
	// MkDocsTabPrefix is the prefix for tab marker lines in MkDocs.
	MkDocsTabPrefix = "=== "
	// MkDocsIndent is the 4-space indentation used in admonition/tab bodies.
	MkDocsIndent = "    "
	// MkDocsIndentWidth is the number of characters to dedent from body lines.
	MkDocsIndentWidth = 4
	// MkDocsFrontmatterDelim is the YAML frontmatter delimiter.
	MkDocsFrontmatterDelim = "---"
)
