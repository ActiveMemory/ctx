//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// MarkdownHeading matches Markdown heading lines (1-6 hashes + space).
//
// Groups:
//   - 1: hash prefix (e.g., "##")
//   - 2: heading text
var MarkdownHeading = regexp.MustCompile(`^(#{1,6}) (.+)$`)

// TurnHeader matches conversation turn headers.
//
// Groups:
//   - 1: turn number
//   - 2: role (e.g. "Assistant", "Tool Output")
//   - 3: timestamp (HH:MM:SS)
var TurnHeader = regexp.MustCompile(`^### (\d+)\. (.+?) \((\d{2}:\d{2}:\d{2})\)$`)

// ListStart matches lines that begin an ordered or unordered list item.
var ListStart = regexp.MustCompile(`^(\d+\.|[-*]) `)

// MarkdownLink matches Markdown links with relative .md targets.
var MarkdownLink = regexp.MustCompile(`\[([^]]+)]\([^)]*\.md[^)]*\)`)

// MarkdownLinkAny matches any Markdown link: [display](target).
//
// Groups:
//   - 1: display text
//   - 2: target URL/path
var MarkdownLinkAny = regexp.MustCompile(`\[([^]]+)]\(([^)]+)\)`)

// MarkdownImage matches Markdown image lines.
var MarkdownImage = regexp.MustCompile(`^\s*!\[.*]\(.*\)\s*$`)

// ToolBold matches tool-use lines like "🔧 **Glob: .context/journal/*.md**".
var ToolBold = regexp.MustCompile(`🔧\s*\*\*(.+?)\*\*`)

// InlineCodeAngle matches single-line inline code spans containing
// angle brackets (e.g., `</com`). Backticks are replaced with quotes and
// angles with HTML entities to prevent broken HTML in rendered output.
var InlineCodeAngle = regexp.MustCompile("`([^`\n]*[<>][^`\n]*)`")

// Phase matches phase headers at any heading level
// (e.g., "## Phase 1", "### Phase").
var Phase = regexp.MustCompile(`^#{1,6}\s+Phase`)

// BulletItem matches any Markdown bullet item (not just tasks).
//
// Groups:
//   - 1: item content
var BulletItem = regexp.MustCompile(`(?m)^-\s*(.+)$`)
