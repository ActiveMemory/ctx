//   /    Context:                     https://ctx.ist
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
