//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package synccmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
	"github.com/ActiveMemory/ctx/internal/bootstrap"
	"github.com/ActiveMemory/ctx/internal/cli/steering/cmd/synccmd"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// The test lives in package `synccmd_test` to avoid an import cycle
// (bootstrap → cli/steering → cli/steering/cmd/synccmd).

// TestMain initializes the embedded text-asset lookup so the write
// helpers and error constructors resolve their DescKey strings.
func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}

// runSync executes `ctx steering sync <args…>` from a fresh project
// directory whose .ctxrc carries rcBody, through a production-shaped
// root (so the persistent --tool flag resolves exactly as it does
// for users). The root's PersistentPreRunE gate is detached: it
// demands an initialized, git-backed project, which is not what
// this test exercises.
func runSync(
	t *testing.T, rcBody string, args ...string,
) (string, error) {
	t.Helper()

	dir := t.TempDir()
	// .ctxrc is only read when $PWD/.context/ exists (cwd-anchored
	// resolution model).
	if mkErr := os.Mkdir(filepath.Join(dir, ".context"), 0o750); mkErr != nil {
		t.Fatalf("mkdir .context: %v", mkErr)
	}
	if rcBody != "" {
		if writeErr := os.WriteFile(
			filepath.Join(dir, ".ctxrc"), []byte(rcBody), 0o600,
		); writeErr != nil {
			t.Fatalf("write .ctxrc: %v", writeErr)
		}
	}
	t.Chdir(dir)
	rc.Reset()
	t.Cleanup(rc.Reset)

	root := bootstrap.RootCmd()
	root.PersistentPreRunE = nil
	steering := &cobra.Command{Use: "steering"}
	steering.AddCommand(synccmd.Cmd())
	root.AddCommand(steering)

	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs(append([]string{"steering", "sync"}, args...))

	execErr := root.Execute()
	return out.String(), execErr
}

// TestRun_CtxrcCodexSkipsPolitely covers the journal-reported bug:
// `tool: codex` in .ctxrc plus a bare `ctx steering sync` used to
// fail with "unsupported sync tool"; it must print the info line
// and exit 0.
func TestRun_CtxrcCodexSkipsPolitely(t *testing.T) {
	out, err := runSync(t, "tool: codex\n")
	if err != nil {
		t.Fatalf("expected exit 0, got error: %v", err)
	}
	assertDirectSkipLine(t, out, cfgHook.ToolCodex)
}

// TestRun_ToolFlagClaudeSkipsPolitely covers the explicit flag
// path for the other direct consumer.
func TestRun_ToolFlagClaudeSkipsPolitely(t *testing.T) {
	out, err := runSync(t, "", "--tool", cfgHook.ToolClaude)
	if err != nil {
		t.Fatalf("expected exit 0, got error: %v", err)
	}
	assertDirectSkipLine(t, out, cfgHook.ToolClaude)
}

// TestRun_ToolFlagCodexOverridesCtxrc asserts the flag wins over
// a syncable .ctxrc tool and still takes the polite-skip path.
func TestRun_ToolFlagCodexOverridesCtxrc(t *testing.T) {
	out, err := runSync(t, "tool: cursor\n", "--tool", cfgHook.ToolCodex)
	if err != nil {
		t.Fatalf("expected exit 0, got error: %v", err)
	}
	assertDirectSkipLine(t, out, cfgHook.ToolCodex)
}

// TestRun_UnknownToolStillErrors asserts the polite skip is scoped
// to the documented direct consumers: an unknown tool keeps the
// existing "unsupported sync tool" error.
func TestRun_UnknownToolStillErrors(t *testing.T) {
	_, err := runSync(t, "", "--tool", "foo")
	if err == nil {
		t.Fatal("expected error for unknown tool, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported sync tool") {
		t.Errorf("error = %q, want it to mention unsupported sync tool", err)
	}
	if !strings.Contains(err.Error(), `"foo"`) {
		t.Errorf("error = %q, want it to name the tool", err)
	}
}

// assertDirectSkipLine checks the polite-skip info line names the
// tool and the ctx agent delivery route, and that no sync report
// summary was printed.
func assertDirectSkipLine(t *testing.T, out, tool string) {
	t.Helper()
	if !strings.Contains(out, tool) {
		t.Errorf("output %q does not name tool %q", out, tool)
	}
	if !strings.Contains(out, "ctx agent") {
		t.Errorf("output %q does not explain the ctx agent route", out)
	}
	if !strings.Contains(out, "nothing to sync") {
		t.Errorf("output %q does not say nothing to sync", out)
	}
	if strings.Contains(out, "written") {
		t.Errorf("output %q printed a sync report; expected skip only", out)
	}
}
