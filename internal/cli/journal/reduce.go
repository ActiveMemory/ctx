//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"encoding/json"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
)

// stripFences removes all code fence markers from content, leaving the inner
// text as-is. This eliminates fence nesting conflicts entirely. Files whose
// fences have been verified (fencesVerified=true) are returned unchanged.
//
// The result is plain text with structural markers preserved (turn headers,
// tool calls, section breaks). Serves as a readable baseline without AI
// reconstruction, or as input for the ctx-journal-normalize skill.
//
// Parameters:
//   - content: Raw Markdown content of a journal entry
//   - fencesVerified: Whether the file's fences have been verified via state
//
// Returns:
//   - string: Content with code fence markers removed
func stripFences(content string, fencesVerified bool) string {
	// Skip files whose fences have been verified by the AI skill
	if fencesVerified {
		return content
	}

	lines := strings.Split(content, config.NewlineLF)
	var out []string
	inFrontmatter := false

	for i, line := range lines {
		// Preserve frontmatter
		if i == 0 && strings.TrimSpace(line) == config.Separator {
			inFrontmatter = true
			out = append(out, line)
			continue
		}
		if inFrontmatter {
			out = append(out, line)
			if strings.TrimSpace(line) == config.Separator {
				inFrontmatter = false
			}
			continue
		}

		// Remove fence markers
		if config.RegExFenceLine.MatchString(line) {
			continue
		}

		out = append(out, line)
	}

	return strings.Join(out, config.NewlineLF)
}

// stripSystemReminders removes internal Claude Code blocks from journal content.
// Handles:
//   - XML-style system reminders: <system-reminder>...</system-reminder>
//   - Bold-style system reminders: **System Reminder**: ... (paragraph until blank line)
//   - Context compaction summaries: multi-line <summary>...</summary> blocks
//     (standalone <summary> on its own line — see config.TagCompactionSummaryOpen)
//   - Compaction continuation boilerplate: "If you need specific details from
//     before compaction..." paragraph
//
// The authoritative JSONL transcripts retain them; the exported Markdown
// doesn't need them.
//
// Parameters:
//   - content: Journal entry content with potential internal blocks
//
// Returns:
//   - string: Content with all internal blocks removed
func stripSystemReminders(content string) string {
	lines := strings.Split(content, config.NewlineLF)
	var out []string
	inTagReminder := false
	inBoldReminder := false
	inCompaction := false
	inBoilerplate := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// XML-style: <system-reminder>...</system-reminder>
		if trimmed == config.TagSystemReminderOpen {
			inTagReminder = true
			continue
		}
		if inTagReminder {
			if trimmed == config.TagSystemReminderClose {
				inTagReminder = false
			}
			continue
		}

		// Bold-style: **System Reminder**: ... (runs until blank line)
		if strings.HasPrefix(trimmed, config.LabelBoldReminder) {
			inBoldReminder = true
			continue
		}
		if inBoldReminder {
			if trimmed == "" {
				inBoldReminder = false
			}
			continue
		}

		// Context compaction: standalone <summary> on its own line.
		// Single-line <summary>N lines</summary> (ours) won't match
		// because trimmed != "<summary>" when there's inline content.
		if trimmed == config.TagCompactionSummaryOpen {
			inCompaction = true
			continue
		}
		if inCompaction {
			if trimmed == config.TagCompactionSummaryClose {
				inCompaction = false
			}
			continue
		}

		// Compaction boilerplate: "If you need specific details from
		// before compaction..." paragraph (runs until blank line)
		if strings.HasPrefix(trimmed, config.CompactionBoilerplatePrefix) {
			inBoilerplate = true
			continue
		}
		if inBoilerplate {
			if trimmed == "" {
				inBoilerplate = false
			}
			continue
		}

		out = append(out, line)
	}

	return strings.Join(out, config.NewlineLF)
}

// cleanToolOutputJSON extracts plain text from Tool Output turns whose body is
// raw JSON from the Claude API (e.g. [{"type":"text","text":"..."}]).
// The JSON text field's \n escapes become real newlines.
//
// Parameters:
//   - content: Journal entry content with potential JSON tool output
//
// Returns:
//   - string: Content with JSON tool output replaced by plain text
func cleanToolOutputJSON(content string) string {
	lines := strings.Split(content, config.NewlineLF)
	var out []string
	i := 0

	for i < len(lines) {
		matches := config.RegExTurnHeader.FindStringSubmatch(
			strings.TrimSpace(lines[i]),
		)
		if matches == nil || matches[2] != config.LabelToolOutput {
			out = append(out, lines[i])
			i++
			continue
		}

		// Tool Output header
		out = append(out, lines[i])
		i++

		// Collect body until next header
		bodyStart := i
		for i < len(lines) {
			if config.RegExTurnHeader.MatchString(strings.TrimSpace(lines[i])) {
				break
			}
			i++
		}
		bodyLines := lines[bodyStart:i]

		// Strip code fences wrapping the body, then rejoin and try JSON parse
		var nonEmpty []string
		for _, l := range bodyLines {
			t := strings.TrimSpace(l)
			if t == "" || strings.HasPrefix(t, "```") {
				continue
			}
			nonEmpty = append(nonEmpty, t)
		}
		body := strings.Join(nonEmpty, " ")

		if strings.HasPrefix(body, "[{") {
			var items []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			}
			if json.Unmarshal([]byte(body), &items) == nil && len(items) > 0 {
				out = append(out, "")
				for _, item := range items {
					out = append(out, item.Text)
				}
				out = append(out, "")
				continue
			}
		}

		// Not JSON or parse failed — keep original
		out = append(out, bodyLines...)
	}

	return strings.Join(out, config.NewlineLF)
}
