//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProjectSlug(t *testing.T) {
	tests := []struct {
		name    string
		absPath string
		want    string
	}{
		{
			name:    "typical workspace path",
			absPath: "/home/jose/WORKSPACE/ctx",
			want:    "-home-jose-WORKSPACE-ctx",
		},
		{
			name:    "root-level project",
			absPath: "/opt/project",
			want:    "-opt-project",
		},
		{
			name:    "deeply nested path",
			absPath: "/home/user/dev/go/src/github.com/org/repo",
			want:    "-home-user-dev-go-src-github.com-org-repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProjectSlug(tt.absPath)
			if got != tt.want {
				t.Errorf("ProjectSlug(%q) = %q, want %q", tt.absPath, got, tt.want)
			}
		})
	}
}

func TestDiscoverMemoryPath_Found(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	projectRoot := filepath.Join(home, "WORKSPACE", "myproject")
	slug := ProjectSlug(projectRoot)

	memDir := filepath.Join(home, ".claude", "projects", slug, "memory")
	if mkErr := os.MkdirAll(memDir, 0o755); mkErr != nil {
		t.Fatal(mkErr)
	}
	memFile := filepath.Join(memDir, "MEMORY.md")
	memData := []byte("# Memory\n")
	if writeErr := os.WriteFile(memFile, memData, 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}

	got, discoverErr := DiscoverPath(projectRoot)
	if discoverErr != nil {
		t.Fatalf("unexpected error: %v", discoverErr)
	}
	if got != memFile {
		t.Errorf("DiscoverPath() = %q, want %q", got, memFile)
	}
}

func TestDiscoverMemoryPath_NotFound(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	projectRoot := filepath.Join(home, "WORKSPACE", "nonexistent")
	_, discoverErr := DiscoverPath(projectRoot)
	if discoverErr == nil {
		t.Fatal("expected error for missing MEMORY.md, got nil")
	}
}
