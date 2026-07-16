//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package checkceremony_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/checkceremony"
	"github.com/ActiveMemory/ctx/internal/config/ceremony"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// writeInitializedProject creates a minimal initialized ctx project
// under tempDir (the files in ctx.FilesRequired) so state.Initialized
// returns true. Returns the .context directory.
func writeInitializedProject(t *testing.T, tempDir string) string {
	t.Helper()
	ctxDir := filepath.Join(tempDir, dir.Context)
	if mkErr := os.MkdirAll(ctxDir, 0o750); mkErr != nil {
		t.Fatalf("MkdirAll(%s): %v", ctxDir, mkErr)
	}
	for _, f := range ctx.FilesRequired {
		p := filepath.Join(ctxDir, f)
		if wErr := os.WriteFile(p, []byte("# "+f+"\n"), 0o600); wErr != nil {
			t.Fatalf("WriteFile(%s): %v", p, wErr)
		}
	}
	return ctxDir
}

// runWithPrompt feeds a UserPromptSubmit hook envelope carrying prompt
// on stdin and returns anything written to the command's stdout (the
// nudge channel) plus the error from Run.
func runWithPrompt(t *testing.T, prompt string) (string, error) {
	t.Helper()

	r, w, pipeErr := os.Pipe()
	if pipeErr != nil {
		t.Fatalf("os.Pipe: %v", pipeErr)
	}
	envelope := `{"session_id":"00000000-0000-0000-0000-000000000000",` +
		`"prompt":"` + prompt + `"}`
	go func() {
		defer func() { _ = w.Close() }()
		_, _ = io.Copy(w, bytes.NewReader([]byte(envelope)))
	}()
	t.Cleanup(func() { _ = r.Close() })

	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&out)
	cmd.SetErr(io.Discard)

	runErr := checkceremony.Run(cmd, r)
	return out.String(), runErr
}

// TestRun_LivePromptSelfSuppressesAndCredits verifies both defects the
// live-session guard closes: running /ctx-remember must (1) emit no
// ceremony nudge on that very prompt (self-suppress) and (2) touch the
// daily marker so the live ceremony is credited despite journal-import
// lag (live-credit). See specs/ceremony-nudge-live-session.md.
func TestRun_LivePromptSelfSuppressesAndCredits(t *testing.T) {
	tempDir := t.TempDir()
	ctxDir := writeInitializedProject(t, tempDir)

	t.Chdir(tempDir)
	rc.Reset()
	t.Cleanup(rc.Reset)

	out, runErr := runWithPrompt(t, "/ctx-remember")
	if runErr != nil {
		t.Fatalf("Run() error = %v, want nil (hooks must never fail)", runErr)
	}
	if out != "" {
		t.Errorf("nudge emitted on a /ctx-remember prompt: %q (want none)", out)
	}

	marker := filepath.Join(ctxDir, dir.State, ceremony.ThrottleID)
	if _, statErr := os.Stat(marker); statErr != nil {
		t.Errorf("day marker not created: stat(%s) err = %v (want nil)",
			marker, statErr)
	}
}
