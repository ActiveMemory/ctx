//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
)

type installedPlugins struct {
	Plugins map[string]json.RawMessage `json:"plugins"`
}

type globalSettings map[string]json.RawMessage

// EnablePluginGlobally enables the ctx plugin in ~/.claude/settings.json.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil if file operations fail
func EnablePluginGlobally(cmd *cobra.Command) error {
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return fmt.Errorf("cannot determine home directory: %w", homeErr)
	}
	claudeDir := filepath.Join(homeDir, ".claude")
	installedPath := filepath.Join(claudeDir, config.FileInstalledPlugins)
	installedData, readErr := os.ReadFile(installedPath) //nolint:gosec // G304: path from os.UserHomeDir
	if readErr != nil {
		cmd.Println("  ○ Plugin enablement skipped (plugin not installed)")
		return nil
	}
	var installed installedPlugins
	if parseErr := json.Unmarshal(installedData, &installed); parseErr != nil {
		return fmt.Errorf("failed to parse %s: %w", installedPath, parseErr)
	}
	if _, found := installed.Plugins[config.PluginID]; !found {
		cmd.Println("  ○ Plugin enablement skipped (plugin not installed)")
		return nil
	}
	settingsPath := filepath.Join(claudeDir, config.FileGlobalSettings)
	var settings globalSettings
	existingData, readErr := os.ReadFile(settingsPath) //nolint:gosec // G304: path from os.UserHomeDir
	if readErr != nil && !os.IsNotExist(readErr) {
		return fmt.Errorf("failed to read %s: %w", settingsPath, readErr)
	}
	if readErr == nil {
		if parseErr := json.Unmarshal(existingData, &settings); parseErr != nil {
			return fmt.Errorf("failed to parse %s: %w", settingsPath, parseErr)
		}
	} else {
		settings = make(globalSettings)
	}
	if raw, ok := settings["enabledPlugins"]; ok {
		var enabled map[string]bool
		if parseErr := json.Unmarshal(raw, &enabled); parseErr == nil {
			if enabled[config.PluginID] {
				cmd.Println("  ○ Plugin already enabled globally")
				return nil
			}
		}
	}
	var enabled map[string]bool
	if raw, ok := settings["enabledPlugins"]; ok {
		if parseErr := json.Unmarshal(raw, &enabled); parseErr != nil {
			enabled = make(map[string]bool)
		}
	} else {
		enabled = make(map[string]bool)
	}
	enabled[config.PluginID] = true
	enabledJSON, marshalErr := json.Marshal(enabled)
	if marshalErr != nil {
		return fmt.Errorf("failed to marshal enabledPlugins: %w", marshalErr)
	}
	settings["enabledPlugins"] = enabledJSON
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if encodeErr := encoder.Encode(settings); encodeErr != nil {
		return fmt.Errorf("failed to marshal settings: %w", encodeErr)
	}
	if writeErr := os.WriteFile(settingsPath, buf.Bytes(), config.PermFile); writeErr != nil {
		return fmt.Errorf("failed to write %s: %w", settingsPath, writeErr)
	}
	cmd.Println(fmt.Sprintf("  ✓ Plugin enabled globally in %s", settingsPath))
	return nil
}

// PluginInstalled reports whether the ctx plugin is registered in
// ~/.claude/plugins/installed_plugins.json.
func PluginInstalled() bool {
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return false
	}
	installedPath := filepath.Join(homeDir, ".claude", config.FileInstalledPlugins)
	data, readErr := os.ReadFile(installedPath) //nolint:gosec // G304: path from os.UserHomeDir
	if readErr != nil {
		return false
	}
	var installed installedPlugins
	if parseErr := json.Unmarshal(data, &installed); parseErr != nil {
		return false
	}
	_, found := installed.Plugins[config.PluginID]
	return found
}

// PluginEnabledGlobally reports whether the ctx plugin is enabled in
// ~/.claude/settings.json.
func PluginEnabledGlobally() bool {
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return false
	}
	settingsPath := filepath.Join(homeDir, ".claude", config.FileGlobalSettings)
	data, readErr := os.ReadFile(settingsPath) //nolint:gosec // G304: path from os.UserHomeDir
	if readErr != nil {
		return false
	}
	var settings globalSettings
	if parseErr := json.Unmarshal(data, &settings); parseErr != nil {
		return false
	}
	raw, ok := settings["enabledPlugins"]
	if !ok {
		return false
	}
	var enabled map[string]bool
	if parseErr := json.Unmarshal(raw, &enabled); parseErr != nil {
		return false
	}
	return enabled[config.PluginID]
}

// PluginEnabledLocally reports whether the ctx plugin is enabled in
// .claude/settings.local.json in the current project.
func PluginEnabledLocally() bool {
	data, readErr := os.ReadFile(config.FileSettings)
	if readErr != nil {
		return false
	}
	var raw map[string]json.RawMessage
	if parseErr := json.Unmarshal(data, &raw); parseErr != nil {
		return false
	}
	epRaw, ok := raw["enabledPlugins"]
	if !ok {
		return false
	}
	var enabled map[string]bool
	if parseErr := json.Unmarshal(epRaw, &enabled); parseErr != nil {
		return false
	}
	return enabled[config.PluginID]
}
