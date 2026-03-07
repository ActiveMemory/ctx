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

	"github.com/ActiveMemory/ctx/internal/cli/config/core"
)

const (
	devContent  = "notify:\n  events:\n    - loop\n"
	baseContent = "# .ctxrc\n# context_dir: .context\n"
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

func TestStatus_Dev(t *testing.T) {
	root := t.TempDir()
	if writeErr := os.WriteFile(
		filepath.Join(root, core.FileCtxRC), []byte(devContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

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
		filepath.Join(root, core.FileCtxRC), []byte(baseContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

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

	cmd := newTestCmd()
	if statusErr := Run(cmd, root); statusErr != nil {
		t.Fatalf("unexpected error: %v", statusErr)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "active: none") {
		t.Errorf("expected 'active: none', got: %s", out)
	}
}
