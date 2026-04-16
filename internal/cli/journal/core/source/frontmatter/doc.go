//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package frontmatter handles heading resolution and
// YAML field generation for journal source content.
//
// Each journal entry file starts with a YAML
// frontmatter block containing metadata such as title,
// session ID, date, and message count. This package
// provides helpers that the source formatter uses to
// build that block.
//
// # Heading Resolution
//
// [ResolveHeading] picks the best available heading
// from three candidates in priority order: title, slug,
// then base filename. The first non-empty value wins.
// The source formatter calls this to produce the
// Markdown H1 heading that follows the frontmatter.
//
// # YAML Field Writers
//
// Three functions write individual YAML fields to a
// strings.Builder:
//
//   - [WriteFmQuoted] writes a quoted string field
//     using the FmQuoted template, suitable for
//     values that may contain special characters.
//   - [WriteFmString] writes a bare string field
//     using the FmString template, for simple values
//     like dates or session IDs.
//   - [WriteFmInt] writes an integer field using the
//     FmInt template, for numeric metadata such as
//     message counts.
//
// Each writer appends a newline after the field. Write
// errors are silently discarded because the builder
// write to an in-memory buffer that does not fail.
//
// # Data Flow
//
// The source/format package calls these functions
// during journal file generation. Templates for field
// formatting live in the assets/tpl package. Newline
// and delimiter tokens come from config/token.
package frontmatter
