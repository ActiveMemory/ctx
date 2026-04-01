//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/mcp/governance"
	"github.com/ActiveMemory/ctx/internal/config/mcp/tool"
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxio "github.com/ActiveMemory/ctx/internal/io"
)

// violation represents a single governance violation recorded by the
// VS Code extension's detection ring.
//
// Fields:
//   - Kind: violation category identifier
//   - Detail: human-readable description of what was violated
//   - Timestamp: ISO-8601 timestamp of when the violation occurred
type violation struct {
	Kind      string `json:"kind"`
	Detail    string `json:"detail"`
	Timestamp string `json:"timestamp"`
}

// violationsData is the JSON structure of the violations file.
//
// Fields:
//   - Entries: list of recorded violations
type violationsData struct {
	Entries []violation `json:"entries"`
}

// readAndClearViolations reads violations from
// .context/state/violations.json and removes the file to prevent
// repeated escalation.
//
// Returns:
//   - []violation: parsed violations, or nil if no file exists or
//     on read error
func (ss *State) readAndClearViolations() []violation {
	if ss.contextDir == "" {
		return nil
	}
	stateDir := filepath.Join(ss.contextDir, dir.State)
	data, readErr := ctxio.SafeReadFile(stateDir, file.Violations)
	if readErr != nil {
		return nil
	}
	// Remove the file immediately to prevent duplicate alerts.
	_ = os.Remove(filepath.Join(stateDir, file.Violations))

	var vd violationsData
	if unmarshalErr := json.Unmarshal(data, &vd); unmarshalErr != nil {
		return nil
	}
	return vd.Entries
}

// RecordSessionStart marks the session as explicitly started.
func (ss *State) RecordSessionStart() {
	ss.sessionStarted = true
	ss.sessionStartedAt = time.Now()
}

// RecordContextLoaded marks context as loaded for this session.
func (ss *State) RecordContextLoaded() {
	ss.contextLoaded = true
}

// RecordDriftCheck records that a drift check was performed.
func (ss *State) RecordDriftCheck() {
	ss.lastDriftCheck = time.Now()
}

// RecordContextWrite records that a .context/ write occurred (add,
// complete, watch_update, compact).
func (ss *State) RecordContextWrite() {
	ss.lastContextWrite = time.Now()
	ss.callsSinceWrite = 0
}

// IncrementCallsSinceWrite bumps the counter used for persist
// nudges.
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
	if violations := ss.readAndClearViolations(); len(violations) > 0 {
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
