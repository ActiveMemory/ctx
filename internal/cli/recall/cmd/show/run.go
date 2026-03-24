//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/cli/recall/core/format"
	"github.com/ActiveMemory/ctx/internal/cli/recall/core/query"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	errSession "github.com/ActiveMemory/ctx/internal/err/session"
	"github.com/ActiveMemory/ctx/internal/write/recall"
)

// Run handles the recall show command.
//
// Displays detailed information about a session including metadata, token
// usage, tool usage summary, and optionally the full conversation.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - args: session ID or slug to show (ignored if latest is true)
//   - latest: if true, show the most recent session
//   - full: if true, show complete conversation instead of preview
//   - allProjects: if true, search sessions from all projects
//
// Returns:
//   - error: non-nil if session not found or scanning fails
func Run(
	cmd *cobra.Command, args []string, latest, full, allProjects bool,
) error {
	sessions, scanErr := query.FindSessions(allProjects)
	if scanErr != nil {
		return errSession.Find(scanErr)
	}

	if len(sessions) == 0 {
		if allProjects {
			return errSession.NoneFound("")
		}
		return errSession.NoneFound(desc.Text(text.DescKeyLabelHintUseAllProjects))
	}

	var session *entity.Session

	switch {
	case latest:
		session = sessions[0]
	case len(args) == 0:
		return errSession.IDRequired()
	default:
		query := strings.ToLower(args[0])
		var matches []*entity.Session
		for _, s := range sessions {
			if strings.HasPrefix(strings.ToLower(s.ID), query) ||
				strings.Contains(strings.ToLower(s.Slug), query) {
				matches = append(matches, s)
			}
		}
		if len(matches) == 0 {
			return errSession.NotFound(args[0])
		}
		if len(matches) > 1 {
			lines := format.SessionMatchLines(matches)
			recall.AmbiguousSessionMatchWithHint(
				cmd, args[0], lines, matches[0].ID[:journal.SessionIDHintLen],
			)
			return errSession.AmbiguousQuery()
		}
		session = matches[0]
	}

	// Print session details.
	recall.SessionMetadata(cmd, recall.SessionInfo{
		Slug:      session.Slug,
		ID:        session.ID,
		Tool:      session.Tool,
		Project:   session.Project,
		Branch:    session.GitBranch,
		Model:     session.Model,
		Started:   session.StartTime.Format(time.DateTimePreciseFormat),
		Duration:  format.Duration(session.Duration),
		Turns:     session.TurnCount,
		Messages:  len(session.Messages),
		TokensIn:  format.Tokens(session.TotalTokensIn),
		TokensOut: format.Tokens(session.TotalTokensOut),
		TokensAll: format.Tokens(session.TotalTokens),
	})

	// Tool usage summary
	tools := session.AllToolUses()
	if len(tools) > 0 {
		toolCounts := make(map[string]int)
		for _, t := range tools {
			toolCounts[t.Name]++
		}

		recall.SectionHeader(cmd, 2, desc.Text(text.DescKeyLabelSectionToolUsage))
		for name, count := range toolCounts {
			recall.ListItem(cmd, desc.Text(text.DescKeyRecallToolCountLine), name, count)
		}
		recall.BlankLine(cmd)
	}

	// Messages
	if full {
		recall.SectionHeader(cmd, 2, desc.Text(text.DescKeyLabelSectionConversation))

		for i, msg := range session.Messages {
			role := desc.Text(text.DescKeyLabelRoleUser)
			if msg.BelongsToAssistant() {
				role = desc.Text(text.DescKeyLabelRoleAssistant)
			} else if len(msg.ToolResults) > 0 && msg.Text == "" {
				role = desc.Text(text.DescKeyLabelToolOutput)
			}

			recall.ConversationTurn(
				cmd, i+1, role, msg.Timestamp.Format(time.Format),
			)

			if msg.Text != "" {
				recall.TextBlock(cmd, msg.Text)
			}

			for _, t := range msg.ToolUses {
				toolInfo := format.ToolUse(t)
				recall.SessionDetail(
					cmd, desc.Text(text.DescKeyLabelInlineTool), toolInfo,
				)
			}

			for _, tr := range msg.ToolResults {
				if tr.IsError {
					recall.Hint(cmd, desc.Text(text.DescKeyLabelInlineError))
				}
				if tr.Content != "" {
					content := format.StripLineNumbers(tr.Content)
					recall.CodeBlock(cmd, content)
				}
			}

			if len(msg.ToolUses) > 0 || len(msg.ToolResults) > 0 {
				recall.BlankLine(cmd)
			}
		}
	} else {
		recall.SectionHeader(
			cmd, 2, desc.Text(text.DescKeyLabelSectionConversationPreview),
		)

		count := 0
		for _, msg := range session.Messages {
			if msg.BelongsToUser() && msg.Text != "" {
				count++
				if count > journal.PreviewMaxTurns {
					recall.MoreTurns(cmd, session.TurnCount-journal.PreviewMaxTurns)
					break
				}
				t := msg.Text
				if len(t) > journal.PreviewMaxTextLen {
					t = t[:journal.PreviewMaxTextLen] + token.Ellipsis
				}
				recall.NumberedItem(cmd, count, t)
			}
		}
		recall.BlankLine(cmd)
		recall.Hint(cmd, desc.Text(text.DescKeyLabelHintUseFull))
	}

	return nil
}
