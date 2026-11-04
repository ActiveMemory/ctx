//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package source

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	srcFmt "github.com/ActiveMemory/ctx/internal/cli/journal/core/source/format"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	errSession "github.com/ActiveMemory/ctx/internal/err/session"
	"github.com/ActiveMemory/ctx/internal/parse"
	writeRecall "github.com/ActiveMemory/ctx/internal/write/journal"
)

// runSelectiveShow handles --turns and --grep show modes for a session.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - session: Target session to inspect
//   - opts: Source command options
//
// Returns:
//   - error: Non-nil on invalid --turns, --grep, or --max-chars
func runSelectiveShow(
	cmd *cobra.Command, session *entity.Session, opts Opts,
) error {
	if opts.MaxChars < 0 {
		return errSession.NegativeMaxChars()
	}
	total := len(session.Messages)
	start := 1
	end := total
	if opts.Turns != "" {
		var rangeErr error
		start, end, rangeErr = ParseTurnRange(opts.Turns, total)
		if rangeErr != nil {
			return rangeErr
		}
	}

	sid := shortSessionID(session.ID)
	if total == 0 || start > total {
		writeRecall.TurnRangePastEnd(cmd, sid, total, start)
		return nil
	}

	if opts.Grep != "" {
		return renderGrepMatches(cmd, session, sid, start, end, opts)
	}

	renderTurnSlice(cmd, session, sid, start, end, opts.MaxChars)
	return nil
}

// renderGrepMatches compiles the --grep pattern, searches [start, end], and
// prints matching turn snippets with a follow-up --turns hint.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - session: Target session
//   - sid: Short session ID for hints
//   - start: 1-based inclusive start turn
//   - end: 1-based inclusive end turn
//   - opts: Source command options
//
// Returns:
//   - error: Non-nil if the regex pattern fails to compile
func renderGrepMatches(
	cmd *cobra.Command,
	session *entity.Session,
	sid string,
	start, end int,
	opts Opts,
) error {
	re, compileErr := regex.CompileDynamic(
		regex.CaseInsensitivePrefix + opts.Grep,
	)
	if compileErr != nil {
		return errSession.InvalidGrepPattern(opts.Grep, compileErr)
	}
	searched := end - start + 1
	matches := MatchTurns(session.Messages, start, end, re)
	if len(matches) == 0 {
		writeRecall.GrepMatchesEmpty(cmd, opts.Grep, sid, searched)
		return nil
	}
	writeRecall.GrepMatchesHeader(cmd, sid, len(matches), len(session.Messages))
	for _, m := range matches {
		writeRecall.GrepMatchItem(cmd, m.Turn, m.Role, m.Timestamp, m.Snippet)
	}
	first := matches[0].Turn
	hintStart := max(1, first-1)
	hintEnd := min(len(session.Messages), first+2)
	writeRecall.BlankLine(cmd)
	writeRecall.GrepMatchesHint(cmd, sid, hintStart, hintEnd)
	return nil
}

// renderTurnSlice formats and prints turns [start, end] subject to maxChars,
// along with top and bottom continuation banners.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - session: Target session
//   - sid: Short session identifier
//   - start: 1-based inclusive start turn
//   - end: 1-based inclusive end turn
//   - maxChars: Maximum character budget (0 = unlimited)
func renderTurnSlice(
	cmd *cobra.Command,
	session *entity.Session,
	sid string,
	start, end, maxChars int,
) {
	total := len(session.Messages)
	windowSize := max(1, end-start+1)
	var out strings.Builder
	actualEnd := start - 1

	for idx := start - 1; idx < end && idx < total; idx++ {
		turnBlock := formatSingleTurn(idx+1, session.Messages[idx])
		if maxChars > 0 && out.Len()+len(turnBlock) > maxChars {
			rem := maxChars - out.Len()
			if rem > 0 {
				truncated := truncateUTF8(turnBlock, rem)
				out.WriteString(truncated)
				out.WriteString(fmt.Sprintf(
					desc.Text(text.DescKeyWriteJournalSourceTruncatedChars),
					maxChars,
				))
				out.WriteString(token.NewlineLF)
				actualEnd = idx + 1
			}
			break
		}
		out.WriteString(turnBlock)
		actualEnd = idx + 1
	}

	if actualEnd < start {
		actualEnd = start
	}
	nextStart := 0
	nextEnd := 0
	if actualEnd < total {
		nextStart = actualEnd + 1
		nextEnd = min(total, nextStart+windowSize-1)
	}

	body := strings.TrimRight(out.String(), token.NewlineLF)
	chars := utf8.RuneCountInString(body)
	writeRecall.SectionHeader(
		cmd, 2, desc.Text(text.DescKeyLabelSectionConversation),
	)
	writeRecall.TurnRangeBanner(
		cmd, sid, start, actualEnd, total, chars, nextStart, nextEnd,
	)
	writeRecall.BlankLine(cmd)
	writeRecall.TextBlock(cmd, body)
	writeRecall.TurnRangeBanner(
		cmd, sid, start, actualEnd, total, chars, nextStart, nextEnd,
	)
}

// formatSingleTurn renders a single message turn into Markdown text.
//
// Parameters:
//   - turnNum: 1-based turn index
//   - msg: Session message to format
//
// Returns:
//   - string: Rendered Markdown block for the turn
func formatSingleTurn(turnNum int, msg entity.Message) string {
	var sb strings.Builder
	nl := token.NewlineLF
	role := turnRoleLabel(msg)
	sb.WriteString(fmt.Sprintf(
		desc.Text(text.DescKeyWriteJournalSourceConversationTurn),
		turnNum, role, msg.Timestamp.Format(time.Format),
	))
	sb.WriteString(nl)
	if msg.Text != "" {
		sb.WriteString(msg.Text + nl + nl)
	}
	for _, tu := range msg.ToolUses {
		sb.WriteString(fmt.Sprintf(
			desc.Text(text.DescKeyWriteJournalSourceDetailString)+nl,
			desc.Text(text.DescKeyLabelInlineTool),
			srcFmt.ToolUse(tu),
		))
	}
	for _, tr := range msg.ToolResults {
		if tr.IsError {
			sb.WriteString(desc.Text(text.DescKeyLabelInlineError) + nl)
		}
		if tr.Content != "" {
			cleaned := parse.StripLineNumbers(tr.Content)
			sb.WriteString(fmt.Sprintf(
				desc.Text(text.DescKeyWriteJournalSourceCodeBlock),
				cleaned,
			))
		}
	}
	if len(msg.ToolUses) > 0 || len(msg.ToolResults) > 0 {
		sb.WriteString(nl)
	}
	return sb.String()
}

// turnRoleLabel returns the localized display label for a message's role.
//
// Parameters:
//   - msg: Session message
//
// Returns:
//   - string: Role label ("User", "Assistant", or "Tool Output")
func turnRoleLabel(msg entity.Message) string {
	if msg.BelongsToAssistant() {
		return desc.Text(text.DescKeyLabelRoleAssistant)
	}
	if len(msg.ToolResults) > 0 && msg.Text == "" {
		return desc.Text(text.DescKeyLabelToolOutput)
	}
	return desc.Text(text.DescKeyLabelRoleUser)
}

// extractMatchSnippet extracts a single-line snippet of at most maxChars around
// the matched byte span loc.
//
// Parameters:
//   - textBody: Full text containing the match
//   - loc: Two-element byte index slice [start, end] from regexp
//   - maxChars: Maximum snippet length in bytes/chars
//
// Returns:
//   - string: Single-line snippet with ellipses when clipped
func extractMatchSnippet(textBody string, loc []int, maxChars int) string {
	half := maxChars / 2
	s := max(0, loc[0]-half)
	e := min(len(textBody), loc[1]+half)
	for s > 0 && !utf8.RuneStart(textBody[s]) {
		s--
	}
	for e < len(textBody) && !utf8.RuneStart(textBody[e]) {
		e++
	}
	raw := strings.ReplaceAll(textBody[s:e], token.NewlineLF, token.Space)
	raw = strings.TrimSpace(raw)
	if s > 0 {
		raw = token.Ellipsis + raw
	}
	if e < len(textBody) {
		raw = raw + token.Ellipsis
	}
	return raw
}

// truncateUTF8 truncates s to at most maxBytes on a valid UTF-8 rune boundary.
//
// Parameters:
//   - s: Input string
//   - maxBytes: Maximum byte length
//
// Returns:
//   - string: Valid UTF-8 prefix of s
func truncateUTF8(s string, maxBytes int) string {
	if len(s) <= maxBytes {
		return s
	}
	end := maxBytes
	for end > 0 && !utf8.RuneStart(s[end]) {
		end--
	}
	return s[:end]
}

// shortSessionID returns the short prefix of a session ID for command hints.
//
// Parameters:
//   - id: Full session ID
//
// Returns:
//   - string: Short session ID prefix
func shortSessionID(id string) string {
	if len(id) > journal.SessionIDShortLen {
		return id[:journal.SessionIDShortLen]
	}
	return id
}
