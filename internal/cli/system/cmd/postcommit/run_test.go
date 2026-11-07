//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package postcommit

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func writeInitializedProject(t *testing.T, tempDir string) string {
	t.Helper()
	ctxDir := filepath.Join(tempDir, dir.Context)
	if mkErr := os.MkdirAll(ctxDir, 0o750); mkErr != nil {
		t.Fatalf("MkdirAll(%s): %v", ctxDir, mkErr)
	}
	for _, f := range cfgCtx.FilesRequired {
		p := filepath.Join(ctxDir, f)
		if wErr := os.WriteFile(p, []byte("# "+f+"\n"), 0o600); wErr != nil {
			t.Fatalf("WriteFile(%s): %v", p, wErr)
		}
	}
	return ctxDir
}

func writeTempFile(t *testing.T, content string) *os.File {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "hook-stdin-*.json")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("WriteString: %v", err)
	}
	if _, err := f.Seek(0, 0); err != nil {
		t.Fatalf("Seek: %v", err)
	}
	return f
}

func TestRun_AgyToolErrorSuppression(t *testing.T) {
	origWd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(origWd)
		rc.Reset()
		coreSession.SetAgyMode(false, "", "", "")
	}()

	workspaceDir := t.TempDir()
	writeInitializedProject(t, workspaceDir)

	// Tool call failed: error is non-empty
	errorJSON := `{
		"conversationId": "sess-post-err",
		"workspacePaths": ["` + filepath.ToSlash(workspaceDir) + `"],
		"error": "command failed with exit code 1",
		"toolCall": {
			"name": "run_command",
			"args": {
				"CommandLine": "git commit -m 'feat: something'"
			}
		}
	}`
	errorFile := writeTempFile(t, errorJSON)
	defer func() { _ = errorFile.Close() }()

	cmd := &cobra.Command{}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := Run(cmd, errorFile)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	out := strings.TrimSpace(buf.String())
	if out != "{}" {
		t.Errorf("expected empty JSON {} on error suppression, got %q", out)
	}

	// Verify nothing was queued
	pending := coreSession.DrainPendingAgy()
	if len(pending) != 0 {
		t.Errorf("expected no pending nudges on tool error, got: %+v", pending)
	}
}

func TestRun_AgySuccessQueuesNudge(t *testing.T) {
	origWd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(origWd)
		rc.Reset()
		coreSession.SetAgyMode(false, "", "", "")
	}()

	workspaceDir := t.TempDir()
	writeInitializedProject(t, workspaceDir)

	// Tool call succeeded: error is empty, command is git commit
	successJSON := `{
		"conversationId": "sess-post-ok",
		"workspacePaths": ["` + filepath.ToSlash(workspaceDir) + `"],
		"error": "",
		"toolCall": {
			"name": "run_command",
			"args": {
				"CommandLine": "git commit -m 'feat: something new'"
			}
		}
	}`
	successFile := writeTempFile(t, successJSON)
	defer func() { _ = successFile.Close() }()

	cmd := &cobra.Command{}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := Run(cmd, successFile)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	out := strings.TrimSpace(buf.String())
	if out != "{}" {
		t.Errorf("expected empty JSON {} on stdout, got %q", out)
	}

	// Verify nudge was queued
	pending := coreSession.DrainPendingAgy()
	if len(pending) == 0 {
		t.Error("expected nudge to be queued in pending queue for PostInvocation")
	}
}
