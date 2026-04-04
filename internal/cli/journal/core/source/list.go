//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package source

import (
	"fmt"
	"strconv"
	"strings"
	goTime "time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/assets/tpl"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/query"
	srcFmt "github.com/ActiveMemory/ctx/internal/cli/journal/core/source/format"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/err/date"
	errSession "github.com/ActiveMemory/ctx/internal/err/session"
	sharedFmt "github.com/ActiveMemory/ctx/internal/format"
	"github.com/ActiveMemory/ctx/internal/parse"
	writeRecall "github.com/ActiveMemory/ctx/internal/write/journal"
)

// RunList finds all sessions, applies optional filters, and
// displays them in a formatted list with project, time, turn
// count, and preview.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - opts: combined flags including limit, project, tool,
//     since, until, and allProjects
//
// Returns:
//   - error: non-nil if date parsing or scanning fails
func RunList(cmd *cobra.Command, opts Opts) error {
	// Parse date filters (only when flags are set).
	var sinceTime, untilTime goTime.Time
	if opts.Since != "" {
		var sinceErr error
		sinceTime, sinceErr = parse.Date(opts.Since)
		if sinceErr != nil {
			return date.Invalid(
				flag.PrefixLong+flag.Since,
				opts.Since, sinceErr,
			)
		}
	}
	if opts.Until != "" {
		parsed, untilErr := parse.Date(opts.Until)
		if untilErr != nil {
			return date.Invalid(
				flag.PrefixLong+flag.Until,
				opts.Until, untilErr,
			)
		}
		// --until is inclusive: advance to end of day
		untilTime = parsed.Add(time.InclusiveUntilOffset)
	}

	sessions, scanErr := query.FindSessions(
		opts.AllProjects,
	)
	if scanErr != nil {
		return errSession.Find(scanErr)
	}

	if len(sessions) == 0 {
		writeRecall.NoSessionsWithHint(
			cmd, opts.AllProjects,
		)
		return nil
	}

	// Apply filters
	var filtered []*entity.Session
	for _, s := range sessions {
		if opts.Project != "" && !strings.Contains(
			strings.ToLower(s.Project),
			strings.ToLower(opts.Project),
		) {
			continue
		}
		if opts.Tool != "" && s.Tool != opts.Tool {
			continue
		}
		if opts.Since != "" &&
			s.StartTime.Before(sinceTime) {
			continue
		}
		if opts.Until != "" &&
			s.StartTime.After(untilTime) {
			continue
		}
		filtered = append(filtered, s)
	}

	if len(filtered) == 0 {
		writeRecall.NoFiltersMatch(cmd)
		return nil
	}

	// Apply limit
	if opts.Limit > 0 && len(filtered) > opts.Limit {
		filtered = filtered[:opts.Limit]
	}

	shown := 0
	if opts.Project != "" || opts.Tool != "" {
		shown = len(filtered)
	}
	writeRecall.SessionListHeader(cmd, len(sessions), shown)

	// Compute dynamic column widths from data.
	slugW, projW := len(
		desc.Text(text.DescKeyLabelColSlug),
	), len(desc.Text(text.DescKeyLabelColProject))
	for _, s := range filtered {
		slug := sharedFmt.Truncate(
			s.Slug, journal.SlugMaxLen,
		)
		if len(slug) > slugW {
			slugW = len(slug)
		}
		if len(s.Project) > projW {
			projW = len(s.Project)
		}
	}

	// Print column header.
	rowFmt := fmt.Sprintf(tpl.RecallListRow, slugW, projW)
	writeRecall.SessionListRow(cmd, rowFmt,
		desc.Text(text.DescKeyLabelColSlug),
		desc.Text(text.DescKeyLabelColProject),
		desc.Text(text.DescKeyLabelColDate),
		desc.Text(text.DescKeyLabelColDuration),
		desc.Text(text.DescKeyLabelColTurns),
		desc.Text(text.DescKeyLabelColUsage),
	)

	// Print sessions.
	for _, s := range filtered {
		slug := sharedFmt.Truncate(
			s.Slug, journal.SlugMaxLen,
		)
		dateStr := s.StartTime.Local().Format(
			time.DateTimeFmt,
		)
		dur := srcFmt.Duration(s.Duration)
		turns := strconv.Itoa(s.TurnCount)
		tokens := ""
		if s.TotalTokens > 0 {
			tokens = sharedFmt.Tokens(s.TotalTokens)
		}
		writeRecall.SessionListRow(cmd, rowFmt,
			slug, s.Project, dateStr, dur, turns, tokens,
		)
	}

	writeRecall.SessionListFooter(
		cmd, len(sessions) > len(filtered),
	)

	return nil
}
