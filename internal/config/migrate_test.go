//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMigrateKeyFile_LegacyRename(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	contextDir := filepath.Join(dir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	// Create legacy key.
	legacyKey := filepath.Join(contextDir, ".context.key")
	if err := os.WriteFile(legacyKey, []byte("legacy-data"), PermSecret); err != nil {
		t.Fatal(err)
	}

	MigrateKeyFile(contextDir, dir)

	// Legacy key should be gone.
	if _, err := os.Stat(legacyKey); err == nil {
		t.Error("legacy key still exists after migration")
	}

	// User-level key should exist with same content.
	userKey := ProjectKeyPath(dir)
	data, readErr := os.ReadFile(userKey) //nolint:gosec // test path
	if readErr != nil {
		t.Fatalf("user-level key not found: %v", readErr)
	}
	if string(data) != "legacy-data" {
		t.Errorf("key content = %q, want %q", string(data), "legacy-data")
	}
}

func TestMigrateKeyFile_PromotesToUserLevel(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	contextDir := filepath.Join(dir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	// Create current-name key at project-local path.
	localKey := filepath.Join(contextDir, FileContextKey)
	if err := os.WriteFile(localKey, []byte("local-key"), PermSecret); err != nil {
		t.Fatal(err)
	}

	MigrateKeyFile(contextDir, dir)

	// Project-local key should be removed.
	if _, err := os.Stat(localKey); err == nil {
		t.Error("project-local key still exists after promotion")
	}

	// User-level key should exist.
	userKey := ProjectKeyPath(dir)
	data, readErr := os.ReadFile(userKey) //nolint:gosec // test path
	if readErr != nil {
		t.Fatalf("user-level key not found: %v", readErr)
	}
	if string(data) != "local-key" {
		t.Errorf("key content = %q, want %q", string(data), "local-key")
	}
}

func TestMigrateKeyFile_UserLevelExists_CleansLocal(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	contextDir := filepath.Join(dir, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	// Create both user-level and project-local keys.
	userKey := ProjectKeyPath(dir)
	if err := os.MkdirAll(filepath.Dir(userKey), PermKeyDir); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(userKey, []byte("user-key"), PermSecret); err != nil {
		t.Fatal(err)
	}

	localKey := filepath.Join(contextDir, FileContextKey)
	if err := os.WriteFile(localKey, []byte("stale-local"), PermSecret); err != nil {
		t.Fatal(err)
	}

	MigrateKeyFile(contextDir, dir)

	// User-level key should be preserved (not overwritten).
	data, readErr := os.ReadFile(userKey) //nolint:gosec // test path
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(data) != "user-key" {
		t.Errorf("user key was overwritten: got %q", string(data))
	}

	// Project-local key should be cleaned up.
	if _, err := os.Stat(localKey); err == nil {
		t.Error("stale project-local key should have been removed")
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
	MigrateKeyFile(contextDir, dir)

	userKey := ProjectKeyPath(dir)
	if _, err := os.Stat(userKey); err == nil {
		t.Error("key was created when none should exist")
	}
}
