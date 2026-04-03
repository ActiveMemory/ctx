//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"os"
	"path/filepath"
	"testing"
)

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
		t.Errorf(
			"DetectBuilder() = %q in empty dir, want nil",
			b.Name(),
		)
	}

	// go.mod -> Go builder.
	goMod := filepath.Join(tmp, "go.mod")
	writeErr := os.WriteFile(
		goMod, []byte("module test\n"), 0o644,
	)
	if writeErr != nil {
		t.Fatal(writeErr)
	}
	if b := DetectBuilder(); b == nil ||
		b.Name() != "go" {
		t.Errorf(
			"DetectBuilder() with go.mod: want 'go', got %v",
			b,
		)
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

	pkgJSON := filepath.Join(tmp, "package.json")
	writeErr := os.WriteFile(
		pkgJSON, []byte(`{"name":"test"}`), 0o644,
	)
	if writeErr != nil {
		t.Fatal(writeErr)
	}
	if b := DetectBuilder(); b == nil ||
		b.Name() != "node" {
		t.Errorf(
			"DetectBuilder() with package.json: want 'node', got %v",
			b,
		)
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

	reqsPath := filepath.Join(tmp, "requirements.txt")
	writeErr := os.WriteFile(
		reqsPath, []byte("flask\n"), 0o644,
	)
	if writeErr != nil {
		t.Fatal(writeErr)
	}
	if b := DetectBuilder(); b == nil ||
		b.Name() != "python" {
		t.Errorf(
			"DetectBuilder() with requirements.txt: want 'python', got %v",
			b,
		)
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

	cargoPath := filepath.Join(tmp, "Cargo.toml")
	cargoContent := "[package]\nname = \"test\"\n"
	writeErr := os.WriteFile(
		cargoPath, []byte(cargoContent), 0o644,
	)
	if writeErr != nil {
		t.Fatal(writeErr)
	}
	if b := DetectBuilder(); b == nil ||
		b.Name() != "rust" {
		t.Errorf(
			"DetectBuilder() with Cargo.toml: want 'rust', got %v",
			b,
		)
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

	goMod := filepath.Join(tmp, "go.mod")
	writeErr := os.WriteFile(
		goMod, []byte("module test\n"), 0o644,
	)
	if writeErr != nil {
		t.Fatal(writeErr)
	}
	pkgJSON := filepath.Join(tmp, "package.json")
	writeErr = os.WriteFile(
		pkgJSON, []byte(`{"name":"test"}`), 0o644,
	)
	if writeErr != nil {
		t.Fatal(writeErr)
	}

	if b := DetectBuilder(); b == nil ||
		b.Name() != "go" {
		t.Errorf(
			"DetectBuilder() with go.mod+package.json: want 'go', got %v",
			b,
		)
	}
}

func TestFindBuilder(t *testing.T) {
	for _, name := range []string{
		"go", "node", "python", "rust",
	} {
		if b := FindBuilder(name); b == nil {
			t.Errorf(
				"FindBuilder(%q) = nil, want builder", name,
			)
		}
	}
	if b := FindBuilder("java"); b != nil {
		t.Errorf(
			"FindBuilder('java') = %v, want nil", b,
		)
	}
}

func TestBuilderNames(t *testing.T) {
	names := BuilderNames()
	if len(names) != 4 {
		t.Fatalf(
			"BuilderNames() returned %d names, want 4",
			len(names),
		)
	}
	expected := []string{
		"go", "node", "python", "rust",
	}
	for i, want := range expected {
		if names[i] != want {
			t.Errorf(
				"BuilderNames()[%d] = %q, want %q",
				i, names[i], want,
			)
		}
	}
}
