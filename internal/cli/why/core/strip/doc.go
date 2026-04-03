//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package strip removes MkDocs-specific syntax from Markdown
// content for terminal display.
//
// [MkDocs] handles frontmatter, images, admonitions,
// tabs, and relative links. [ExtractAdmonitionTitle] and
// [ExtractTabTitle] parse quoted titles from MkDocs markers.
package strip
