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

func TestRustBuilder_Detect(t *testing.T) {
	orig, getErr := os.Getwd()
	if getErr != nil {
		t.Fatal(getErr)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	tmp := t.TempDir()
	if chdirErr := os.Chdir(tmp); chdirErr != nil {
		t.Fatal(chdirErr)
	}

	b := &RustBuilder{}
	if b.Detect() {
		t.Error("RustBuilder.Detect() = true in empty dir")
	}

	cargoContent := "[package]\nname = \"test\"\n" +
		"version = \"0.1.0\"\n"
	cargoPath := filepath.Join(tmp, "Cargo.toml")
	if writeErr := os.WriteFile(
		cargoPath, []byte(cargoContent), 0o644,
	); writeErr != nil {
		t.Fatal(writeErr)
	}
	if !b.Detect() {
		t.Error("RustBuilder.Detect() = false with Cargo.toml")
	}
}

func TestRustBuilder_Name(t *testing.T) {
	b := &RustBuilder{}
	if got := b.Name(); got != "rust" {
		t.Errorf("RustBuilder.Name() = %q, want 'rust'", got)
	}
}
