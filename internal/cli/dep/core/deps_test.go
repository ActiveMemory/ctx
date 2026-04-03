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

	"github.com/ActiveMemory/ctx/internal/cli/dep/core/builder"
)

func TestDetect(t *testing.T) {
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
	if b := builder.Detect(); b != nil {
		t.Errorf(
			"Detect() = %q in empty dir, want nil",
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
	if b := builder.Detect(); b == nil ||
		b.Name() != "go" {
		t.Errorf(
			"Detect() with go.mod: want 'go', got %v",
			b,
		)
	}
}

func TestDetect_Node(t *testing.T) {
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
	if b := builder.Detect(); b == nil ||
		b.Name() != "node" {
		t.Errorf(
			"Detect() with package.json: want 'node', got %v",
			b,
		)
	}
}

func TestDetect_Python(t *testing.T) {
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
	if b := builder.Detect(); b == nil ||
		b.Name() != "python" {
		t.Errorf(
			"Detect() with requirements.txt: want 'python', got %v",
			b,
		)
	}
}

func TestDetect_Rust(t *testing.T) {
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
	if b := builder.Detect(); b == nil ||
		b.Name() != "rust" {
		t.Errorf(
			"Detect() with Cargo.toml: want 'rust', got %v",
			b,
		)
	}
}

func TestDetect_PriorityOrder(t *testing.T) {
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

	if b := builder.Detect(); b == nil ||
		b.Name() != "go" {
		t.Errorf(
			"Detect() with go.mod+package.json: want 'go', got %v",
			b,
		)
	}
}

func TestFind(t *testing.T) {
	for _, name := range []string{
		"go", "node", "python", "rust",
	} {
		if b := builder.Find(name); b == nil {
			t.Errorf(
				"Find(%q) = nil, want builder", name,
			)
		}
	}
	if b := builder.Find("java"); b != nil {
		t.Errorf(
			"Find('java') = %v, want nil", b,
		)
	}
}

func TestNames(t *testing.T) {
	names := builder.Names()
	if len(names) != 4 {
		t.Fatalf(
			"Names() returned %d names, want 4",
			len(names),
		)
	}
	expected := []string{
		"go", "node", "python", "rust",
	}
	for i, want := range expected {
		if names[i] != want {
			t.Errorf(
				"Names()[%d] = %q, want %q",
				i, names[i], want,
			)
		}
	}
}
