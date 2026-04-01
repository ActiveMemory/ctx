//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package format

import (
	"encoding/json"
	"fmt"
	"html"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/assets/tpl"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/source/frontmatter"
	"github.com/ActiveMemory/ctx/internal/config/box"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	sharedFmt "github.com/ActiveMemory/ctx/internal/format"
	"github.com/ActiveMemory/ctx/internal/parse"
)

// PartNavigation generates previous/next navigation links for
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
func PartNavigation(part, totalParts int, baseName string) string {
	var sb strings.Builder
	nl := token.NewlineLF

	sb.WriteString(fmt.Sprintf(tpl.RecallPartOf, part, totalParts))

	if part > 1 || part < totalParts {
		sb.WriteString(box.PipeSeparator)
	}

	// Previous link
	if part > 1 {
		prevFile := baseName + file.ExtMarkdown
		if part > 2 {
			prevFile = fmt.Sprintf(tpl.RecallPartFilename, baseName, part-1)
		}
		sb.WriteString(fmt.Sprintf(tpl.RecallNavPrev, prevFile))
	}

	// Separator between prev and next
	if part > 1 && part < totalParts {
		sb.WriteString(box.PipeSeparator)
	}

	// Next link
	if part < totalParts {
		nextFile := fmt.Sprintf(tpl.RecallPartFilename, baseName, part+1)
		sb.WriteString(fmt.Sprintf(tpl.RecallNavNext, nextFile))
	}

	sb.WriteString(nl)
	return sb.String()
}

// Duration formats a duration in a human-readable way.
//
// Parameters:
//   - d: Duration with Minutes() method
//
// Returns:
//   - string: Human-readable duration (e.g., "<1m", "5m", "1h30m")
func Duration(d interface{ Minutes() float64 }) string {
	mins := d.Minutes()
	if mins < 1 {
		return desc.Text(text.DescKeyWriteFormatDurationLTMin)
	}
	if mins < 60 {
		return fmt.Sprintf(desc.Text(text.DescKeyWriteFormatDurationMin), int(mins))
	}
	hours := int(mins) / 60
	remainMins := int(mins) % 60
	if remainMins == 0 {
		return fmt.Sprintf(desc.Text(text.DescKeyWriteFormatDurationHour), hours)
	}
	return fmt.Sprintf(
		desc.Text(text.DescKeyWriteFormatDurationHourMin),
		hours, remainMins,
	)
}

// ToolUse formats a tool invocation with its key parameters.
//
// Parameters:
//   - t: Tool use to format
//
// Returns:
//   - string: Formatted string like "Read: /path/to/file" or just tool name
func ToolUse(t entity.ToolUse) string {
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
	if t.Name == session.ToolBash && len(val) > session.ToolDisplayMaxLen {
		val = val[:session.ToolDisplayMaxLen] + token.Ellipsis
	}
	return fmt.Sprintf(tpl.ToolDisplay, t.Name, val)
}

// SessionMatchLines formats session matches for ambiguous query output.
//
// Parameters:
//   - matches: sessions that matched the query.
//
// Returns:
//   - []string: pre-formatted lines, one per match.
func SessionMatchLines(matches []*entity.Session) []string {
	lines := make([]string, 0, len(matches))
	for _, m := range matches {
		lines = append(lines, fmt.Sprintf(
			tpl.SessionMatch,
			m.Slug,
			m.ID[:journal.SessionIDShortLen],
			m.StartTime.Format(time.DateTimeFmt)),
		)
	}
	return lines
}

// JournalFilename generates the filename for a journal entry.
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
func JournalFilename(s *entity.Session, slugOverride string) string {
	date := s.StartTime.Local().Format(time.DateFormat)
	shortID := s.ID
	if len(shortID) > journal.ShortIDLen {
		shortID = shortID[:journal.ShortIDLen]
	}
	slug := s.Slug
	if slugOverride != "" {
		slug = slugOverride
	}
	return fmt.Sprintf(tpl.RecallFilename, date, slug, shortID)
}

// JournalEntryPart generates Markdown content for a part of a journal entry.
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
func JournalEntryPart(
	s *entity.Session,
	messages []entity.Message,
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
		durationStr := Duration(s.Duration)

		// Basic YAML frontmatter
		sb.WriteString(sep + nl)
		frontmatter.WriteFmQuoted(&sb, session.FrontmatterDate, dateStr)
		frontmatter.WriteFmQuoted(&sb, session.FmKeyTime, timeStr)
		frontmatter.WriteFmString(&sb, session.FmKeyProject, s.Project)
		if s.GitBranch != "" {
			frontmatter.WriteFmString(&sb, session.FmKeyBranch, s.GitBranch)
		}
		if s.Model != "" {
			frontmatter.WriteFmString(&sb, session.FmKeyModel, s.Model)
		}
		if s.TotalTokensIn > 0 {
			frontmatter.WriteFmInt(&sb, session.FmKeyTokensIn, s.TotalTokensIn)
		}
		if s.TotalTokensOut > 0 {
			frontmatter.WriteFmInt(&sb, session.FmKeyTokensOut, s.TotalTokensOut)
		}
		frontmatter.WriteFmQuoted(&sb, session.FmKeyID, s.ID)
		if title != "" {
			frontmatter.WriteFmQuoted(&sb, session.FrontmatterTitle, title)
		}
		sb.WriteString(sep + nl + nl)

		// Header: prefer title, fall back to slug, then baseName.
		heading := frontmatter.ResolveHeading(title, s.Slug, baseName)
		sb.WriteString(fmt.Sprintf(tpl.JournalPageHeading+nl+nl, heading))

		// Navigation header for multipart sessions
		if totalParts > 1 {
			sb.WriteString(PartNavigation(part, totalParts, baseName))
			sb.WriteString(nl + sep + nl + nl)
		}

		// Session metadata as collapsible HTML table
		// (Markdown tables don't render inside <details> in Zensical)
		summaryText := fmt.Sprintf(
			desc.Text(text.DescKeyJournalSourceMetaSummary), dateStr, durationStr, s.Model,
		)
		sb.WriteString(fmt.Sprintf(tpl.MetaDetailsOpen, summaryText))
		sb.WriteString(fmt.Sprintf(
			tpl.MetaRow+nl, desc.Text(text.DescKeyLabelMetaID), s.ID))
		sb.WriteString(
			fmt.Sprintf(
				tpl.MetaRow+nl, desc.Text(text.DescKeyLabelMetaDate), dateStr,
			),
		)
		sb.WriteString(
			fmt.Sprintf(
				tpl.MetaRow+nl, desc.Text(text.DescKeyLabelMetaTime), timeStr,
			),
		)
		sb.WriteString(
			fmt.Sprintf(
				tpl.MetaRow+nl, desc.Text(text.DescKeyLabelMetaDuration), durationStr,
			),
		)
		sb.WriteString(
			fmt.Sprintf(
				tpl.MetaRow+nl, desc.Text(text.DescKeyLabelMetaTool), s.Tool,
			),
		)
		sb.WriteString(
			fmt.Sprintf(
				tpl.MetaRow+nl, desc.Text(text.DescKeyLabelMetaProject), s.Project,
			),
		)
		if s.GitBranch != "" {
			sb.WriteString(
				fmt.Sprintf(
					tpl.MetaRow+nl, desc.Text(text.DescKeyLabelMetaBranch), s.GitBranch,
				),
			)
		}
		if s.Model != "" {
			sb.WriteString(
				fmt.Sprintf(
					tpl.MetaRow+nl, desc.Text(text.DescKeyLabelMetaModel), s.Model,
				),
			)
		}
		sb.WriteString(tpl.MetaDetailsClose + nl + nl)

		// Token stats as collapsible HTML table
		turnStr := fmt.Sprintf("%d", s.TurnCount)
		sb.WriteString(fmt.Sprintf(tpl.MetaDetailsOpen, turnStr))
		sb.WriteString(
			fmt.Sprintf(
				tpl.MetaRow+nl, desc.Text(text.DescKeyLabelMetaTurns), turnStr,
			),
		)
		tokenSummary := fmt.Sprintf(desc.Text(text.DescKeyJournalSourceTokenSummary),
			sharedFmt.Tokens(s.TotalTokens),
			sharedFmt.Tokens(s.TotalTokensIn),
			sharedFmt.Tokens(s.TotalTokensOut))
		sb.WriteString(
			fmt.Sprintf(
				tpl.MetaRow+nl, desc.Text(text.DescKeyLabelMetaTokens), tokenSummary,
			),
		)
		if totalParts > 1 {
			sb.WriteString(
				fmt.Sprintf(tpl.MetaRow+nl, desc.Text(text.DescKeyLabelMetaParts),
					fmt.Sprintf("%d", totalParts)),
			)
		}
		sb.WriteString(tpl.MetaDetailsClose + nl + nl)

		sb.WriteString(sep + nl + nl)

		// Tool usage summary
		tools := s.AllToolUses()
		if len(tools) > 0 {
			sb.WriteString(desc.Text(text.DescKeyHeadingToolUsage) + nl + nl)
			toolCounts := make(map[string]int)
			for _, t := range tools {
				toolCounts[t.Name]++
			}
			for name, count := range toolCounts {
				sb.WriteString(fmt.Sprintf(
					tpl.RecallToolCount+nl, name, count),
				)
			}
			sb.WriteString(nl + sep + nl + nl)
		}
	} else {
		// Header (non-part-1) - the same fallback as part 1.
		heading := frontmatter.ResolveHeading(title, s.Slug, baseName)
		sb.WriteString(fmt.Sprintf(tpl.JournalPageHeading+nl+nl, heading))

		// Navigation header for multipart sessions
		if totalParts > 1 {
			sb.WriteString(PartNavigation(part, totalParts, baseName))
			sb.WriteString(nl + sep + nl + nl)
		}
	}

	// Conversation section
	if part == 1 {
		sb.WriteString(desc.Text(text.DescKeyHeadingConversation) + nl + nl)
	} else {
		sb.WriteString(fmt.Sprintf(
			tpl.RecallConversationContinued+nl+nl, part-1),
		)
	}

	for i, msg := range messages {
		msgNum := startMsgIdx + i + 1
		role := desc.Text(text.DescKeyLabelRoleUser)
		if msg.BelongsToAssistant() {
			role = desc.Text(text.DescKeyLabelRoleAssistant)
		} else if len(msg.ToolResults) > 0 && msg.Text == "" {
			role = desc.Text(text.DescKeyLabelToolOutput)
		}

		localTime := msg.Timestamp.Local()
		sb.WriteString(fmt.Sprintf(tpl.RecallTurnHeader+nl+nl,
			msgNum, role, localTime.Format(time.Format)))

		if msg.Text != "" {
			t := msg.Text
			// Normalize code fences in user messages
			// (users often type "text: ```code")
			if !msg.BelongsToAssistant() {
				t = parse.NormalizeCodeFences(t)
			}
			sb.WriteString(t + nl + nl)
		}

		// Tool uses
		for _, t := range msg.ToolUses {
			sb.WriteString(fmt.Sprintf(tpl.RecallToolUse+nl, ToolUse(t)))
		}

		// Tool results
		for _, tr := range msg.ToolResults {
			if tr.IsError {
				sb.WriteString(tpl.RecallErrorMarker + nl)
			}
			if tr.Content != "" {
				stripped := parse.StripLineNumbers(tr.Content)
				content, reminders := parse.ExtractSystemReminders(stripped)
				fence := parse.FenceForContent(content)
				lines := strings.Count(content, nl)

				if lines > journal.DetailsThreshold {
					summary := fmt.Sprintf(tpl.RecallDetailsSummary, lines)
					sb.WriteString(fmt.Sprintf(tpl.RecallDetailsOpen+nl+nl, summary))
					sb.WriteString(marker.TagPre + nl)
					sb.WriteString(html.EscapeString(content) + nl)
					sb.WriteString(marker.TagPreClose + nl)
					sb.WriteString(tpl.RecallDetailsClose + nl)
				} else {
					sb.WriteString(fmt.Sprintf(
						tpl.RecallFencedBlock+nl, fence, content, fence),
					)
				}

				// Render system reminders as Markdown outside the code fence
				for _, reminder := range reminders {
					sb.WriteString(
						fmt.Sprintf(
							nl+desc.Text(text.DescKeyLabelBoldReminder)+" %s"+nl,
							reminder,
						),
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
		sb.WriteString(PartNavigation(part, totalParts, baseName))
	}

	return sb.String()
}
