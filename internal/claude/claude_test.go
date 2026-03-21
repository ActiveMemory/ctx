//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude

import (
	"embed"
	"errors"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/skill"
)

func TestSkills(t *testing.T) {
	skills, err := SkillList()
	if err != nil {
		t.Fatalf("SkillList() unexpected error: %v", err)
	}

	if len(skills) == 0 {
		t.Error("SkillList() returned empty list")
	}

	// Check that all entries are skill directory names (no extension)
	for _, skill := range skills {
		if strings.Contains(skill, ".") {
			t.Errorf("SkillList() returned name with extension: %s", skill)
		}
	}
}

func TestSkillContent(t *testing.T) {
	// First get the list of skills to test with
	skills, err := SkillList()
	if err != nil {
		t.Fatalf("SkillList() failed: %v", err)
	}

	if len(skills) == 0 {
		t.Skip("no skills available to test")
	}

	// Test getting the first skill
	content, err := SkillContent(skills[0])
	if err != nil {
		t.Errorf("SkillContent(%q) unexpected error: %v", skills[0], err)
	}
	if len(content) == 0 {
		t.Errorf("SkillContent(%q) returned empty content", skills[0])
	}

	// Verify it's a valid SKILL.md with frontmatter
	contentStr := string(content)
	if !strings.HasPrefix(contentStr, "---") {
		t.Errorf("SkillContent(%q) missing frontmatter", skills[0])
	}

	// Test getting nonexistent skill
	_, err = SkillContent("nonexistent-skill")
	if err == nil {
		t.Error("SkillContent(nonexistent) expected error, got nil")
	}
}

func TestSkillContentAllSkills(t *testing.T) {
	skills, err := SkillList()
	if err != nil {
		t.Fatalf("SkillList() failed: %v", err)
	}
	for _, name := range skills {
		content, err := SkillContent(name)
		if err != nil {
			t.Errorf("SkillContent(%q) error: %v", name, err)
			continue
		}
		if len(content) == 0 {
			t.Errorf("SkillContent(%q) returned empty content", name)
		}
	}
}

func TestSettingsStructure(t *testing.T) {
	settings := Settings{
		Permissions: PermissionsConfig{
			Allow: []string{"Bash(ctx status:*)", "Bash(ctx agent:*)"},
		},
	}

	if len(settings.Permissions.Allow) == 0 {
		t.Error("Settings.Permissions.Allow should not be empty")
	}
}

func TestErrSkillList(t *testing.T) {
	cause := errors.New("read dir failed")
	err := ctxerr.List(cause)
	if err == nil {
		t.Fatal("ctxerr.List() returned nil")
	}
	if !strings.Contains(err.Error(), "failed to list skills") {
		t.Errorf("ctxerr.List() error missing prefix: %v", err)
	}
	if !errors.Is(err, cause) {
		t.Error("ctxerr.List() does not wrap the cause error")
	}
}

func TestErrSkillRead(t *testing.T) {
	cause := errors.New("not found")
	err := ctxerr.Read("my-skill", cause)
	if err == nil {
		t.Fatal("ctxerr.Read() returned nil")
	}
	if !strings.Contains(err.Error(), "my-skill") {
		t.Errorf("ctxerr.Read() error missing skill name: %v", err)
	}
	if !errors.Is(err, cause) {
		t.Error("ctxerr.Read() does not wrap the cause error")
	}
}

// TestScriptErrorPaths swaps assets.FS with an empty embed.FS to trigger
// error branches in skill functions.
func TestScriptErrorPaths(t *testing.T) {
	orig := assets.FS
	defer func() { assets.FS = orig }()
	assets.FS = embed.FS{} // empty FS causes all reads to fail

	if _, err := SkillList(); err == nil {
		t.Error("SkillList() expected error with empty FS")
	}
}
