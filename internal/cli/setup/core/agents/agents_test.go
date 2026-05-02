//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agents

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/agent"
	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/spf13/cobra"
)

func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}

func testCmd(buf *bytes.Buffer) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	return cmd
}

func withTempProjectDir(t *testing.T) string {
	t.Helper()

	tmp := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })
	return tmp
}

func TestDeploy_CreatesFileWhenAbsent(t *testing.T) {
	withTempProjectDir(t)

	var buf bytes.Buffer
	if err := Deploy(testCmd(&buf)); err != nil {
		t.Fatalf("Deploy() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Clean("AGENTS.md"))
	if err != nil {
		t.Fatalf("read AGENTS.md: %v", err)
	}
	if !strings.Contains(string(data), marker.AgentsStart) {
		t.Fatalf("AGENTS.md missing ctx marker: %q", string(data))
	}
	if !strings.Contains(buf.String(), "✓") {
		t.Fatalf("expected create output, got %q", buf.String())
	}
}

func TestDeploy_MergesExistingFileWithoutMarkers(t *testing.T) {
	withTempProjectDir(t)

	const existing = `# Local Instructions

Do local things.`
	if err := os.WriteFile("AGENTS.md", []byte(existing), 0o644); err != nil {
		t.Fatalf("seed AGENTS.md: %v", err)
	}

	template, err := agent.AgentsMd()
	if err != nil {
		t.Fatalf("read embedded AGENTS.md: %v", err)
	}

	var buf bytes.Buffer
	if err := Deploy(testCmd(&buf)); err != nil {
		t.Fatalf("Deploy() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Clean("AGENTS.md"))
	if err != nil {
		t.Fatalf("read AGENTS.md: %v", err)
	}
	got := string(data)
	if !strings.Contains(got, existing) {
		t.Fatalf("merged file lost existing content: %q", got)
	}
	if !strings.Contains(got, string(template)) {
		t.Fatalf("merged file missing ctx template")
	}
	if !strings.Contains(buf.String(), "merged") {
		t.Fatalf("expected merge output, got %q", buf.String())
	}
}

func TestDeploy_SkipsExistingFileWithMarkers(t *testing.T) {
	withTempProjectDir(t)

	existing := marker.AgentsStart + `
custom ctx-managed section
`
	if err := os.WriteFile("AGENTS.md", []byte(existing), 0o644); err != nil {
		t.Fatalf("seed AGENTS.md: %v", err)
	}

	var buf bytes.Buffer
	if err := Deploy(testCmd(&buf)); err != nil {
		t.Fatalf("Deploy() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Clean("AGENTS.md"))
	if err != nil {
		t.Fatalf("read AGENTS.md: %v", err)
	}
	if string(data) != existing {
		t.Fatalf("AGENTS.md should be unchanged; got %q", string(data))
	}
	if !strings.Contains(buf.String(), "skipped") {
		t.Fatalf("expected skip output, got %q", buf.String())
	}
}

func TestDeploy_RejectsSymlinkTarget(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink behavior varies on Windows in this environment")
	}

	tmp := withTempProjectDir(t)
	realFile := filepath.Join(t.TempDir(), "outside.md")
	if err := os.WriteFile(realFile, []byte("secret"), 0o644); err != nil {
		t.Fatalf("seed real file: %v", err)
	}
	if err := os.Symlink(realFile, filepath.Join(tmp, "AGENTS.md")); err != nil {
		t.Fatalf("create symlink: %v", err)
	}

	err := Deploy(testCmd(&bytes.Buffer{}))
	if err == nil {
		t.Fatal("expected symlink rejection, got nil")
	}
}

func TestDeploy_RejectsNonRegularTarget(t *testing.T) {
	tmp := withTempProjectDir(t)
	if err := os.Mkdir(filepath.Join(tmp, "AGENTS.md"), 0o755); err != nil {
		t.Fatalf("mkdir AGENTS.md: %v", err)
	}

	err := Deploy(testCmd(&bytes.Buffer{}))
	if err == nil {
		t.Fatal("expected non-regular file rejection, got nil")
	}
}
