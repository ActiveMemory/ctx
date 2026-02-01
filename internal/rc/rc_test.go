//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
)

func TestDefaultRC(t *testing.T) {
	rc := Default()

	if rc.ContextDir != config.DirContext {
		t.Errorf("ContextDir = %q, want %q", rc.ContextDir, config.DirContext)
	}
	if rc.TokenBudget != DefaultTokenBudget {
		t.Errorf("TokenBudget = %d, want %d", rc.TokenBudget, DefaultTokenBudget)
	}
	if rc.PriorityOrder != nil {
		t.Errorf("PriorityOrder = %v, want nil", rc.PriorityOrder)
	}
	if !rc.AutoArchive {
		t.Error("AutoArchive = false, want true")
	}
	if rc.ArchiveAfterDays != DefaultArchiveAfterDays {
		t.Errorf("ArchiveAfterDays = %d, want %d", rc.ArchiveAfterDays, DefaultArchiveAfterDays)
	}
}

func TestGetRC_NoFile(t *testing.T) {
	// Change to temp directory with no .contextrc
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origDir)

	Reset()

	rc := RC()

	if rc.ContextDir != config.DirContext {
		t.Errorf("ContextDir = %q, want %q", rc.ContextDir, config.DirContext)
	}
	if rc.TokenBudget != DefaultTokenBudget {
		t.Errorf("TokenBudget = %d, want %d", rc.TokenBudget, DefaultTokenBudget)
	}
}

func TestGetRC_WithFile(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origDir)

	// Create .contextrc file
	rcContent := `context_dir: custom-context
token_budget: 4000
priority_order:
  - TASKS.md
  - DECISIONS.md
auto_archive: false
archive_after_days: 14
`
	os.WriteFile(filepath.Join(tempDir, ".contextrc"), []byte(rcContent), 0644)

	Reset()

	rc := RC()

	if rc.ContextDir != "custom-context" {
		t.Errorf("ContextDir = %q, want %q", rc.ContextDir, "custom-context")
	}
	if rc.TokenBudget != 4000 {
		t.Errorf("TokenBudget = %d, want %d", rc.TokenBudget, 4000)
	}
	if len(rc.PriorityOrder) != 2 || rc.PriorityOrder[0] != "TASKS.md" {
		t.Errorf("PriorityOrder = %v, want [TASKS.md DECISIONS.md]", rc.PriorityOrder)
	}
	if rc.AutoArchive {
		t.Error("AutoArchive = true, want false")
	}
	if rc.ArchiveAfterDays != 14 {
		t.Errorf("ArchiveAfterDays = %d, want %d", rc.ArchiveAfterDays, 14)
	}
}

func TestGetRC_EnvOverrides(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origDir)

	// Create .contextrc file
	rcContent := `context_dir: file-context
token_budget: 4000
`
	os.WriteFile(filepath.Join(tempDir, ".contextrc"), []byte(rcContent), 0644)

	// Set environment variables (t.Setenv auto-restores after test)
	t.Setenv(config.EnvCtxDir, "env-context")
	t.Setenv(config.EnvCtxTokenBudget, "2000")

	Reset()

	rc := RC()

	// Env should override file
	if rc.ContextDir != "env-context" {
		t.Errorf("ContextDir = %q, want %q (env override)", rc.ContextDir, "env-context")
	}
	if rc.TokenBudget != 2000 {
		t.Errorf("TokenBudget = %d, want %d (env override)", rc.TokenBudget, 2000)
	}
}

func TestGetContextDir_CLIOverride(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origDir)

	// Create .contextrc file
	rcContent := `context_dir: file-context`
	os.WriteFile(filepath.Join(tempDir, ".contextrc"), []byte(rcContent), 0644)

	// Set env override (t.Setenv auto-restores after test)
	t.Setenv(config.EnvCtxDir, "env-context")

	Reset()

	// CLI override takes precedence over all
	OverrideContextDir("cli-context")
	defer Reset()

	dir := ContextDir()
	if dir != "cli-context" {
		t.Errorf("ContextDir() = %q, want %q (CLI override)", dir, "cli-context")
	}
}

func TestGetTokenBudget(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origDir)

	Reset()

	// Default value
	budget := TokenBudget()
	if budget != DefaultTokenBudget {
		t.Errorf("TokenBudget() = %d, want %d", budget, DefaultTokenBudget)
	}
}

func TestGetRC_InvalidYAML(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origDir)

	// Create invalid .contextrc file
	os.WriteFile(filepath.Join(tempDir, ".contextrc"), []byte("invalid: [yaml: content"), 0644)

	Reset()

	// Should return defaults on invalid YAML
	rc := RC()
	if rc.TokenBudget != DefaultTokenBudget {
		t.Errorf("TokenBudget = %d, want %d (defaults on invalid YAML)", rc.TokenBudget, DefaultTokenBudget)
	}
}

func TestGetRC_PartialConfig(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origDir)

	// Create .contextrc with only some fields
	rcContent := `token_budget: 5000`
	os.WriteFile(filepath.Join(tempDir, ".contextrc"), []byte(rcContent), 0644)

	Reset()

	rc := RC()

	// Specified value should be used
	if rc.TokenBudget != 5000 {
		t.Errorf("TokenBudget = %d, want %d", rc.TokenBudget, 5000)
	}
	// Unspecified values should use defaults
	if rc.ContextDir != config.DirContext {
		t.Errorf("ContextDir = %q, want %q (default)", rc.ContextDir, config.DirContext)
	}
}

func TestGetRC_InvalidEnvBudget(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origDir)

	t.Setenv(config.EnvCtxTokenBudget, "not-a-number")

	Reset()

	// Invalid env should be ignored, use default
	rc := RC()
	if rc.TokenBudget != DefaultTokenBudget {
		t.Errorf("TokenBudget = %d, want %d (default on invalid env)", rc.TokenBudget, DefaultTokenBudget)
	}
}

func TestGetRC_Singleton(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(origDir)

	Reset()

	rc1 := RC()
	rc2 := RC()

	if rc1 != rc2 {
		t.Error("RC() should return same instance")
	}
}
