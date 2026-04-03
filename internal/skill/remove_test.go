//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package skill

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRemove_ExistingSkill(t *testing.T) {
	dir := t.TempDir()
	writeSkillManifest(t, dir, "to-remove", `---
name: to-remove
description: Will be removed
---
Body
`)

	if err := Remove(dir, "to-remove"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify directory was deleted.
	if _, err := os.Stat(filepath.Join(dir, "to-remove")); !os.IsNotExist(err) {
		t.Error("expected skill directory to be deleted")
	}
}

func TestRemove_NonExistentSkill(t *testing.T) {
	dir := t.TempDir()

	err := Remove(dir, "does-not-exist")
	if err == nil {
		t.Fatal("expected error for non-existent skill")
	}
}

func TestRemove_NotADirectory(t *testing.T) {
	dir := t.TempDir()
	// Create a regular file instead of a directory.
	if err := os.WriteFile(filepath.Join(dir, "not-a-dir"), []byte("file"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Remove(dir, "not-a-dir")
	if err == nil {
		t.Fatal("expected error when target is not a directory")
	}
}
