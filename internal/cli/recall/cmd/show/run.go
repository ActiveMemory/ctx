//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/recall/core"
	"github.com/ActiveMemory/ctx/internal/config"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/recall/parser"
	"github.com/ActiveMemory/ctx/internal/write"
)

// runShow handles the recall show command.
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
func runShow(
	cmd *cobra.Command, args []string, latest, full, allProjects bool,
) error {
	sessions, scanErr := core.FindSessions(allProjects)
	if scanErr != nil {
		return ctxerr.FindSessions(scanErr)
	}

	if len(sessions) == 0 {
		if allProjects {
			return ctxerr.NoSessionsFound("")
		}
		return ctxerr.NoSessionsFound(config.HintUseAllProjects)
	}

	var session *parser.Session

	switch {
	case latest:
		session = sessions[0]
	case len(args) == 0:
		return ctxerr.SessionIDRequired()
	default:
		query := strings.ToLower(args[0])
		var matches []*parser.Session
		for _, s := range sessions {
			if strings.HasPrefix(strings.ToLower(s.ID), query) ||
				strings.Contains(strings.ToLower(s.Slug), query) {
				matches = append(matches, s)
			}
		}
		if len(matches) == 0 {
			return ctxerr.SessionNotFound(args[0])
		}
		if len(matches) > 1 {
			lines := core.FormatSessionMatchLines(matches)
			write.AmbiguousSessionMatchWithHint(
				cmd, args[0], lines, matches[0].ID[:config.SessionIDHintLen],
			)
			return ctxerr.AmbiguousQuery()
		}
		session = matches[0]
	}

	// Print session details
	write.SectionHeader(cmd, 1, session.Slug)

	write.SessionDetail(cmd, config.MetadataID, session.ID)
	write.SessionDetail(cmd, config.MetadataTool, session.Tool)
	write.SessionDetail(cmd, config.MetadataProject, session.Project)
	if session.GitBranch != "" {
		write.SessionDetail(cmd, config.MetadataBranch, session.GitBranch)
	}
	if session.Model != "" {
		write.SessionDetail(cmd, config.MetadataModel, session.Model)
	}
	write.BlankLine(cmd)

	write.SessionDetail(
		cmd, config.MetadataStarted,
		session.StartTime.Format(config.DateTimePreciseFormat),
	)
	write.SessionDetail(
		cmd, config.MetadataDuration, core.FormatDuration(session.Duration),
	)
	write.SessionDetailInt(cmd, config.MetadataTurns, session.TurnCount)
	write.SessionDetailInt(cmd, config.MetadataMessages, len(session.Messages))
	write.BlankLine(cmd)

	write.SessionDetail(
		cmd, config.MetadataInputUsage, core.FormatTokens(session.TotalTokensIn),
	)
	write.SessionDetail(
		cmd, config.MetadataOutputUsage, core.FormatTokens(session.TotalTokensOut),
	)
	write.SessionDetail(
		cmd, config.MetadataTotal, core.FormatTokens(session.TotalTokens),
	)
	write.BlankLine(cmd)

	// Tool usage summary
	tools := session.AllToolUses()
	if len(tools) > 0 {
		toolCounts := make(map[string]int)
		for _, t := range tools {
			toolCounts[t.Name]++
		}

		write.SectionHeader(cmd, 2, config.SectionToolUsage)
		for name, count := range toolCounts {
			write.ListItem(cmd, "%s: %d", name, count)
		}
		write.BlankLine(cmd)
	}

	// Messages
	if full {
		write.SectionHeader(cmd, 2, config.SectionConversation)

		for i, msg := range session.Messages {
			role := config.LabelRoleUser
			if msg.BelongsToAssistant() {
				role = config.LabelRoleAssistant
			} else if len(msg.ToolResults) > 0 && msg.Text == "" {
				role = config.LabelToolOutput
			}

			write.ConversationTurn(
				cmd, i+1, role, msg.Timestamp.Format(config.TimeFormat),
			)

			if msg.Text != "" {
				write.TextBlock(cmd, msg.Text)
			}

			for _, t := range msg.ToolUses {
				toolInfo := core.FormatToolUse(t)
				write.SessionDetail(cmd, config.LabelTool, toolInfo)
			}

			for _, tr := range msg.ToolResults {
				if tr.IsError {
					write.Hint(cmd, config.LabelError)
				}
				if tr.Content != "" {
					content := core.StripLineNumbers(tr.Content)
					write.CodeBlock(cmd, content)
				}
			}

			if len(msg.ToolUses) > 0 || len(msg.ToolResults) > 0 {
				write.BlankLine(cmd)
			}
		}
	} else {
		write.SectionHeader(cmd, 2, config.SectionConversationPreview)

		count := 0
		for _, msg := range session.Messages {
			if msg.BelongsToUser() && msg.Text != "" {
				count++
				if count > config.PreviewMaxTurns {
					write.MoreTurns(cmd, session.TurnCount-config.PreviewMaxTurns)
					break
				}
				text := msg.Text
				if len(text) > config.PreviewMaxTextLen {
					text = text[:config.PreviewMaxTextLen] + config.Ellipsis
				}
				write.NumberedItem(cmd, count, text)
			}
		}
		write.BlankLine(cmd)
		write.Hint(cmd, config.HintUseFullFlag)
	}

	return nil
}
