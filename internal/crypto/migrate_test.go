//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package crypto

import (
	"os"
	"path/filepath"
	"testing"

	cryptocfg "github.com/ActiveMemory/ctx/internal/config/crypto"
	"github.com/ActiveMemory/ctx/internal/config/fs"
)

func TestMigrateKeyFile_GlobalExists_Noop(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	// Create global key.
	globalDir := filepath.Join(dir, ".ctx")
	if err := os.MkdirAll(globalDir, fs.PermKeyDir); err != nil {
		t.Fatal(err)
	}
	globalKey := filepath.Join(globalDir, cryptocfg.ContextKey)
	if err := os.WriteFile(globalKey, []byte("global-key"), fs.PermSecret); err != nil {
		t.Fatal(err)
	}

	contextDir := filepath.Join(dir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	// Should not panic or error — just a noop.
	MigrateKeyFile(contextDir)

	// Global key should be untouched.
	data, readErr := os.ReadFile(globalKey) //nolint:gosec // test path
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(data) != "global-key" {
		t.Errorf("global key was modified: got %q", string(data))
	}
}

func TestMigrateKeyFile_LegacyLocal_WarnsOnly(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	contextDir := filepath.Join(dir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	// Create legacy project-local key.
	localKey := filepath.Join(contextDir, cryptocfg.ContextKey)
	if err := os.WriteFile(localKey, []byte("local-key"), fs.PermSecret); err != nil {
		t.Fatal(err)
	}

	// Should warn but NOT auto-migrate.
	MigrateKeyFile(contextDir)

	// Local key should still exist (not moved).
	if _, err := os.Stat(localKey); err != nil {
		t.Error("local key was removed — should only warn, not migrate")
	}

	// Global key should NOT have been created.
	globalKey := GlobalKeyPath()
	if _, err := os.Stat(globalKey); err == nil {
		t.Error("global key was created — should only warn, not migrate")
	}
}

func TestMigrateKeyFile_LegacyUserLevel_WarnsOnly(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	// Create a legacy user-level key at ~/.local/ctx/keys/.
	legacyKeyDir := filepath.Join(dir, ".local", "ctx", "keys")
	if err := os.MkdirAll(legacyKeyDir, fs.PermKeyDir); err != nil {
		t.Fatal(err)
	}
	legacyKey := filepath.Join(legacyKeyDir, "some-project--abcd1234.key")
	if err := os.WriteFile(legacyKey, []byte("user-level-data"), fs.PermSecret); err != nil {
		t.Fatal(err)
	}

	contextDir := filepath.Join(dir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	// Should warn but NOT auto-migrate.
	MigrateKeyFile(contextDir)

	// Legacy key should still exist.
	if _, err := os.Stat(legacyKey); err != nil {
		t.Error("legacy key was removed — should only warn, not migrate")
	}

	// Global key should NOT have been created.
	globalKey := GlobalKeyPath()
	if _, err := os.Stat(globalKey); err == nil {
		t.Error("global key was created — should only warn, not migrate")
	}
}

func TestMigrateKeyFile_NothingToDo(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	contextDir := filepath.Join(dir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	// No keys anywhere — should be a noop.
	MigrateKeyFile(contextDir)

	globalKey := GlobalKeyPath()
	if _, err := os.Stat(globalKey); err == nil {
		t.Error("key was created when none should exist")
	}
}
