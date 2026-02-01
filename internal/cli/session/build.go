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
	sb.WriteString(fmt.Sprintf("# Session: %s"+nl+nl, topic))
	sb.WriteString(fmt.Sprintf("**Date**: %s"+nl, timestamp.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("**Time**: %s"+nl, timestamp.Format("15:04:05")))
	sb.WriteString(fmt.Sprintf("**Type**: %s"+nl, sessionType))

	// Session correlation timestamps
	// (YYYY-MM-DD-HHMM format matches ctx add timestamps)
	// start_time: When session began
	// (use CTX_SESSION_START env var if available, else save time)
	startTime := timestamp
	if envStart := os.Getenv("CTX_SESSION_START"); envStart != "" {
		if parsed, err := time.Parse("2006-01-02-1504", envStart); err == nil {
			startTime = parsed
		}
	}
	sb.WriteString(
		fmt.Sprintf("**start_time**: %s"+nl, startTime.Format("2006-01-02-1504")),
	)
	sb.WriteString(
		fmt.Sprintf("**end_time**: %s"+nl, timestamp.Format("2006-01-02-1504")),
	)
	sb.WriteString(nl + sep + nl + nl)

	// Summary section (placeholder for the user to fill in)
	sb.WriteString("## Summary" + nl + nl)
	sb.WriteString("[Describe what was accomplished in this session]" + nl + nl)
	sb.WriteString(sep + nl + nl)

	// Current Tasks
	sb.WriteString("## Current Tasks" + nl + nl)
	tasks, err := readContextSection(
		"TASKS.md", "## In Progress", "## Next Up",
	)
	if err == nil && tasks != "" {
		sb.WriteString("### In Progress" + nl + nl)
		sb.WriteString(tasks)
		sb.WriteString(nl)
	}
	nextTasks, err := readContextSection(
		"TASKS.md", "## Next Up", "## Completed",
	)
	if err == nil && nextTasks != "" {
		sb.WriteString("### Next Up" + nl + nl)
		sb.WriteString(nextTasks)
		sb.WriteString(nl)
	}
	sb.WriteString(sep + nl + nl)

	// Recent Decisions
	sb.WriteString("## Recent Decisions" + nl + nl)
	decisions, err := readRecentDecisions()
	if err == nil && decisions != "" {
		sb.WriteString(decisions)
	} else {
		sb.WriteString("[No recent decisions found]" + nl)
	}
	sb.WriteString(nl + sep + nl + nl)

	// Recent Learnings
	sb.WriteString("## Recent Learnings" + nl + nl)
	learnings, err := readRecentLearnings()
	if err == nil && learnings != "" {
		sb.WriteString(learnings)
	} else {
		sb.WriteString("[No recent learnings found]" + nl)
	}
	sb.WriteString(nl + sep + nl + nl)

	// Tasks for Next Session
	sb.WriteString("## Tasks for Next Session" + nl + nl)
	sb.WriteString("[List tasks to continue in the next session]" + nl + nl)

	return sb.String(), nil
}
