//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestKeyDir(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	got := KeyDir()
	want := filepath.Join(dir, ".local", "ctx", "keys")
	if got != want {
		t.Errorf("KeyDir() = %q, want %q", got, want)
	}
}

func TestProjectKeySlug(t *testing.T) {
	slug := ProjectKeySlug("/home/jose/WORKSPACE/ctx")

	if !strings.HasPrefix(slug, "-home-jose-WORKSPACE-ctx--") {
		t.Errorf("slug = %q, want prefix -home-jose-WORKSPACE-ctx--", slug)
	}
	if !strings.HasSuffix(slug, ".key") {
		t.Errorf("slug = %q, want .key suffix", slug)
	}
	// SHA portion: 8 hex chars between -- and .key
	parts := strings.Split(slug, "--")
	if len(parts) != 2 {
		t.Fatalf("slug = %q, want exactly one -- separator", slug)
	}
	hashPart := strings.TrimSuffix(parts[1], ".key")
	if len(hashPart) != 8 {
		t.Errorf("hash = %q, want 8 hex chars", hashPart)
	}
}

func TestProjectKeySlug_Deterministic(t *testing.T) {
	a := ProjectKeySlug("/home/jose/project")
	b := ProjectKeySlug("/home/jose/project")
	if a != b {
		t.Errorf("non-deterministic: %q != %q", a, b)
	}
}

func TestProjectKeySlug_UniquePerProject(t *testing.T) {
	a := ProjectKeySlug("/home/jose/project-a")
	b := ProjectKeySlug("/home/jose/project-b")
	if a == b {
		t.Errorf("different projects produced same slug: %q", a)
	}
}

func TestProjectKeyPath(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	path := ProjectKeyPath("/some/project")
	if !strings.HasPrefix(path, filepath.Join(dir, ".local", "ctx", "keys")) {
		t.Errorf("path = %q, want prefix in ~/.local/ctx/keys/", path)
	}
	if !strings.HasSuffix(path, ".key") {
		t.Errorf("path = %q, want .key suffix", path)
	}
}

func TestResolveKeyPath_OverrideTakesPrecedence(t *testing.T) {
	override := "/custom/path/my.key"
	got := ResolveKeyPath(".context", "/project", override)
	if got != override {
		t.Errorf("ResolveKeyPath() = %q, want override %q", got, override)
	}
}

func TestResolveKeyPath_UserLevelBeforeLegacy(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	projectRoot := filepath.Join(dir, "project")

	// Create both user-level and legacy keys.
	userKey := ProjectKeyPath(projectRoot)
	if err := os.MkdirAll(filepath.Dir(userKey), PermKeyDir); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(userKey, []byte("user-key"), PermSecret); err != nil {
		t.Fatal(err)
	}

	contextDir := filepath.Join(projectRoot, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}
	legacyKey := filepath.Join(contextDir, FileContextKey)
	if err := os.WriteFile(legacyKey, []byte("legacy-key"), PermSecret); err != nil {
		t.Fatal(err)
	}

	got := ResolveKeyPath(contextDir, projectRoot, "")
	if got != userKey {
		t.Errorf("ResolveKeyPath() = %q, want user-level %q", got, userKey)
	}
}

func TestResolveKeyPath_FallbackToLegacy(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	projectRoot := filepath.Join(dir, "project")

	contextDir := filepath.Join(projectRoot, ".context")
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}
	legacyKey := filepath.Join(contextDir, FileContextKey)
	if err := os.WriteFile(legacyKey, []byte("legacy-key"), PermSecret); err != nil {
		t.Fatal(err)
	}

	got := ResolveKeyPath(contextDir, projectRoot, "")
	if got != legacyKey {
		t.Errorf("ResolveKeyPath() = %q, want legacy %q", got, legacyKey)
	}
}

func TestResolveKeyPath_DefaultsToUserLevel(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	projectRoot := filepath.Join(dir, "project")
	contextDir := filepath.Join(projectRoot, ".context")

	// Neither key exists — should default to user-level.
	got := ResolveKeyPath(contextDir, projectRoot, "")
	want := ProjectKeyPath(projectRoot)
	if got != want {
		t.Errorf("ResolveKeyPath() = %q, want default user-level %q", got, want)
	}
}
