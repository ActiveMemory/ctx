//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/ctxctl/cli/audit"
	"github.com/ActiveMemory/ctx/internal/ctxctl/cli/audit/core/store"
	cfgAudit "github.com/ActiveMemory/ctx/internal/ctxctl/config/audit"
	errAudit "github.com/ActiveMemory/ctx/internal/ctxctl/err/audit"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// setup creates a temp dir with a .context/ directory and resets
// rc so the audit dir resolves to the temp tree.
func setup(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	origDir, getErr := os.Getwd()
	if getErr != nil {
		t.Fatal(getErr)
	}
	if chErr := os.Chdir(tmpDir); chErr != nil {
		t.Fatal(chErr)
	}
	t.Cleanup(func() {
		if chErr := os.Chdir(origDir); chErr != nil {
			t.Error(chErr)
		}
		rc.Reset()
	})
	ctxDir := filepath.Join(tmpDir, dir.Context)
	if mkErr := os.MkdirAll(ctxDir, 0750); mkErr != nil {
		t.Fatal(mkErr)
	}
	rc.Reset()
	return tmpDir
}

// dropReport writes a fixture audit report at
// .context/audit/<id>.md with frontmatter + body.
func dropReport(
	t *testing.T, projectRoot, id, status, body string,
) {
	t.Helper()
	dirPath := filepath.Join(
		projectRoot, dir.Context, cfgAudit.DirName,
	)
	if mkErr := os.MkdirAll(dirPath, 0750); mkErr != nil {
		t.Fatal(mkErr)
	}
	content := "---" + "\n" +
		"kind: " + id + "\n" +
		"status: " + status + "\n" +
		"commit-range: main..HEAD\n" +
		"generated-at: 2026-05-24T14:30:12Z\n" +
		"generator: /ctx-surface-audit\n" +
		"digest: abc123\n" +
		"---" + "\n" +
		body + "\n"
	if writeErr := os.WriteFile(
		filepath.Join(dirPath, id+cfgAudit.ReportExt),
		[]byte(content), 0600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}
}

// runCmd executes a cobra command and captures combined output.
func runCmd(cmd *cobra.Command) (string, error) {
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	err := cmd.Execute()
	return buf.String(), err
}

// newAuditCmd builds a fresh audit command (wired with the
// shipped English strings) with the given args.
func newAuditCmd(args ...string) *cobra.Command {
	c := audit.Cmd(auditStrings())
	c.SetArgs(args)
	return c
}

// runRoot executes ctxctl's real root command (the same wiring
// main() uses), capturing combined stdout+stderr and applying
// printErr on error exactly as main() does. This exercises the
// SilenceErrors + sole-printer path end to end, so error-path
// assertions see the user-facing English rendered from the typed
// errors.
func runRoot(args ...string) (string, error) {
	var buf bytes.Buffer
	root := newRoot()
	root.SetArgs(args)
	root.SetOut(&buf)
	root.SetErr(&buf)
	err := root.Execute()
	if err != nil {
		printErr(root, err)
	}
	return buf.String(), err
}

func TestList_EmptyDir(t *testing.T) {
	setup(t)

	out, err := runCmd(newAuditCmd("list"))
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if !strings.Contains(out, "No audit reports.") {
		t.Errorf("output = %q, want 'No audit reports.'", out)
	}
}

func TestList_ShowsFindings(t *testing.T) {
	tmp := setup(t)
	dropReport(t, tmp, "surface", cfgAudit.StatusFindings,
		"Surface audit body line 1\nSurface audit body line 2",
	)

	out, err := runCmd(newAuditCmd("list"))
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if !strings.Contains(out, "surface") {
		t.Errorf("output = %q, missing report id", out)
	}
	if !strings.Contains(out, "findings") {
		t.Errorf("output = %q, missing status", out)
	}
	if !strings.Contains(out, "main..HEAD") {
		t.Errorf("output = %q, missing commit-range", out)
	}
}

func TestShow_PrintsBodyVerbatim(t *testing.T) {
	tmp := setup(t)
	body := "line one\nline two\nline three"
	dropReport(t, tmp, "surface", cfgAudit.StatusFindings, body)

	out, err := runCmd(newAuditCmd("show", "surface"))
	if err != nil {
		t.Fatalf("show error: %v", err)
	}
	if !strings.Contains(out, body) {
		t.Errorf("output = %q, want body %q", out, body)
	}
}

func TestShow_UnknownID(t *testing.T) {
	setup(t)
	// Drive the real root so the typed error renders through
	// printErr in ctxctl's English voice (the format strings now
	// live in tools/ctxctl, not the err package).
	out, err := runRoot("audit", "show", "ghost")
	if err == nil {
		t.Fatal("expected error for unknown id")
	}
	if _, ok := errors.AsType[*errAudit.UnknownIDError](err); !ok {
		t.Errorf("err type = %T, want *errAudit.UnknownIDError", err)
	}
	if !strings.Contains(out, "unknown audit id: ghost") {
		t.Errorf("output = %q, want 'unknown audit id: ghost'", out)
	}
}

func TestShow_MalformedReport(t *testing.T) {
	tmp := setup(t)
	// A report file with no YAML frontmatter delimiter. loadOne
	// wraps parse.ErrNoFrontmatter in a *ParseReportError; printErr
	// must resolve the wrapped sentinel to its English prose
	// (auditCause), not leak the diagnostic code.
	dirPath := filepath.Join(tmp, dir.Context, cfgAudit.DirName)
	if mkErr := os.MkdirAll(dirPath, 0750); mkErr != nil {
		t.Fatal(mkErr)
	}
	if writeErr := os.WriteFile(
		filepath.Join(dirPath, "surface"+cfgAudit.ReportExt),
		[]byte("no frontmatter here\n"), 0600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	out, err := runRoot("audit", "show", "surface")
	if err == nil {
		t.Fatal("expected parse error for malformed report")
	}
	if _, ok := errors.AsType[*errAudit.ParseReportError](err); !ok {
		t.Errorf("err type = %T, want *errAudit.ParseReportError", err)
	}
	if !errors.Is(err, errAudit.ErrNoFrontmatter) {
		t.Errorf("err = %v, want wrapped ErrNoFrontmatter", err)
	}
	want := "parse audit report surface: " +
		"audit report missing yaml frontmatter"
	if !strings.Contains(out, want) {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestDismiss_StopsRelay(t *testing.T) {
	tmp := setup(t)
	dropReport(t, tmp, "surface", cfgAudit.StatusFindings, "body")

	if _, err := runCmd(newAuditCmd("dismiss", "surface")); err != nil {
		t.Fatalf("dismiss error: %v", err)
	}

	// Re-read; IsDismissed should now be true.
	reports, readErr := store.Read()
	if readErr != nil {
		t.Fatalf("store.Read: %v", readErr)
	}
	led, ledErr := store.ReadDismissals()
	if ledErr != nil {
		t.Fatalf("store.ReadDismissals: %v", ledErr)
	}
	if len(reports) != 1 {
		t.Fatalf("got %d reports, want 1", len(reports))
	}
	if !store.IsDismissed(reports[0], led) {
		t.Error("dismissed report not marked in ledger")
	}

	// Ledger file should exist on disk too.
	path := filepath.Join(
		tmp, dir.Context, cfgAudit.DirName, cfgAudit.DismissedLedger,
	)
	data, readErr := os.ReadFile(path) //nolint:gosec // test path
	if readErr != nil {
		t.Fatalf("read ledger: %v", readErr)
	}
	var led2 store.DismissalLedger
	if parseErr := json.Unmarshal(data, &led2); parseErr != nil {
		t.Fatalf("parse ledger: %v", parseErr)
	}
	if _, ok := led2.Entries["surface"]; !ok {
		t.Error("ledger missing surface entry")
	}
}

func TestDismissAll_DismissesEvery(t *testing.T) {
	tmp := setup(t)
	dropReport(t, tmp, "surface", cfgAudit.StatusFindings, "a")
	dropReport(t, tmp, "specs", cfgAudit.StatusFindings, "b")
	dropReport(t, tmp, "capture", cfgAudit.StatusFindings, "c")

	out, err := runCmd(newAuditCmd("dismiss", "--all"))
	if err != nil {
		t.Fatalf("dismiss --all error: %v", err)
	}
	if !strings.Contains(out, "Dismissed 3 audit") {
		t.Errorf("output = %q, want '3 audit'", out)
	}

	reports, readErr := store.Read()
	if readErr != nil {
		t.Fatalf("store.Read: %v", readErr)
	}
	led, ledErr := store.ReadDismissals()
	if ledErr != nil {
		t.Fatalf("store.ReadDismissals: %v", ledErr)
	}
	for _, r := range reports {
		if !store.IsDismissed(r, led) {
			t.Errorf("report %s not dismissed", r.ID)
		}
	}
}

func TestDismiss_NoIDsNoAll_Errors(t *testing.T) {
	setup(t)
	out, err := runRoot("audit", "dismiss")
	if err == nil {
		t.Fatal("expected error when no ids and no --all")
	}
	if !errors.Is(err, errAudit.ErrIDRequired) {
		t.Errorf("err = %v, want ErrIDRequired", err)
	}
	if !strings.Contains(out, "audit id required") {
		t.Errorf("output = %q, want 'audit id required'", out)
	}
}

func TestDismiss_FreshDigestResurfaces(t *testing.T) {
	tmp := setup(t)
	dropReport(t, tmp, "surface", cfgAudit.StatusFindings, "body v1")

	// Dismiss against the v1 digest (abc123 per dropReport fixture).
	if _, err := runCmd(newAuditCmd("dismiss", "surface")); err != nil {
		t.Fatalf("dismiss error: %v", err)
	}

	// Drop a new report with a different digest by hand-writing
	// a fixture with a changed digest field.
	newContent := "---" + "\n" +
		"kind: surface\nstatus: findings\n" +
		"commit-range: main..HEAD\n" +
		"generated-at: 2026-05-24T15:00:00Z\n" +
		"generator: /ctx-surface-audit\n" +
		"digest: NEW-digest-xyz\n" +
		"---\nbody v2\n"
	path := filepath.Join(
		tmp, dir.Context, cfgAudit.DirName,
		"surface"+cfgAudit.ReportExt,
	)
	if writeErr := os.WriteFile(path, []byte(newContent), 0600); writeErr != nil {
		t.Fatal(writeErr)
	}

	reports, readErr := store.Read()
	if readErr != nil {
		t.Fatalf("store.Read: %v", readErr)
	}
	led, ledErr := store.ReadDismissals()
	if ledErr != nil {
		t.Fatalf("store.ReadDismissals: %v", ledErr)
	}
	if len(reports) != 1 {
		t.Fatalf("got %d reports, want 1", len(reports))
	}
	if store.IsDismissed(reports[0], led) {
		t.Error(
			"fresh-digest report should re-surface (not be dismissed)",
		)
	}
}
