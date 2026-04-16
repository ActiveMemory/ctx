//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package skill discovers and lists available skills for the
// ctx guide output.
//
// Skills are SKILL.md files with YAML frontmatter containing
// a name and description. The [List] function enumerates all
// installed skills via [internal/claude.SkillList], parses
// each file's frontmatter with [ParseFrontmatter], and
// prints a summary line per skill.
//
// # Frontmatter Parsing
//
// [ParseFrontmatter] extracts the YAML block delimited by
// "---" fences at the top of a SKILL.md file. If no valid
// frontmatter is found (missing fences or malformed YAML),
// it returns a zero [Meta] without error. Only YAML parse
// failures produce an error.
//
// # Description Truncation
//
// [TruncateDescription] shortens long descriptions for
// compact display. It prefers a natural sentence break
// (first ". " within the limit) over a hard cut. If the
// text is shorter than the limit, it passes through
// unchanged; otherwise it is cut at the limit with an
// ellipsis appended.
package skill
