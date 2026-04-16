//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package strip removes MkDocs-specific syntax from
// Markdown content so embedded documents read cleanly
// in the terminal. The why command uses this to convert
// documentation authored for MkDocs into plain Markdown.
//
// # MkDocs Syntax Handling
//
// [MkDocs] processes the following constructs:
//
//   - YAML frontmatter: delimited by "---" lines at the
//     start of the document; removed entirely
//   - Image references: lines matching ![alt](path);
//     removed entirely since terminals cannot display
//     images
//   - Admonitions: lines starting with "!!!" followed by
//     a type and quoted title; converted to a bold title
//     with the body lines prefixed as blockquotes
//   - Tab markers: lines starting with "===" followed by
//     a quoted name; converted to a bold name with body
//     lines dedented by 4 spaces
//   - Relative .md links: Markdown links pointing to
//     local .md files; replaced with just the display
//     text, since the link target is not navigable
//
// # Title Extraction
//
// [ExtractAdmonitionTitle] and [ExtractTabTitle] parse
// the quoted title from admonition and tab marker lines
// respectively. Both use simple string indexing to find
// the first and last double-quote characters.
package strip
