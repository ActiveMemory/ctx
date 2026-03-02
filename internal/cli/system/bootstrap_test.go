//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func TestBootstrap_TextOutput(t *testing.T) {
	origDir, _ := os.Getwd()
	_ = os.Chdir(t.TempDir())
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	cmd := newTestCmd()
	if err := runBootstrap(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)

	if !strings.Contains(out, "context_dir:") {
		t.Errorf("expected output to contain 'context_dir:', got: %s", out)
	}
	if !strings.Contains(out, "ctx bootstrap") {
		t.Errorf("expected output to contain 'ctx bootstrap', got: %s", out)
	}
	if !strings.Contains(out, "CONSTITUTION.md") {
		t.Errorf("expected output to contain 'CONSTITUTION.md', got: %s", out)
	}
	if !strings.Contains(out, "TASKS.md") {
		t.Errorf("expected output to contain 'TASKS.md', got: %s", out)
	}
	if !strings.Contains(out, "DECISIONS.md") {
		t.Errorf("expected output to contain 'DECISIONS.md', got: %s", out)
	}
	if !strings.Contains(out, "Rules:") {
		t.Errorf("expected output to contain 'Rules:', got: %s", out)
	}
	if !strings.Contains(out, "Next steps:") {
		t.Errorf("expected output to contain 'Next steps:', got: %s", out)
	}
	if !strings.Contains(out, "AGENT_PLAYBOOK") {
		t.Errorf("expected output to contain 'AGENT_PLAYBOOK', got: %s", out)
	}
}

func TestBootstrap_JSONOutput(t *testing.T) {
	origDir, _ := os.Getwd()
	_ = os.Chdir(t.TempDir())
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	cmd := newTestCmd()
	cmd.Flags().Bool("json", true, "")
	_ = cmd.Flags().Set("json", "true")

	if err := runBootstrap(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)

	var result struct {
		ContextDir string   `json:"context_dir"`
		Files      []string `json:"files"`
		Rules      []string `json:"rules"`
		NextSteps  []string `json:"next_steps"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON output: %v\noutput: %s", err, out)
	}

	if result.ContextDir == "" {
		t.Error("expected context_dir to be non-empty")
	}
	if len(result.Files) == 0 {
		t.Error("expected files to be non-empty")
	}
	if len(result.Rules) == 0 {
		t.Error("expected rules to be non-empty")
	}
	if len(result.NextSteps) == 0 {
		t.Error("expected next_steps to be non-empty")
	}
	foundPlaybook := false
	for _, s := range result.NextSteps {
		if strings.Contains(s, "AGENT_PLAYBOOK") {
			foundPlaybook = true
			break
		}
	}
	if !foundPlaybook {
		t.Errorf("expected next_steps to mention AGENT_PLAYBOOK, got: %v", result.NextSteps)
	}

	// Verify known files are present
	fileSet := make(map[string]bool)
	for _, f := range result.Files {
		fileSet[f] = true
	}
	for _, expected := range []string{"CONSTITUTION.md", "TASKS.md", "DECISIONS.md"} {
		if !fileSet[expected] {
			t.Errorf("expected files to contain %q, got: %v", expected, result.Files)
		}
	}
}

func TestBootstrap_QuietOutput(t *testing.T) {
	origDir, _ := os.Getwd()
	_ = os.Chdir(t.TempDir())
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	cmd := newTestCmd()
	cmd.Flags().BoolP("quiet", "q", false, "")
	_ = cmd.Flags().Set("quiet", "true")

	if err := runBootstrap(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	trimmed := strings.TrimSpace(out)

	// Quiet output should be exactly the context dir path, nothing else.
	expectedDir := rc.ContextDir()
	if trimmed != expectedDir {
		t.Errorf("expected quiet output to be %q, got %q", expectedDir, trimmed)
	}

	// Should NOT contain verbose markers.
	if strings.Contains(out, "Rules:") {
		t.Error("quiet output should not contain 'Rules:'")
	}
	if strings.Contains(out, "ctx bootstrap") {
		t.Error("quiet output should not contain 'ctx bootstrap'")
	}
}

func TestBootstrap_CustomDir(t *testing.T) {
	origDir, _ := os.Getwd()
	workDir := t.TempDir()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	customDir := "my-custom-context"
	rc.OverrideContextDir(customDir)
	defer rc.Reset()

	// Create the custom context dir with a file
	if err := os.MkdirAll(customDir, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(customDir+"/CONSTITUTION.md", []byte("# test\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := runBootstrap(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, customDir) {
		t.Errorf("expected output to contain custom dir %q, got: %s", customDir, out)
	}
}

func TestBootstrap_MissingDir(t *testing.T) {
	origDir, _ := os.Getwd()
	_ = os.Chdir(t.TempDir())
	defer func() { _ = os.Chdir(origDir) }()

	rc.OverrideContextDir("nonexistent-dir")
	defer rc.Reset()

	cmd := newTestCmd()
	err := runBootstrap(cmd)
	if err == nil {
		t.Fatal("expected error for missing directory, got nil")
	}
	if !strings.Contains(err.Error(), "context directory not found") {
		t.Errorf("expected error to contain 'context directory not found', got: %v", err)
	}
	if !strings.Contains(err.Error(), "ctx init") {
		t.Errorf("expected error to mention 'ctx init', got: %v", err)
	}
}

// writeTestJSON writes a JSON file at the given path, creating parent dirs.
func writeTestJSON(t *testing.T, path string, v any) {
	t.Helper()
	if mkErr := os.MkdirAll(filepath.Dir(path), 0o750); mkErr != nil {
		t.Fatalf("mkdir: %v", mkErr)
	}
	data, marshalErr := json.Marshal(v)
	if marshalErr != nil {
		t.Fatalf("marshal: %v", marshalErr)
	}
	if writeErr := os.WriteFile(path, data, 0o600); writeErr != nil {
		t.Fatalf("write: %v", writeErr)
	}
}

func TestPluginWarning_EmptyWhenNotInstalled(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	if w := pluginWarning(); w != "" {
		t.Errorf("expected empty warning, got: %s", w)
	}
}

func TestPluginWarning_EmptyWhenEnabledGlobally(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// Create installed_plugins.json.
	writeTestJSON(t,
		filepath.Join(tmpHome, ".claude", config.FileInstalledPlugins),
		map[string]any{
			"version": 2,
			"plugins": map[string]any{
				config.PluginID: []map[string]string{
					{"scope": "user"},
				},
			},
		},
	)

	// Create settings.json with plugin enabled.
	writeTestJSON(t,
		filepath.Join(tmpHome, ".claude", config.FileGlobalSettings),
		map[string]any{
			"enabledPlugins": map[string]bool{
				config.PluginID: true,
			},
		},
	)

	if w := pluginWarning(); w != "" {
		t.Errorf("expected empty warning, got: %s", w)
	}
}

func TestPluginWarning_WarnsWhenInstalledButNotEnabled(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// Create installed_plugins.json.
	writeTestJSON(t,
		filepath.Join(tmpHome, ".claude", config.FileInstalledPlugins),
		map[string]any{
			"version": 2,
			"plugins": map[string]any{
				config.PluginID: []map[string]string{
					{"scope": "user"},
				},
			},
		},
	)

	// No settings.json — plugin not enabled.

	w := pluginWarning()
	if w == "" {
		t.Fatal("expected warning, got empty string")
	}
	if !strings.Contains(w, "not enabled") {
		t.Errorf("warning should mention 'not enabled', got: %s", w)
	}
	if !strings.Contains(w, "ctx init") {
		t.Errorf("warning should mention 'ctx init', got: %s", w)
	}
}
