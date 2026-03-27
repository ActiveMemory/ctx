//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/cli/journal/core/query"
	sourceFormat "github.com/ActiveMemory/ctx/internal/cli/journal/core/source/format"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/assets/tpl"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/err/date"
	errSession "github.com/ActiveMemory/ctx/internal/err/session"
	"github.com/ActiveMemory/ctx/internal/parse"
	"github.com/ActiveMemory/ctx/internal/write/recall"
)

// Run handles the recall list command.
//
// Finds all sessions, applies optional filters, and displays them in a
// formatted list with project, time, turn count, and preview.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - limit: maximum sessions to display (0 for unlimited)
//   - project: filter by project name (case-insensitive substring match)
//   - tool: filter by tool identifier (exact match)
//   - since: inclusive start date filter (YYYY-MM-DD)
//   - until: inclusive end date filter (YYYY-MM-DD)
//   - allProjects: if true, include sessions from all projects
//
// Returns:
//   - error: non-nil if date parsing or session scanning fails
func Run(
	cmd *cobra.Command, limit int, project, tool,
	since, until string,
	allProjects bool,
) error {
	// Parse date filters
	sinceTime, sinceErr := parse.Date(since)
	if since != "" && sinceErr != nil {
		return date.InvalidDate(flag.PrefixLong+flag.Since, since, sinceErr)
	}
	untilTime, untilErr := parse.Date(until)
	if until != "" && untilErr != nil {
		return date.InvalidDate(flag.PrefixLong+flag.Until, until, untilErr)
	}
	// --until is inclusive: advance to the end of the day
	if until != "" {
		untilTime = untilTime.Add(time.InclusiveUntilOffset)
	}

	sessions, scanErr := query.FindSessions(allProjects)
	if scanErr != nil {
		return errSession.Find(scanErr)
	}

	if len(sessions) == 0 {
		recall.NoSessionsWithHint(cmd, allProjects)
		return nil
	}

	// Apply filters
	var filtered []*entity.Session
	for _, s := range sessions {
		if project != "" && !strings.Contains(
			strings.ToLower(s.Project), strings.ToLower(project),
		) {
			continue
		}
		if tool != "" && s.Tool != tool {
			continue
		}
		if since != "" && s.StartTime.Before(sinceTime) {
			continue
		}
		if until != "" && s.StartTime.After(untilTime) {
			continue
		}
		filtered = append(filtered, s)
	}

	if len(filtered) == 0 {
		recall.NoFiltersMatch(cmd)
		return nil
	}

	// Apply limit
	if limit > 0 && len(filtered) > limit {
		filtered = filtered[:limit]
	}

	shown := 0
	if project != "" || tool != "" {
		shown = len(filtered)
	}
	recall.SessionListHeader(cmd, len(sessions), shown)

	// Compute dynamic column widths from data.
	slugW, projW := len(desc.Text(text.DescKeyLabelColSlug)),
		len(desc.Text(text.DescKeyLabelColProject))
	for _, s := range filtered {
		slug := sourceFormat.Truncate(s.Slug, journal.SlugMaxLen)
		if len(slug) > slugW {
			slugW = len(slug)
		}
		if len(s.Project) > projW {
			projW = len(s.Project)
		}
	}

	// Print column header.
	rowFmt := fmt.Sprintf(tpl.RecallListRow, slugW, projW)
	recall.SessionListRow(cmd, rowFmt,
		desc.Text(text.DescKeyLabelColSlug),
		desc.Text(text.DescKeyLabelColProject),
		desc.Text(text.DescKeyLabelColDate),
		desc.Text(text.DescKeyLabelColDuration),
		desc.Text(text.DescKeyLabelColTurns),
		desc.Text(text.DescKeyLabelColUsage),
	)

	// Print sessions.
	for _, s := range filtered {
		slug := sourceFormat.Truncate(s.Slug, journal.SlugMaxLen)
		dateStr := s.StartTime.Local().Format(time.DateTimeFormat)
		dur := sourceFormat.Duration(s.Duration)
		turns := fmt.Sprintf("%d", s.TurnCount)
		tokens := ""
		if s.TotalTokens > 0 {
			tokens = sourceFormat.Tokens(s.TotalTokens)
		}
		recall.SessionListRow(cmd, rowFmt,
			slug, s.Project, dateStr, dur, turns, tokens)
	}

	recall.SessionListFooter(cmd, len(sessions) > len(filtered))

	return nil
}
