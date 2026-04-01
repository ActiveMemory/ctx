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

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/mcp/governance"
	"github.com/ActiveMemory/ctx/internal/config/mcp/tool"
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxio "github.com/ActiveMemory/ctx/internal/io"
)

// violation represents a single governance violation recorded by the
// VS Code extension's detection ring.
type violation struct {
	Kind      string `json:"kind"`
	Detail    string `json:"detail"`
	Timestamp string `json:"timestamp"`
}

// violationsData is the JSON structure of the violations file.
type violationsData struct {
	Entries []violation `json:"entries"`
}

// readAndClearViolations reads violations from .context/state/violations.json
// and removes the file to prevent repeated escalation. Returns nil if
// no file exists or on read error.
func (ss *State) readAndClearViolations() []violation {
	if ss.contextDir == "" {
		return nil
	}
	stateDir := filepath.Join(ss.contextDir, dir.State)
	data, err := ctxio.SafeReadFile(stateDir, file.Violations)
	if err != nil {
		return nil
	}
	// Remove the file immediately to prevent duplicate alerts.
	_ = os.Remove(filepath.Join(stateDir, file.Violations))

	var vd violationsData
	if err := json.Unmarshal(data, &vd); err != nil {
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

// IncrementCallsSinceWrite bumps the counter used for persist nudges.
func (ss *State) IncrementCallsSinceWrite() {
	ss.callsSinceWrite++
}

// CheckGovernance returns governance warnings that should be appended
// to the current tool response. Returns an empty string when no action
// is warranted.
//
// The caller (toolName) is used to suppress redundant warnings — for
// example, a drift warning is not appended to a ctx_drift response.
func (ss *State) CheckGovernance(toolName string) string {
	var warnings []string

	// 1. Session not started
	if !ss.sessionStarted && toolName != tool.SessionEvent {
		warnings = append(warnings,
			"⚠ Session not started. "+
				"Call ctx_session_event(type=\"start\") to enable tracking.")
	}

	// 2. Context not loaded
	if !ss.contextLoaded && toolName != "ctx_status" &&
		toolName != tool.SessionEvent {
		warnings = append(warnings,
			"⚠ Context not loaded. "+
				"Call ctx_status() to load context before proceeding.")
	}

	// 3. Drift not checked recently
	if ss.sessionStarted && toolName != "ctx_drift" &&
		toolName != tool.SessionEvent {
		if !ss.lastDriftCheck.IsZero() {
			if time.Since(ss.lastDriftCheck) > governance.DriftCheckInterval {
				warnings = append(warnings, fmt.Sprintf(
					"⚠ Drift not checked in %d minutes. Consider calling ctx_drift().",
					int(time.Since(ss.lastDriftCheck).Minutes())))
			}
		} else if ss.ToolCalls > 5 {
			// Never checked drift and already 5+ calls in
			warnings = append(warnings,
				"⚠ Drift has not been checked this session. Consider calling ctx_drift().")
		}
	}

	// 4. Persist nudge — no context writes in a while
	if ss.sessionStarted && ss.callsSinceWrite >= governance.PersistNudgeAfter &&
		toolName != "ctx_add" && toolName != "ctx_watch_update" &&
		toolName != "ctx_complete" && toolName != "ctx_compact" &&
		toolName != tool.SessionEvent {
		// Fire at threshold, then every governance.PersistNudgeRepeat calls after
		if ss.callsSinceWrite == governance.PersistNudgeAfter ||
			(ss.callsSinceWrite-governance.PersistNudgeAfter)%governance.PersistNudgeRepeat == 0 {
			warnings = append(warnings, fmt.Sprintf(
				"⚠ %d tool calls since last context write. "+
					"Persist decisions, learnings, or completed tasks with ctx_add() or ctx_complete().",
				ss.callsSinceWrite))
		}
	}

	// 5. Violations from extension detection ring
	if violations := ss.readAndClearViolations(); len(violations) > 0 {
		for _, v := range violations {
			detail := v.Detail
			if len(detail) > 120 {
				detail = detail[:120] + "..."
			}
			warnings = append(warnings, fmt.Sprintf(
				"🚨 CRITICAL: %s — %s (at %s). "+
					"Review this action immediately. If unintended, revert it.",
				v.Kind, detail, v.Timestamp))
		}
	}

	if len(warnings) == 0 {
		return ""
	}

	nl := token.NewlineLF
	return nl + nl + "---" + nl + strings.Join(warnings, nl)
}
