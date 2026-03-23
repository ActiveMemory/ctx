//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
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
		{"internal/cli/dep", "internal_cli_dep"},
		{"github.com/foo/bar", "github_com_foo_bar"},
		{"my-pkg", "my_pkg"},
	}
	for _, tt := range tests {
		if got := MermaidID(tt.input); got != tt.want {
			t.Errorf("MermaidID(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestRenderMermaid(t *testing.T) {
	graph := map[string][]string{
		"cmd":          {"internal/cli"},
		"internal/cli": {"internal/config"},
	}

	out := RenderMermaid(graph)
	if !strings.HasPrefix(out, "graph TD\n") {
		t.Errorf("RenderMermaid should start with 'graph TD\\n', got: %s", out)
	}
	if !strings.Contains(out, `cmd["cmd"] --> internal_cli["internal/cli"]`) {
		t.Errorf("RenderMermaid missing expected edge, got: %s", out)
	}
}

func TestRenderTable(t *testing.T) {
	graph := map[string][]string{
		"cmd": {"internal/cli"},
	}

	out := RenderTable(graph)
	if !strings.Contains(out, "Package") {
		t.Errorf("RenderTable should contain header 'Package', got: %s", out)
	}
	if !strings.Contains(out, "cmd") {
		t.Errorf("RenderTable should contain 'cmd', got: %s", out)
	}
}

func TestRenderJSON(t *testing.T) {
	graph := map[string][]string{
		"cmd": {"internal/cli"},
	}

	out := RenderJSON(graph)
	var parsed map[string][]string
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("RenderJSON produced invalid JSON: %v", err)
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

	// Empty dir - no builder detected.
	tmp := t.TempDir()
	if chdirErr := os.Chdir(tmp); chdirErr != nil {
		t.Fatal(chdirErr)
	}
	if b := DetectBuilder(); b != nil {
		t.Errorf("DetectBuilder() = %q in empty dir, want nil", b.Name())
	}

	// go.mod → Go builder.
	if writeErr := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module test\n"), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}
	if b := DetectBuilder(); b == nil || b.Name() != "go" {
		t.Errorf("DetectBuilder() with go.mod: want 'go', got %v", b)
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
	if b := DetectBuilder(); b == nil || b.Name() != "node" {
		t.Errorf("DetectBuilder() with package.json: want 'node', got %v", b)
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
	if b := DetectBuilder(); b == nil || b.Name() != "python" {
		t.Errorf("DetectBuilder() with requirements.txt: want 'python', got %v", b)
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
	if b := DetectBuilder(); b == nil || b.Name() != "rust" {
		t.Errorf("DetectBuilder() with Cargo.toml: want 'rust', got %v", b)
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

	// Create both go.mod and package.json - Go should win (first in registry).
	if writeErr := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module test\n"), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}
	if writeErr := os.WriteFile(filepath.Join(tmp, "package.json"), []byte(`{"name":"test"}`), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}

	if b := DetectBuilder(); b == nil || b.Name() != "go" {
		t.Errorf("DetectBuilder() with go.mod+package.json: want 'go', got %v", b)
	}
}

func TestFindBuilder(t *testing.T) {
	for _, name := range []string{"go", "node", "python", "rust"} {
		if b := FindBuilder(name); b == nil {
			t.Errorf("FindBuilder(%q) = nil, want builder", name)
		}
	}
	if b := FindBuilder("java"); b != nil {
		t.Errorf("FindBuilder('java') = %v, want nil", b)
	}
}

func TestBuilderNames(t *testing.T) {
	names := BuilderNames()
	if len(names) != 4 {
		t.Fatalf("BuilderNames() returned %d names, want 4", len(names))
	}
	expected := []string{"go", "node", "python", "rust"}
	for i, want := range expected {
		if names[i] != want {
			t.Errorf("BuilderNames()[%d] = %q, want %q", i, names[i], want)
		}
	}
}
