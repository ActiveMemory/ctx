//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

const (
	devContent  = "notify:\n  events:\n    - loop\n"
	baseContent = "# .ctxrc\n# context_dir: .context\n"
)

func setupProfiles(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	if writeErr := os.WriteFile(
		filepath.Join(root, fileCtxRCDev), []byte(devContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}
	if writeErr := os.WriteFile(
		filepath.Join(root, fileCtxRCBase), []byte(baseContent), 0o600,
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

func TestDetectProfile_Dev(t *testing.T) {
	root := setupProfiles(t)
	if writeErr := os.WriteFile(
		filepath.Join(root, fileCtxRC), []byte(devContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	got := detectProfile(root)
	if got != profileDev {
		t.Errorf("expected dev, got %q", got)
	}
}

func TestDetectProfile_Base(t *testing.T) {
	root := setupProfiles(t)
	if writeErr := os.WriteFile(
		filepath.Join(root, fileCtxRC), []byte(baseContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	got := detectProfile(root)
	if got != profileBase {
		t.Errorf("expected base, got %q", got)
	}
}

func TestDetectProfile_Missing(t *testing.T) {
	root := t.TempDir()
	got := detectProfile(root)
	if got != "" {
		t.Errorf("expected empty for missing file, got %q", got)
	}
}

func TestSwitch_DevToBase(t *testing.T) {
	root := setupProfiles(t)
	// Start on dev.
	if writeErr := os.WriteFile(
		filepath.Join(root, fileCtxRC), []byte(devContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := newTestCmd()
	if switchErr := runSwitch(cmd, root, []string{"base"}); switchErr != nil {
		t.Fatalf("unexpected error: %v", switchErr)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "switched to base") {
		t.Errorf("expected 'switched to base', got: %s", out)
	}

	if got := detectProfile(root); got != profileBase {
		t.Errorf("profile should be base after switch, got %q", got)
	}
}

func TestSwitch_BaseToDev(t *testing.T) {
	root := setupProfiles(t)
	if writeErr := os.WriteFile(
		filepath.Join(root, fileCtxRC), []byte(baseContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := newTestCmd()
	if switchErr := runSwitch(cmd, root, []string{"dev"}); switchErr != nil {
		t.Fatalf("unexpected error: %v", switchErr)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "switched to dev") {
		t.Errorf("expected 'switched to dev', got: %s", out)
	}

	if got := detectProfile(root); got != profileDev {
		t.Errorf("profile should be dev after switch, got %q", got)
	}
}

func TestSwitch_AlreadyOnProfile(t *testing.T) {
	root := setupProfiles(t)
	if writeErr := os.WriteFile(
		filepath.Join(root, fileCtxRC), []byte(devContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := newTestCmd()
	if switchErr := runSwitch(cmd, root, []string{"dev"}); switchErr != nil {
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
		filepath.Join(root, fileCtxRC), []byte(devContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := newTestCmd()
	if switchErr := runSwitch(cmd, root, []string{"prod"}); switchErr != nil {
		t.Fatalf("unexpected error: %v", switchErr)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "switched to base") {
		t.Errorf("expected 'switched to base' (prod alias), got: %s", out)
	}
}

func TestSwap_Toggle_DevToBase(t *testing.T) {
	root := setupProfiles(t)
	if writeErr := os.WriteFile(
		filepath.Join(root, fileCtxRC), []byte(devContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := newTestCmd()
	// No args = toggle.
	if switchErr := runSwitch(cmd, root, nil); switchErr != nil {
		t.Fatalf("unexpected error: %v", switchErr)
	}

	if got := detectProfile(root); got != profileBase {
		t.Errorf("toggle from dev should go to base, got %q", got)
	}
}

func TestSwap_Toggle_BaseToDev(t *testing.T) {
	root := setupProfiles(t)
	if writeErr := os.WriteFile(
		filepath.Join(root, fileCtxRC), []byte(baseContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := newTestCmd()
	if switchErr := runSwitch(cmd, root, nil); switchErr != nil {
		t.Fatalf("unexpected error: %v", switchErr)
	}

	if got := detectProfile(root); got != profileDev {
		t.Errorf("toggle from base should go to dev, got %q", got)
	}
}

func TestSwap_Toggle_MissingCtxrc(t *testing.T) {
	root := setupProfiles(t)
	// No .ctxrc — toggle should create dev (missing = base-like → dev).

	cmd := newTestCmd()
	if switchErr := runSwitch(cmd, root, nil); switchErr != nil {
		t.Fatalf("unexpected error: %v", switchErr)
	}

	if got := detectProfile(root); got != profileDev {
		t.Errorf("toggle from missing should go to dev, got %q", got)
	}
}

func TestStatus_Dev(t *testing.T) {
	root := setupProfiles(t)
	if writeErr := os.WriteFile(
		filepath.Join(root, fileCtxRC), []byte(devContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := newTestCmd()
	if statusErr := runStatus(cmd, root); statusErr != nil {
		t.Fatalf("unexpected error: %v", statusErr)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "active: dev") {
		t.Errorf("expected 'active: dev', got: %s", out)
	}
}

func TestStatus_Base(t *testing.T) {
	root := setupProfiles(t)
	if writeErr := os.WriteFile(
		filepath.Join(root, fileCtxRC), []byte(baseContent), 0o600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := newTestCmd()
	if statusErr := runStatus(cmd, root); statusErr != nil {
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
	if statusErr := runStatus(cmd, root); statusErr != nil {
		t.Fatalf("unexpected error: %v", statusErr)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, "active: none") {
		t.Errorf("expected 'active: none', got: %s", out)
	}
}

func TestSwitch_InvalidProfile(t *testing.T) {
	root := setupProfiles(t)

	cmd := newTestCmd()
	switchErr := runSwitch(cmd, root, []string{"invalid"})
	if switchErr == nil {
		t.Fatal("expected error for invalid profile")
	}
}
