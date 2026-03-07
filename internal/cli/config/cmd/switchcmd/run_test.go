//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package switchcmd

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

func setupProfiles(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	if writeErr := os.WriteFile(
		filepath.Join(root, core.FileCtxRCDev), []byte(devContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}
	if writeErr := os.WriteFile(
		filepath.Join(root, core.FileCtxRCBase), []byte(baseContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}
	return root
}

func newTestCmd() *cobra.Command {
	buf := new(bytes.Buffer)
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	return cmd
}

func cmdOutput(cmd *cobra.Command) string {
	return cmd.OutOrStdout().(*bytes.Buffer).String()
}

func TestSwitch_DevToBase(t *testing.T) {
	root := setupProfiles(t)
	if writeErr := os.WriteFile(
		filepath.Join(root, core.FileCtxRC), []byte(devContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := newTestCmd()
	if switchErr := Run(cmd, root, []string{"base"}); switchErr != nil {
		t.Fatalf("unexpected error: %v", switchErr)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "switched to base") {
		t.Errorf("expected 'switched to base', got: %s", out)
	}

	if got := core.DetectProfile(root); got != core.ProfileBase {
		t.Errorf("profile should be base after switch, got %q", got)
	}
}

func TestSwitch_BaseToDev(t *testing.T) {
	root := setupProfiles(t)
	if writeErr := os.WriteFile(
		filepath.Join(root, core.FileCtxRC), []byte(baseContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := newTestCmd()
	if switchErr := Run(cmd, root, []string{"dev"}); switchErr != nil {
		t.Fatalf("unexpected error: %v", switchErr)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "switched to dev") {
		t.Errorf("expected 'switched to dev', got: %s", out)
	}

	if got := core.DetectProfile(root); got != core.ProfileDev {
		t.Errorf("profile should be dev after switch, got %q", got)
	}
}

func TestSwitch_AlreadyOnProfile(t *testing.T) {
	root := setupProfiles(t)
	if writeErr := os.WriteFile(
		filepath.Join(root, core.FileCtxRC), []byte(devContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := newTestCmd()
	if switchErr := Run(cmd, root, []string{"dev"}); switchErr != nil {
		t.Fatalf("unexpected error: %v", switchErr)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "already on dev") {
		t.Errorf("expected 'already on dev', got: %s", out)
	}
}

func TestSwitch_ProdAlias(t *testing.T) {
	root := setupProfiles(t)
	if writeErr := os.WriteFile(
		filepath.Join(root, core.FileCtxRC), []byte(devContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := newTestCmd()
	if switchErr := Run(cmd, root, []string{"prod"}); switchErr != nil {
		t.Fatalf("unexpected error: %v", switchErr)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "switched to base") {
		t.Errorf("expected 'switched to base' (prod alias), got: %s", out)
	}
}

func TestSwitch_Toggle_DevToBase(t *testing.T) {
	root := setupProfiles(t)
	if writeErr := os.WriteFile(
		filepath.Join(root, core.FileCtxRC), []byte(devContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := newTestCmd()
	if switchErr := Run(cmd, root, nil); switchErr != nil {
		t.Fatalf("unexpected error: %v", switchErr)
	}

	if got := core.DetectProfile(root); got != core.ProfileBase {
		t.Errorf("toggle from dev should go to base, got %q", got)
	}
}

func TestSwitch_Toggle_BaseToDev(t *testing.T) {
	root := setupProfiles(t)
	if writeErr := os.WriteFile(
		filepath.Join(root, core.FileCtxRC), []byte(baseContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := newTestCmd()
	if switchErr := Run(cmd, root, nil); switchErr != nil {
		t.Fatalf("unexpected error: %v", switchErr)
	}

	if got := core.DetectProfile(root); got != core.ProfileDev {
		t.Errorf("toggle from base should go to dev, got %q", got)
	}
}

func TestSwitch_Toggle_MissingCtxrc(t *testing.T) {
	root := setupProfiles(t)

	cmd := newTestCmd()
	if switchErr := Run(cmd, root, nil); switchErr != nil {
		t.Fatalf("unexpected error: %v", switchErr)
	}

	if got := core.DetectProfile(root); got != core.ProfileDev {
		t.Errorf("toggle from missing should go to dev, got %q", got)
	}
}

func TestSwitch_InvalidProfile(t *testing.T) {
	root := setupProfiles(t)

	cmd := newTestCmd()
	switchErr := Run(cmd, root, []string{"invalid"})
	if switchErr == nil {
		t.Fatal("expected error for invalid profile")
	}
}
