//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
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
