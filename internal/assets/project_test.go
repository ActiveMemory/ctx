//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package assets

import (
	"path"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/asset"
)

func TestClaudeMd(t *testing.T) {
	content, err := FS.ReadFile(asset.PathCLAUDEMd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(content), "Context") {
		t.Error("CLAUDE.md does not contain 'Context'")
	}
}

func TestProjectFile(t *testing.T) {
	tests := []struct {
		name        string
		file        string
		wantContain string
		wantErr     bool
	}{
		{"Makefile.ctx exists", "Makefile.ctx", "ctx", false},
		{"nonexistent returns error", "NONEXISTENT.md", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := FS.ReadFile(path.Join(asset.DirProject, tt.file))
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for %q", tt.file)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for %q: %v", tt.file, err)
				return
			}
			if !strings.Contains(string(content), tt.wantContain) {
				t.Errorf("content of %q does not contain %q", tt.file, tt.wantContain)
			}
		})
	}
}

func TestMakefileCtx(t *testing.T) {
	content, readErr := FS.ReadFile(asset.PathMakefileCtx)
	if readErr != nil {
		t.Fatalf("unexpected error: %v", readErr)
	}
	if len(content) == 0 {
		t.Fatal("returned empty content")
	}
	if !strings.Contains(string(content), "ctx") {
		t.Error("content does not contain 'ctx'")
	}
}
