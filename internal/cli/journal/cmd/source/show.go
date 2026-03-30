//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package source

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/query"
	srcFmt "github.com/ActiveMemory/ctx/internal/cli/journal/core/source/format"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	errSession "github.com/ActiveMemory/ctx/internal/err/session"
	sharedFmt "github.com/ActiveMemory/ctx/internal/format"
	"github.com/ActiveMemory/ctx/internal/parse"
	writeRecall "github.com/ActiveMemory/ctx/internal/write/journal"
)

// runShow displays detailed information about a session including metadata,
// token usage, tool usage summary, and optionally the full conversation.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - args: positional arguments (session ID triggers show mode)
//   - opts: combined flags including ShowID, Latest, Full, and AllProjects
//
// Returns:
//   - error: non-nil if session not found or scanning fails
func runShow(cmd *cobra.Command, args []string, opts Opts) error {
	// If --show <id> was used, pass the ID as a positional arg.
	showArgs := args
	if opts.ShowID != "" {
		showArgs = []string{opts.ShowID}
	}

	sessions, scanErr := query.FindSessions(opts.AllProjects)
	if scanErr != nil {
		return errSession.Find(scanErr)
	}

	if len(sessions) == 0 {
		if opts.AllProjects {
			return errSession.NoneFound("")
		}
		return errSession.NoneFound(
			desc.Text(text.DescKeyLabelHintUseAllProjects),
		)
	}

	var session *entity.Session

	switch {
	case opts.Latest:
		session = sessions[0]
	case len(showArgs) == 0:
		return errSession.IDRequired()
	default:
		q := strings.ToLower(showArgs[0])
		var matches []*entity.Session
		for _, s := range sessions {
			if strings.HasPrefix(strings.ToLower(s.ID), q) ||
				strings.Contains(strings.ToLower(s.Slug), q) {
				matches = append(matches, s)
			}
		}
		if len(matches) == 0 {
			return errSession.NotFound(showArgs[0])
		}
		if len(matches) > 1 {
			lines := srcFmt.SessionMatchLines(matches)
			writeRecall.AmbiguousSessionMatchWithHint(
				cmd, showArgs[0], lines,
				matches[0].ID[:journal.SessionIDHintLen],
			)
			return errSession.AmbiguousQuery()
		}
		session = matches[0]
	}

	// Print session details.
	writeRecall.SessionMetadata(cmd, writeRecall.SessionInfo{
		Slug:      session.Slug,
		ID:        session.ID,
		Tool:      session.Tool,
		Project:   session.Project,
		Branch:    session.GitBranch,
		Model:     session.Model,
		Started:   session.StartTime.Format(time.DateTimePreciseFmt),
		Duration:  srcFmt.Duration(session.Duration),
		Turns:     session.TurnCount,
		Messages:  len(session.Messages),
		TokensIn:  sharedFmt.Tokens(session.TotalTokensIn),
		TokensOut: sharedFmt.Tokens(session.TotalTokensOut),
		TokensAll: sharedFmt.Tokens(session.TotalTokens),
	})

	// Tool usage summary
	tools := session.AllToolUses()
	if len(tools) > 0 {
		toolCounts := make(map[string]int)
		for _, t := range tools {
			toolCounts[t.Name]++
		}

		writeRecall.SectionHeader(
			cmd, 2, desc.Text(text.DescKeyLabelSectionToolUsage),
		)
		for name, count := range toolCounts {
			writeRecall.ListItem(
				cmd, desc.Text(text.DescKeyJournalSourceToolCountLine),
				name, count,
			)
		}
		writeRecall.BlankLine(cmd)
	}

	// Messages
	if opts.Full {
		writeRecall.SectionHeader(
			cmd, 2, desc.Text(text.DescKeyLabelSectionConversation),
		)

		for i, msg := range session.Messages {
			role := desc.Text(text.DescKeyLabelRoleUser)
			if msg.BelongsToAssistant() {
				role = desc.Text(text.DescKeyLabelRoleAssistant)
			} else if len(msg.ToolResults) > 0 && msg.Text == "" {
				role = desc.Text(text.DescKeyLabelToolOutput)
			}

			writeRecall.ConversationTurn(
				cmd, i+1, role, msg.Timestamp.Format(time.Format),
			)

			if msg.Text != "" {
				writeRecall.TextBlock(cmd, msg.Text)
			}

			for _, t := range msg.ToolUses {
				toolInfo := srcFmt.ToolUse(t)
				writeRecall.SessionDetail(
					cmd,
					desc.Text(text.DescKeyLabelInlineTool),
					toolInfo,
				)
			}

			for _, tr := range msg.ToolResults {
				if tr.IsError {
					writeRecall.Hint(
						cmd, desc.Text(text.DescKeyLabelInlineError),
					)
				}
				if tr.Content != "" {
					content := parse.StripLineNumbers(tr.Content)
					writeRecall.CodeBlock(cmd, content)
				}
			}

			if len(msg.ToolUses) > 0 || len(msg.ToolResults) > 0 {
				writeRecall.BlankLine(cmd)
			}
		}
	} else {
		writeRecall.SectionHeader(
			cmd, 2,
			desc.Text(text.DescKeyLabelSectionConversationPreview),
		)

		count := 0
		for _, msg := range session.Messages {
			if msg.BelongsToUser() && msg.Text != "" {
				count++
				if count > journal.PreviewMaxTurns {
					writeRecall.MoreTurns(
						cmd,
						session.TurnCount-journal.PreviewMaxTurns,
					)
					break
				}
				t := msg.Text
				if len(t) > journal.PreviewMaxTextLen {
					t = t[:journal.PreviewMaxTextLen] + token.Ellipsis
				}
				writeRecall.NumberedItem(cmd, count, t)
			}
		}
		writeRecall.BlankLine(cmd)
		writeRecall.Hint(cmd, desc.Text(text.DescKeyLabelHintUseFull))
	}

	return nil
}
