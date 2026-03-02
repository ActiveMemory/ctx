//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
)

// withFakeHome sets HOME to a temp dir for the test and restores it after.
func withFakeHome(t *testing.T) string {
	t.Helper()
	tmpHome := t.TempDir()
	orig := os.Getenv("HOME")
	t.Setenv("HOME", tmpHome)
	t.Cleanup(func() { _ = os.Setenv("HOME", orig) })
	return tmpHome
}

// writeJSON writes a JSON file at the given path, creating parent dirs.
func writeJSON(t *testing.T, path string, v any) {
	t.Helper()
	if mkErr := os.MkdirAll(filepath.Dir(path), 0o750); mkErr != nil {
		t.Fatalf("mkdir: %v", mkErr)
	}
	data, marshalErr := json.MarshalIndent(v, "", "  ")
	if marshalErr != nil {
		t.Fatalf("marshal: %v", marshalErr)
	}
	if writeErr := os.WriteFile(path, data, 0o600); writeErr != nil {
		t.Fatalf("write: %v", writeErr)
	}
}

// fakeInstalledPlugins writes a minimal installed_plugins.json with the
// ctx plugin present.
func fakeInstalledPlugins(t *testing.T, home string) {
	t.Helper()
	data := map[string]any{
		"version": 2,
		"plugins": map[string]any{
			config.PluginID: []map[string]string{
				{"scope": "user", "version": "0.7.2"},
			},
		},
	}
	writeJSON(t, filepath.Join(home, ".claude", config.FileInstalledPlugins), data)
}

// testCmd returns a cobra command that captures output.
func testCmd() (*cobra.Command, *bytes.Buffer) {
	var buf bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&buf)
	return cmd, &buf
}

func TestEnablePluginGlobally_SkipsWhenNotInstalled(t *testing.T) {
	_ = withFakeHome(t)
	cmd, buf := testCmd()

	if err := enablePluginGlobally(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("not installed")) {
		t.Errorf("expected 'not installed' message, got: %s", buf.String())
	}
}

func TestEnablePluginGlobally_EnablesWhenInstalled(t *testing.T) {
	home := withFakeHome(t)
	fakeInstalledPlugins(t, home)
	cmd, buf := testCmd()

	if err := enablePluginGlobally(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("enabled globally")) {
		t.Errorf("expected 'enabled globally', got: %s", buf.String())
	}

	// Verify settings.json was written correctly.
	data, readErr := os.ReadFile(
		filepath.Join(home, ".claude", config.FileGlobalSettings),
	)
	if readErr != nil {
		t.Fatalf("settings.json not created: %v", readErr)
	}

	var settings map[string]json.RawMessage
	if parseErr := json.Unmarshal(data, &settings); parseErr != nil {
		t.Fatalf("invalid JSON: %v", parseErr)
	}

	var enabled map[string]bool
	if parseErr := json.Unmarshal(settings["enabledPlugins"], &enabled); parseErr != nil {
		t.Fatalf("enabledPlugins parse error: %v", parseErr)
	}

	if !enabled[config.PluginID] {
		t.Error("plugin not enabled in settings.json")
	}
}

func TestEnablePluginGlobally_PreservesExistingSettings(t *testing.T) {
	home := withFakeHome(t)
	fakeInstalledPlugins(t, home)

	// Write existing settings with another plugin and custom setting.
	existing := map[string]any{
		"cleanupPeriodDays": 365,
		"enabledPlugins": map[string]bool{
			"other@plugin": true,
		},
	}
	writeJSON(t, filepath.Join(home, ".claude", config.FileGlobalSettings), existing)

	cmd, _ := testCmd()

	if err := enablePluginGlobally(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, readErr := os.ReadFile(
		filepath.Join(home, ".claude", config.FileGlobalSettings),
	)
	if readErr != nil {
		t.Fatalf("settings.json not found: %v", readErr)
	}

	var settings map[string]json.RawMessage
	if parseErr := json.Unmarshal(data, &settings); parseErr != nil {
		t.Fatalf("invalid JSON: %v", parseErr)
	}

	// Check cleanupPeriodDays preserved.
	var days int
	if parseErr := json.Unmarshal(settings["cleanupPeriodDays"], &days); parseErr != nil {
		t.Fatalf("cleanupPeriodDays missing: %v", parseErr)
	}
	if days != 365 {
		t.Errorf("cleanupPeriodDays = %d, want 365", days)
	}

	// Check both plugins present.
	var enabled map[string]bool
	if parseErr := json.Unmarshal(settings["enabledPlugins"], &enabled); parseErr != nil {
		t.Fatalf("enabledPlugins parse error: %v", parseErr)
	}
	if !enabled["other@plugin"] {
		t.Error("existing plugin was removed")
	}
	if !enabled[config.PluginID] {
		t.Error("ctx plugin not added")
	}
}

func TestEnablePluginGlobally_SkipsWhenAlreadyEnabled(t *testing.T) {
	home := withFakeHome(t)
	fakeInstalledPlugins(t, home)

	existing := map[string]any{
		"enabledPlugins": map[string]bool{
			config.PluginID: true,
		},
	}
	writeJSON(t, filepath.Join(home, ".claude", config.FileGlobalSettings), existing)

	cmd, buf := testCmd()

	if err := enablePluginGlobally(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("already enabled")) {
		t.Errorf("expected 'already enabled', got: %s", buf.String())
	}
}

func TestPluginInstalled_ReturnsFalseWhenMissing(t *testing.T) {
	_ = withFakeHome(t)
	if PluginInstalled() {
		t.Error("expected false when no installed_plugins.json")
	}
}

func TestPluginInstalled_ReturnsTrueWhenPresent(t *testing.T) {
	home := withFakeHome(t)
	fakeInstalledPlugins(t, home)
	if !PluginInstalled() {
		t.Error("expected true when plugin is in installed_plugins.json")
	}
}

func TestPluginEnabledGlobally_ReturnsFalseWhenMissing(t *testing.T) {
	_ = withFakeHome(t)
	if PluginEnabledGlobally() {
		t.Error("expected false when no settings.json")
	}
}

func TestPluginEnabledGlobally_ReturnsTrueWhenEnabled(t *testing.T) {
	home := withFakeHome(t)
	existing := map[string]any{
		"enabledPlugins": map[string]bool{
			config.PluginID: true,
		},
	}
	writeJSON(t, filepath.Join(home, ".claude", config.FileGlobalSettings), existing)

	if !PluginEnabledGlobally() {
		t.Error("expected true when plugin is enabled")
	}
}

func TestPluginEnabledLocally_ReturnsFalseWhenMissing(t *testing.T) {
	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	if chErr := os.Chdir(tmpDir); chErr != nil {
		t.Fatal(chErr)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	if PluginEnabledLocally() {
		t.Error("expected false when no settings.local.json")
	}
}

func TestPluginEnabledLocally_ReturnsTrueWhenEnabled(t *testing.T) {
	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	if chErr := os.Chdir(tmpDir); chErr != nil {
		t.Fatal(chErr)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	existing := map[string]any{
		"enabledPlugins": map[string]bool{
			config.PluginID: true,
		},
	}
	writeJSON(t, filepath.Join(tmpDir, config.FileSettings), existing)

	if !PluginEnabledLocally() {
		t.Error("expected true when plugin is enabled locally")
	}
}
