//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package codex

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	cfgCodex "github.com/ActiveMemory/ctx/internal/config/codex"
	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// Detect returns the combined state of the Codex CLI and the ctx
// plugin. It never errors: every read failure is treated as
// "not yet installed".
//
// Returns:
//   - State: the current combined state
func Detect() State {
	if _, lookErr := exec.LookPath(cfgCodex.Binary); lookErr != nil {
		return StateAbsent
	}
	home := Home()
	if !PluginInstalled(home) {
		return StatePluginNotInstalled
	}
	if !PluginEnabled(home) {
		return StatePluginInstalledNotEnabled
	}
	return StatePluginReady
}

// Unwired reports whether Codex is installed but ctx is not yet
// integrated with it: the `codex` binary is on PATH, the project
// has no `.codex/hooks.json`, and the ctx plugin is not enabled in
// the user's config.toml. `ctx init` prints its Codex hint when
// this is true.
//
// Returns:
//   - bool: true when `ctx setup codex --write` is still needed
func Unwired() bool {
	if _, lookErr := exec.LookPath(cfgCodex.Binary); lookErr != nil {
		return false
	}
	if _, statErr := os.Stat(cfgSetup.HooksPathCodex); statErr == nil {
		return false
	}
	home := Home()
	return !PluginEnabled(home) || !PluginNativeVariant(home)
}

// PluginNativeVariant reports whether the installed ctx plugin is
// the Codex variant: at least one cached version root contains
// `.codex-plugin/`. A GitHub install from a revision that predates
// the Codex marketplace silently delivers the legacy Claude Code
// variant (`.claude-plugin/`), whose hooks cannot run under Codex;
// the deployer must not treat that as a working plugin.
//
// Parameters:
//   - home: Codex home directory (see [Home])
//
// Returns:
//   - bool: true when a cached copy carries `.codex-plugin/`
func PluginNativeVariant(home string) bool {
	if home == "" {
		return false
	}
	pluginDir := filepath.Join(
		home, cfgCodex.DirPlugins, cfgCodex.DirPluginCache,
		cfgCodex.MarketplaceID, cfgCodex.PluginName,
	)
	entries, readErr := os.ReadDir(pluginDir)
	if readErr != nil {
		return false
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		manifest := filepath.Join(
			pluginDir, entry.Name(), cfgCodex.DirPluginManifest,
		)
		if _, statErr := os.Stat(manifest); statErr == nil {
			return true
		}
	}
	return false
}

// ProjectConfigured reports whether the project-local Codex
// integration is already deployed in the current working
// directory (`.codex/hooks.json` exists). `ctx setup codex`
// prefers this state label over the plugin detection states.
//
// Returns:
//   - bool: true when `.codex/hooks.json` exists at $PWD
func ProjectConfigured() bool {
	_, statErr := os.Stat(cfgSetup.HooksPathCodex)
	return statErr == nil
}

// PluginInstalled reports whether the ctx plugin is present in
// the Codex plugin cache
// (<home>/plugins/cache/<marketplace>/<plugin>/<version>/).
//
// Parameters:
//   - home: Codex home directory (see [Home])
//
// Returns:
//   - bool: true when the plugin directory exists and holds at
//     least one version subdirectory
func PluginInstalled(home string) bool {
	if home == "" {
		return false
	}
	pluginDir := filepath.Join(
		home, cfgCodex.DirPlugins, cfgCodex.DirPluginCache,
		cfgCodex.MarketplaceID, cfgCodex.PluginName,
	)
	entries, readErr := os.ReadDir(pluginDir)
	if readErr != nil {
		return false
	}
	for _, entry := range entries {
		if entry.IsDir() {
			return true
		}
	}
	return false
}

// PluginEnabled reports whether <home>/config.toml enables the
// ctx plugin. The file is scanned line by line, never parsed:
// the `[plugins."ctx@activememory-ctx"]` header must be present
// and, until the next table header, no line may set
// `enabled = false`. A header with no `enabled` key counts as
// enabled (Codex's default).
//
// Parameters:
//   - home: Codex home directory (see [Home])
//
// Returns:
//   - bool: true when the plugin table is present and not
//     explicitly disabled
func PluginEnabled(home string) bool {
	if home == "" {
		return false
	}
	data, readErr := ctxIo.SafeReadUserFile(
		filepath.Join(home, cfgCodex.FileConfigTOML),
	)
	if readErr != nil {
		return false
	}

	inTable := false
	found := false
	for _, raw := range splitLines(data) {
		line := strings.TrimSpace(raw)
		if headerLine(line, cfgCodex.TOMLHeaderPluginCtx) {
			inTable = true
			found = true
			continue
		}
		if !inTable {
			continue
		}
		if strings.HasPrefix(line, cfgCodex.TOMLBracketOpen) {
			inTable = false
			continue
		}
		key, value, hasAssign := strings.Cut(line, token.KeyValueSep)
		if !hasAssign || strings.TrimSpace(key) != cfgCodex.TOMLKeyEnabled {
			continue
		}
		value, _, _ = strings.Cut(value, token.Hash)
		if strings.TrimSpace(value) != cfgCodex.TOMLTrue {
			return false
		}
	}
	return found
}
