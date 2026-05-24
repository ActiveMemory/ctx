//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package checkaudit_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkaudit"
	cfgAudit "github.com/ActiveMemory/ctx/internal/config/audit"
	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}

// TestRun_NoLeakInUninitializedProject mirrors the
// check-reminder invariant: hooks must not materialize
// `.context/` when called against an uninitialized project.
func TestRun_NoLeakInUninitializedProject(t *testing.T) {
	tempDir := t.TempDir()
	ctxDir := filepath.Join(tempDir, dir.Context)
	stateDir := filepath.Join(ctxDir, dir.State)

	t.Chdir(tempDir)
	rc.Reset()
	t.Cleanup(rc.Reset)

	r, w, pipeErr := os.Pipe()
	if pipeErr != nil {
		t.Fatalf("os.Pipe: %v", pipeErr)
	}
	go func() {
		defer func() { _ = w.Close() }()
		_, _ = io.Copy(w, bytes.NewReader([]byte(
			`{"session_id":"00000000-0000-0000-0000-000000000000"}`,
		)))
	}()
	t.Cleanup(func() { _ = r.Close() })

	cmd := &cobra.Command{}
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)

	if err := checkaudit.Run(cmd, r); err != nil {
		t.Fatalf("Run() error = %v, want nil", err)
	}
	if _, statErr := os.Stat(ctxDir); !os.IsNotExist(statErr) {
		t.Errorf(".context/ leaked: stat err = %v", statErr)
	}
	if _, statErr := os.Stat(stateDir); !os.IsNotExist(statErr) {
		t.Errorf(".context/state/ leaked: stat err = %v", statErr)
	}
}

// TestRun_SilentOnCleanReports verifies that a status:
// clean report does not produce a relay box.
func TestRun_SilentOnCleanReports(t *testing.T) {
	tempDir := setupInitializedProject(t)
	dropReport(t, tempDir, "surface", cfgAudit.StatusClean, "no findings")

	out, err := runHook(t)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	// Provenance is always emitted; absence of the audit
	// box title is the signal.
	if bytes.Contains(out, []byte("Audit Reports")) {
		t.Errorf(
			"clean report produced relay box; output:\n%s", out,
		)
	}
}

// TestRun_RelaysFindingsBody verifies that a findings
// report is wrapped in the verbatim-relay box with its
// body visible inside.
func TestRun_RelaysFindingsBody(t *testing.T) {
	tempDir := setupInitializedProject(t)
	body := "Surface drift detected: ctx pad undo missing from SKILL.md"
	dropReport(t, tempDir, "surface", cfgAudit.StatusFindings, body)

	out, err := runHook(t)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if !bytes.Contains(out, []byte("Audit Reports")) {
		t.Errorf("missing box title; output:\n%s", out)
	}
	if !bytes.Contains(out, []byte(body)) {
		t.Errorf("missing report body; output:\n%s", out)
	}
	if !bytes.Contains(out, []byte("ctx audit dismiss")) {
		t.Errorf("missing dismiss hint; output:\n%s", out)
	}
}

// TestRun_SilentWhenDismissed verifies that a dismissed
// report (against current digest) is suppressed.
func TestRun_SilentWhenDismissed(t *testing.T) {
	tempDir := setupInitializedProject(t)
	dropReport(t, tempDir, "surface", cfgAudit.StatusFindings, "body")

	// Write a dismissal ledger that matches the report's
	// digest (abc123 from the fixture).
	ledgerPath := filepath.Join(
		tempDir, dir.Context, cfgAudit.DirName,
		cfgAudit.DismissedLedger,
	)
	led := `{"entries":{"surface":{"digest":"abc123","at":"2026-05-24T16:00:00Z"}}}`
	if writeErr := os.WriteFile(ledgerPath, []byte(led), 0600); writeErr != nil {
		t.Fatal(writeErr)
	}

	out, err := runHook(t)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if bytes.Contains(out, []byte("Audit Reports")) {
		t.Errorf(
			"dismissed report still relayed; output:\n%s", out,
		)
	}
}

// TestRun_StalePrefixOnOldReport verifies that a report
// older than StalenessAge gets a STALE prefix in the box.
func TestRun_StalePrefixOnOldReport(t *testing.T) {
	tempDir := setupInitializedProject(t)
	// generated-at in dropReport's fixture is 2026-05-24;
	// today's date in CI may be later. To force staleness
	// independent of wall clock, drop a much older report.
	dropOldReport(t, tempDir, "surface", "ancient finding")

	out, err := runHook(t)
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if !bytes.Contains(out, []byte("STALE")) {
		t.Errorf(
			"old report missing STALE prefix; output:\n%s", out,
		)
	}
}

// setupInitializedProject creates a temp dir with the
// canonical required ctx files so [state.Initialized]
// returns true (the hook bails early otherwise).
func setupInitializedProject(t *testing.T) string {
	t.Helper()
	tempDir := t.TempDir()
	ctxDir := filepath.Join(tempDir, dir.Context)
	if mkErr := os.MkdirAll(ctxDir, 0750); mkErr != nil {
		t.Fatal(mkErr)
	}
	for _, fname := range cfgCtx.FilesRequired {
		if writeErr := os.WriteFile(
			filepath.Join(ctxDir, fname),
			[]byte("# placeholder\n"), 0600,
		); writeErr != nil {
			t.Fatal(writeErr)
		}
	}
	t.Chdir(tempDir)
	rc.Reset()
	t.Cleanup(rc.Reset)
	return tempDir
}

// dropReport writes a fresh-timestamp findings report.
func dropReport(
	t *testing.T, projectRoot, id, status, body string,
) {
	t.Helper()
	auditDir := filepath.Join(
		projectRoot, dir.Context, cfgAudit.DirName,
	)
	if mkErr := os.MkdirAll(auditDir, 0750); mkErr != nil {
		t.Fatal(mkErr)
	}
	// Use a generated-at far in the future relative to
	// production but stable for the test: matches the
	// audit_test dropReport fixture.
	content := "---\n" +
		"kind: " + id + "\nstatus: " + status + "\n" +
		"commit-range: main..HEAD\n" +
		"generated-at: 2099-05-24T14:30:12Z\n" +
		"generator: /ctx-surface-audit\n" +
		"digest: abc123\n---\n" + body + "\n"
	if writeErr := os.WriteFile(
		filepath.Join(auditDir, id+cfgAudit.ReportExt),
		[]byte(content), 0600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}
}

// dropOldReport writes a deliberately-stale findings
// report (year 1990 generated-at) so the StalenessAge
// branch fires.
func dropOldReport(
	t *testing.T, projectRoot, id, body string,
) {
	t.Helper()
	auditDir := filepath.Join(
		projectRoot, dir.Context, cfgAudit.DirName,
	)
	if mkErr := os.MkdirAll(auditDir, 0750); mkErr != nil {
		t.Fatal(mkErr)
	}
	content := "---\nkind: " + id + "\nstatus: findings\n" +
		"commit-range: main..ancient\n" +
		"generated-at: 1990-01-01T00:00:00Z\n" +
		"generator: /ctx-surface-audit\n" +
		"digest: old-digest\n---\n" + body + "\n"
	if writeErr := os.WriteFile(
		filepath.Join(auditDir, id+cfgAudit.ReportExt),
		[]byte(content), 0600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}
}

// runHook feeds a minimal hook envelope and captures the
// combined output.
func runHook(t *testing.T) ([]byte, error) {
	t.Helper()
	r, w, pipeErr := os.Pipe()
	if pipeErr != nil {
		t.Fatalf("os.Pipe: %v", pipeErr)
	}
	go func() {
		defer func() { _ = w.Close() }()
		_, _ = io.Copy(w, bytes.NewReader([]byte(
			`{"session_id":"00000000-0000-0000-0000-000000000000"}`,
		)))
	}()
	t.Cleanup(func() { _ = r.Close() })

	cmd := &cobra.Command{}
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	err := checkaudit.Run(cmd, r)
	return buf.Bytes(), err
}
