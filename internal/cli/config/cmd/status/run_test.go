//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/testutil/testctx"
)

const (
	devContent  = "profile: dev\nnotify:\n  events:\n    - loop\n"
	baseContent = "profile: base\n"
)

func newTestCmd() *cobra.Command {
	buf := new(bytes.Buffer)
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	return cmd
}

func cmdOutput(cmd *cobra.Command) string {
	return cmd.OutOrStdout().(*bytes.Buffer).String()
}

func chdirWithCleanup(t *testing.T, dir string) {
	t.Helper()
	origDir, _ := os.Getwd()
	_ = os.Chdir(dir)
	// Under the explicit-context-dir model, .ctxrc is read from
	// `filepath.Dir(CTX_DIR)/.ctxrc`. Declaring CTX_DIR at
	// `<dir>/.context` keeps this test's root-adjacent .ctxrc
	// visible to the loader.
	testctx.Declare(t, dir)
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
		rc.Reset()
	})
}

func TestStatus_Dev(t *testing.T) {
	root := t.TempDir()
	if writeErr := os.WriteFile(
		filepath.Join(root, file.CtxRC), []byte(devContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}
	chdirWithCleanup(t, root)

	cmd := newTestCmd()
	if statusErr := Run(cmd, root); statusErr != nil {
		t.Fatalf("unexpected error: %v", statusErr)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "active: dev") {
		t.Errorf("expected 'active: dev', got: %s", out)
	}
}

func TestStatus_Base(t *testing.T) {
	root := t.TempDir()
	if writeErr := os.WriteFile(
		filepath.Join(root, file.CtxRC), []byte(baseContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}
	chdirWithCleanup(t, root)

	cmd := newTestCmd()
	if statusErr := Run(cmd, root); statusErr != nil {
		t.Fatalf("unexpected error: %v", statusErr)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "active: base") {
		t.Errorf("expected 'active: base', got: %s", out)
	}
}

func TestStatus_Missing(t *testing.T) {
	root := t.TempDir()
	chdirWithCleanup(t, root)

	cmd := newTestCmd()
	if statusErr := Run(cmd, root); statusErr != nil {
		t.Fatalf("unexpected error: %v", statusErr)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "active: none") {
		t.Errorf("expected 'active: none', got: %s", out)
	}
}
