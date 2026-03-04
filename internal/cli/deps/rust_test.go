//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package deps

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

	b := &rustBuilder{}
	if b.Detect() {
		t.Error("rustBuilder.Detect() = true in empty dir")
	}

	if writeErr := os.WriteFile(filepath.Join(tmp, "Cargo.toml"),
		[]byte("[package]\nname = \"test\"\nversion = \"0.1.0\"\n"), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}
	if !b.Detect() {
		t.Error("rustBuilder.Detect() = false with Cargo.toml")
	}
}

func TestRustBuilder_Name(t *testing.T) {
	b := &rustBuilder{}
	if got := b.Name(); got != "rust" {
		t.Errorf("rustBuilder.Name() = %q, want 'rust'", got)
	}
}
