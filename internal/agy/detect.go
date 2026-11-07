//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agy

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	cfgAgy "github.com/ActiveMemory/ctx/internal/config/agy"
	cfgAsset "github.com/ActiveMemory/ctx/internal/config/asset"
	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// Detect returns the combined state of the Antigravity CLI and the ctx
// Antigravity plugin.
//
// Returns:
//   - State: current install and enablement state
func Detect() State {
	if _, lookErr := exec.LookPath(cfgAgy.Binary); lookErr != nil {
		return StateAbsent
	}
	configDir := ConfigDir()
	if !PluginInstalled(configDir) {
		return StatePluginNotInstalled
	}
	if !PluginEnabled(configDir) {
		return StatePluginInstalledNotEnabled
	}
	return StatePluginReady
}

// Unwired reports whether Antigravity is installed on PATH but neither
// project-local hooks nor the ctx plugin are configured.
//
// Returns:
//   - bool: true when `ctx setup agy --write` is still needed
func Unwired() bool {
	if _, lookErr := exec.LookPath(cfgAgy.Binary); lookErr != nil {
		return false
	}
	if ProjectConfigured() {
		return false
	}
	return !PluginEnabled(ConfigDir())
}

// ProjectConfigured reports whether the project-local Antigravity
// integration is deployed: `.agents/hooks.json` exists and contains
// a top-level `"ctx"` key.
//
// Returns:
//   - bool: true when `.agents/hooks.json` has a `"ctx"` key
func ProjectConfigured() bool {
	data, readErr := ctxIo.SafeReadUserFile(cfgSetup.HooksPathAgy)
	if readErr != nil {
		return false
	}
	var manifest map[string]json.RawMessage
	if parseErr := json.Unmarshal(data, &manifest); parseErr != nil {
		return false
	}
	_, hasCtx := manifest[cfgAgy.KeyCtx]
	return hasCtx
}

// PluginInstalled reports whether the ctx plugin manifest exists
// under `<configDir>/plugins/ctx/plugin.json`.
//
// Parameters:
//   - configDir: Antigravity configuration directory
//
// Returns:
//   - bool: true when plugin.json exists
func PluginInstalled(configDir string) bool {
	if configDir == "" {
		return false
	}
	manifest := filepath.Join(
		configDir, cfgAgy.DirPlugins, cfgAgy.PluginName,
		cfgAsset.FilePluginJSON,
	)
	_, statErr := os.Stat(manifest)
	return statErr == nil
}

// PluginEnabled reports whether the ctx plugin is installed and enabled
// (i.e. plugin.json does not set `"disabled": true`).
//
// Parameters:
//   - configDir: Antigravity configuration directory
//
// Returns:
//   - bool: true when plugin.json is present and not disabled
func PluginEnabled(configDir string) bool {
	if !PluginInstalled(configDir) {
		return false
	}
	manifest := filepath.Join(
		configDir, cfgAgy.DirPlugins, cfgAgy.PluginName,
		cfgAsset.FilePluginJSON,
	)
	data, readErr := ctxIo.SafeReadUserFile(manifest)
	if readErr != nil {
		return false
	}
	var fields map[string]any
	if parseErr := json.Unmarshal(data, &fields); parseErr != nil {
		return false
	}
	if disabledVal, hasDisabled := fields[cfgAgy.KeyDisabled]; hasDisabled {
		if d, ok := disabledVal.(bool); ok && d {
			return false
		}
	}
	return true
}
