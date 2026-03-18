//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/err/config"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errInit "github.com/ActiveMemory/ctx/internal/err/initialize"
	errParser "github.com/ActiveMemory/ctx/internal/err/parser"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/write/add"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
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
		return errInit.HomeDir(homeErr)
	}
	claudeDir := filepath.Join(homeDir, ".claude")
	installedPath := filepath.Join(claudeDir, claude.InstalledPlugins)
	installedData, readErr := io.SafeReadUserFile(installedPath)
	if readErr != nil {
		initialize.PluginSkipped(cmd)
		return nil
	}
	var installed installedPlugins
	if parseErr := json.Unmarshal(installedData, &installed); parseErr != nil {
		return errParser.ParseFile(installedPath, parseErr)
	}
	if _, found := installed.Plugins[claude.PluginID]; !found {
		initialize.PluginSkipped(cmd)
		return nil
	}
	settingsPath := filepath.Join(claudeDir, claude.GlobalSettings)
	var settings globalSettings
	existingData, readErr := io.SafeReadUserFile(settingsPath)
	if readErr != nil && !os.IsNotExist(readErr) {
		return add.ErrFileRead(settingsPath, readErr)
	}
	if readErr == nil {
		if parseErr := json.Unmarshal(existingData, &settings); parseErr != nil {
			return errParser.ParseFile(settingsPath, parseErr)
		}
	} else {
		settings = make(globalSettings)
	}
	if raw, ok := settings["enabledPlugins"]; ok {
		var enabled map[string]bool
		if parseErr := json.Unmarshal(raw, &enabled); parseErr == nil {
			if enabled[claude.PluginID] {
				initialize.PluginAlreadyEnabled(cmd)
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
	enabled[claude.PluginID] = true
	enabledJSON, marshalErr := json.Marshal(enabled)
	if marshalErr != nil {
		return config.MarshalPlugins(marshalErr)
	}
	settings["enabledPlugins"] = enabledJSON
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if encodeErr := encoder.Encode(settings); encodeErr != nil {
		return config.MarshalSettings(encodeErr)
	}
	if writeErr := os.WriteFile(settingsPath, buf.Bytes(), fs.PermFile); writeErr != nil {
		return errFs.FileWrite(settingsPath, writeErr)
	}
	initialize.PluginEnabled(cmd, settingsPath)
	return nil
}

// PluginInstalled reports whether the ctx plugin is registered in
// ~/.claude/plugins/installed_plugins.json.
func PluginInstalled() bool {
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return false
	}
	installedPath := filepath.Join(homeDir, ".claude", claude.InstalledPlugins)
	data, readErr := io.SafeReadUserFile(installedPath)
	if readErr != nil {
		return false
	}
	var installed installedPlugins
	if parseErr := json.Unmarshal(data, &installed); parseErr != nil {
		return false
	}
	_, found := installed.Plugins[claude.PluginID]
	return found
}

// PluginEnabledGlobally reports whether the ctx plugin is enabled in
// ~/.claude/settings.json.
func PluginEnabledGlobally() bool {
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return false
	}
	settingsPath := filepath.Join(homeDir, ".claude", claude.GlobalSettings)
	data, readErr := io.SafeReadUserFile(settingsPath)
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
	return enabled[claude.PluginID]
}

// PluginEnabledLocally reports whether the ctx plugin is enabled in
// .claude/settings.local.json in the current project.
func PluginEnabledLocally() bool {
	data, readErr := os.ReadFile(claude.Settings)
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
	return enabled[claude.PluginID]
}
