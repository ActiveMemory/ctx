//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"fmt"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/mcp/governance"
	"github.com/ActiveMemory/ctx/internal/config/mcp/tool"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// RecordSessionStart marks the session as explicitly started and
// resets the session start timestamp.
//
// Called by the session_event tool when the agent reports a "start"
// event. Sets sessionStarted to true and captures the current wall
// time so governance checks can measure elapsed time.
func (ss *State) RecordSessionStart() {
	ss.sessionStarted = true
	ss.sessionStartedAt = time.Now()
}

// RecordContextLoaded marks context as loaded for this session.
//
// Called after the agent successfully loads context files (TASKS.md,
// DECISIONS.md, etc.). Suppresses the "context not loaded" governance
// warning that would otherwise appear on every tool response.
func (ss *State) RecordContextLoaded() {
	ss.contextLoaded = true
}

// RecordDriftCheck records that a drift check was performed.
//
// Called after the agent runs ctx_drift. Updates the last-drift-check
// timestamp so CheckGovernance can determine whether a follow-up drift
// check is overdue based on governance.DriftCheckInterval.
func (ss *State) RecordDriftCheck() {
	ss.lastDriftCheck = time.Now()
}

// RecordContextWrite records that a .context/ write occurred.
//
// Called after successful ctx_add, ctx_complete, ctx_watch_update, or
// ctx_compact invocations. Captures the current wall time and resets
// the calls-since-write counter to zero, which suppresses persist
// nudges until governance.PersistNudgeAfter more tool calls elapse.
func (ss *State) RecordContextWrite() {
	ss.lastContextWrite = time.Now()
	ss.callsSinceWrite = 0
}

// IncrementCallsSinceWrite bumps the counter used for persist nudges.
//
// Called by the MCP server after every tool dispatch regardless of tool
// type. When the counter reaches governance.PersistNudgeAfter,
// CheckGovernance begins emitting persist nudge warnings.
func (ss *State) IncrementCallsSinceWrite() {
	ss.callsSinceWrite++
}

// CheckGovernance returns governance warnings that should be
// appended to the current tool response. Returns an empty string
// when no action is warranted.
//
// Parameters:
//   - toolName: the MCP tool that was just called, used to
//     suppress redundant warnings (e.g. drift warning is not
//     appended to a ctx_drift response)
//
// Returns:
//   - string: newline-separated warnings preceded by a separator,
//     or empty string when no warnings apply
func (ss *State) CheckGovernance(toolName string) string {
	var warnings []string

	// 1. Session not started
	if !ss.sessionStarted && toolName != tool.SessionEvent {
		warnings = append(warnings,
			desc.Text(text.DescKeyGovSessionNotStarted))
	}

	// 2. Context not loaded
	if !ss.contextLoaded && toolName != tool.Status &&
		toolName != tool.SessionEvent {
		warnings = append(warnings,
			desc.Text(text.DescKeyGovContextNotLoaded))
	}

	// 3. Drift not checked recently
	if ss.sessionStarted && toolName != tool.Drift &&
		toolName != tool.SessionEvent {
		if !ss.lastDriftCheck.IsZero() {
			if time.Since(ss.lastDriftCheck) > governance.DriftCheckInterval {
				warnings = append(warnings, fmt.Sprintf(
					desc.Text(text.DescKeyGovDriftNotChecked),
					int(time.Since(ss.lastDriftCheck).Minutes())))
			}
		} else if ss.ToolCalls > 5 {
			// Never checked drift and already 5+ calls in
			warnings = append(warnings,
				desc.Text(text.DescKeyGovDriftNeverChecked))
		}
	}

	// 4. Persist nudge — no context writes in a while
	if ss.sessionStarted && ss.callsSinceWrite >= governance.PersistNudgeAfter &&
		toolName != tool.Add && toolName != tool.WatchUpdate &&
		toolName != tool.Complete && toolName != tool.Compact &&
		toolName != tool.SessionEvent {
		// Fire at threshold, then every governance.PersistNudgeRepeat
		// calls after.
		if ss.callsSinceWrite == governance.PersistNudgeAfter ||
			(ss.callsSinceWrite-governance.PersistNudgeAfter)%governance.PersistNudgeRepeat == 0 {
			warnings = append(warnings, fmt.Sprintf(
				desc.Text(text.DescKeyGovPersistNudge),
				ss.callsSinceWrite))
		}
	}

	// 5. Violations from extension detection ring
	if violations := readAndClearViolations(ss.contextDir); len(violations) > 0 {
		for _, v := range violations {
			detail := v.Detail
			if len(detail) > 120 {
				detail = detail[:120] + token.Ellipsis
			}
			warnings = append(warnings, fmt.Sprintf(
				desc.Text(text.DescKeyGovViolationCritical),
				v.Kind, detail, v.Timestamp))
		}
	}

	if len(warnings) == 0 {
		return ""
	}

	nl := token.NewlineLF
	return nl + nl + token.Separator + nl + strings.Join(warnings, nl)
}
