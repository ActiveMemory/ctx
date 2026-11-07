//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package heartbeat

import (
	"bytes"
	"encoding/json"
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

func TestRun_AgyHeartbeatDrainsPending(t *testing.T) {
	origWd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(origWd)
		rc.Reset()
		coreSession.SetAgyMode(false, "", "", "")
	}()

	workspaceDir := t.TempDir()
	writeInitializedProject(t, workspaceDir)

	stdinJSON := `{
		"conversationId": "sess-hb-drain",
		"workspacePaths": ["` + filepath.ToSlash(workspaceDir) + `"]
	}`
	stdinFile := writeTempFile(t, stdinJSON)
	defer stdinFile.Close()

	// Pretend a previous PreToolUse/PostToolUse queued a pending nudge
	_ = os.Chdir(workspaceDir)
	rc.Reset()
	coreSession.SetAgyMode(true, "sess-hb-drain", "", "")
	coreSession.QueuePendingAgy("queued qa reminder nudge")

	cmd := &cobra.Command{}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := Run(cmd, stdinFile)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	out := strings.TrimSpace(buf.String())
	var resp struct {
		InjectSteps []struct {
			SystemMessage struct {
				SystemMessage string `json:"systemMessage"`
			} `json:"systemMessage"`
		} `json:"injectSteps"`
	}
	if unmarshalErr := json.Unmarshal([]byte(out), &resp); unmarshalErr != nil {
		t.Fatalf("expected valid JSON injectSteps, got %q (err: %v)", out, unmarshalErr)
	}
	if len(resp.InjectSteps) != 1 {
		t.Fatalf("len(InjectSteps) = %d, want 1", len(resp.InjectSteps))
	}
	if !strings.Contains(resp.InjectSteps[0].SystemMessage.SystemMessage, "queued qa reminder nudge") {
		t.Errorf("unexpected injectSteps content: %q", resp.InjectSteps[0].SystemMessage.SystemMessage)
	}

	// Second run: pending queue was drained, should emit "{}"
	stdinFile2 := writeTempFile(t, stdinJSON)
	defer stdinFile2.Close()
	buf.Reset()
	cmd.SetOut(&buf)

	err = Run(cmd, stdinFile2)
	if err != nil {
		t.Fatalf("Run 2 returned error: %v", err)
	}
	out2 := strings.TrimSpace(buf.String())
	if out2 != "{}" {
		t.Errorf("expected empty JSON {} when queue is empty, got %q", out2)
	}
}

func TestRun_ClaudeProducesNoStdout(t *testing.T) {
	origWd, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(origWd)
		rc.Reset()
		coreSession.SetAgyMode(false, "", "", "")
	}()

	workspaceDir := t.TempDir()
	writeInitializedProject(t, workspaceDir)

	stdinJSON := `{"session_id": "sess-claude-hb"}`
	stdinFile := writeTempFile(t, stdinJSON)
	defer stdinFile.Close()

	_ = os.Chdir(workspaceDir)
	rc.Reset()

	cmd := &cobra.Command{}
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	err := Run(cmd, stdinFile)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if buf.Len() != 0 {
		t.Errorf("expected no stdout in Claude mode, got: %q", buf.String())
	}
}
