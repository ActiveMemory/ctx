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

	cfgCrypto "github.com/ActiveMemory/ctx/internal/config/crypto"
	"github.com/ActiveMemory/ctx/internal/config/fs"
)

func TestGlobalKeyPath(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	got := GlobalKeyPath()
	want := filepath.Join(dir, ".ctx", cfgCrypto.ContextKey)
	if got != want {
		t.Errorf("GlobalKeyPath() = %q, want %q", got, want)
	}
}

func TestExpandHome_Tilde(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	got := ExpandHome("~/foo")
	want := filepath.Join(dir, "foo")
	if got != want {
		t.Errorf("ExpandHome(~/foo) = %q, want %q", got, want)
	}
}

func TestExpandHome_NoTilde(t *testing.T) {
	got := ExpandHome("/abs/path")
	if got != "/abs/path" {
		t.Errorf("ExpandHome(/abs/path) = %q, want /abs/path", got)
	}
}

func TestExpandHome_TildeOnly(t *testing.T) {
	got := ExpandHome("~")
	if got != "~" {
		t.Errorf("ExpandHome(~) = %q, want ~ (no trailing /)", got)
	}
}

func TestResolveKeyPath_OverrideTakesPrecedence(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	got := ResolveKeyPath(".context", "~/custom/my.key")
	want := filepath.Join(dir, "custom", "my.key")
	if got != want {
		t.Errorf("ResolveKeyPath() = %q, want override %q", got, want)
	}
}

func TestResolveKeyPath_ProjectLocalIgnored(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	// A project-local key exists, but it must NOT be auto-detected:
	// the global key wins. The implicit project-local tier was removed
	// (it broke worktrees and was a security antipattern) — see
	// specs/notify-resolution-hardening.md.
	contextDir := filepath.Join(dir, "project", ".context")
	if err := os.MkdirAll(contextDir, fs.PermKeyDir); err != nil {
		t.Fatal(err)
	}
	localKey := filepath.Join(contextDir, cfgCrypto.ContextKey)
	if err := os.WriteFile(localKey, []byte("local-key"), fs.PermSecret); err != nil {
		t.Fatal(err)
	}

	got := ResolveKeyPath(contextDir, "")
	want := GlobalKeyPath()
	if got != want {
		t.Errorf("ResolveKeyPath() = %q, want global %q (project-local must be ignored)", got, want)
	}
	if got == localKey {
		t.Errorf("ResolveKeyPath() returned the project-local key %q; it must not be auto-detected", localKey)
	}
}

func TestResolveKeyPath_HomeUnavailableFallsBackToLocal(t *testing.T) {
	// With no home dir, GlobalKeyPath() returns "" and resolution
	// falls back to the project-local path as a last resort.
	t.Setenv("HOME", "")

	contextDir := filepath.Join("project", ".context")
	got := ResolveKeyPath(contextDir, "")
	want := filepath.Join(contextDir, cfgCrypto.ContextKey)
	if got != want {
		t.Errorf("ResolveKeyPath() = %q, want local fallback %q", got, want)
	}
}

func TestResolveKeyPath_FallbackToGlobal(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	// Create global key only: no project-local.
	globalDir := filepath.Join(dir, ".ctx")
	if err := os.MkdirAll(globalDir, fs.PermKeyDir); err != nil {
		t.Fatal(err)
	}
	globalKey := filepath.Join(globalDir, cfgCrypto.ContextKey)
	gData := []byte("global-key")
	if err := os.WriteFile(globalKey, gData, fs.PermSecret); err != nil {
		t.Fatal(err)
	}

	contextDir := filepath.Join(dir, "project", ".context")
	got := ResolveKeyPath(contextDir, "")
	if got != globalKey {
		t.Errorf("ResolveKeyPath() = %q, want global %q", got, globalKey)
	}
}

func TestResolveKeyPath_DefaultsToGlobal(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	contextDir := filepath.Join(dir, "project", ".context")

	// Neither key exists: should default to global path.
	got := ResolveKeyPath(contextDir, "")
	want := GlobalKeyPath()
	if got != want {
		t.Errorf("ResolveKeyPath() = %q, want global default %q", got, want)
	}
}
