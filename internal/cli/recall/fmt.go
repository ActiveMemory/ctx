//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/recall/parser"
)

// fenceForContent returns the appropriate code fence for content.
//
// Uses longer fences when content contains backticks to avoid
// nested Markdown rendering issues. Starts with ``` and adds
// more backticks as needed.
//
// Parameters:
//   - content: The content to be fenced
//
// Returns:
//   - string: A fence string (e.g., "```", "````", "`````")
func fenceForContent(content string) string {
	fence := config.CodeFence
	for strings.Contains(content, fence) {
		fence += config.Backtick
	}
	return fence
}

// formatJournalFilename generates the filename for a journal entry.
//
// Format: YYYY-MM-DD-slug-shortid.md
// Uses local time for the date.
//
// Parameters:
//   - s: Session to generate filename for
//
// Returns:
//   - string: Filename like "2026-01-15-gleaming-wobbling-sutherland-abc12345.md"
func formatJournalFilename(s *parser.Session) string {
	date := s.StartTime.Local().Format("2006-01-02")
	shortID := s.ID
	if len(shortID) > config.RecallShortIDLen {
		shortID = shortID[:config.RecallShortIDLen]
	}
	return fmt.Sprintf(config.TplRecallFilename, date, s.Slug, shortID)
}

// formatJournalEntryPart generates Markdown content for a part of a journal entry.
//
// Includes metadata, tool usage summary (on part 1 only), navigation links,
// and the conversation subset for this part.
//
// Parameters:
//   - s: Session to format
//   - messages: Subset of messages for this part
//   - startMsgIdx: Starting message index (for numbering)
//   - part: Current part number (1-indexed)
//   - totalParts: Total number of parts
//   - baseName: Base filename without extension (for navigation links)
//
// Returns:
//   - string: Markdown content for this part
func formatJournalEntryPart(
	s *parser.Session,
	messages []parser.Message,
	startMsgIdx, part, totalParts int,
	baseName string,
) string {
	var sb strings.Builder
	nl := config.NewlineLF
	sep := config.Separator

	// Header
	if s.Slug != "" {
		sb.WriteString(fmt.Sprintf(config.TplJournalPageHeading+nl+nl, s.Slug))
	} else {
		sb.WriteString(fmt.Sprintf(config.TplJournalPageHeading+nl+nl, baseName))
	}

	// Navigation header for multipart sessions
	if totalParts > 1 {
		sb.WriteString(formatPartNavigation(part, totalParts, baseName))
		sb.WriteString(nl + sep + nl + nl)
	}

	// Metadata (use local time) - only on part 1
	if part == 1 {
		localStart := s.StartTime.Local()
		sb.WriteString(fmt.Sprintf(config.MetadataID+" %s"+nl, s.ID))
		sb.WriteString(fmt.Sprintf(
			config.MetadataDate+" %s"+nl, localStart.Format("2006-01-02")),
		)
		sb.WriteString(fmt.Sprintf(
			config.MetadataTime+" %s"+nl, localStart.Format("15:04:05")),
		)
		sb.WriteString(fmt.Sprintf(
			config.MetadataDuration+" %s"+nl, formatDuration(s.Duration)),
		)
		sb.WriteString(fmt.Sprintf(config.MetadataTool+" %s"+nl, s.Tool))
		sb.WriteString(fmt.Sprintf(config.MetadataProject+" %s"+nl, s.Project))
		if s.GitBranch != "" {
			sb.WriteString(fmt.Sprintf(config.MetadataBranch+" %s"+nl, s.GitBranch))
		}
		if s.Model != "" {
			sb.WriteString(fmt.Sprintf(config.MetadataModel+" %s"+nl, s.Model))
		}
		sb.WriteString(nl)

		// Token stats
		sb.WriteString(fmt.Sprintf(config.MetadataTurns+" %d"+nl, s.TurnCount))
		sb.WriteString(fmt.Sprintf(config.TplRecallTokens+nl,
			formatTokens(s.TotalTokens),
			formatTokens(s.TotalTokensIn),
			formatTokens(s.TotalTokensOut)))
		if totalParts > 1 {
			sb.WriteString(fmt.Sprintf(config.MetadataParts+" %d"+nl, totalParts))
		}
		sb.WriteString(nl + sep + nl + nl)

		// Summary section (placeholder for the user to fill in)
		sb.WriteString(config.RecallHeadingSummary + nl + nl)
		sb.WriteString(config.TplRecallSummaryPlaceholder + nl + nl)
		sb.WriteString(sep + nl + nl)

		// Tool usage summary
		tools := s.AllToolUses()
		if len(tools) > 0 {
			sb.WriteString(config.RecallHeadingToolUsage + nl + nl)
			toolCounts := make(map[string]int)
			for _, t := range tools {
				toolCounts[t.Name]++
			}
			for name, count := range toolCounts {
				sb.WriteString(fmt.Sprintf(
					config.TplRecallToolCount+nl, name, count),
				)
			}
			sb.WriteString(nl + sep + nl + nl)
		}
	}

	// Conversation section
	if part == 1 {
		sb.WriteString(config.RecallHeadingConversation + nl + nl)
	} else {
		sb.WriteString(fmt.Sprintf(
			config.TplRecallConversationContinued+nl+nl, part-1),
		)
	}

	for i, msg := range messages {
		msgNum := startMsgIdx + i + 1
		role := config.LabelRoleUser
		if msg.BelongsToAssistant() {
			role = config.LabelRoleAssistant
		} else if len(msg.ToolResults) > 0 && msg.Text == "" {
			role = config.LabelToolOutput
		}

		localTime := msg.Timestamp.Local()
		sb.WriteString(fmt.Sprintf(config.TplRecallTurnHeader+nl+nl,
			msgNum, role, localTime.Format("15:04:05")))

		if msg.Text != "" {
			text := msg.Text
			// Normalize code fences in user messages
			// (users often type "text: ```code")
			if !msg.BelongsToAssistant() {
				text = normalizeCodeFences(text)
			}
			sb.WriteString(text + nl + nl)
		}

		// Tool uses
		for _, t := range msg.ToolUses {
			sb.WriteString(fmt.Sprintf(config.TplRecallToolUse+nl, formatToolUse(t)))
		}

		// Tool results
		for _, tr := range msg.ToolResults {
			if tr.IsError {
				sb.WriteString(config.TplRecallErrorMarker + nl)
			}
			if tr.Content != "" {
				content := stripLineNumbers(tr.Content)
				content, reminders := extractSystemReminders(content)
				fence := fenceForContent(content)
				lines := strings.Count(content, nl)

				if lines > config.RecallDetailsThreshold {
					summary := fmt.Sprintf(config.TplRecallDetailsSummary, lines)
					sb.WriteString(fmt.Sprintf(config.TplRecallDetailsOpen+nl+nl, summary))
					sb.WriteString(fmt.Sprintf(
						config.TplRecallFencedBlock+nl, fence, content, fence),
					)
					sb.WriteString(config.TplRecallDetailsClose + nl)
				} else {
					sb.WriteString(fmt.Sprintf(
						config.TplRecallFencedBlock+nl, fence, content, fence),
					)
				}

				// Render system reminders as Markdown outside the code fence
				for _, reminder := range reminders {
					sb.WriteString(
						fmt.Sprintf(nl+config.LabelBoldReminder+" %s"+nl, reminder),
					)
				}
			}
		}

		if len(msg.ToolUses) > 0 || len(msg.ToolResults) > 0 {
			sb.WriteString(nl)
		}
	}

	// Navigation footer for multipart sessions
	if totalParts > 1 {
		sb.WriteString(nl + sep + nl + nl)
		sb.WriteString(formatPartNavigation(part, totalParts, baseName))
	}

	return sb.String()
}

// formatPartNavigation generates previous/next navigation links for
// multipart sessions.
//
// Parameters:
//   - part: Current part number (1-indexed)
//   - totalParts: Total number of parts
//   - baseName: Base filename without extension
//
// Returns:
//   - string: Formatted navigation line
//     (e.g., "**Part 2 of 3** | [← Previous](...) | [Next →](...)")
func formatPartNavigation(part, totalParts int, baseName string) string {
	var sb strings.Builder
	nl := config.NewlineLF

	sb.WriteString(fmt.Sprintf(config.TplRecallPartOf, part, totalParts))

	if part > 1 || part < totalParts {
		sb.WriteString(config.PipeSeparator)
	}

	// Previous link
	if part > 1 {
		prevFile := baseName + config.ExtMarkdown
		if part > 2 {
			prevFile = fmt.Sprintf(config.TplRecallPartFilename, baseName, part-1)
		}
		sb.WriteString(fmt.Sprintf(config.TplRecallNavPrev, prevFile))
	}

	// Separator between prev and next
	if part > 1 && part < totalParts {
		sb.WriteString(config.PipeSeparator)
	}

	// Next link
	if part < totalParts {
		nextFile := fmt.Sprintf(config.TplRecallPartFilename, baseName, part+1)
		sb.WriteString(fmt.Sprintf(config.TplRecallNavNext, nextFile))
	}

	sb.WriteString(nl)
	return sb.String()
}

// formatDuration formats a duration in a human-readable way.
//
// Parameters:
//   - d: Duration with Minutes() method
//
// Returns:
//   - string: Human-readable duration (e.g., "<1m", "5m", "1h30m")
func formatDuration(d interface{ Minutes() float64 }) string {
	mins := d.Minutes()
	if mins < 1 {
		return "<1m"
	}
	if mins < 60 {
		return fmt.Sprintf("%dm", int(mins))
	}
	hours := int(mins) / 60
	remainMins := int(mins) % 60
	if remainMins == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh%dm", hours, remainMins)
}

// formatTokens formats token counts in a human-readable way.
//
// Parameters:
//   - tokens: Token count to format
//
// Returns:
//   - string: Human-readable count (e.g., "500", "1.5K", "2.3M")
func formatTokens(tokens int) string {
	if tokens < 1000 {
		return fmt.Sprintf("%d", tokens)
	}
	if tokens < 1000000 {
		return fmt.Sprintf("%.1fK", float64(tokens)/1000)
	}
	return fmt.Sprintf("%.1fM", float64(tokens)/1000000)
}

// stripLineNumbers removes Claude Code's line number prefixes from content.
//
// Parameters:
//   - content: Text potentially containing "    1→" style prefixes
//
// Returns:
//   - string: Content with line number prefixes removed
func stripLineNumbers(content string) string {
	return config.RegExLineNumber.ReplaceAllString(content, "")
}

// extractSystemReminders separates system-reminder content from tool output.
//
// Claude Code injects <system-reminder> tags into tool results. This function
// extracts them so they can be rendered as Markdown outside code fences.
//
// Parameters:
//   - content: Tool result content potentially containing system-reminder tags
//
// Returns:
//   - string: Content with system-reminder tags removed
//   - []string: Extracted reminder texts (may be empty)
func extractSystemReminders(content string) (string, []string) {
	matches := config.RegExSystemReminder.FindAllStringSubmatch(content, -1)
	var reminders []string
	for _, m := range matches {
		if len(m) > 1 && m[1] != "" {
			reminders = append(reminders, m[1])
		}
	}
	cleaned := config.RegExSystemReminder.ReplaceAllString(content, "")
	return cleaned, reminders
}

// normalizeCodeFences ensures code fences are on their own lines with proper spacing.
//
// Users often type "text: ```code" without proper line breaks. Markdown requires
// code fences to be on their own lines with blank lines separating them from
// surrounding content.
//
// Parameters:
//   - content: Text that may contain inline code fences
//
// Returns:
//   - string: Content with code fences properly separated by blank lines
func normalizeCodeFences(content string) string {
	// Add newlines before code fences that follow text on the same line
	result := config.RegExCodeFenceInline.ReplaceAllString(content, "$1\n\n$2")
	// Add newlines after code fences that are followed by text on the same line
	result = config.RegExCodeFenceClose.ReplaceAllString(result, "$1\n\n$2")
	return result
}

// formatToolUse formats a tool invocation with its key parameters.
//
// Extracts the most relevant parameter based on tool type (e.g., file path
// for Read/Write, command for Bash, pattern for Grep).
//
// Parameters:
//   - t: Tool use to format
//
// Returns:
//   - string: Formatted string like "Read: /path/to/file" or just tool name
func formatToolUse(t parser.ToolUse) string {
	// Parse the JSON input to extract meaningful parameters
	var input map[string]any
	if err := json.Unmarshal([]byte(t.Input), &input); err != nil {
		return t.Name
	}

	// Extract the most relevant parameter based on tool type
	switch t.Name {
	case "Read":
		if path, ok := input["file_path"].(string); ok {
			return fmt.Sprintf("Read: %s", path)
		}
	case "Write":
		if path, ok := input["file_path"].(string); ok {
			return fmt.Sprintf("Write: %s", path)
		}
	case "Edit":
		if path, ok := input["file_path"].(string); ok {
			return fmt.Sprintf("Edit: %s", path)
		}
	case "Bash":
		if cmd, ok := input["command"].(string); ok {
			// Truncate long commands for readability
			if len(cmd) > 100 {
				cmd = cmd[:100] + "..."
			}
			return fmt.Sprintf("Bash: %s", cmd)
		}
	case "Grep":
		if pattern, ok := input["pattern"].(string); ok {
			return fmt.Sprintf("Grep: %s", pattern)
		}
	case "Glob":
		if pattern, ok := input["pattern"].(string); ok {
			return fmt.Sprintf("Glob: %s", pattern)
		}
	case "WebFetch":
		if url, ok := input["url"].(string); ok {
			return fmt.Sprintf("WebFetch: %s", url)
		}
	case "WebSearch":
		if query, ok := input["query"].(string); ok {
			return fmt.Sprintf("WebSearch: %s", query)
		}
	case "Task":
		if desc, ok := input["description"].(string); ok {
			return fmt.Sprintf("Task: %s", desc)
		}
	}

	// Default: just show the tool name
	return t.Name
}
