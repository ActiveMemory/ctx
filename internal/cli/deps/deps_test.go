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

func TestRenderMermaid(t *testing.T) {
	graph := map[string][]string{
		"cmd":          {"internal/cli"},
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

func TestDetectBuilder(t *testing.T) {
	orig, getErr := os.Getwd()
	if getErr != nil {
		t.Fatal(getErr)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	// Empty dir — no builder detected.
	tmp := t.TempDir()
	if chdirErr := os.Chdir(tmp); chdirErr != nil {
		t.Fatal(chdirErr)
	}
	if b := detectBuilder(); b != nil {
		t.Errorf("detectBuilder() = %q in empty dir, want nil", b.Name())
	}

	// go.mod → Go builder.
	if writeErr := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module test\n"), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}
	if b := detectBuilder(); b == nil || b.Name() != "go" {
		t.Errorf("detectBuilder() with go.mod: want 'go', got %v", b)
	}
}

func TestDetectBuilder_Node(t *testing.T) {
	orig, getErr := os.Getwd()
	if getErr != nil {
		t.Fatal(getErr)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	tmp := t.TempDir()
	if chdirErr := os.Chdir(tmp); chdirErr != nil {
		t.Fatal(chdirErr)
	}

	if writeErr := os.WriteFile(filepath.Join(tmp, "package.json"), []byte(`{"name":"test"}`), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}
	if b := detectBuilder(); b == nil || b.Name() != "node" {
		t.Errorf("detectBuilder() with package.json: want 'node', got %v", b)
	}
}

func TestDetectBuilder_Python(t *testing.T) {
	orig, getErr := os.Getwd()
	if getErr != nil {
		t.Fatal(getErr)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	tmp := t.TempDir()
	if chdirErr := os.Chdir(tmp); chdirErr != nil {
		t.Fatal(chdirErr)
	}

	if writeErr := os.WriteFile(filepath.Join(tmp, "requirements.txt"), []byte("flask\n"), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}
	if b := detectBuilder(); b == nil || b.Name() != "python" {
		t.Errorf("detectBuilder() with requirements.txt: want 'python', got %v", b)
	}
}

func TestDetectBuilder_Rust(t *testing.T) {
	orig, getErr := os.Getwd()
	if getErr != nil {
		t.Fatal(getErr)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	tmp := t.TempDir()
	if chdirErr := os.Chdir(tmp); chdirErr != nil {
		t.Fatal(chdirErr)
	}

	if writeErr := os.WriteFile(filepath.Join(tmp, "Cargo.toml"), []byte("[package]\nname = \"test\"\n"), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}
	if b := detectBuilder(); b == nil || b.Name() != "rust" {
		t.Errorf("detectBuilder() with Cargo.toml: want 'rust', got %v", b)
	}
}

func TestDetectBuilder_PriorityOrder(t *testing.T) {
	orig, getErr := os.Getwd()
	if getErr != nil {
		t.Fatal(getErr)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	tmp := t.TempDir()
	if chdirErr := os.Chdir(tmp); chdirErr != nil {
		t.Fatal(chdirErr)
	}

	// Create both go.mod and package.json — Go should win (first in registry).
	if writeErr := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module test\n"), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}
	if writeErr := os.WriteFile(filepath.Join(tmp, "package.json"), []byte(`{"name":"test"}`), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}

	if b := detectBuilder(); b == nil || b.Name() != "go" {
		t.Errorf("detectBuilder() with go.mod+package.json: want 'go', got %v", b)
	}
}

func TestFindBuilder(t *testing.T) {
	for _, name := range []string{"go", "node", "python", "rust"} {
		if b := findBuilder(name); b == nil {
			t.Errorf("findBuilder(%q) = nil, want builder", name)
		}
	}
	if b := findBuilder("java"); b != nil {
		t.Errorf("findBuilder('java') = %v, want nil", b)
	}
}

func TestBuilderNames(t *testing.T) {
	names := builderNames()
	if len(names) != 4 {
		t.Fatalf("builderNames() returned %d names, want 4", len(names))
	}
	expected := []string{"go", "node", "python", "rust"}
	for i, want := range expected {
		if names[i] != want {
			t.Errorf("builderNames()[%d] = %q, want %q", i, names[i], want)
		}
	}
}

func TestRunDeps_TypeFlag(t *testing.T) {
	cmd := Cmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"--type", "invalid"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid --type, got nil")
	}
	if !strings.Contains(err.Error(), "unknown project type") {
		t.Errorf("unexpected error: %v", err)
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
