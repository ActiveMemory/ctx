//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
)

// buildSessionContent creates the Markdown content for a session file.
//
// Assembles a session document with metadata, current tasks, recent decisions,
// and learnings. Uses CTX_SESSION_START environment variable for session
// correlation if available.
//
// Parameters:
//   - topic: Session topic used as the document title
//   - sessionType: Type of session (e.g., "manual", "auto-save")
//   - timestamp: Time used for end_time and fallback start_time
//
// Returns:
//   - string: Complete Markdown content for the session file
//   - error: Currently always nil (reserved for future validation)
func buildSessionContent(
	topic, sessionType string, timestamp time.Time,
) (string, error) {
	var sb strings.Builder
	nl := config.NewlineLF
	sep := config.Separator

	// Header with timestamp fields for session correlation
	sb.WriteString(fmt.Sprintf(config.TplSessionHeading+nl+nl, topic))
	sb.WriteString(
		fmt.Sprintf(config.MetadataDate+" %s"+nl, timestamp.Format("2006-01-02")),
	)
	sb.WriteString(
		fmt.Sprintf(config.MetadataTime+" %s"+nl, timestamp.Format("15:04:05")),
	)
	sb.WriteString(fmt.Sprintf(config.MetadataType+" %s"+nl, sessionType))

	// Session correlation timestamps
	// (YYYY-MM-DD-HHMM format matches ctx add timestamps)
	// start_time: When session began
	// (use CTX_SESSION_START env var if available, else save time)
	startTime := timestamp
	if envStart := os.Getenv(config.EnvCtxSessionStart); envStart != "" {
		if parsed, parseErr := time.Parse(
			"2006-01-02-1504", envStart,
		); parseErr == nil {
			startTime = parsed
		}
	}
	sb.WriteString(
		fmt.Sprintf(
			config.MetadataStartTime+" %s"+nl, startTime.Format("2006-01-02-1504"),
		),
	)
	sb.WriteString(
		fmt.Sprintf(
			config.MetadataEndTime+" %s"+nl, timestamp.Format("2006-01-02-1504"),
		),
	)
	sb.WriteString(nl + sep + nl + nl)

	// Summary section (placeholder for the user to fill in)
	sb.WriteString(config.RecallHeadingSummary + nl + nl)
	sb.WriteString(config.TplSessionPlaceholderSummary + nl + nl)
	sb.WriteString(sep + nl + nl)

	// Current Tasks
	sb.WriteString(config.SessionHeadingCurrentTasks + nl + nl)
	tasks, readErr := readContextSection(
		config.FileTask, config.HeadingInProgress, config.HeadingNextUp,
	)
	if readErr == nil && tasks != "" {
		sb.WriteString(config.SessionHeadingInProgress + nl + nl)
		sb.WriteString(tasks)
		sb.WriteString(nl)
	}
	nextTasks, readErr := readContextSection(
		config.FileTask, config.HeadingNextUp, config.HeadingCompleted,
	)
	if readErr == nil && nextTasks != "" {
		sb.WriteString(config.SessionHeadingNextUp + nl + nl)
		sb.WriteString(nextTasks)
		sb.WriteString(nl)
	}
	sb.WriteString(sep + nl + nl)

	// Recent Decisions
	sb.WriteString(config.SessionHeadingRecentDecisions + nl + nl)
	decisions, readErr := readRecentDecisions()
	if readErr == nil && decisions != "" {
		sb.WriteString(decisions)
	} else {
		sb.WriteString(config.TplSessionPlaceholderNoDecisions + nl)
	}
	sb.WriteString(nl + sep + nl + nl)

	// Recent Learnings
	sb.WriteString(config.SessionHeadingRecentLearnings + nl + nl)
	learnings, readErr := readRecentLearnings()
	if readErr == nil && learnings != "" {
		sb.WriteString(learnings)
	} else {
		sb.WriteString(config.TplSessionPlaceholderNoLearnings + nl)
	}
	sb.WriteString(nl + sep + nl + nl)

	// Tasks for Next Session
	sb.WriteString(config.SessionHeadingNextSessionTasks + nl + nl)
	sb.WriteString(config.TplSessionPlaceholderNextTasks + nl + nl)

	return sb.String(), nil
}
