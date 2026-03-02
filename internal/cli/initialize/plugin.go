//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
)

// installedPlugins represents the structure of installed_plugins.json.
type installedPlugins struct {
	Plugins map[string]json.RawMessage `json:"plugins"`
}

// globalSettings represents the structure of ~/.claude/settings.json.
//
// Only the fields ctx cares about are modeled; unknown fields are
// preserved via the raw map approach.
type globalSettings map[string]json.RawMessage

// enablePluginGlobally adds the ctx plugin to enabledPlugins in
// ~/.claude/settings.json if the plugin is installed but not yet enabled.
//
// Checks ~/.claude/plugins/installed_plugins.json first to confirm the
// plugin is actually installed. Merges into the existing file, preserving
// all other settings.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if JSON parsing or file operations fail
func enablePluginGlobally(cmd *cobra.Command) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return fmt.Errorf("cannot determine home directory: %w", homeErr)
	}

	claudeDir := filepath.Join(homeDir, ".claude")

	// Check if plugin is installed.
	installedPath := filepath.Join(claudeDir, config.FileInstalledPlugins)
	installedData, readErr := os.ReadFile(installedPath) //nolint:gosec // G304: path from os.UserHomeDir
	if readErr != nil {
		cmd.Printf(
			"  %s Plugin enablement skipped (plugin not installed)\n",
			yellow("○"),
		)
		return nil
	}

	var installed installedPlugins
	if parseErr := json.Unmarshal(installedData, &installed); parseErr != nil {
		return fmt.Errorf("failed to parse %s: %w", installedPath, parseErr)
	}

	if _, found := installed.Plugins[config.PluginID]; !found {
		cmd.Printf(
			"  %s Plugin enablement skipped (plugin not installed)\n",
			yellow("○"),
		)
		return nil
	}

	// Read existing global settings.
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

	// Check if already enabled.
	if raw, ok := settings["enabledPlugins"]; ok {
		var enabled map[string]bool
		if parseErr := json.Unmarshal(raw, &enabled); parseErr == nil {
			if enabled[config.PluginID] {
				cmd.Printf(
					"  %s Plugin already enabled globally\n",
					yellow("○"),
				)
				return nil
			}
		}
	}

	// Merge the plugin into enabledPlugins.
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

	// Write back.
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

	cmd.Printf("  %s Plugin enabled globally in %s\n", green("✓"), settingsPath)
	return nil
}

// PluginInstalled reports whether the ctx plugin is registered in
// ~/.claude/plugins/installed_plugins.json.
func PluginInstalled() bool {
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return false
	}

	installedPath := filepath.Join(
		homeDir, ".claude", config.FileInstalledPlugins,
	)
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

	settingsPath := filepath.Join(
		homeDir, ".claude", config.FileGlobalSettings,
	)
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

	// settings.local.json uses a different shape (claude.Settings) but
	// enabledPlugins sits at the top level in both files.
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
