//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"encoding/json"
	"fmt"
	"html"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/recall/parser"
)

// Claude Code tool names used in session transcripts.
const (
	toolRead      = "Read"
	toolWrite     = "Write"
	toolEdit      = "Edit"
	toolBash      = "Bash"
	toolGrep      = "Grep"
	toolGlob      = "Glob"
	toolWebFetch  = "WebFetch"
	toolWebSearch = "WebSearch"
	toolTask      = "Task"
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
//   - string: A fence string (e.g., "```", "````", "```“")
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
// When slugOverride is non-empty it replaces s.Slug in the filename,
// allowing title-derived slugs to be used instead of Claude Code's
// random slug.
//
// Parameters:
//   - s: Session to generate filename for
//   - slugOverride: If non-empty, used instead of s.Slug
//
// Returns:
//   - string: Filename like "2026-01-15-fix-auth-bug-abc12345.md"
func formatJournalFilename(s *parser.Session, slugOverride string) string {
	date := s.StartTime.Local().Format("2006-01-02")
	shortID := s.ID
	if len(shortID) > config.RecallShortIDLen {
		shortID = shortID[:config.RecallShortIDLen]
	}
	slug := s.Slug
	if slugOverride != "" {
		slug = slugOverride
	}
	return fmt.Sprintf(config.TplRecallFilename, date, slug, shortID)
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
//   - title: Human-readable title for frontmatter and H1 heading (may be empty)
//
// Returns:
//   - string: Markdown content for this part
func formatJournalEntryPart(
	s *parser.Session,
	messages []parser.Message,
	startMsgIdx, part, totalParts int,
	baseName, title string,
) string {
	var sb strings.Builder
	nl := config.NewlineLF
	sep := config.Separator

	// Metadata (YAML frontmatter + HTML details) - only on part 1
	if part == 1 {
		localStart := s.StartTime.Local()
		dateStr := localStart.Format("2006-01-02")
		timeStr := localStart.Format("15:04:05")
		durationStr := formatDuration(s.Duration)

		// Basic YAML frontmatter
		sb.WriteString(sep + nl)
		sb.WriteString(fmt.Sprintf("date: %q"+nl, dateStr))
		sb.WriteString(fmt.Sprintf("time: %q"+nl, timeStr))
		sb.WriteString(fmt.Sprintf("project: %s"+nl, s.Project))
		if s.GitBranch != "" {
			sb.WriteString(fmt.Sprintf("branch: %s"+nl, s.GitBranch))
		}
		if s.Model != "" {
			sb.WriteString(fmt.Sprintf("model: %s"+nl, s.Model))
		}
		sb.WriteString(fmt.Sprintf("session_id: %q"+nl, s.ID))
		if title != "" {
			sb.WriteString(fmt.Sprintf("title: %q"+nl, title))
		}
		sb.WriteString(sep + nl + nl)

		// Header — prefer title, fall back to slug, then baseName.
		heading := title
		if heading == "" {
			heading = s.Slug
		}
		if heading == "" {
			heading = baseName
		}
		sb.WriteString(fmt.Sprintf(config.TplJournalPageHeading+nl+nl, heading))

		// Navigation header for multipart sessions
		if totalParts > 1 {
			sb.WriteString(formatPartNavigation(part, totalParts, baseName))
			sb.WriteString(nl + sep + nl + nl)
		}

		// Session metadata as collapsible HTML table
		// (Markdown tables don't render inside <details> in Zensical)
		summaryText := fmt.Sprintf("%s · %s · %s", dateStr, durationStr, s.Model)
		sb.WriteString(fmt.Sprintf(config.TplMetaDetailsOpen, summaryText))
		sb.WriteString(fmt.Sprintf(config.TplMetaRow+nl, "ID", s.ID))
		sb.WriteString(fmt.Sprintf(config.TplMetaRow+nl, "Date", dateStr))
		sb.WriteString(fmt.Sprintf(config.TplMetaRow+nl, "Time", timeStr))
		sb.WriteString(fmt.Sprintf(config.TplMetaRow+nl, "Duration", durationStr))
		sb.WriteString(fmt.Sprintf(config.TplMetaRow+nl, "Tool", s.Tool))
		sb.WriteString(fmt.Sprintf(config.TplMetaRow+nl, "Project", s.Project))
		if s.GitBranch != "" {
			sb.WriteString(fmt.Sprintf(config.TplMetaRow+nl, "Branch", s.GitBranch))
		}
		if s.Model != "" {
			sb.WriteString(fmt.Sprintf(config.TplMetaRow+nl, "Model", s.Model))
		}
		sb.WriteString(config.TplMetaDetailsClose + nl + nl)

		// Token stats as collapsible HTML table
		turnStr := fmt.Sprintf("%d", s.TurnCount)
		sb.WriteString(fmt.Sprintf(config.TplMetaDetailsOpen, turnStr))
		sb.WriteString(fmt.Sprintf(config.TplMetaRow+nl, "Turns", turnStr))
		tokenSummary := fmt.Sprintf("%s (in: %s, out: %s)",
			formatTokens(s.TotalTokens),
			formatTokens(s.TotalTokensIn),
			formatTokens(s.TotalTokensOut))
		sb.WriteString(fmt.Sprintf(config.TplMetaRow+nl, "Tokens", tokenSummary))
		if totalParts > 1 {
			sb.WriteString(fmt.Sprintf(config.TplMetaRow+nl, "Parts",
				fmt.Sprintf("%d", totalParts)))
		}
		sb.WriteString(config.TplMetaDetailsClose + nl + nl)

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
	} else {
		// Header (non-part-1) — same fallback as part 1.
		heading := title
		if heading == "" {
			heading = s.Slug
		}
		if heading == "" {
			heading = baseName
		}
		sb.WriteString(fmt.Sprintf(config.TplJournalPageHeading+nl+nl, heading))

		// Navigation header for multipart sessions
		if totalParts > 1 {
			sb.WriteString(formatPartNavigation(part, totalParts, baseName))
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
					sb.WriteString("<pre>" + nl + html.EscapeString(content) + nl + "</pre>" + nl)
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

// truncate shortens s to max characters, appending "…" if truncated.
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
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
//
// toolDisplayKey maps tool names to the JSON input key that best
// describes each invocation.
var toolDisplayKey = map[string]string{
	toolRead:      "file_path",
	toolWrite:     "file_path",
	toolEdit:      "file_path",
	toolBash:      "command",
	toolGrep:      "pattern",
	toolGlob:      "pattern",
	toolWebFetch:  "url",
	toolWebSearch: "query",
	toolTask:      "description",
}

func formatToolUse(t parser.ToolUse) string {
	key, ok := toolDisplayKey[t.Name]
	if !ok {
		return t.Name
	}
	var input map[string]any
	if err := json.Unmarshal([]byte(t.Input), &input); err != nil {
		return t.Name
	}
	val, ok := input[key].(string)
	if !ok {
		return t.Name
	}
	if t.Name == toolBash && len(val) > 100 {
		val = val[:100] + "..."
	}
	return fmt.Sprintf("%s: %s", t.Name, val)
}
