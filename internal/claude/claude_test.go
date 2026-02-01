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

func TestAutoSaveScript(t *testing.T) {
	content, err := AutoSaveScript()
	if err != nil {
		t.Fatalf("AutoSaveScript() unexpected error: %v", err)
	}

	if len(content) == 0 {
		t.Error("AutoSaveScript() returned empty content")
	}

	// Check for expected script content
	script := string(content)
	if !strings.Contains(script, "#!/") {
		t.Error("AutoSaveScript() script missing shebang")
	}
}

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

func TestPromptCoachScript(t *testing.T) {
	content, err := PromptCoachScript()
	if err != nil {
		t.Fatalf("PromptCoachScript() unexpected error: %v", err)
	}

	if len(content) == 0 {
		t.Error("PromptCoachScript() returned empty content")
	}

	// Check for expected script content
	script := string(content)
	if !strings.Contains(script, "#!/") {
		t.Error("PromptCoachScript() script missing shebang")
	}

	// Check that it contains pattern detection logic
	if !strings.Contains(script, "idiomatic") {
		t.Error("PromptCoachScript() should contain anti-pattern detection")
	}
}

func TestCommands(t *testing.T) {
	commands, err := Commands()
	if err != nil {
		t.Fatalf("Commands() unexpected error: %v", err)
	}

	if len(commands) == 0 {
		t.Error("Commands() returned empty list")
	}

	// Check that all entries are .md files
	for _, cmd := range commands {
		if !strings.HasSuffix(cmd, ".md") {
			t.Errorf("Commands() returned non-.md file: %s", cmd)
		}
	}
}

func TestCommandByName(t *testing.T) {
	// First get the list of commands to test with
	commands, err := Commands()
	if err != nil {
		t.Fatalf("Commands() failed: %v", err)
	}

	if len(commands) == 0 {
		t.Skip("no commands available to test")
	}

	// Test getting the first command
	content, err := CommandByName(commands[0])
	if err != nil {
		t.Errorf("CommandByName(%q) unexpected error: %v", commands[0], err)
	}
	if len(content) == 0 {
		t.Errorf("CommandByName(%q) returned empty content", commands[0])
	}

	// Test getting nonexistent command
	_, err = CommandByName("nonexistent-command.md")
	if err == nil {
		t.Error("CommandByName(nonexistent) expected error, got nil")
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

			// Check SessionEnd hooks
			if len(hooks.SessionEnd) == 0 {
				t.Error("DefaultHooks() SessionEnd is empty")
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

