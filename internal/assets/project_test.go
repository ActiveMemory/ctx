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
	"github.com/ActiveMemory/ctx/internal/config/marker"
)

// TestClaudeMd verifies the CLAUDE.md template is the bridge that
// imports AGENTS.md: ctx markers around an @AGENTS.md line and none of
// the instructions themselves, which live in the AGENTS.md template.
func TestClaudeMd(t *testing.T) {
	content, err := FS.ReadFile(asset.PathCLAUDEMd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	text := string(content)
	if !strings.Contains(text, "Context") {
		t.Error("CLAUDE.md does not contain 'Context'")
	}
	for _, want := range []string{
		marker.CtxStart, marker.CtxEnd, "\n@AGENTS.md\n",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("CLAUDE.md template missing %q", want)
		}
	}
	if strings.Contains(text, "ctx system bootstrap") {
		t.Error("CLAUDE.md template should import AGENTS.md, not repeat it")
	}
}

// TestAgentsMd verifies the AGENTS.md template carries the ctx block
// markers the init merge and the setup deployers key on, plus the
// bootstrap instruction every agent must follow.
func TestAgentsMd(t *testing.T) {
	content, err := FS.ReadFile(asset.PathAgentsMd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	text := string(content)
	for _, want := range []string{
		marker.AgentsStart, marker.AgentsEnd, "ctx system bootstrap",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("AGENTS.md template missing %q", want)
		}
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
