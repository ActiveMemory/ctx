//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"encoding/json"
	"fmt"
	"html"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/box"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/recall/parser"
)

// FenceForContent returns the appropriate code fence for content.
//
// Uses longer fences when content contains backticks to avoid
// nested Markdown rendering issues. Starts with ``` and adds
// more backticks as needed.
//
// Parameters:
//   - content: The content to be fenced
//
// Returns:
//   - string: A fence string (e.g., "```", "````", "```"")
func FenceForContent(content string) string {
	fence := token.CodeFence
	for strings.Contains(content, fence) {
		fence += token.Backtick
	}
	return fence
}

// FormatJournalFilename generates the filename for a journal entry.
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
func FormatJournalFilename(s *parser.Session, slugOverride string) string {
	date := s.StartTime.Local().Format(time.DateFormat)
	shortID := s.ID
	if len(shortID) > journal.ShortIDLen {
		shortID = shortID[:journal.ShortIDLen]
	}
	slug := s.Slug
	if slugOverride != "" {
		slug = slugOverride
	}
	return fmt.Sprintf(assets.TplRecallFilename, date, slug, shortID)
}

// FormatJournalEntryPart generates Markdown content for a part of a journal entry.
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
func FormatJournalEntryPart(
	s *parser.Session,
	messages []parser.Message,
	startMsgIdx, part, totalParts int,
	baseName, title string,
) string {
	var sb strings.Builder
	nl := token.NewlineLF
	sep := token.Separator

	// Metadata (YAML frontmatter + HTML details) - only on part 1
	if part == 1 {
		localStart := s.StartTime.Local()
		dateStr := localStart.Format(time.DateFormat)
		timeStr := localStart.Format(time.Format)
		durationStr := FormatDuration(s.Duration)

		// Basic YAML frontmatter
		sb.WriteString(sep + nl)
		writeFmQuoted(&sb, assets.FmKeyDate, dateStr)
		writeFmQuoted(&sb, assets.FmKeyTime, timeStr)
		writeFmString(&sb, assets.FmKeyProject, s.Project)
		if s.GitBranch != "" {
			writeFmString(&sb, assets.FmKeyBranch, s.GitBranch)
		}
		if s.Model != "" {
			writeFmString(&sb, assets.FmKeyModel, s.Model)
		}
		if s.TotalTokensIn > 0 {
			writeFmInt(&sb, assets.FmKeyTokensIn, s.TotalTokensIn)
		}
		if s.TotalTokensOut > 0 {
			writeFmInt(&sb, assets.FmKeyTokensOut, s.TotalTokensOut)
		}
		writeFmQuoted(&sb, assets.FmKeySessionID, s.ID)
		if title != "" {
			writeFmQuoted(&sb, assets.FmKeyTitle, title)
		}
		sb.WriteString(sep + nl + nl)

		// Header — prefer title, fall back to slug, then baseName.
		heading := resolveHeading(title, s.Slug, baseName)
		sb.WriteString(fmt.Sprintf(assets.TplJournalPageHeading+nl+nl, heading))

		// Navigation header for multipart sessions
		if totalParts > 1 {
			sb.WriteString(FormatPartNavigation(part, totalParts, baseName))
			sb.WriteString(nl + sep + nl + nl)
		}

		// Session metadata as collapsible HTML table
		// (Markdown tables don't render inside <details> in Zensical)
		summaryText := fmt.Sprintf("%s · %s · %s", dateStr, durationStr, s.Model)
		sb.WriteString(fmt.Sprintf(assets.TplMetaDetailsOpen, summaryText))
		sb.WriteString(fmt.Sprintf(assets.TplMetaRow+nl, assets.MetaLabelID, s.ID))
		sb.WriteString(fmt.Sprintf(assets.TplMetaRow+nl, assets.MetaLabelDate, dateStr))
		sb.WriteString(fmt.Sprintf(assets.TplMetaRow+nl, assets.MetaLabelTime, timeStr))
		sb.WriteString(fmt.Sprintf(assets.TplMetaRow+nl, assets.MetaLabelDuration, durationStr))
		sb.WriteString(fmt.Sprintf(assets.TplMetaRow+nl, assets.MetaLabelTool, s.Tool))
		sb.WriteString(fmt.Sprintf(assets.TplMetaRow+nl, assets.MetaLabelProject, s.Project))
		if s.GitBranch != "" {
			sb.WriteString(fmt.Sprintf(assets.TplMetaRow+nl, assets.MetaLabelBranch, s.GitBranch))
		}
		if s.Model != "" {
			sb.WriteString(fmt.Sprintf(assets.TplMetaRow+nl, assets.MetaLabelModel, s.Model))
		}
		sb.WriteString(assets.TplMetaDetailsClose + nl + nl)

		// Token stats as collapsible HTML table
		turnStr := fmt.Sprintf("%d", s.TurnCount)
		sb.WriteString(fmt.Sprintf(assets.TplMetaDetailsOpen, turnStr))
		sb.WriteString(fmt.Sprintf(assets.TplMetaRow+nl, assets.MetaLabelTurns, turnStr))
		tokenSummary := fmt.Sprintf("%s (in: %s, out: %s)",
			FormatTokens(s.TotalTokens),
			FormatTokens(s.TotalTokensIn),
			FormatTokens(s.TotalTokensOut))
		sb.WriteString(fmt.Sprintf(assets.TplMetaRow+nl, assets.MetaLabelTokens, tokenSummary))
		if totalParts > 1 {
			sb.WriteString(fmt.Sprintf(assets.TplMetaRow+nl, assets.MetaLabelParts,
				fmt.Sprintf("%d", totalParts)))
		}
		sb.WriteString(assets.TplMetaDetailsClose + nl + nl)

		sb.WriteString(sep + nl + nl)

		// Tool usage summary
		tools := s.AllToolUses()
		if len(tools) > 0 {
			sb.WriteString(assets.RecallHeadingToolUsage + nl + nl)
			toolCounts := make(map[string]int)
			for _, t := range tools {
				toolCounts[t.Name]++
			}
			for name, count := range toolCounts {
				sb.WriteString(fmt.Sprintf(
					assets.TplRecallToolCount+nl, name, count),
				)
			}
			sb.WriteString(nl + sep + nl + nl)
		}
	} else {
		// Header (non-part-1) — same fallback as part 1.
		heading := resolveHeading(title, s.Slug, baseName)
		sb.WriteString(fmt.Sprintf(assets.TplJournalPageHeading+nl+nl, heading))

		// Navigation header for multipart sessions
		if totalParts > 1 {
			sb.WriteString(FormatPartNavigation(part, totalParts, baseName))
			sb.WriteString(nl + sep + nl + nl)
		}
	}

	// Conversation section
	if part == 1 {
		sb.WriteString(assets.RecallHeadingConversation + nl + nl)
	} else {
		sb.WriteString(fmt.Sprintf(
			assets.TplRecallConversationContinued+nl+nl, part-1),
		)
	}

	for i, msg := range messages {
		msgNum := startMsgIdx + i + 1
		role := assets.RoleUser
		if msg.BelongsToAssistant() {
			role = assets.LabelRoleAssistant
		} else if len(msg.ToolResults) > 0 && msg.Text == "" {
			role = assets.ToolOutput
		}

		localTime := msg.Timestamp.Local()
		sb.WriteString(fmt.Sprintf(assets.TplRecallTurnHeader+nl+nl,
			msgNum, role, localTime.Format(time.Format)))

		if msg.Text != "" {
			text := msg.Text
			// Normalize code fences in user messages
			// (users often type "text: ```code")
			if !msg.BelongsToAssistant() {
				text = NormalizeCodeFences(text)
			}
			sb.WriteString(text + nl + nl)
		}

		// Tool uses
		for _, t := range msg.ToolUses {
			sb.WriteString(fmt.Sprintf(assets.TplRecallToolUse+nl, FormatToolUse(t)))
		}

		// Tool results
		for _, tr := range msg.ToolResults {
			if tr.IsError {
				sb.WriteString(assets.TplRecallErrorMarker + nl)
			}
			if tr.Content != "" {
				content := StripLineNumbers(tr.Content)
				content, reminders := ExtractSystemReminders(content)
				fence := FenceForContent(content)
				lines := strings.Count(content, nl)

				if lines > journal.DetailsThreshold {
					summary := fmt.Sprintf(assets.TplRecallDetailsSummary, lines)
					sb.WriteString(fmt.Sprintf(assets.TplRecallDetailsOpen+nl+nl, summary))
					sb.WriteString("<pre>" + nl + html.EscapeString(content) + nl + "</pre>" + nl)
					sb.WriteString(assets.TplRecallDetailsClose + nl)
				} else {
					sb.WriteString(fmt.Sprintf(
						assets.TplRecallFencedBlock+nl, fence, content, fence),
					)
				}

				// Render system reminders as Markdown outside the code fence
				for _, reminder := range reminders {
					sb.WriteString(
						fmt.Sprintf(nl+assets.BoldReminder+" %s"+nl, reminder),
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
		sb.WriteString(FormatPartNavigation(part, totalParts, baseName))
	}

	return sb.String()
}

// resolveHeading returns the first non-empty value among title, slug, baseName.
func resolveHeading(title, slug, baseName string) string {
	if title != "" {
		return title
	}
	if slug != "" {
		return slug
	}
	return baseName
}

// writeFmQuoted writes a YAML frontmatter quoted string field.
func writeFmQuoted(sb *strings.Builder, key, value string) {
	fmt.Fprintf(sb, assets.TplFmQuoted+token.NewlineLF, key, value)
}

// writeFmString writes a YAML frontmatter bare string field.
func writeFmString(sb *strings.Builder, key, value string) {
	fmt.Fprintf(sb, assets.TplFmString+token.NewlineLF, key, value)
}

// writeFmInt writes a YAML frontmatter integer field.
func writeFmInt(sb *strings.Builder, key string, value int) {
	fmt.Fprintf(sb, assets.TplFmInt+token.NewlineLF, key, value)
}

// FormatPartNavigation generates previous/next navigation links for
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
func FormatPartNavigation(part, totalParts int, baseName string) string {
	var sb strings.Builder
	nl := token.NewlineLF

	sb.WriteString(fmt.Sprintf(assets.TplRecallPartOf, part, totalParts))

	if part > 1 || part < totalParts {
		sb.WriteString(box.PipeSeparator)
	}

	// Previous link
	if part > 1 {
		prevFile := baseName + file.ExtMarkdown
		if part > 2 {
			prevFile = fmt.Sprintf(assets.TplRecallPartFilename, baseName, part-1)
		}
		sb.WriteString(fmt.Sprintf(assets.TplRecallNavPrev, prevFile))
	}

	// Separator between prev and next
	if part > 1 && part < totalParts {
		sb.WriteString(box.PipeSeparator)
	}

	// Next link
	if part < totalParts {
		nextFile := fmt.Sprintf(assets.TplRecallPartFilename, baseName, part+1)
		sb.WriteString(fmt.Sprintf(assets.TplRecallNavNext, nextFile))
	}

	sb.WriteString(nl)
	return sb.String()
}

// FormatDuration formats a duration in a human-readable way.
//
// Parameters:
//   - d: Duration with Minutes() method
//
// Returns:
//   - string: Human-readable duration (e.g., "<1m", "5m", "1h30m")
func FormatDuration(d interface{ Minutes() float64 }) string {
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

// FormatTokens formats token counts in a human-readable way.
//
// Parameters:
//   - tokens: Token count to format
//
// Returns:
//   - string: Human-readable count (e.g., "500", "1.5K", "2.3M")
func FormatTokens(tokens int) string {
	if tokens < 1000 {
		return fmt.Sprintf("%d", tokens)
	}
	if tokens < 1000000 {
		return fmt.Sprintf("%.1fK", float64(tokens)/1000)
	}
	return fmt.Sprintf("%.1fM", float64(tokens)/1000000)
}

// Truncate shortens s to max characters, appending "…" if truncated.
//
// Parameters:
//   - s: String to truncate
//   - max: Maximum length
//
// Returns:
//   - string: Truncated string
func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

// StripLineNumbers removes Claude Code's line number prefixes from content.
//
// Parameters:
//   - content: Text potentially containing "    1→" style prefixes
//
// Returns:
//   - string: Content with line number prefixes removed
func StripLineNumbers(content string) string {
	return regex.LineNumber.ReplaceAllString(content, "")
}

// ExtractSystemReminders separates system-reminder content from tool output.
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
func ExtractSystemReminders(content string) (string, []string) {
	matches := regex.SystemReminder.FindAllStringSubmatch(content, -1)
	var reminders []string
	for _, m := range matches {
		if len(m) > 1 && m[1] != "" {
			reminders = append(reminders, m[1])
		}
	}
	cleaned := regex.SystemReminder.ReplaceAllString(content, "")
	return cleaned, reminders
}

// NormalizeCodeFences ensures code fences are on their own lines with proper spacing.
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
func NormalizeCodeFences(content string) string {
	// Add newlines before code fences that follow text on the same line
	result := regex.CodeFenceInline.ReplaceAllString(content, "$1\n\n$2")
	// Add newlines after code fences that are followed by text on the same line
	result = regex.CodeFenceClose.ReplaceAllString(result, "$1\n\n$2")
	return result
}

// toolDisplayKey maps tool names to the JSON input key that best
// describes each invocation.
var toolDisplayKey = map[string]string{
	assets.ToolRead:      assets.ToolInputFilePath,
	assets.ToolWrite:     assets.ToolInputFilePath,
	assets.ToolEdit:      assets.ToolInputFilePath,
	assets.ToolBash:      assets.ToolInputCommand,
	assets.ToolGrep:      assets.ToolInputPattern,
	assets.ToolGlob:      assets.ToolInputPattern,
	assets.ToolWebFetch:  assets.ToolInputURL,
	assets.ToolWebSearch: assets.ToolInputQuery,
	assets.ToolTask:      assets.ToolInputDescription,
}

// FormatToolUse formats a tool invocation with its key parameters.
//
// Parameters:
//   - t: Tool use to format
//
// Returns:
//   - string: Formatted string like "Read: /path/to/file" or just tool name
func FormatToolUse(t parser.ToolUse) string {
	key, ok := toolDisplayKey[t.Name]
	if !ok {
		return t.Name
	}
	var input map[string]any
	if jsonErr := json.Unmarshal([]byte(t.Input), &input); jsonErr != nil {
		return t.Name
	}
	val, ok := input[key].(string)
	if !ok {
		return t.Name
	}
	if t.Name == assets.ToolBash && len(val) > assets.ToolDisplayMaxLen {
		val = val[:assets.ToolDisplayMaxLen] + token.Ellipsis
	}
	return fmt.Sprintf(assets.TplToolDisplay, t.Name, val)
}
