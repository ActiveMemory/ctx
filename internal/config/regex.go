//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import "regexp"

// RegExEntryHeader matches entry headers like "## [2026-01-28-051426] Title here".
//
// Groups:
//   - 1: date (YYYY-MM-DD)
//   - 2: time (HHMMSS)
//   - 3: title
var RegExEntryHeader = regexp.MustCompile(
	`## \[(\d{4}-\d{2}-\d{2})-(\d{6})] (.+)`,
)

// RegExLineNumber matches Claude Code's line number prefixes like "     1→".
var RegExLineNumber = regexp.MustCompile(`(?m)^\s*\d+→`)

// RegExSystemReminder matches <system-reminder>...</system-reminder> blocks.
// These are injected by Claude Code into tool results.
// Groups:
//   - 1: content between tags
var RegExSystemReminder = regexp.MustCompile(`(?s)<system-reminder>\s*(.*?)\s*</system-reminder>`)

// RegExCodeFenceInline matches code fences that appear inline after text.
// E.g., "some text: ```code" where fence should be on its own line.
// Groups:
//   - 1: preceding non-whitespace character
//   - 2: the code fence (3+ backticks)
var RegExCodeFenceInline = regexp.MustCompile("(\\S) *(```+)")

// RegExCodeFenceClose matches code fences immediately followed by text.
// E.g., "```text" where text should be on its own line after the fence.
// Groups:
//   - 1: the code fence (3+ backticks)
//   - 2: following non-whitespace character
var RegExCodeFenceClose = regexp.MustCompile("(```+) *(\\S)")

// RegExPhase matches phase headers at any heading level (e.g., "## Phase 1", "### Phase").
var RegExPhase = regexp.MustCompile(`^#{1,6}\s+Phase`)

// RegExBulletItem matches any Markdown bullet item (not just tasks).
//
// Groups:
//   - 1: item content
var RegExBulletItem = regexp.MustCompile(`(?m)^-\s*(.+)$`)

// RegExDecision matches decision entry headers in multiline content.
// Use for finding decision positions without capturing groups.
var RegExDecision = regexp.MustCompile(`(?m)^## \[\d{4}-\d{2}-\d{2}-\d{6}].*$`)

// RegExLearning matches learning entry headers in multiline content.
// Use for finding learning positions without capturing groups.
var RegExLearning = regexp.MustCompile(`(?m)^- \*\*\[\d{4}-\d{2}-\d{2}]\*\*.*$`)

// RegExNonFileNameChar matches characters not allowed in file names.
var RegExNonFileNameChar = regexp.MustCompile(`[^a-zA-Z0-9-]+`)

// RegExEntryHeading matches any entry heading (## [timestamp]).
// Use for counting entries without capturing groups.
var RegExEntryHeading = regexp.MustCompile(`(?m)^## \[`)

// RegExPath matches file paths in Markdown backticks.
//
// Groups:
//   - 1: file path
var RegExPath = regexp.MustCompile("`([^`]+\\.[a-zA-Z]{1,5})`")

// RegExContextUpdate matches context-update XML tags.
//
// Groups:
//   - 1: opening tag attributes (e.g., ` type="task" context="..."`)
//   - 2: content between tags
var RegExContextUpdate = regexp.MustCompile(`<context-update(\s+[^>]+)>([^<]+)</context-update>`)

// RegExGlossary matches glossary definition entries (lines with **term**).
var RegExGlossary = regexp.MustCompile(`(?m)(?:^|\n)\s*(?:-\s*)?\*\*[^*]+\*\*`)

// RegExDecisionPatterns detects decision-like phrases in text.
var RegExDecisionPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)decided to\s+(.{20,100})`),
	regexp.MustCompile(`(?i)decision:\s*(.{20,100})`),
	regexp.MustCompile(`(?i)we('ll| will) use\s+(.{10,80})`),
	regexp.MustCompile(`(?i)going with\s+(.{10,80})`),
	regexp.MustCompile(`(?i)chose\s+(.{10,80})\s+(over|instead)`),
}

// RegExLearningPatterns detects learning-like phrases in text.
var RegExLearningPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)learned that\s+(.{20,100})`),
	regexp.MustCompile(`(?i)gotcha:\s*(.{20,100})`),
	regexp.MustCompile(`(?i)lesson:\s*(.{20,100})`),
	regexp.MustCompile(`(?i)TIL:?\s*(.{20,100})`),
	regexp.MustCompile(`(?i)turns out\s+(.{20,100})`),
	regexp.MustCompile(`(?i)important to (note|remember):\s*(.{20,100})`),
}

// regExTaskPattern captures indent, checkbox state, and content.
//
// Pattern: ^(\s*)-\s*\[([x ]?)]\s*(.+)$
//
// Groups:
//   - 1: indent (leading whitespace, may be empty)
//   - 2: state ("x" for completed, " " or "" for pending)
//   - 3: content (task text)
const regExTaskPattern = `^(\s*)-\s*\[([x ]?)]\s*(.+)$`

// RegExTask matches a task item on a single line.
//
// Use with MatchString or FindStringSubmatch on individual lines.
// For multiline content, use RegExTaskMultiline.
var RegExTask = regexp.MustCompile(regExTaskPattern)

// RegExTaskMultiline matches task items across multiple lines.
//
// Use with FindAllStringSubmatch on multiline content.
var RegExTaskMultiline = regexp.MustCompile(`(?m)` + regExTaskPattern)

// RegExTaskDoneTimestamp extracts the #done: timestamp from a task line.
//
// Groups:
//   - 1: timestamp (YYYY-MM-DD-HHMMSS)
var RegExTaskDoneTimestamp = regexp.MustCompile(`#done:(\d{4}-\d{2}-\d{2}-\d{6})`)

// RegExClaudeTag matches Claude Code internal markup tags that leak into
// session titles via the first user message. This MUST remain an allowlist
// of known Claude Code tags — do NOT replace with a blanket regex.
var RegExClaudeTag = regexp.MustCompile(`</?(?:command-message|command-name|local-command-caveat)>`)

// Journal site pipeline patterns.

// RegExMultiPart matches session part files like "...-p2.md", "...-p3.md", etc.
var RegExMultiPart = regexp.MustCompile(`-p\d+\.md$`)

// RegExGlobStar matches glob-like wildcards: *.ext, */, *) etc.
var RegExGlobStar = regexp.MustCompile(`\*(\.\w+|[/)])`)

// RegExToolBold matches tool-use lines like "🔧 **Glob: .context/journal/*.md**".
var RegExToolBold = regexp.MustCompile(`🔧\s*\*\*(.+?)\*\*`)

// RegExInlineCodeAngle matches single-line inline code spans containing
// angle brackets (e.g., `</com`). Backticks are replaced with quotes and
// angles with HTML entities to prevent broken HTML in rendered output.
var RegExInlineCodeAngle = regexp.MustCompile("`([^`\n]*[<>][^`\n]*)`")

// RegExTurnHeader matches conversation turn headers.
//
// Groups:
//   - 1: turn number
//   - 2: role (e.g. "Assistant", "Tool Output")
//   - 3: timestamp (HH:MM:SS)
var RegExTurnHeader = regexp.MustCompile(`^### (\d+)\. (.+?) \((\d{2}:\d{2}:\d{2})\)$`)

// RegExFenceLine matches lines that are code fence markers (3+ backticks or
// tildes, optionally followed by a language tag).
var RegExFenceLine = regexp.MustCompile("^\\s*(`{3,}|~{3,})(.*)$")

// RegExMarkdownHeading matches Markdown heading lines (1-6 hashes + space).
//
// Groups:
//   - 1: hash prefix (e.g., "##")
//   - 2: heading text
var RegExMarkdownHeading = regexp.MustCompile(`^(#{1,6}) (.+)$`)

// RegExListStart matches lines that begin an ordered or unordered list item.
var RegExListStart = regexp.MustCompile(`^(\d+\.|[-*]) `)

// RegExFromAttrName creates a regex to extract an XML attribute value by name.
//
// Parameters:
//   - name: The attribute name to match
//
// Returns:
//   - *regexp.Regexp: Pattern matching name="value" with value in group 1
func RegExFromAttrName(name string) *regexp.Regexp {
	return regexp.MustCompile(name + `="([^"]*)"`)
}
