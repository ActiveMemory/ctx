//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/claude"
	cfgClaude "github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/env"
)

// TestInitCommand tests the init command creates the .context directory.
func TestInitCommand(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-init-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()
	t.Setenv("HOME", tmpDir)
	t.Setenv(env.CtxDir, filepath.Join(tmpDir, ".context"))
	t.Setenv(env.SkipPathCheck, env.True)

	cmd := Cmd()
	cmd.SetArgs([]string{})
	if err = cmd.Execute(); err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	ctxDir := filepath.Join(tmpDir, ".context")
	info, err := os.Stat(ctxDir)
	if err != nil {
		t.Fatalf(".context directory was not created: %v", err)
	}
	if !info.IsDir() {
		t.Fatal(".context should be a directory")
	}

	requiredFiles := []string{
		"CONSTITUTION.md",
		"TASKS.md",
		"DECISIONS.md",
		"CONVENTIONS.md",
		"ARCHITECTURE.md",
	}
	for _, name := range requiredFiles {
		path := filepath.Join(ctxDir, name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("required file %s was not created", name)
		}
	}
}

func TestInitCreatesSteeringHooksSkillsDirs(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-init-dirs-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()
	t.Setenv("HOME", tmpDir)
	t.Setenv(env.CtxDir, filepath.Join(tmpDir, ".context"))
	t.Setenv(env.SkipPathCheck, env.True)

	cmd := Cmd()
	cmd.SetArgs([]string{})
	if err = cmd.Execute(); err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	for _, sub := range []string{"steering", "hooks", "skills"} {
		dirPath := filepath.Join(tmpDir, ".context", sub)
		info, statErr := os.Stat(dirPath)
		if statErr != nil {
			t.Errorf(".context/%s was not created: %v", sub, statErr)
			continue
		}
		if !info.IsDir() {
			t.Errorf(".context/%s should be a directory", sub)
		}
	}
}

func TestInitSkipsExistingSteeringHooksSkillsDirs(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-init-dirs-exist-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()
	t.Setenv("HOME", tmpDir)
	t.Setenv(env.CtxDir, filepath.Join(tmpDir, ".context"))
	t.Setenv(env.SkipPathCheck, env.True)

	// Pre-create the directories with a marker file inside each.
	for _, sub := range []string{"steering", "hooks", "skills"} {
		dirPath := filepath.Join(tmpDir, ".context", sub)
		if mkErr := os.MkdirAll(dirPath, 0755); mkErr != nil {
			t.Fatalf("failed to pre-create %s: %v", sub, mkErr)
		}
		marker := filepath.Join(dirPath, "marker.txt")
		if wErr := os.WriteFile(marker, []byte("keep"), 0644); wErr != nil {
			t.Fatalf("failed to write marker in %s: %v", sub, wErr)
		}
	}

	cmd := Cmd()
	cmd.SetArgs([]string{"--force"})
	if err = cmd.Execute(); err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	// Verify directories still exist and marker files are preserved.
	for _, sub := range []string{"steering", "hooks", "skills"} {
		marker := filepath.Join(tmpDir, ".context", sub, "marker.txt")
		content, readErr := os.ReadFile(marker)
		if readErr != nil {
			t.Errorf(".context/%s/marker.txt was lost: %v", sub, readErr)
			continue
		}
		if string(content) != "keep" {
			t.Errorf(".context/%s/marker.txt content changed", sub)
		}
	}
}

func TestInitMergeInsertsAfterH1(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-init-merge-h1-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()
	t.Setenv("HOME", tmpDir)
	t.Setenv(env.CtxDir, filepath.Join(tmpDir, ".context"))
	t.Setenv(env.SkipPathCheck, env.True)

	existingContent := "# My Amazing Project\n\n" +
		"This is the project description.\n\n" +
		"## Build Instructions\n\nRun make build.\n"
	if err = os.WriteFile("CLAUDE.md", []byte(existingContent), 0600); err != nil {
		t.Fatalf("failed to create CLAUDE.md: %v", err)
	}

	initCmd := Cmd()
	initCmd.SetArgs([]string{"--merge"})
	if err = initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	content, err := os.ReadFile("CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read CLAUDE.md: %v", err)
	}
	contentStr := string(content)

	if !strings.HasPrefix(contentStr, "# My Amazing Project") {
		t.Error("H1 heading should remain at the start")
	}
	ctxIdx := strings.Index(contentStr, "ctx:context")
	buildIdx := strings.Index(contentStr, "## Build Instructions")
	if ctxIdx == -1 {
		t.Fatal("ctx:context marker not found")
	}
	if buildIdx == -1 {
		t.Fatal("Build Instructions section not found")
	}
	if ctxIdx > buildIdx {
		t.Error("ctx content should appear before Build Instructions, not after")
	}
	if !strings.Contains(contentStr, "Run make build") {
		t.Error("original content was lost")
	}
}

func TestInitMergeInsertsAtTopWhenNoH1(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-init-merge-no-h1-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()
	t.Setenv("HOME", tmpDir)
	t.Setenv(env.CtxDir, filepath.Join(tmpDir, ".context"))
	t.Setenv(env.SkipPathCheck, env.True)

	existingContent := "## Build Instructions\n\nRun make build.\n\n" +
		"## Testing\n\nRun make test.\n"
	if err = os.WriteFile("CLAUDE.md", []byte(existingContent), 0600); err != nil {
		t.Fatalf("failed to create CLAUDE.md: %v", err)
	}

	initCmd := Cmd()
	initCmd.SetArgs([]string{"--merge"})
	if err = initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	content, err := os.ReadFile("CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read CLAUDE.md: %v", err)
	}
	contentStr := string(content)

	ctxIdx := strings.Index(contentStr, "ctx:context")
	buildIdx := strings.Index(contentStr, "## Build Instructions")
	if ctxIdx == -1 {
		t.Fatal("ctx:context marker not found")
	}
	if buildIdx == -1 {
		t.Fatal("Build Instructions section not found")
	}
	if ctxIdx > buildIdx {
		t.Error("ctx content should appear at top, before Build Instructions")
	}
	if !strings.Contains(contentStr, "Run make test") {
		t.Error("original content was lost")
	}
}

func TestInitCreatesPermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-init-perms-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()
	t.Setenv("HOME", tmpDir)
	t.Setenv(env.CtxDir, filepath.Join(tmpDir, ".context"))
	t.Setenv(env.SkipPathCheck, env.True)

	cmd := Cmd()
	cmd.SetArgs([]string{})
	if err = cmd.Execute(); err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	settingsPath := filepath.Join(tmpDir, ".claude", "settings.local.json")
	content, err := os.ReadFile(filepath.Clean(settingsPath))
	if err != nil {
		t.Fatalf("failed to read settings.local.json: %v", err)
	}

	var settings claude.Settings
	if err := json.Unmarshal(content, &settings); err != nil {
		t.Fatalf("failed to parse settings.local.json: %v", err)
	}

	permSet := make(map[string]bool)
	for _, p := range settings.Permissions.Allow {
		permSet[p] = true
	}
	requiredPerms := []string{
		"Bash(ctx:*)", "Skill(ctx-agent)",
		"Skill(ctx-learning-add)",
	}
	for _, p := range requiredPerms {
		if !permSet[p] {
			t.Errorf("missing required permission: %s", p)
		}
	}

	denySet := make(map[string]bool)
	for _, d := range settings.Permissions.Deny {
		denySet[d] = true
	}
	if !denySet["Bash(sudo *)"] {
		t.Error("missing deny rule: Bash(sudo *)")
	}
	if !denySet["Bash(git push *)"] {
		t.Error("missing deny rule: Bash(git push *)")
	}
}

func TestInitMergesPermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-init-merge-perms-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()
	t.Setenv("HOME", tmpDir)
	t.Setenv(env.CtxDir, filepath.Join(tmpDir, ".context"))
	t.Setenv(env.SkipPathCheck, env.True)

	if err = os.MkdirAll(".claude", 0750); err != nil {
		t.Fatalf("failed to create .claude: %v", err)
	}

	existingSettings := claude.Settings{
		Permissions: claude.PermissionsConfig{
			Allow: []string{"Bash(git status:*)", "Bash(make build:*)", "Bash(ctx:*)"},
			Deny:  []string{"Bash(custom-block *)"},
		},
	}
	existingJSON, _ := json.MarshalIndent(existingSettings, "", "  ")
	settingsPath := ".claude/settings.local.json"
	if err = os.WriteFile(settingsPath, existingJSON, 0600); err != nil {
		t.Fatalf("failed to write settings: %v", err)
	}

	cmd := Cmd()
	cmd.SetArgs([]string{})
	if err = cmd.Execute(); err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	content, err := os.ReadFile(".claude/settings.local.json")
	if err != nil {
		t.Fatalf("failed to read settings: %v", err)
	}

	var settings claude.Settings
	if err := json.Unmarshal(content, &settings); err != nil {
		t.Fatalf("failed to parse settings: %v", err)
	}

	permSet := make(map[string]bool)
	for _, p := range settings.Permissions.Allow {
		permSet[p] = true
	}
	if !permSet["Bash(git status:*)"] {
		t.Error("existing permission 'Bash(git status:*)' was removed")
	}
	if !permSet["Skill(ctx-agent)"] {
		t.Error("missing new permission 'Skill(ctx-agent)'")
	}

	denySet := make(map[string]bool)
	for _, d := range settings.Permissions.Deny {
		denySet[d] = true
	}
	if !denySet["Bash(custom-block *)"] {
		t.Error("existing deny rule 'Bash(custom-block *)' was removed")
	}
	if !denySet["Bash(sudo *)"] {
		t.Error("missing default deny rule 'Bash(sudo *)'")
	}
}

func TestInitWithExistingClaudeMdWithCtxMarker(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-init-existing-claude-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()
	t.Setenv("HOME", tmpDir)
	t.Setenv(env.CtxDir, filepath.Join(tmpDir, ".context"))
	t.Setenv(env.SkipPathCheck, env.True)

	existingContent := "# My Project\n\n" +
		"This is my existing CLAUDE.md content.\n\n" +
		"<!-- ctx:context -->\nOld ctx content here\n" +
		"<!-- ctx:end -->\n\n## Custom Section\n\n" +
		"Some custom content here.\n"
	if err = os.WriteFile("CLAUDE.md", []byte(existingContent), 0600); err != nil {
		t.Fatalf("failed to create CLAUDE.md: %v", err)
	}

	initCmd := Cmd()
	initCmd.SetArgs([]string{})
	if err = initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	content, err := os.ReadFile("CLAUDE.md")
	if err != nil {
		t.Fatalf("failed to read CLAUDE.md: %v", err)
	}
	if !strings.Contains(string(content), "ctx:context") {
		t.Error("CLAUDE.md missing ctx:context marker")
	}
	if !strings.Contains(string(content), "Custom Section") {
		t.Error("CLAUDE.md lost custom section")
	}
}

func TestCmd_Flags(t *testing.T) {
	cmd := Cmd()
	if cmd == nil {
		t.Fatal("Cmd() returned nil")
	}
	if cmd.Use != "init" {
		t.Errorf("Cmd().Use = %q, want %q", cmd.Use, "init")
	}
	flags := []string{"force", "minimal", "merge"}
	for _, f := range flags {
		if cmd.Flags().Lookup(f) == nil {
			t.Errorf("missing --%s flag", f)
		}
	}
	if cmd.Flags().ShorthandLookup("f") == nil {
		t.Error("missing -f shorthand for --force")
	}
	if cmd.Flags().ShorthandLookup("m") == nil {
		t.Error("missing -m shorthand for --minimal")
	}
}

func TestRunInit_Minimal(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ctx-init-minimal-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()
	t.Setenv("HOME", tmpDir)
	t.Setenv(env.CtxDir, filepath.Join(tmpDir, ".context"))
	t.Setenv(env.SkipPathCheck, env.True)

	cmd := Cmd()
	cmd.SetArgs([]string{"--minimal"})
	if err = cmd.Execute(); err != nil {
		t.Fatalf("init --minimal failed: %v", err)
	}

	for _, name := range ctx.FilesRequired {
		path := filepath.Join(".context", name)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("required file %s missing with --minimal: %v", name, err)
		}
	}

	glossaryPath := filepath.Join(".context", ctx.Glossary)
	if _, err := os.Stat(glossaryPath); err == nil {
		t.Error("GLOSSARY.md should not exist with --minimal")
	}
}

func TestRunInit_Force(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ctx-init-force-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()
	t.Setenv("HOME", tmpDir)
	t.Setenv(env.CtxDir, filepath.Join(tmpDir, ".context"))
	t.Setenv(env.SkipPathCheck, env.True)

	cmd := Cmd()
	cmd.SetArgs([]string{})
	if err = cmd.Execute(); err != nil {
		t.Fatalf("first init failed: %v", err)
	}

	cmd2 := Cmd()
	cmd2.SetArgs([]string{"--force"})
	if err = cmd2.Execute(); err != nil {
		t.Fatalf("init --force failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(".context", ctx.Constitution)); err != nil {
		t.Error("CONSTITUTION.md missing after force reinit")
	}
}

func TestRunInit_Merge(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ctx-init-merge-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()
	t.Setenv("HOME", tmpDir)
	t.Setenv(env.CtxDir, filepath.Join(tmpDir, ".context"))
	t.Setenv(env.SkipPathCheck, env.True)

	mdContent := "# My Project\n\nExisting.\n"
	if err = os.WriteFile(cfgClaude.Md, []byte(mdContent), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := Cmd()
	cmd.SetArgs([]string{"--merge"})
	if err = cmd.Execute(); err != nil {
		t.Fatalf("init --merge failed: %v", err)
	}

	content, _ := os.ReadFile(cfgClaude.Md)
	if !strings.Contains(string(content), "My Project") {
		t.Error("original content lost with --merge")
	}
}

// TestInitScaffoldsFoundationSteeringFiles verifies that
// `ctx init` (without --no-steering-init) populates
// .context/steering/ with the four foundation files, each
// containing the inline inclusion-mode guidance comment.
func TestInitScaffoldsFoundationSteeringFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-init-steering-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()
	t.Setenv("HOME", tmpDir)
	t.Setenv(env.CtxDir, filepath.Join(tmpDir, ".context"))
	t.Setenv(env.SkipPathCheck, env.True)

	cmd := Cmd()
	cmd.SetArgs([]string{})
	if err = cmd.Execute(); err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	steeringDir := filepath.Join(tmpDir, ".context", "steering")
	foundationFiles := []string{
		"product.md", "tech.md",
		"structure.md", "workflow.md",
	}

	for _, name := range foundationFiles {
		path := filepath.Join(steeringDir, name)
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			t.Errorf("%s was not created: %v", name, readErr)
			continue
		}
		if !strings.Contains(string(data), "inclusion: always") {
			t.Errorf("%s missing inclusion: always frontmatter", name)
		}
		if !strings.Contains(string(data), "This is a ctx steering file") {
			t.Errorf("%s missing inline mode guidance comment", name)
		}
	}
}

// TestInitNoSteeringInitFlagSkipsScaffold verifies that
// --no-steering-init suppresses foundation file scaffolding
// while still creating the empty directory.
func TestInitNoSteeringInitFlagSkipsScaffold(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cli-init-no-steering-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	origDir, _ := os.Getwd()
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()
	t.Setenv("HOME", tmpDir)
	t.Setenv(env.CtxDir, filepath.Join(tmpDir, ".context"))
	t.Setenv(env.SkipPathCheck, env.True)

	cmd := Cmd()
	cmd.SetArgs([]string{"--no-steering-init"})
	if err = cmd.Execute(); err != nil {
		t.Fatalf("init --no-steering-init failed: %v", err)
	}

	steeringDir := filepath.Join(tmpDir, ".context", "steering")
	info, statErr := os.Stat(steeringDir)
	if statErr != nil || !info.IsDir() {
		t.Errorf(".context/steering should still be created")
	}
	entries, readErr := os.ReadDir(steeringDir)
	if readErr != nil {
		t.Fatalf("failed to read steering dir: %v", readErr)
	}
	if len(entries) != 0 {
		t.Errorf(
			"expected empty steering dir with --no-steering-init, got %d entries",
			len(entries),
		)
	}
}
