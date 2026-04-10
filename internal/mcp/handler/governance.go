//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgFmt "github.com/ActiveMemory/ctx/internal/config/format"
	"github.com/ActiveMemory/ctx/internal/config/mcp/governance"
	"github.com/ActiveMemory/ctx/internal/config/mcp/tool"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// CheckGovernance returns governance warnings that should be
// appended to the current tool response. Returns an empty string
// when no action is warranted.
//
// It also drains the VS Code extension's violations file for the
// given context directory, which is why this is a free function in
// the handler package (I/O) rather than a method on
// [entity.MCPSession].
//
// Parameters:
//   - d: runtime dependencies carrying the session state and
//     context directory
//   - toolName: the MCP tool that was just called, used to
//     suppress redundant warnings (e.g. drift warning is not
//     appended to a ctx_drift response)
//
// Returns:
//   - string: newline-separated warnings preceded by a separator,
//     or empty string when no warnings apply
func CheckGovernance(d *entity.MCPDeps, toolName string) string {
	ss := d.Session
	var warnings []string

	// 1. Session not started
	if !ss.SessionStarted && toolName != tool.SessionEvent {
		warnings = append(warnings,
			desc.Text(text.DescKeyGovSessionNotStarted))
	}

	// 2. Context not loaded
	if !ss.ContextLoaded && toolName != tool.Status &&
		toolName != tool.SessionEvent {
		warnings = append(warnings,
			desc.Text(text.DescKeyGovContextNotLoaded))
	}

	// 3. Drift not checked recently
	if ss.SessionStarted && toolName != tool.Drift &&
		toolName != tool.SessionEvent {
		if !ss.LastDriftCheck.IsZero() {
			if time.Since(ss.LastDriftCheck) > governance.DriftCheckInterval {
				warnings = append(warnings, fmt.Sprintf(
					desc.Text(text.DescKeyGovDriftNotChecked),
					int(time.Since(ss.LastDriftCheck).Minutes())))
			}
		} else if ss.ToolCalls > governance.DriftCheckMinCalls {
			// Never checked drift and already past threshold
			warnings = append(warnings,
				desc.Text(text.DescKeyGovDriftNeverChecked))
		}
	}

	// 4. Persist nudge — no context writes in a while
	if ss.SessionStarted && ss.CallsSinceWrite >= governance.PersistNudgeAfter &&
		toolName != tool.Add && toolName != tool.WatchUpdate &&
		toolName != tool.Complete && toolName != tool.Compact &&
		toolName != tool.SessionEvent {
		// Fire at threshold, then every governance.PersistNudgeRepeat
		// calls after.
		if ss.CallsSinceWrite == governance.PersistNudgeAfter ||
			(ss.CallsSinceWrite-governance.PersistNudgeAfter)%
				governance.PersistNudgeRepeat == 0 {
			warnings = append(warnings, fmt.Sprintf(
				desc.Text(text.DescKeyGovPersistNudge),
				ss.CallsSinceWrite))
		}
	}

	// 5. Violations from extension detection ring
	if violations := readAndClearViolations(d.ContextDir); len(violations) > 0 {
		for _, v := range violations {
			detail := v.Detail
			if len(detail) > cfgFmt.TruncateDetail {
				detail = detail[:cfgFmt.TruncateDetail] + token.Ellipsis
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
