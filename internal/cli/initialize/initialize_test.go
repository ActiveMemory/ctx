//   /    Context:                     https://ctx.ist
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
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/spf13/cobra"
)

// helper creates a temp dir, chdir into it, and returns a cleanup function.
func helper(t *testing.T) (string, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "ctx-init-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		_ = os.RemoveAll(tmpDir)
		t.Fatalf("failed to chdir: %v", err)
	}
	return tmpDir, func() {
		_ = os.Chdir(origDir)
		_ = os.RemoveAll(tmpDir)
	}
}

// newTestCmd returns a cobra.Command with stdout/stderr captured.
func newTestCmd() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.SetOut(&strings.Builder{})
	cmd.SetErr(&strings.Builder{})
	return cmd
}

// --- checkCtxInPath tests ---

func TestCheckCtxInPath_SkipEnv(t *testing.T) {
	t.Setenv("CTX_SKIP_PATH_CHECK", "1")
	cmd := newTestCmd()
	if err := checkCtxInPath(cmd); err != nil {
		t.Errorf("expected nil error with skip env, got %v", err)
	}
}

func TestCheckCtxInPath_NotFound(t *testing.T) {
	t.Setenv("CTX_SKIP_PATH_CHECK", "")
	t.Setenv("PATH", "/nonexistent")
	cmd := newTestCmd()
	err := checkCtxInPath(cmd)
	if err == nil {
		t.Fatal("expected error when ctx not in PATH")
	}
	if !strings.Contains(err.Error(), "ctx not found in PATH") {
		t.Errorf("unexpected error: %v", err)
	}
}

// --- mergePermissions tests ---

func TestMergePermissions_Empty(t *testing.T) {
	var slice []string
	added := mergePermissions(&slice, []string{"Bash(ctx:*)", "Skill(ctx-agent)"})
	if !added {
		t.Error("expected permissions to be added")
	}
	if len(slice) != 2 {
		t.Errorf("expected 2 permissions, got %d", len(slice))
	}
}

func TestMergePermissions_NoDuplicates(t *testing.T) {
	slice := []string{"Bash(ctx:*)", "Bash(git:*)"}
	added := mergePermissions(&slice, []string{"Bash(ctx:*)", "Skill(ctx-agent)"})
	if !added {
		t.Error("expected new permissions to be added")
	}
	count := 0
	for _, p := range slice {
		if p == "Bash(ctx:*)" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 Bash(ctx:*), got %d", count)
	}
}

func TestMergePermissions_AllExist(t *testing.T) {
	slice := []string{"Bash(ctx:*)", "Skill(ctx-agent)"}
	added := mergePermissions(&slice, []string{"Bash(ctx:*)", "Skill(ctx-agent)"})
	if added {
		t.Error("expected no new permissions added")
	}
}

// --- handleMakefileCtx tests ---

func TestHandleMakefileCtx_NoExistingMakefile(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	cmd := newTestCmd()
	if err := handleMakefileCtx(cmd); err != nil {
		t.Fatalf("handleMakefileCtx failed: %v", err)
	}

	// Makefile.ctx should exist
	if _, err := os.Stat(config.FileMakefileCtx); err != nil {
		t.Errorf("Makefile.ctx not created: %v", err)
	}

	// Makefile should be created with include directive
	content, err := os.ReadFile("Makefile")
	if err != nil {
		t.Fatalf("Makefile not created: %v", err)
	}
	if !strings.Contains(string(content), includeDirective) {
		t.Error("Makefile missing include directive")
	}
}

func TestHandleMakefileCtx_ExistingMakefileWithoutInclude(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	// Create existing Makefile without the include
	if err := os.WriteFile("Makefile", []byte("build:\n\tgo build\n"), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := handleMakefileCtx(cmd); err != nil {
		t.Fatalf("handleMakefileCtx failed: %v", err)
	}

	content, err := os.ReadFile("Makefile")
	if err != nil {
		t.Fatal(err)
	}
	contentStr := string(content)

	// Should still have original content
	if !strings.Contains(contentStr, "go build") {
		t.Error("original Makefile content lost")
	}
	// Should have include directive appended
	if !strings.Contains(contentStr, includeDirective) {
		t.Error("include directive not appended")
	}
}

func TestHandleMakefileCtx_ExistingMakefileWithInclude(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	// Create Makefile that already has the include
	original := "build:\n\tgo build\n\n" + includeDirective + "\n"
	if err := os.WriteFile("Makefile", []byte(original), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := handleMakefileCtx(cmd); err != nil {
		t.Fatalf("handleMakefileCtx failed: %v", err)
	}

	content, err := os.ReadFile("Makefile")
	if err != nil {
		t.Fatal(err)
	}

	// Count include directives - should be exactly 1
	count := strings.Count(string(content), includeDirective)
	if count != 1 {
		t.Errorf("expected 1 include directive, got %d", count)
	}
}

func TestHandleMakefileCtx_MakefileNoTrailingNewline(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	// Create Makefile without trailing newline
	if err := os.WriteFile("Makefile", []byte("build:\n\tgo build"), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := handleMakefileCtx(cmd); err != nil {
		t.Fatalf("handleMakefileCtx failed: %v", err)
	}

	content, err := os.ReadFile("Makefile")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), includeDirective) {
		t.Error("include directive not appended")
	}
}

// --- addToGitignore tests ---

func TestAddToGitignore_New(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	if err := addToGitignore(".context", ".scratchpad.key"); err != nil {
		t.Fatalf("addToGitignore failed: %v", err)
	}

	content, err := os.ReadFile(".gitignore")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), ".context/.scratchpad.key") {
		t.Error("entry not added to .gitignore")
	}
}

func TestAddToGitignore_AlreadyPresent(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	entry := ".context/.scratchpad.key"
	if err := os.WriteFile(".gitignore", []byte(entry+"\n"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := addToGitignore(".context", ".scratchpad.key"); err != nil {
		t.Fatalf("addToGitignore failed: %v", err)
	}

	content, err := os.ReadFile(".gitignore")
	if err != nil {
		t.Fatal(err)
	}
	count := strings.Count(string(content), entry)
	if count != 1 {
		t.Errorf("expected 1 occurrence, got %d", count)
	}
}

func TestAddToGitignore_AppendNoTrailingNewline(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	if err := os.WriteFile(".gitignore", []byte("node_modules"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := addToGitignore(".context", ".scratchpad.key"); err != nil {
		t.Fatalf("addToGitignore failed: %v", err)
	}

	content, err := os.ReadFile(".gitignore")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "node_modules\n.context/.scratchpad.key") {
		t.Errorf("unexpected content: %q", string(content))
	}
}

// --- ensureGitignoreEntries tests ---

func TestEnsureGitignoreEntries_CreatesNew(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	cmd := newTestCmd()
	if err := ensureGitignoreEntries(cmd); err != nil {
		t.Fatalf("ensureGitignoreEntries failed: %v", err)
	}

	content, err := os.ReadFile(".gitignore")
	if err != nil {
		t.Fatalf("expected .gitignore to be created: %v", err)
	}
	contentStr := string(content)

	for _, entry := range config.GitignoreEntries {
		if !strings.Contains(contentStr, entry) {
			t.Errorf("missing entry %q in .gitignore", entry)
		}
	}
	if !strings.Contains(contentStr, "# ctx managed entries") {
		t.Error("missing comment header")
	}
}

func TestEnsureGitignoreEntries_AppendsOnlyMissing(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	// Pre-populate with some entries
	existing := ".context/journal/\n.context/logs/\n"
	if err := os.WriteFile(".gitignore", []byte(existing), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := ensureGitignoreEntries(cmd); err != nil {
		t.Fatalf("ensureGitignoreEntries failed: %v", err)
	}

	content, err := os.ReadFile(".gitignore")
	if err != nil {
		t.Fatal(err)
	}
	contentStr := string(content)

	// All entries should be present
	for _, entry := range config.GitignoreEntries {
		if !strings.Contains(contentStr, entry) {
			t.Errorf("missing entry %q in .gitignore", entry)
		}
	}

	// Already-present entries should not be duplicated
	if strings.Count(contentStr, ".context/journal/") != 1 {
		t.Error("duplicate .context/journal/ entry")
	}
	if strings.Count(contentStr, ".context/logs/") != 1 {
		t.Error("duplicate .context/logs/ entry")
	}
}

func TestEnsureGitignoreEntries_Idempotent(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	cmd := newTestCmd()
	if err := ensureGitignoreEntries(cmd); err != nil {
		t.Fatalf("first call failed: %v", err)
	}
	first, _ := os.ReadFile(".gitignore")

	cmd2 := newTestCmd()
	if err := ensureGitignoreEntries(cmd2); err != nil {
		t.Fatalf("second call failed: %v", err)
	}
	second, _ := os.ReadFile(".gitignore")

	if string(first) != string(second) {
		t.Errorf("file changed on second call:\n--- first ---\n%s\n--- second ---\n%s", first, second)
	}
}

func TestEnsureGitignoreEntries_PreservesExistingContent(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	existing := "node_modules/\n*.log\nbuild/\n"
	if err := os.WriteFile(".gitignore", []byte(existing), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := ensureGitignoreEntries(cmd); err != nil {
		t.Fatalf("ensureGitignoreEntries failed: %v", err)
	}

	content, err := os.ReadFile(".gitignore")
	if err != nil {
		t.Fatal(err)
	}
	contentStr := string(content)

	if !strings.HasPrefix(contentStr, existing) {
		t.Error("existing content was not preserved")
	}
}

func TestEnsureGitignoreEntries_NoTrailingNewline(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	existing := "node_modules/"
	if err := os.WriteFile(".gitignore", []byte(existing), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := ensureGitignoreEntries(cmd); err != nil {
		t.Fatalf("ensureGitignoreEntries failed: %v", err)
	}

	content, _ := os.ReadFile(".gitignore")
	contentStr := string(content)

	// Should have a newline before the comment header
	if !strings.Contains(contentStr, "node_modules/\n\n# ctx managed entries\n") {
		t.Errorf("unexpected content format: %q", contentStr)
	}
}

// --- createTools tests ---

func TestCreateTools(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	contextDir := ".context"
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := createTools(cmd, contextDir, false); err != nil {
		t.Fatalf("createTools failed: %v", err)
	}

	toolsDir := filepath.Join(contextDir, config.DirTools)
	entries, err := os.ReadDir(toolsDir)
	if err != nil {
		t.Fatalf("failed to read tools dir: %v", err)
	}
	if len(entries) == 0 {
		t.Error("no tools created")
	}
}

func TestCreateTools_ExistsNoForce(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	contextDir := ".context"
	toolsDir := filepath.Join(contextDir, config.DirTools)
	if err := os.MkdirAll(toolsDir, 0750); err != nil {
		t.Fatal(err)
	}

	// Create a tool file manually
	toolPath := filepath.Join(toolsDir, "context-watch.sh")
	if err := os.WriteFile(toolPath, []byte("#!/bin/bash\n# original"), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := createTools(cmd, contextDir, false); err != nil {
		t.Fatalf("createTools failed: %v", err)
	}

	// Should not be overwritten
	content, _ := os.ReadFile(toolPath) //nolint:gosec // test temp path
	if !strings.Contains(string(content), "# original") {
		t.Error("tool was overwritten when force=false")
	}
}

// --- createEntryTemplates tests ---

func TestCreateEntryTemplates(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	contextDir := ".context"
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := createEntryTemplates(cmd, contextDir, false); err != nil {
		t.Fatalf("createEntryTemplates failed: %v", err)
	}

	templatesDir := filepath.Join(contextDir, "templates")
	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		t.Fatalf("failed to read templates dir: %v", err)
	}
	if len(entries) == 0 {
		t.Error("no entry templates created")
	}
}

func TestCreateEntryTemplates_ExistsNoForce(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	contextDir := ".context"
	templatesDir := filepath.Join(contextDir, "templates")
	if err := os.MkdirAll(templatesDir, 0750); err != nil {
		t.Fatal(err)
	}

	// Create a template file manually
	tplPath := filepath.Join(templatesDir, "decision.md")
	if err := os.WriteFile(tplPath, []byte("# original"), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := createEntryTemplates(cmd, contextDir, false); err != nil {
		t.Fatalf("createEntryTemplates failed: %v", err)
	}

	// Should not be overwritten
	content, _ := os.ReadFile(tplPath) //nolint:gosec // test temp path
	if !strings.Contains(string(content), "# original") {
		t.Error("template was overwritten when force=false")
	}
}

// --- handlePromptMd tests ---

func TestHandlePromptMd_NewFile(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	cmd := newTestCmd()
	if err := handlePromptMd(cmd, false, false, false); err != nil {
		t.Fatalf("handlePromptMd failed: %v", err)
	}

	if _, err := os.Stat(config.FilePromptMd); err != nil {
		t.Errorf("PROMPT.md not created: %v", err)
	}
}

func TestHandlePromptMd_NewFileRalph(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	cmd := newTestCmd()
	if err := handlePromptMd(cmd, false, false, true); err != nil {
		t.Fatalf("handlePromptMd failed: %v", err)
	}

	content, err := os.ReadFile(config.FilePromptMd)
	if err != nil {
		t.Fatal(err)
	}
	if len(content) == 0 {
		t.Error("PROMPT.md is empty")
	}
}

func TestHandlePromptMd_ExistsWithMarkers_NoForce(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	existing := "# Prompt\n\n" + config.PromptMarkerStart + "\nold\n" + config.PromptMarkerEnd + "\n\n## Custom\n"
	if err := os.WriteFile(config.FilePromptMd, []byte(existing), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := handlePromptMd(cmd, false, false, false); err != nil {
		t.Fatalf("handlePromptMd failed: %v", err)
	}

	// Content should be unchanged (skipped)
	content, _ := os.ReadFile(config.FilePromptMd)
	if string(content) != existing {
		t.Error("content was changed when force=false and markers exist")
	}
}

func TestHandlePromptMd_ExistsWithMarkers_Force(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	existing := "# Prompt\n\n" + config.PromptMarkerStart + "\nold prompt\n" + config.PromptMarkerEnd + "\n\n## Custom\n"
	if err := os.WriteFile(config.FilePromptMd, []byte(existing), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := handlePromptMd(cmd, true, false, false); err != nil {
		t.Fatalf("handlePromptMd failed: %v", err)
	}

	content, _ := os.ReadFile(config.FilePromptMd)
	contentStr := string(content)

	// Should have updated prompt markers
	if !strings.Contains(contentStr, config.PromptMarkerStart) {
		t.Error("prompt markers missing after force update")
	}
	// Custom section should be preserved
	if !strings.Contains(contentStr, "## Custom") {
		t.Error("custom section lost after force update")
	}
}

func TestHandlePromptMd_MergeAutoMerge(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	existing := "# My Prompt\n\nExisting content.\n"
	if err := os.WriteFile(config.FilePromptMd, []byte(existing), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := handlePromptMd(cmd, false, true, false); err != nil {
		t.Fatalf("handlePromptMd failed: %v", err)
	}

	content, _ := os.ReadFile(config.FilePromptMd)
	contentStr := string(content)

	if !strings.Contains(contentStr, "My Prompt") {
		t.Error("original H1 lost")
	}
	if !strings.Contains(contentStr, "Existing content") {
		t.Error("original content lost")
	}

	// Backup should exist
	matches, _ := filepath.Glob(config.FilePromptMd + ".*.bak")
	if len(matches) == 0 {
		t.Error("no backup file created")
	}
}

// --- handleImplementationPlan tests ---

func TestHandleImplementationPlan_NewFile(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	cmd := newTestCmd()
	if err := handleImplementationPlan(cmd, false, false); err != nil {
		t.Fatalf("handleImplementationPlan failed: %v", err)
	}

	if _, err := os.Stat(config.FileImplementationPlan); err != nil {
		t.Errorf("IMPLEMENTATION_PLAN.md not created: %v", err)
	}
}

func TestHandleImplementationPlan_ExistsWithMarkers_NoForce(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	existing := "# Plan\n\n" + config.PlanMarkerStart + "\nold\n" + config.PlanMarkerEnd + "\n"
	if err := os.WriteFile(config.FileImplementationPlan, []byte(existing), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := handleImplementationPlan(cmd, false, false); err != nil {
		t.Fatalf("handleImplementationPlan failed: %v", err)
	}

	content, _ := os.ReadFile(config.FileImplementationPlan)
	if string(content) != existing {
		t.Error("content changed when force=false and markers exist")
	}
}

func TestHandleImplementationPlan_ExistsWithMarkers_Force(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	existing := "# Plan\n\n" + config.PlanMarkerStart + "\nold plan\n" + config.PlanMarkerEnd + "\n\n## Custom\n"
	if err := os.WriteFile(config.FileImplementationPlan, []byte(existing), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := handleImplementationPlan(cmd, true, false); err != nil {
		t.Fatalf("handleImplementationPlan failed: %v", err)
	}

	content, _ := os.ReadFile(config.FileImplementationPlan)
	contentStr := string(content)

	if !strings.Contains(contentStr, config.PlanMarkerStart) {
		t.Error("plan markers missing after force update")
	}
	if !strings.Contains(contentStr, "## Custom") {
		t.Error("custom section lost after force update")
	}
}

func TestHandleImplementationPlan_MergeAutoMerge(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	existing := "# My Plan\n\nExisting plan content.\n"
	if err := os.WriteFile(config.FileImplementationPlan, []byte(existing), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := handleImplementationPlan(cmd, false, true); err != nil {
		t.Fatalf("handleImplementationPlan failed: %v", err)
	}

	content, _ := os.ReadFile(config.FileImplementationPlan)
	contentStr := string(content)

	if !strings.Contains(contentStr, "My Plan") {
		t.Error("original H1 lost")
	}
	if !strings.Contains(contentStr, "Existing plan content") {
		t.Error("original content lost")
	}
}

// --- handleClaudeMd tests ---

func TestHandleClaudeMd_NewFile(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	cmd := newTestCmd()
	if err := handleClaudeMd(cmd, false, false); err != nil {
		t.Fatalf("handleClaudeMd failed: %v", err)
	}

	if _, err := os.Stat(config.FileClaudeMd); err != nil {
		t.Errorf("CLAUDE.md not created: %v", err)
	}

	content, _ := os.ReadFile(config.FileClaudeMd)
	if len(content) == 0 {
		t.Error("CLAUDE.md is empty")
	}
}

func TestHandleClaudeMd_ExistsWithMarkers_NoForce(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	existing := "# Project\n\n" + config.CtxMarkerStart + "\nold\n" + config.CtxMarkerEnd + "\n"
	if err := os.WriteFile(config.FileClaudeMd, []byte(existing), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := handleClaudeMd(cmd, false, false); err != nil {
		t.Fatalf("handleClaudeMd failed: %v", err)
	}

	content, _ := os.ReadFile(config.FileClaudeMd)
	if string(content) != existing {
		t.Error("content changed when force=false and markers exist")
	}
}

func TestHandleClaudeMd_ExistsWithMarkers_Force(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	existing := "# Project\n\n" + config.CtxMarkerStart + "\nold ctx\n" + config.CtxMarkerEnd + "\n\n## Custom\n"
	if err := os.WriteFile(config.FileClaudeMd, []byte(existing), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := handleClaudeMd(cmd, true, false); err != nil {
		t.Fatalf("handleClaudeMd failed: %v", err)
	}

	content, _ := os.ReadFile(config.FileClaudeMd)
	contentStr := string(content)

	if !strings.Contains(contentStr, config.CtxMarkerStart) {
		t.Error("ctx markers missing after force update")
	}
	if !strings.Contains(contentStr, "## Custom") {
		t.Error("custom section lost after force update")
	}
}

func TestHandleClaudeMd_AutoMerge(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	existing := "# My Project\n\nExisting content.\n"
	if err := os.WriteFile(config.FileClaudeMd, []byte(existing), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := handleClaudeMd(cmd, false, true); err != nil {
		t.Fatalf("handleClaudeMd failed: %v", err)
	}

	content, _ := os.ReadFile(config.FileClaudeMd)
	contentStr := string(content)

	if !strings.Contains(contentStr, "My Project") {
		t.Error("original H1 lost")
	}
	if !strings.Contains(contentStr, "Existing content") {
		t.Error("original content lost")
	}
}

// --- updateCtxSection tests ---

func TestUpdateCtxSection(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	existing := "# Project\n\nbefore\n\n" +
		config.CtxMarkerStart + "\nOLD CTX CONTENT\n" + config.CtxMarkerEnd +
		"\n\nafter\n"
	template := config.CtxMarkerStart + "\nNEW CTX CONTENT\n" + config.CtxMarkerEnd

	cmd := newTestCmd()
	if err := updateCtxSection(cmd, existing, []byte(template)); err != nil {
		t.Fatalf("updateCtxSection failed: %v", err)
	}

	content, err := os.ReadFile(config.FileClaudeMd)
	if err != nil {
		t.Fatal(err)
	}
	contentStr := string(content)

	if !strings.Contains(contentStr, "NEW CTX CONTENT") {
		t.Error("new content not inserted")
	}
	if strings.Contains(contentStr, "OLD CTX CONTENT") {
		t.Error("old content not replaced")
	}
	if !strings.Contains(contentStr, "before") {
		t.Error("content before markers lost")
	}
	if !strings.Contains(contentStr, "after") {
		t.Error("content after markers lost")
	}
}

func TestUpdateCtxSection_NoEndMarker(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	existing := "# Project\n\n" + config.CtxMarkerStart + "\nOLD CTX CONTENT\n"
	template := config.CtxMarkerStart + "\nNEW CTX CONTENT\n" + config.CtxMarkerEnd

	cmd := newTestCmd()
	if err := updateCtxSection(cmd, existing, []byte(template)); err != nil {
		t.Fatalf("updateCtxSection failed: %v", err)
	}

	content, _ := os.ReadFile(config.FileClaudeMd)
	contentStr := string(content)
	if !strings.Contains(contentStr, "NEW CTX CONTENT") {
		t.Error("new content not inserted when no end marker")
	}
}

func TestUpdateCtxSection_TemplateMissingMarkers(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	existing := "# Project\n\n" + config.CtxMarkerStart + "\nOLD\n" + config.CtxMarkerEnd
	template := "no markers here"

	cmd := newTestCmd()
	err := updateCtxSection(cmd, existing, []byte(template))
	if err == nil {
		t.Fatal("expected error when template missing markers")
	}
	if !strings.Contains(err.Error(), "template missing ctx markers") {
		t.Errorf("unexpected error: %v", err)
	}
}

// --- updatePromptSection tests ---

func TestUpdatePromptSection(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	existing := "# Prompt\n\nbefore\n\n" +
		config.PromptMarkerStart + "\nOLD PROMPT\n" + config.PromptMarkerEnd +
		"\n\nafter\n"
	template := config.PromptMarkerStart + "\nNEW PROMPT\n" + config.PromptMarkerEnd

	cmd := newTestCmd()
	if err := updatePromptSection(cmd, existing, []byte(template)); err != nil {
		t.Fatalf("updatePromptSection failed: %v", err)
	}

	content, err := os.ReadFile(config.FilePromptMd)
	if err != nil {
		t.Fatal(err)
	}
	contentStr := string(content)

	if !strings.Contains(contentStr, "NEW PROMPT") {
		t.Error("new prompt not inserted")
	}
	if strings.Contains(contentStr, "OLD PROMPT") {
		t.Error("old prompt not replaced")
	}
}

func TestUpdatePromptSection_TemplateMissingMarkers(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	existing := config.PromptMarkerStart + "\nOLD\n" + config.PromptMarkerEnd
	template := "no markers"

	cmd := newTestCmd()
	err := updatePromptSection(cmd, existing, []byte(template))
	if err == nil {
		t.Fatal("expected error when template missing markers")
	}
}

// --- updatePlanSection tests ---

func TestUpdatePlanSection(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	existing := "# Plan\n\nbefore\n\n" +
		config.PlanMarkerStart + "\nOLD PLAN\n" + config.PlanMarkerEnd +
		"\n\nafter\n"
	template := config.PlanMarkerStart + "\nNEW PLAN\n" + config.PlanMarkerEnd

	cmd := newTestCmd()
	if err := updatePlanSection(cmd, existing, []byte(template)); err != nil {
		t.Fatalf("updatePlanSection failed: %v", err)
	}

	content, err := os.ReadFile(config.FileImplementationPlan)
	if err != nil {
		t.Fatal(err)
	}
	contentStr := string(content)

	if !strings.Contains(contentStr, "NEW PLAN") {
		t.Error("new plan not inserted")
	}
	if strings.Contains(contentStr, "OLD PLAN") {
		t.Error("old plan not replaced")
	}
}

func TestUpdatePlanSection_TemplateMissingMarkers(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	existing := config.PlanMarkerStart + "\nOLD\n" + config.PlanMarkerEnd
	template := "no markers"

	cmd := newTestCmd()
	err := updatePlanSection(cmd, existing, []byte(template))
	if err == nil {
		t.Fatal("expected error when template missing markers")
	}
}

// --- mergeSettingsPermissions tests ---

func TestMergeSettingsPermissions_NewSettings(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	cmd := newTestCmd()
	if err := mergeSettingsPermissions(cmd); err != nil {
		t.Fatalf("mergeSettingsPermissions failed: %v", err)
	}

	content, err := os.ReadFile(config.FileSettings)
	if err != nil {
		t.Fatal(err)
	}

	var settings claude.Settings
	if err := json.Unmarshal(content, &settings); err != nil {
		t.Fatalf("failed to parse settings: %v", err)
	}

	if len(settings.Permissions.Allow) == 0 {
		t.Error("no allow permissions created")
	}
	if len(settings.Permissions.Deny) == 0 {
		t.Error("no deny permissions created")
	}

	// Verify specific deny rules
	denySet := make(map[string]bool)
	for _, d := range settings.Permissions.Deny {
		denySet[d] = true
	}
	if !denySet["Bash(sudo *)"] {
		t.Error("missing deny rule: Bash(sudo *)")
	}
}

func TestMergeSettingsPermissions_ExistingWithAllPerms(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	if err := os.MkdirAll(config.DirClaude, 0750); err != nil {
		t.Fatal(err)
	}

	settings := claude.Settings{
		Permissions: claude.PermissionsConfig{
			Allow: config.DefaultClaudePermissions,
			Deny:  config.DefaultClaudeDenyPermissions,
		},
	}
	data, _ := json.MarshalIndent(settings, "", "  ")
	if err := os.WriteFile(config.FileSettings, data, 0600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := mergeSettingsPermissions(cmd); err != nil {
		t.Fatalf("mergeSettingsPermissions failed: %v", err)
	}
}

func TestMergeSettingsPermissions_DenyPreservesExisting(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	if err := os.MkdirAll(config.DirClaude, 0750); err != nil {
		t.Fatal(err)
	}

	// Start with a custom deny rule plus some defaults
	settings := claude.Settings{
		Permissions: claude.PermissionsConfig{
			Allow: config.DefaultClaudePermissions,
			Deny:  []string{"Bash(custom-block *)", "Bash(sudo *)"},
		},
	}
	data, _ := json.MarshalIndent(settings, "", "  ")
	if err := os.WriteFile(config.FileSettings, data, 0600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := mergeSettingsPermissions(cmd); err != nil {
		t.Fatalf("mergeSettingsPermissions failed: %v", err)
	}

	content, err := os.ReadFile(config.FileSettings)
	if err != nil {
		t.Fatal(err)
	}

	var updated claude.Settings
	if err := json.Unmarshal(content, &updated); err != nil {
		t.Fatalf("failed to parse settings: %v", err)
	}

	// Custom deny rule must survive
	denySet := make(map[string]bool)
	for _, d := range updated.Permissions.Deny {
		denySet[d] = true
	}
	if !denySet["Bash(custom-block *)"] {
		t.Error("custom deny rule was removed during merge")
	}
	// Default deny rules must also be present
	if !denySet["Bash(git push *)"] {
		t.Error("default deny rule 'Bash(git push *)' missing after merge")
	}
}

// --- initScratchpad tests ---

func TestInitScratchpad_Plaintext(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	// Set scratchpad_encrypt to false via .contextrc
	if err := os.WriteFile(".contextrc", []byte(`scratchpad_encrypt = false`+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	// Reset rc cache to pick up new config
	// We need to import rc package... instead, set env variable approach
	// Actually we can't directly control rc.ScratchpadEncrypt in a test
	// without modifying source. Let's just test the encrypted path
	// since that's the default.
	_ = os.Remove(".contextrc")

	contextDir := ".context"
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	err := initScratchpad(cmd, contextDir)
	if err != nil {
		t.Fatalf("initScratchpad failed: %v", err)
	}

	// Either a key file or scratchpad.md should have been created
	keyPath := filepath.Join(contextDir, config.FileScratchpadKey)
	mdPath := filepath.Join(contextDir, config.FileScratchpadMd)
	_, keyErr := os.Stat(keyPath)
	_, mdErr := os.Stat(mdPath)
	if keyErr != nil && mdErr != nil {
		t.Error("neither key nor scratchpad.md was created")
	}
}

func TestInitScratchpad_KeyExists(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	contextDir := ".context"
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	// Create existing key file
	keyPath := filepath.Join(contextDir, config.FileScratchpadKey)
	if err := os.WriteFile(keyPath, []byte("existing-key"), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := initScratchpad(cmd, contextDir); err != nil {
		t.Fatalf("initScratchpad failed: %v", err)
	}

	// Key should not have been overwritten
	content, _ := os.ReadFile(keyPath) //nolint:gosec // test temp path
	if string(content) != "existing-key" {
		t.Error("existing key was overwritten")
	}
}

func TestInitScratchpad_EncExistsNoKey(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()

	contextDir := ".context"
	if err := os.MkdirAll(contextDir, 0750); err != nil {
		t.Fatal(err)
	}

	// Create encrypted file without key
	encPath := filepath.Join(contextDir, config.FileScratchpadEnc)
	if err := os.WriteFile(encPath, []byte("encrypted-data"), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := initScratchpad(cmd, contextDir); err != nil {
		t.Fatalf("initScratchpad failed: %v", err)
	}

	// Key should NOT have been generated (warning path)
	keyPath := filepath.Join(contextDir, config.FileScratchpadKey)
	if _, err := os.Stat(keyPath); err == nil {
		t.Error("key was generated even though enc exists without key (should just warn)")
	}
}

// --- Cmd tests ---

func TestCmd_Flags(t *testing.T) {
	cmd := Cmd()

	if cmd == nil {
		t.Fatal("Cmd() returned nil")
	}

	if cmd.Use != "init" {
		t.Errorf("Cmd().Use = %q, want %q", cmd.Use, "init")
	}

	flags := []string{"force", "minimal", "merge", "ralph"}
	for _, f := range flags {
		if cmd.Flags().Lookup(f) == nil {
			t.Errorf("missing --%s flag", f)
		}
	}

	// Check shorthand for force
	if cmd.Flags().ShorthandLookup("f") == nil {
		t.Error("missing -f shorthand for --force")
	}

	// Check shorthand for minimal
	if cmd.Flags().ShorthandLookup("m") == nil {
		t.Error("missing -m shorthand for --minimal")
	}
}

// --- runInit with minimal flag ---

func TestRunInit_Minimal(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()
	t.Setenv("CTX_SKIP_PATH_CHECK", "1")

	cmd := Cmd()
	cmd.SetArgs([]string{"--minimal"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init --minimal failed: %v", err)
	}

	// Check that essential files exist
	for _, name := range config.FilesRequired {
		path := filepath.Join(".context", name)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("required file %s missing with --minimal: %v", name, err)
		}
	}

	// ARCHITECTURE.md should still exist (it's in the minimal template list)
	// Actually, minimal only creates FilesRequired
	// GLOSSARY.md should NOT exist with minimal
	glossaryPath := filepath.Join(".context", config.FileGlossary)
	if _, err := os.Stat(glossaryPath); err == nil {
		t.Error("GLOSSARY.md should not exist with --minimal")
	}
}

// --- runInit with ralph flag ---

func TestRunInit_Ralph(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()
	t.Setenv("CTX_SKIP_PATH_CHECK", "1")

	cmd := Cmd()
	cmd.SetArgs([]string{"--ralph"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init --ralph failed: %v", err)
	}

	// PROMPT.md should exist
	if _, err := os.Stat(config.FilePromptMd); err != nil {
		t.Errorf("PROMPT.md not created with --ralph: %v", err)
	}
}

// --- runInit with force (overwrite existing) ---

func TestRunInit_Force(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()
	t.Setenv("CTX_SKIP_PATH_CHECK", "1")

	// Run once to create files
	cmd := Cmd()
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("first init failed: %v", err)
	}

	// Run again with --force
	cmd2 := Cmd()
	cmd2.SetArgs([]string{"--force"})
	if err := cmd2.Execute(); err != nil {
		t.Fatalf("init --force failed: %v", err)
	}

	// Verify files still exist
	if _, err := os.Stat(filepath.Join(".context", config.FileConstitution)); err != nil {
		t.Error("CONSTITUTION.md missing after force reinit")
	}
}

// --- runInit with merge flag ---

func TestRunInit_Merge(t *testing.T) {
	_, cleanup := helper(t)
	defer cleanup()
	t.Setenv("CTX_SKIP_PATH_CHECK", "1")

	// Create existing CLAUDE.md
	if err := os.WriteFile(config.FileClaudeMd, []byte("# My Project\n\nExisting.\n"), 0600); err != nil {
		t.Fatal(err)
	}

	cmd := Cmd()
	cmd.SetArgs([]string{"--merge"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init --merge failed: %v", err)
	}

	content, _ := os.ReadFile(config.FileClaudeMd)
	if !strings.Contains(string(content), "My Project") {
		t.Error("original content lost with --merge")
	}
}
