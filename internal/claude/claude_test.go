//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude

import (
	"strings"
	"testing"
)

func TestBlockNonPathCtxScript(t *testing.T) {
	content, err := BlockNonPathCtxScript()
	if err != nil {
		t.Fatalf("BlockNonPathCtxScript() unexpected error: %v", err)
	}

	if len(content) == 0 {
		t.Error("BlockNonPathCtxScript() returned empty content")
	}

	// Check for expected script content
	script := string(content)
	if !strings.Contains(script, "#!/") {
		t.Error("BlockNonPathCtxScript() script missing shebang")
	}
}

func TestSkills(t *testing.T) {
	skills, err := Skills()
	if err != nil {
		t.Fatalf("Skills() unexpected error: %v", err)
	}

	if len(skills) == 0 {
		t.Error("Skills() returned empty list")
	}

	// Check that all entries are skill directory names (no extension)
	for _, skill := range skills {
		if strings.Contains(skill, ".") {
			t.Errorf("Skills() returned name with extension: %s", skill)
		}
	}
}

func TestSkillContent(t *testing.T) {
	// First get the list of skills to test with
	skills, err := Skills()
	if err != nil {
		t.Fatalf("Skills() failed: %v", err)
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

func TestDefaultHooks(t *testing.T) {
	tests := []struct {
		name       string
		projectDir string
	}{
		{
			name:       "empty project dir",
			projectDir: "",
		},
		{
			name:       "with project dir",
			projectDir: "/home/user/myproject",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hooks := DefaultHooks(tt.projectDir)

			// Check PreToolUse hooks
			if len(hooks.PreToolUse) == 0 {
				t.Error("DefaultHooks() PreToolUse is empty")
			}

			// Check that project dir is used in paths when provided
			if tt.projectDir != "" {
				found := false
				for _, matcher := range hooks.PreToolUse {
					for _, hook := range matcher.Hooks {
						if strings.Contains(hook.Command, tt.projectDir) {
							found = true
							break
						}
					}
				}
				if !found {
					t.Error("DefaultHooks() project dir not found in hook commands")
				}
			}
		})
	}
}

func TestSettingsStructure(t *testing.T) {
	// Test that Settings struct can be instantiated correctly
	settings := Settings{
		Hooks: DefaultHooks(""),
		Permissions: PermissionsConfig{
			Allow: []string{"Bash(ctx status:*)", "Bash(ctx agent:*)"},
		},
	}

	if len(settings.Hooks.PreToolUse) == 0 {
		t.Error("Settings.Hooks.PreToolUse should not be empty")
	}

	if len(settings.Permissions.Allow) == 0 {
		t.Error("Settings.Permissions.Allow should not be empty")
	}
}
