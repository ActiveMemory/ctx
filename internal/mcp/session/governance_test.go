//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/file"
	cfgGov "github.com/ActiveMemory/ctx/internal/config/mcp/governance"
)

func newTestState() *State {
	return NewState("/tmp/test/.context")
}

func TestCheckGovernance_SessionNotStarted(t *testing.T) {
	ss := newTestState()
	got := ss.CheckGovernance("ctx_status")
	if !strings.Contains(got, "Session not started") {
		t.Errorf("expected session-not-started warning, got: %q", got)
	}
}

func TestCheckGovernance_SessionNotStarted_SuppressedForSessionEvent(t *testing.T) {
	ss := newTestState()
	got := ss.CheckGovernance("ctx_session_event")
	if strings.Contains(got, "Session not started") {
		t.Errorf("session-not-started should be suppressed for ctx_session_event, got: %q", got)
	}
}

func TestCheckGovernance_ContextNotLoaded(t *testing.T) {
	ss := newTestState()
	ss.RecordSessionStart()
	got := ss.CheckGovernance("ctx_add")
	if !strings.Contains(got, "Context not loaded") {
		t.Errorf("expected context-not-loaded warning, got: %q", got)
	}
}

func TestCheckGovernance_ContextNotLoaded_SuppressedForStatus(t *testing.T) {
	ss := newTestState()
	ss.RecordSessionStart()
	got := ss.CheckGovernance("ctx_status")
	if strings.Contains(got, "Context not loaded") {
		t.Errorf("context-not-loaded should be suppressed for ctx_status, got: %q", got)
	}
}

func TestCheckGovernance_DriftNeverChecked(t *testing.T) {
	ss := newTestState()
	ss.RecordSessionStart()
	ss.RecordContextLoaded()
	ss.ToolCalls = 6 // above the 5-call threshold

	got := ss.CheckGovernance("ctx_add")
	if !strings.Contains(got, "Drift has not been checked") {
		t.Errorf("expected drift-never-checked warning, got: %q", got)
	}
}

func TestCheckGovernance_DriftNeverChecked_BelowThreshold(t *testing.T) {
	ss := newTestState()
	ss.RecordSessionStart()
	ss.RecordContextLoaded()
	ss.ToolCalls = 3 // below 5

	got := ss.CheckGovernance("ctx_add")
	if strings.Contains(got, "Drift") {
		t.Errorf("drift warning should not fire below 5 calls, got: %q", got)
	}
}

func TestCheckGovernance_DriftStale(t *testing.T) {
	ss := newTestState()
	ss.RecordSessionStart()
	ss.RecordContextLoaded()
	ss.lastDriftCheck = time.Now().Add(-20 * time.Minute) // 20 min ago

	got := ss.CheckGovernance("ctx_add")
	if !strings.Contains(got, "Drift not checked in") {
		t.Errorf("expected stale-drift warning, got: %q", got)
	}
}

func TestCheckGovernance_DriftStale_SuppressedForDrift(t *testing.T) {
	ss := newTestState()
	ss.RecordSessionStart()
	ss.RecordContextLoaded()
	ss.lastDriftCheck = time.Now().Add(-20 * time.Minute)

	got := ss.CheckGovernance("ctx_drift")
	if strings.Contains(got, "Drift") {
		t.Errorf("drift warning should be suppressed for ctx_drift, got: %q", got)
	}
}

func TestCheckGovernance_PersistNudge_AtThreshold(t *testing.T) {
	ss := newTestState()
	ss.RecordSessionStart()
	ss.RecordContextLoaded()
	ss.RecordDriftCheck()
	ss.callsSinceWrite = cfgGov.PersistNudgeAfter // exactly at threshold

	got := ss.CheckGovernance("ctx_status")
	if !strings.Contains(got, "tool calls since last context write") {
		t.Errorf("expected persist-nudge at threshold, got: %q", got)
	}
}

func TestCheckGovernance_PersistNudge_BelowThreshold(t *testing.T) {
	ss := newTestState()
	ss.RecordSessionStart()
	ss.RecordContextLoaded()
	ss.RecordDriftCheck()
	ss.callsSinceWrite = cfgGov.PersistNudgeAfter - 1

	got := ss.CheckGovernance("ctx_status")
	if strings.Contains(got, "tool calls since last context write") {
		t.Errorf("persist-nudge should not fire below threshold, got: %q", got)
	}
}

func TestCheckGovernance_PersistNudge_Repeat(t *testing.T) {
	ss := newTestState()
	ss.RecordSessionStart()
	ss.RecordContextLoaded()
	ss.RecordDriftCheck()
	ss.callsSinceWrite = cfgGov.PersistNudgeAfter + cfgGov.PersistNudgeRepeat

	got := ss.CheckGovernance("ctx_status")
	if !strings.Contains(got, "tool calls since last context write") {
		t.Errorf("expected persist-nudge at repeat interval, got: %q", got)
	}
}

func TestCheckGovernance_PersistNudge_SuppressedForWriteTools(t *testing.T) {
	ss := newTestState()
	ss.RecordSessionStart()
	ss.callsSinceWrite = cfgGov.PersistNudgeAfter

	for _, tool := range []string{"ctx_add", "ctx_complete", "ctx_watch_update", "ctx_compact"} {
		got := ss.CheckGovernance(tool)
		if strings.Contains(got, "tool calls since last context write") {
			t.Errorf("persist-nudge should be suppressed for %s, got: %q", tool, got)
		}
	}
}

func TestCheckGovernance_NoWarnings(t *testing.T) {
	ss := newTestState()
	ss.RecordSessionStart()
	ss.RecordContextLoaded()
	ss.RecordDriftCheck()
	ss.RecordContextWrite()

	got := ss.CheckGovernance("ctx_status")
	if got != "" {
		t.Errorf("expected no warnings, got: %q", got)
	}
}

func TestRecordSessionStart(t *testing.T) {
	ss := newTestState()
	if ss.sessionStarted {
		t.Fatal("sessionStarted should be false initially")
	}
	ss.RecordSessionStart()
	if !ss.sessionStarted {
		t.Fatal("sessionStarted should be true after RecordSessionStart")
	}
}

func TestRecordContextWrite_ResetsCounter(t *testing.T) {
	ss := newTestState()
	ss.callsSinceWrite = 15
	ss.RecordContextWrite()
	if ss.callsSinceWrite != 0 {
		t.Errorf("callsSinceWrite should be 0 after RecordContextWrite, got %d", ss.callsSinceWrite)
	}
}

func TestIncrementCallsSinceWrite(t *testing.T) {
	ss := newTestState()
	ss.IncrementCallsSinceWrite()
	ss.IncrementCallsSinceWrite()
	ss.IncrementCallsSinceWrite()
	if ss.callsSinceWrite != 3 {
		t.Errorf("expected 3, got %d", ss.callsSinceWrite)
	}
}

func TestCheckGovernance_WarningFormat(t *testing.T) {
	ss := newTestState()
	got := ss.CheckGovernance("ctx_add")
	if got != "" && !strings.HasPrefix(got, "\n\n---\n") {
		t.Errorf("warnings should start with separator, got: %q", got)
	}
}

func newTestStateWithDir(t *testing.T) *State {
	t.Helper()
	contextDir := filepath.Join(t.TempDir(), ".context")
	if err := os.MkdirAll(filepath.Join(contextDir, dir.State), 0o755); err != nil {
		t.Fatal(err)
	}
	return NewState(contextDir)
}

func writeViolations(t *testing.T, contextDir string, entries []violation) {
	t.Helper()
	data, err := json.Marshal(violationsData{Entries: entries})
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(contextDir, dir.State, file.Violations)
	if err := os.WriteFile(p, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestCheckGovernance_ViolationsDetected(t *testing.T) {
	ss := newTestStateWithDir(t)
	ss.RecordSessionStart()
	ss.RecordContextLoaded()
	ss.RecordDriftCheck()
	ss.RecordContextWrite()

	writeViolations(t, ss.contextDir, []violation{
		{Kind: "dangerous_command", Detail: "sudo rm -rf /tmp", Timestamp: "2026-03-17T10:00:00Z"},
	})

	got := ss.CheckGovernance("ctx_status")
	if !strings.Contains(got, "CRITICAL") {
		t.Errorf("expected CRITICAL warning, got: %q", got)
	}
	if !strings.Contains(got, "dangerous_command") {
		t.Errorf("expected violation kind in warning, got: %q", got)
	}
}

func TestCheckGovernance_ViolationsFileRemovedAfterRead(t *testing.T) {
	ss := newTestStateWithDir(t)
	writeViolations(t, ss.contextDir, []violation{
		{Kind: "sensitive_file_read", Detail: ".env", Timestamp: "2026-03-17T10:00:00Z"},
	})

	p := filepath.Join(ss.contextDir, dir.State, file.Violations)
	if _, err := os.Stat(p); err != nil {
		t.Fatal("violations file should exist before read")
	}

	ss.CheckGovernance("ctx_status")

	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Error("violations file should be removed after read")
	}
}

func TestCheckGovernance_NoViolationsFile(t *testing.T) {
	ss := newTestStateWithDir(t)
	ss.RecordSessionStart()
	ss.RecordContextLoaded()
	ss.RecordDriftCheck()
	ss.RecordContextWrite()

	got := ss.CheckGovernance("ctx_status")
	if strings.Contains(got, "CRITICAL") {
		t.Errorf("no violations should mean no CRITICAL warning, got: %q", got)
	}
}

func TestCheckGovernance_ViolationDetailTruncated(t *testing.T) {
	ss := newTestStateWithDir(t)
	ss.RecordSessionStart()
	ss.RecordContextLoaded()
	ss.RecordDriftCheck()
	ss.RecordContextWrite()

	longDetail := strings.Repeat("x", 200)
	writeViolations(t, ss.contextDir, []violation{
		{Kind: "hack_script", Detail: longDetail, Timestamp: "2026-03-17T10:00:00Z"},
	})

	got := ss.CheckGovernance("ctx_status")
	if strings.Contains(got, longDetail) {
		t.Error("full 200-char detail should be truncated")
	}
	if !strings.Contains(got, "...") {
		t.Errorf("truncated detail should contain ellipsis, got: %q", got)
	}
}

func TestCheckGovernance_MultipleViolations(t *testing.T) {
	ss := newTestStateWithDir(t)
	ss.RecordSessionStart()
	ss.RecordContextLoaded()
	ss.RecordDriftCheck()
	ss.RecordContextWrite()

	writeViolations(t, ss.contextDir, []violation{
		{Kind: "dangerous_command", Detail: "git push --force", Timestamp: "2026-03-17T10:00:00Z"},
		{Kind: "sensitive_file_read", Detail: ".env.local", Timestamp: "2026-03-17T10:00:01Z"},
	})

	got := ss.CheckGovernance("ctx_status")
	count := strings.Count(got, "CRITICAL")
	if count != 2 {
		t.Errorf("expected 2 CRITICAL warnings, got %d in: %q", count, got)
	}
}

func TestReadAndClearViolations_EmptyContextDir(t *testing.T) {
	ss := &State{contextDir: ""}
	violations := readAndClearViolations(ss.contextDir)
	if violations != nil {
		t.Errorf("expected nil for empty contextDir, got: %v", violations)
	}
}

func TestReadAndClearViolations_CorruptFile(t *testing.T) {
	ss := newTestStateWithDir(t)
	p := filepath.Join(ss.contextDir, dir.State, file.Violations)
	if err := os.WriteFile(p, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	violations := readAndClearViolations(ss.contextDir)
	if violations != nil {
		t.Errorf("expected nil for corrupt file, got: %v", violations)
	}
}
