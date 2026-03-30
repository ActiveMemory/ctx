//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package frontmatter handles heading and YAML field generation
// for session source content.
//
// [ResolveHeading] picks the best heading from title, slug, or
// base name. [WriteFmQuoted] and [WriteFmString] write individual
// YAML frontmatter fields to a string builder.
package frontmatter
