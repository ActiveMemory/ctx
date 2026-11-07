//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

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

func TestAgyPreambleAndCWDAnchor(t *testing.T) {
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	defer func() {
		_ = os.Chdir(origWd)
		rc.Reset()
		coreSession.SetAgyMode(false, "", "", "")
	}()

	// Outer directory where hook command was ostensibly spawned
	spawnDir := t.TempDir()
	_ = os.Chdir(spawnDir)

	// Target workspace with initialized .context
	workspaceDir := t.TempDir()
	writeInitializedProject(t, workspaceDir)

	stdinPayload := `{
		"conversationId": "sess-anchor-test",
		"workspacePaths": ["` + filepath.ToSlash(workspaceDir) + `"],
		"transcriptPath": "/logs/transcript.jsonl",
		"modelName": "gemini-2.0-flash",
		"lastUserInput": "test prompt",
		"toolCall": {
			"name": "run_command",
			"args": {
				"CommandLine": "echo hello"
			}
		}
	}`

	stdinFile := writeTempFile(t, stdinPayload)
	defer stdinFile.Close()

	input, sessionID, ctxDir, stateDir, ok := FullPreamble(stdinFile)
	if !ok {
		t.Fatal("expected FullPreamble to succeed after CWD anchor to workspace")
	}

	if sessionID != "sess-anchor-test" {
		t.Errorf("sessionID = %q, want %q", sessionID, "sess-anchor-test")
	}
	if input.ToolName != "run_command" {
		t.Errorf("input.ToolName = %q, want %q", input.ToolName, "run_command")
	}
	if input.ToolInput.Command != "echo hello" {
		t.Errorf(
			"input.ToolInput.Command = %q, want %q",
			input.ToolInput.Command, "echo hello",
		)
	}

	currentWd, _ := os.Getwd()
	currentWdCanonical, _ := filepath.EvalSymlinks(currentWd)
	workspaceCanonical, _ := filepath.EvalSymlinks(workspaceDir)
	if !strings.EqualFold(currentWdCanonical, workspaceCanonical) {
		t.Errorf(
			"cwd = %q, want workspace %q",
			currentWdCanonical, workspaceCanonical,
		)
	}

	if ctxDir == "" || stateDir == "" {
		t.Errorf(
			"expected non-empty ctxDir and stateDir, got ctxDir=%q stateDir=%q",
			ctxDir, stateDir,
		)
	}
}
