//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/recall"
	"github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/recall/core"
	"github.com/ActiveMemory/ctx/internal/config"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/recall/parser"
	"github.com/ActiveMemory/ctx/internal/write"
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
				cmd, args[0], lines, matches[0].ID[:recall.SessionIDHintLen],
			)
			return ctxerr.AmbiguousQuery()
		}
		session = matches[0]
	}

	// Print session details.
	write.SessionMetadata(cmd, write.SessionInfo{
		Slug:      session.Slug,
		ID:        session.ID,
		Tool:      session.Tool,
		Project:   session.Project,
		Branch:    session.GitBranch,
		Model:     session.Model,
		Started:   session.StartTime.Format(time.DateTimePreciseFormat),
		Duration:  core.FormatDuration(session.Duration),
		Turns:     session.TurnCount,
		Messages:  len(session.Messages),
		TokensIn:  core.FormatTokens(session.TotalTokensIn),
		TokensOut: core.FormatTokens(session.TotalTokensOut),
		TokensAll: core.FormatTokens(session.TotalTokens),
	})

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
				cmd, i+1, role, msg.Timestamp.Format(time.Format),
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
				if count > recall.PreviewMaxTurns {
					write.MoreTurns(cmd, session.TurnCount-recall.PreviewMaxTurns)
					break
				}
				text := msg.Text
				if len(text) > recall.PreviewMaxTextLen {
					text = text[:recall.PreviewMaxTextLen] + config.Ellipsis
				}
				write.NumberedItem(cmd, count, text)
			}
		}
		write.BlankLine(cmd)
		write.Hint(cmd, config.HintUseFullFlag)
	}

	return nil
}
