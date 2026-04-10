//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handler

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
	"github.com/ActiveMemory/ctx/internal/entity"
)

func newTestDeps() *entity.MCPDeps {
	return &entity.MCPDeps{
		ContextDir: "/tmp/test/.context",
		Session:    entity.NewMCPSession(),
	}
}

func TestCheckGovernance_SessionNotStarted(t *testing.T) {
	d := newTestDeps()
	got := CheckGovernance(d, "ctx_status")
	if !strings.Contains(got, "Session not started") {
		t.Errorf("expected session-not-started warning, got: %q", got)
	}
}

func TestCheckGovernance_SessionNotStarted_SuppressedForSessionEvent(t *testing.T) {
	d := newTestDeps()
	got := CheckGovernance(d, "ctx_session_event")
	if strings.Contains(got, "Session not started") {
		t.Errorf("session-not-started should be suppressed for ctx_session_event, got: %q", got)
	}
}

func TestCheckGovernance_ContextNotLoaded(t *testing.T) {
	d := newTestDeps()
	d.Session.RecordSessionStart()
	got := CheckGovernance(d, "ctx_add")
	if !strings.Contains(got, "Context not loaded") {
		t.Errorf("expected context-not-loaded warning, got: %q", got)
	}
}

func TestCheckGovernance_ContextNotLoaded_SuppressedForStatus(t *testing.T) {
	d := newTestDeps()
	d.Session.RecordSessionStart()
	got := CheckGovernance(d, "ctx_status")
	if strings.Contains(got, "Context not loaded") {
		t.Errorf("context-not-loaded should be suppressed for ctx_status, got: %q", got)
	}
}

func TestCheckGovernance_DriftNeverChecked(t *testing.T) {
	d := newTestDeps()
	d.Session.RecordSessionStart()
	d.Session.RecordContextLoaded()
	d.Session.ToolCalls = 6 // above the 5-call threshold

	got := CheckGovernance(d, "ctx_add")
	if !strings.Contains(got, "Drift has not been checked") {
		t.Errorf("expected drift-never-checked warning, got: %q", got)
	}
}

func TestCheckGovernance_DriftNeverChecked_BelowThreshold(t *testing.T) {
	d := newTestDeps()
	d.Session.RecordSessionStart()
	d.Session.RecordContextLoaded()
	d.Session.ToolCalls = 3 // below 5

	got := CheckGovernance(d, "ctx_add")
	if strings.Contains(got, "Drift") {
		t.Errorf("drift warning should not fire below 5 calls, got: %q", got)
	}
}

func TestCheckGovernance_DriftStale(t *testing.T) {
	d := newTestDeps()
	d.Session.RecordSessionStart()
	d.Session.RecordContextLoaded()
	d.Session.LastDriftCheck = time.Now().Add(-20 * time.Minute) // 20 min ago

	got := CheckGovernance(d, "ctx_add")
	if !strings.Contains(got, "Drift not checked in") {
		t.Errorf("expected stale-drift warning, got: %q", got)
	}
}

func TestCheckGovernance_DriftStale_SuppressedForDrift(t *testing.T) {
	d := newTestDeps()
	d.Session.RecordSessionStart()
	d.Session.RecordContextLoaded()
	d.Session.LastDriftCheck = time.Now().Add(-20 * time.Minute)

	got := CheckGovernance(d, "ctx_drift")
	if strings.Contains(got, "Drift") {
		t.Errorf("drift warning should be suppressed for ctx_drift, got: %q", got)
	}
}

func TestCheckGovernance_PersistNudge_AtThreshold(t *testing.T) {
	d := newTestDeps()
	d.Session.RecordSessionStart()
	d.Session.RecordContextLoaded()
	d.Session.RecordDriftCheck()
	d.Session.CallsSinceWrite = cfgGov.PersistNudgeAfter // exactly at threshold

	got := CheckGovernance(d, "ctx_status")
	if !strings.Contains(got, "tool calls since last context write") {
		t.Errorf("expected persist-nudge at threshold, got: %q", got)
	}
}

func TestCheckGovernance_PersistNudge_BelowThreshold(t *testing.T) {
	d := newTestDeps()
	d.Session.RecordSessionStart()
	d.Session.RecordContextLoaded()
	d.Session.RecordDriftCheck()
	d.Session.CallsSinceWrite = cfgGov.PersistNudgeAfter - 1

	got := CheckGovernance(d, "ctx_status")
	if strings.Contains(got, "tool calls since last context write") {
		t.Errorf("persist-nudge should not fire below threshold, got: %q", got)
	}
}

func TestCheckGovernance_PersistNudge_Repeat(t *testing.T) {
	d := newTestDeps()
	d.Session.RecordSessionStart()
	d.Session.RecordContextLoaded()
	d.Session.RecordDriftCheck()
	d.Session.CallsSinceWrite = cfgGov.PersistNudgeAfter + cfgGov.PersistNudgeRepeat

	got := CheckGovernance(d, "ctx_status")
	if !strings.Contains(got, "tool calls since last context write") {
		t.Errorf("expected persist-nudge at repeat interval, got: %q", got)
	}
}

func TestCheckGovernance_PersistNudge_SuppressedForWriteTools(t *testing.T) {
	d := newTestDeps()
	d.Session.RecordSessionStart()
	d.Session.CallsSinceWrite = cfgGov.PersistNudgeAfter

	for _, tool := range []string{"ctx_add", "ctx_complete", "ctx_watch_update", "ctx_compact"} {
		got := CheckGovernance(d, tool)
		if strings.Contains(got, "tool calls since last context write") {
			t.Errorf("persist-nudge should be suppressed for %s, got: %q", tool, got)
		}
	}
}

func TestCheckGovernance_NoWarnings(t *testing.T) {
	d := newTestDeps()
	d.Session.RecordSessionStart()
	d.Session.RecordContextLoaded()
	d.Session.RecordDriftCheck()
	d.Session.RecordContextWrite()

	got := CheckGovernance(d, "ctx_status")
	if got != "" {
		t.Errorf("expected no warnings, got: %q", got)
	}
}

func TestRecordSessionStart(t *testing.T) {
	d := newTestDeps()
	if d.Session.SessionStarted {
		t.Fatal("SessionStarted should be false initially")
	}
	d.Session.RecordSessionStart()
	if !d.Session.SessionStarted {
		t.Fatal("SessionStarted should be true after RecordSessionStart")
	}
}

func TestRecordContextWrite_ResetsCounter(t *testing.T) {
	d := newTestDeps()
	d.Session.CallsSinceWrite = 15
	d.Session.RecordContextWrite()
	if d.Session.CallsSinceWrite != 0 {
		t.Errorf("CallsSinceWrite should be 0 after RecordContextWrite, got %d", d.Session.CallsSinceWrite)
	}
}

func TestIncrementCallsSinceWrite(t *testing.T) {
	d := newTestDeps()
	d.Session.IncrementCallsSinceWrite()
	d.Session.IncrementCallsSinceWrite()
	d.Session.IncrementCallsSinceWrite()
	if d.Session.CallsSinceWrite != 3 {
		t.Errorf("expected 3, got %d", d.Session.CallsSinceWrite)
	}
}

func TestCheckGovernance_WarningFormat(t *testing.T) {
	d := newTestDeps()
	got := CheckGovernance(d, "ctx_add")
	if got != "" && !strings.HasPrefix(got, "\n\n---\n") {
		t.Errorf("warnings should start with separator, got: %q", got)
	}
}

func newTestDepsWithDir(t *testing.T) *entity.MCPDeps {
	t.Helper()
	contextDir := filepath.Join(t.TempDir(), ".context")
	if err := os.MkdirAll(filepath.Join(contextDir, dir.State), 0o755); err != nil {
		t.Fatal(err)
	}
	return &entity.MCPDeps{
		ContextDir: contextDir,
		Session:    entity.NewMCPSession(),
	}
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
	d := newTestDepsWithDir(t)
	d.Session.RecordSessionStart()
	d.Session.RecordContextLoaded()
	d.Session.RecordDriftCheck()
	d.Session.RecordContextWrite()

	writeViolations(t, d.ContextDir, []violation{
		{Kind: "dangerous_command", Detail: "sudo rm -rf /tmp", Timestamp: "2026-03-17T10:00:00Z"},
	})

	got := CheckGovernance(d, "ctx_status")
	if !strings.Contains(got, "CRITICAL") {
		t.Errorf("expected CRITICAL warning, got: %q", got)
	}
	if !strings.Contains(got, "dangerous_command") {
		t.Errorf("expected violation kind in warning, got: %q", got)
	}
}

func TestCheckGovernance_ViolationsFileRemovedAfterRead(t *testing.T) {
	d := newTestDepsWithDir(t)
	writeViolations(t, d.ContextDir, []violation{
		{Kind: "sensitive_file_read", Detail: ".env", Timestamp: "2026-03-17T10:00:00Z"},
	})

	p := filepath.Join(d.ContextDir, dir.State, file.Violations)
	if _, err := os.Stat(p); err != nil {
		t.Fatal("violations file should exist before read")
	}

	CheckGovernance(d, "ctx_status")

	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Error("violations file should be removed after read")
	}
}

func TestCheckGovernance_NoViolationsFile(t *testing.T) {
	d := newTestDepsWithDir(t)
	d.Session.RecordSessionStart()
	d.Session.RecordContextLoaded()
	d.Session.RecordDriftCheck()
	d.Session.RecordContextWrite()

	got := CheckGovernance(d, "ctx_status")
	if strings.Contains(got, "CRITICAL") {
		t.Errorf("no violations should mean no CRITICAL warning, got: %q", got)
	}
}

func TestCheckGovernance_ViolationDetailTruncated(t *testing.T) {
	d := newTestDepsWithDir(t)
	d.Session.RecordSessionStart()
	d.Session.RecordContextLoaded()
	d.Session.RecordDriftCheck()
	d.Session.RecordContextWrite()

	longDetail := strings.Repeat("x", 200)
	writeViolations(t, d.ContextDir, []violation{
		{Kind: "hack_script", Detail: longDetail, Timestamp: "2026-03-17T10:00:00Z"},
	})

	got := CheckGovernance(d, "ctx_status")
	if strings.Contains(got, longDetail) {
		t.Error("full 200-char detail should be truncated")
	}
	if !strings.Contains(got, "...") {
		t.Errorf("truncated detail should contain ellipsis, got: %q", got)
	}
}

func TestCheckGovernance_MultipleViolations(t *testing.T) {
	d := newTestDepsWithDir(t)
	d.Session.RecordSessionStart()
	d.Session.RecordContextLoaded()
	d.Session.RecordDriftCheck()
	d.Session.RecordContextWrite()

	writeViolations(t, d.ContextDir, []violation{
		{Kind: "dangerous_command", Detail: "git push --force", Timestamp: "2026-03-17T10:00:00Z"},
		{Kind: "sensitive_file_read", Detail: ".env.local", Timestamp: "2026-03-17T10:00:01Z"},
	})

	got := CheckGovernance(d, "ctx_status")
	count := strings.Count(got, "CRITICAL")
	if count != 2 {
		t.Errorf("expected 2 CRITICAL warnings, got %d in: %q", count, got)
	}
}

func TestReadAndClearViolations_EmptyContextDir(t *testing.T) {
	violations := readAndClearViolations("")
	if violations != nil {
		t.Errorf("expected nil for empty contextDir, got: %v", violations)
	}
}

func TestReadAndClearViolations_CorruptFile(t *testing.T) {
	d := newTestDepsWithDir(t)
	p := filepath.Join(d.ContextDir, dir.State, file.Violations)
	if err := os.WriteFile(p, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	violations := readAndClearViolations(d.ContextDir)
	if violations != nil {
		t.Errorf("expected nil for corrupt file, got: %v", violations)
	}
}
