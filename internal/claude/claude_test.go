//   /    Context:                     https://ctx.ist
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

	"github.com/ActiveMemory/ctx/internal/tpl"
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

func TestCheckContextSizeScript(t *testing.T) {
	content, err := CheckContextSizeScript()
	if err != nil {
		t.Fatalf("CheckContextSizeScript() unexpected error: %v", err)
	}
	if len(content) == 0 {
		t.Error("CheckContextSizeScript() returned empty content")
	}
	script := string(content)
	if !strings.Contains(script, "#!/") {
		t.Error("CheckContextSizeScript() script missing shebang")
	}
}

func TestCleanupTmpScript(t *testing.T) {
	content, err := CleanupTmpScript()
	if err != nil {
		t.Fatalf("CleanupTmpScript() unexpected error: %v", err)
	}
	if len(content) == 0 {
		t.Error("CleanupTmpScript() returned empty content")
	}
	script := string(content)
	if !strings.Contains(script, "#!/") {
		t.Error("CleanupTmpScript() script missing shebang")
	}
}

func TestCheckPersistenceScript(t *testing.T) {
	content, err := CheckPersistenceScript()
	if err != nil {
		t.Fatalf("CheckPersistenceScript() unexpected error: %v", err)
	}
	if len(content) == 0 {
		t.Error("CheckPersistenceScript() returned empty content")
	}
	script := string(content)
	if !strings.Contains(script, "#!/") {
		t.Error("CheckPersistenceScript() script missing shebang")
	}
}

func TestErrFileRead(t *testing.T) {
	cause := errors.New("permission denied")
	err := errFileRead("/some/path", cause)
	if err == nil {
		t.Fatal("errFileRead() returned nil")
	}
	if !strings.Contains(err.Error(), "/some/path") {
		t.Errorf("errFileRead() error missing path: %v", err)
	}
	if !errors.Is(err, cause) {
		t.Error("errFileRead() does not wrap the cause error")
	}
}

func TestErrSkillList(t *testing.T) {
	cause := errors.New("read dir failed")
	err := errSkillList(cause)
	if err == nil {
		t.Fatal("errSkillList() returned nil")
	}
	if !strings.Contains(err.Error(), "failed to list skills") {
		t.Errorf("errSkillList() error missing prefix: %v", err)
	}
	if !errors.Is(err, cause) {
		t.Error("errSkillList() does not wrap the cause error")
	}
}

func TestErrSkillRead(t *testing.T) {
	cause := errors.New("not found")
	err := errSkillRead("my-skill", cause)
	if err == nil {
		t.Fatal("errSkillRead() returned nil")
	}
	if !strings.Contains(err.Error(), "my-skill") {
		t.Errorf("errSkillRead() error missing skill name: %v", err)
	}
	if !errors.Is(err, cause) {
		t.Error("errSkillRead() does not wrap the cause error")
	}
}

func TestNewHook(t *testing.T) {
	h := NewHook(HookTypeCommand, "echo hello")
	if h.Type != HookTypeCommand {
		t.Errorf("NewHook().Type = %q, want %q", h.Type, HookTypeCommand)
	}
	if h.Command != "echo hello" {
		t.Errorf("NewHook().Command = %q, want %q", h.Command, "echo hello")
	}
}

func TestDefaultHooks_SessionEnd(t *testing.T) {
	hooks := DefaultHooks("")
	if len(hooks.SessionEnd) == 0 {
		t.Error("DefaultHooks() SessionEnd is empty")
	}
	found := false
	for _, matcher := range hooks.SessionEnd {
		for _, hook := range matcher.Hooks {
			if strings.Contains(hook.Command, "cleanup-tmp") {
				found = true
			}
		}
	}
	if !found {
		t.Error("DefaultHooks() SessionEnd missing cleanup-tmp hook")
	}
}

func TestDefaultHooks_UserPromptSubmit(t *testing.T) {
	hooks := DefaultHooks("")
	if len(hooks.UserPromptSubmit) == 0 {
		t.Error("DefaultHooks() UserPromptSubmit is empty")
	}
	var hasContextSize, hasPersistence bool
	for _, matcher := range hooks.UserPromptSubmit {
		for _, hook := range matcher.Hooks {
			if strings.Contains(hook.Command, "check-context-size") {
				hasContextSize = true
			}
			if strings.Contains(hook.Command, "check-persistence") {
				hasPersistence = true
			}
		}
	}
	if !hasContextSize {
		t.Error("DefaultHooks() UserPromptSubmit missing check-context-size hook")
	}
	if !hasPersistence {
		t.Error("DefaultHooks() UserPromptSubmit missing check-persistence hook")
	}
}

func TestDefaultHooks_PreToolUseMatcherTypes(t *testing.T) {
	hooks := DefaultHooks("/project")
	if len(hooks.PreToolUse) < 2 {
		t.Fatalf("DefaultHooks() PreToolUse has %d matchers, want at least 2", len(hooks.PreToolUse))
	}
	if hooks.PreToolUse[0].Matcher != MatcherBash {
		t.Errorf("PreToolUse[0].Matcher = %q, want %q", hooks.PreToolUse[0].Matcher, MatcherBash)
	}
	if hooks.PreToolUse[1].Matcher != MatcherAll {
		t.Errorf("PreToolUse[1].Matcher = %q, want %q", hooks.PreToolUse[1].Matcher, MatcherAll)
	}
}

func TestSkillContentAllSkills(t *testing.T) {
	skills, err := Skills()
	if err != nil {
		t.Fatalf("Skills() failed: %v", err)
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

func TestHookTypeAndMatcherVars(t *testing.T) {
	if HookTypeCommand != "command" {
		t.Errorf("HookTypeCommand = %q, want %q", HookTypeCommand, "command")
	}
	if MatcherBash != "Bash" {
		t.Errorf("MatcherBash = %q, want %q", MatcherBash, "Bash")
	}
	if MatcherAll != ".*" {
		t.Errorf("MatcherAll = %q, want %q", MatcherAll, ".*")
	}
}

func TestCheckJournalScript(t *testing.T) {
	content, err := CheckJournalScript()
	if err != nil {
		t.Fatalf("CheckJournalScript() unexpected error: %v", err)
	}
	if len(content) == 0 {
		t.Error("CheckJournalScript() returned empty content")
	}
	script := string(content)
	if !strings.Contains(script, "#!/") {
		t.Error("CheckJournalScript() script missing shebang")
	}
}

// TestScriptErrorPaths swaps tpl.FS with an empty embed.FS to trigger
// error branches in all script and skill functions.
func TestScriptErrorPaths(t *testing.T) {
	orig := tpl.FS
	defer func() { tpl.FS = orig }()
	tpl.FS = embed.FS{} // empty FS causes all reads to fail

	if _, err := BlockNonPathCtxScript(); err == nil {
		t.Error("BlockNonPathCtxScript() expected error with empty FS")
	}
	if _, err := CheckContextSizeScript(); err == nil {
		t.Error("CheckContextSizeScript() expected error with empty FS")
	}
	if _, err := CleanupTmpScript(); err == nil {
		t.Error("CleanupTmpScript() expected error with empty FS")
	}
	if _, err := CheckPersistenceScript(); err == nil {
		t.Error("CheckPersistenceScript() expected error with empty FS")
	}
	if _, err := CheckJournalScript(); err == nil {
		t.Error("CheckJournalScript() expected error with empty FS")
	}
	if _, err := Skills(); err == nil {
		t.Error("Skills() expected error with empty FS")
	}
}
