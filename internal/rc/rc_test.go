//   /    ctx:                         https://ctx.ist
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
	// Change to temp directory with no .ctxrc
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

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
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	// Create .ctxrc file
	rcContent := `context_dir: custom-context
token_budget: 4000
priority_order:
  - TASKS.md
  - DECISIONS.md
auto_archive: false
archive_after_days: 14
`
	_ = os.WriteFile(filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0600)

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
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	// Create .ctxrc file
	rcContent := `context_dir: file-context
token_budget: 4000
`
	_ = os.WriteFile(filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0600)

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
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	// Create .ctxrc file
	rcContent := `context_dir: file-context`
	_ = os.WriteFile(filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0600)

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
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

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
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	// Create invalid .ctxrc file
	_ = os.WriteFile(filepath.Join(tempDir, ".ctxrc"), []byte("invalid: [yaml: content"), 0600)

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
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	// Create .ctxrc with only some fields
	rcContent := `token_budget: 5000`
	_ = os.WriteFile(filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0600)

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
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

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
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	Reset()

	rc1 := RC()
	rc2 := RC()

	if rc1 != rc2 {
		t.Error("RC() should return same instance")
	}
}

func TestPriorityOrder(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	Reset()

	// Default has nil PriorityOrder
	order := PriorityOrder()
	if order != nil {
		t.Errorf("PriorityOrder() = %v, want nil", order)
	}
}

func TestPriorityOrder_Custom(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	rcContent := `priority_order:
  - TASKS.md
  - DECISIONS.md
  - LEARNINGS.md
`
	_ = os.WriteFile(filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0600)

	Reset()

	order := PriorityOrder()
	if len(order) != 3 {
		t.Fatalf("PriorityOrder() len = %d, want 3", len(order))
	}
	if order[0] != "TASKS.md" {
		t.Errorf("PriorityOrder()[0] = %q, want %q", order[0], "TASKS.md")
	}
}

func TestAutoArchive(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	Reset()

	// Default is true
	if !AutoArchive() {
		t.Error("AutoArchive() = false, want true")
	}
}

func TestAutoArchive_Disabled(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	rcContent := `auto_archive: false`
	_ = os.WriteFile(filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0600)

	Reset()

	if AutoArchive() {
		t.Error("AutoArchive() = true, want false")
	}
}

func TestArchiveAfterDays(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	Reset()

	days := ArchiveAfterDays()
	if days != DefaultArchiveAfterDays {
		t.Errorf("ArchiveAfterDays() = %d, want %d", days, DefaultArchiveAfterDays)
	}
}

func TestArchiveAfterDays_Custom(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	rcContent := `archive_after_days: 30`
	_ = os.WriteFile(filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0600)

	Reset()

	days := ArchiveAfterDays()
	if days != 30 {
		t.Errorf("ArchiveAfterDays() = %d, want %d", days, 30)
	}
}

func TestScratchpadEncrypt_Default(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	Reset()

	// Default (nil pointer) should return true
	if !ScratchpadEncrypt() {
		t.Error("ScratchpadEncrypt() = false, want true (default)")
	}
}

func TestScratchpadEncrypt_Explicit(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	rcContent := `scratchpad_encrypt: false`
	_ = os.WriteFile(filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0600)

	Reset()

	if ScratchpadEncrypt() {
		t.Error("ScratchpadEncrypt() = true, want false")
	}
}

func TestScratchpadEncrypt_ExplicitTrue(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	rcContent := `scratchpad_encrypt: true`
	_ = os.WriteFile(filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0600)

	Reset()

	if !ScratchpadEncrypt() {
		t.Error("ScratchpadEncrypt() = false, want true")
	}
}

func TestFilePriority_DefaultOrder(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	Reset()

	// CONSTITUTION.md should be first in default FileReadOrder
	p := FilePriority(config.FileConstitution)
	if p != 1 {
		t.Errorf("FilePriority(%q) = %d, want 1", config.FileConstitution, p)
	}

	// TASKS.md should be second
	p = FilePriority(config.FileTask)
	if p != 2 {
		t.Errorf("FilePriority(%q) = %d, want 2", config.FileTask, p)
	}

	// Unknown file gets 100
	p = FilePriority("UNKNOWN.md")
	if p != 100 {
		t.Errorf("FilePriority(%q) = %d, want 100", "UNKNOWN.md", p)
	}
}

func TestFilePriority_CustomOrder(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	rcContent := `priority_order:
  - DECISIONS.md
  - TASKS.md
`
	_ = os.WriteFile(filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0600)

	Reset()

	// DECISIONS.md should be first in custom order
	p := FilePriority(config.FileDecision)
	if p != 1 {
		t.Errorf("FilePriority(%q) = %d, want 1", config.FileDecision, p)
	}

	// TASKS.md should be second
	p = FilePriority(config.FileTask)
	if p != 2 {
		t.Errorf("FilePriority(%q) = %d, want 2", config.FileTask, p)
	}

	// File not in custom order gets 100
	p = FilePriority("UNKNOWN.md")
	if p != 100 {
		t.Errorf("FilePriority(%q) = %d, want 100", "UNKNOWN.md", p)
	}
}

func TestContextDir_NoOverride(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	Reset()

	dir := ContextDir()
	if dir != config.DirContext {
		t.Errorf("ContextDir() = %q, want %q", dir, config.DirContext)
	}
}

func TestAllowOutsideCwd_Default(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	Reset()

	// Default is false
	if AllowOutsideCwd() {
		t.Error("AllowOutsideCwd() = true, want false (default)")
	}
}

func TestAllowOutsideCwd_Enabled(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	rcContent := `allow_outside_cwd: true`
	_ = os.WriteFile(filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0600)

	Reset()

	if !AllowOutsideCwd() {
		t.Error("AllowOutsideCwd() = false, want true")
	}
}

func TestNotifyEvents_Default(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	Reset()

	// Default (nil Notify) returns nil
	events := NotifyEvents()
	if events != nil {
		t.Errorf("NotifyEvents() = %v, want nil", events)
	}
}

func TestNotifyEvents_Configured(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	rcContent := `notify:
  events:
    - loop
    - nudge
`
	_ = os.WriteFile(filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0600)

	Reset()

	events := NotifyEvents()
	if len(events) != 2 || events[0] != "loop" || events[1] != "nudge" {
		t.Errorf("NotifyEvents() = %v, want [loop nudge]", events)
	}
}

func TestKeyRotationDays_Default(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	Reset()

	days := KeyRotationDays()
	if days != DefaultKeyRotationDays {
		t.Errorf("KeyRotationDays() = %d, want %d", days, DefaultKeyRotationDays)
	}
}

func TestKeyRotationDays_Custom(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	rcContent := `key_rotation_days: 30
`
	_ = os.WriteFile(filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0600)

	Reset()

	days := KeyRotationDays()
	if days != 30 {
		t.Errorf("KeyRotationDays() = %d, want %d", days, 30)
	}
}

func TestKeyRotationDays_LegacyNotify(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	rcContent := `notify:
  key_rotation_days: 45
`
	_ = os.WriteFile(filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0600)

	Reset()

	days := KeyRotationDays()
	if days != 45 {
		t.Errorf("KeyRotationDays() = %d, want %d (legacy notify fallback)", days, 45)
	}
}

func TestKeyRotationDays_TopLevelTakesPrecedence(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	rcContent := `key_rotation_days: 60
notify:
  key_rotation_days: 45
`
	_ = os.WriteFile(filepath.Join(tempDir, ".ctxrc"), []byte(rcContent), 0600)

	Reset()

	days := KeyRotationDays()
	if days != 60 {
		t.Errorf("KeyRotationDays() = %d, want %d (top-level takes precedence)", days, 60)
	}
}

func TestGetRC_NegativeEnvBudget(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tempDir)
	defer func() { _ = os.Chdir(origDir) }()

	t.Setenv(config.EnvCtxTokenBudget, "-100")

	Reset()

	// Negative budget should be ignored (budget > 0 check)
	rc := RC()
	if rc.TokenBudget != DefaultTokenBudget {
		t.Errorf("TokenBudget = %d, want %d (default on negative env)", rc.TokenBudget, DefaultTokenBudget)
	}
}
