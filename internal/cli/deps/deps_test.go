//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package deps

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMermaidID(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"internal/cli/deps", "internal_cli_deps"},
		{"github.com/foo/bar", "github_com_foo_bar"},
		{"my-pkg", "my_pkg"},
	}
	for _, tt := range tests {
		if got := mermaidID(tt.input); got != tt.want {
			t.Errorf("mermaidID(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestIsStdlib(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"fmt", true},
		{"os/exec", true},
		{"encoding/json", true},
		{"github.com/foo/bar", false},
		{"golang.org/x/tools", false},
	}
	for _, tt := range tests {
		if got := isStdlib(tt.input); got != tt.want {
			t.Errorf("isStdlib(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestShortPkgName(t *testing.T) {
	mod := "github.com/ActiveMemory/ctx"
	tests := []struct {
		input string
		want  string
	}{
		{"github.com/ActiveMemory/ctx/internal/cli/deps", "internal/cli/deps"},
		{"github.com/ActiveMemory/ctx", "github.com/ActiveMemory/ctx"},
		{"github.com/other/pkg", "github.com/other/pkg"},
	}
	for _, tt := range tests {
		if got := shortPkgName(tt.input, mod); got != tt.want {
			t.Errorf("shortPkgName(%q, %q) = %q, want %q", tt.input, mod, got, tt.want)
		}
	}
}

func TestRenderMermaid(t *testing.T) {
	graph := map[string][]string{
		"cmd": {"internal/cli"},
		"internal/cli": {"internal/config"},
	}

	out := renderMermaid(graph)
	if !strings.HasPrefix(out, "graph TD\n") {
		t.Errorf("renderMermaid should start with 'graph TD\\n', got: %s", out)
	}
	if !strings.Contains(out, `cmd["cmd"] --> internal_cli["internal/cli"]`) {
		t.Errorf("renderMermaid missing expected edge, got: %s", out)
	}
}

func TestRenderTable(t *testing.T) {
	graph := map[string][]string{
		"cmd": {"internal/cli"},
	}

	out := renderTable(graph)
	if !strings.Contains(out, "Package") {
		t.Errorf("renderTable should contain header 'Package', got: %s", out)
	}
	if !strings.Contains(out, "cmd") {
		t.Errorf("renderTable should contain 'cmd', got: %s", out)
	}
}

func TestRenderJSON(t *testing.T) {
	graph := map[string][]string{
		"cmd": {"internal/cli"},
	}

	out := renderJSON(graph)
	var parsed map[string][]string
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("renderJSON produced invalid JSON: %v", err)
	}
	if len(parsed["cmd"]) != 1 || parsed["cmd"][0] != "internal/cli" {
		t.Errorf("unexpected parsed result: %v", parsed)
	}
}

func TestDetectProjectType(t *testing.T) {
	// Save and restore working directory.
	orig, getErr := os.Getwd()
	if getErr != nil {
		t.Fatal(getErr)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	// Temp dir without go.mod.
	tmp := t.TempDir()
	if chdirErr := os.Chdir(tmp); chdirErr != nil {
		t.Fatal(chdirErr)
	}
	if got := detectProjectType(); got != "" {
		t.Errorf("detectProjectType() = %q in empty dir, want empty", got)
	}

	// Add go.mod.
	if writeErr := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module test\n"), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}
	if got := detectProjectType(); got != "go" {
		t.Errorf("detectProjectType() = %q with go.mod, want 'go'", got)
	}
}

func TestRunDeps_GoProject(t *testing.T) {
	// Create a mini Go project with two packages and an import relationship.
	orig, getErr := os.Getwd()
	if getErr != nil {
		t.Fatal(getErr)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	tmp := t.TempDir()
	if chdirErr := os.Chdir(tmp); chdirErr != nil {
		t.Fatal(chdirErr)
	}

	// go.mod
	if writeErr := os.WriteFile(filepath.Join(tmp, "go.mod"),
		[]byte("module example.com/testmod\n\ngo 1.21\n"), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Package A: no imports.
	pkgA := filepath.Join(tmp, "pkga")
	if mkErr := os.MkdirAll(pkgA, 0o755); mkErr != nil {
		t.Fatal(mkErr)
	}
	if writeErr := os.WriteFile(filepath.Join(pkgA, "a.go"),
		[]byte("package pkga\n\nfunc Hello() string { return \"hello\" }\n"), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Package B: imports A.
	pkgB := filepath.Join(tmp, "pkgb")
	if mkErr := os.MkdirAll(pkgB, 0o755); mkErr != nil {
		t.Fatal(mkErr)
	}
	if writeErr := os.WriteFile(filepath.Join(pkgB, "b.go"),
		[]byte("package pkgb\n\nimport \"example.com/testmod/pkga\"\n\nvar _ = pkga.Hello\n"), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}

	// Run deps command.
	cmd := Cmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("runDeps failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "graph TD") {
		t.Errorf("expected mermaid output, got: %s", out)
	}
	if !strings.Contains(out, "pkgb") || !strings.Contains(out, "pkga") {
		t.Errorf("expected pkgb -> pkga edge in output, got: %s", out)
	}
}
