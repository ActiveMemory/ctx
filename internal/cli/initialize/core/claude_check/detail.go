//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude_check

import (
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/cli/initialize/core/plugin"
	cfgClaude "github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/dir"
)

// Details loads rich metadata about the currently-installed
// ctx plugin by cross-referencing installed_plugins.json
// and known_marketplaces.json under ~/.claude/plugins/.
//
// Details returns a zero-value PluginDetails and ok=false
// if either file is missing, unreadable, or doesn't mention
// the ctx plugin. Callers should fall back to a minimal
// confirmation message in that case — a metadata parse
// failure must never break the `ctx init` tail.
//
// Returns:
//   - PluginDetails: populated plugin metadata
//   - bool: true iff detail extraction succeeded
func Details() (PluginDetails, bool) {
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return PluginDetails{}, false
	}

	installedPath := filepath.Join(
		homeDir, dir.Claude, cfgClaude.InstalledPlugins,
	)
	installed, installedOk := readInstalled(installedPath)
	if !installedOk {
		return PluginDetails{}, false
	}

	marketplacesPath := filepath.Join(
		homeDir, dir.Claude, cfgClaude.KnownMarketplaces,
	)
	// Marketplace metadata is optional: if it's missing we
	// still return the installed-plugin fields and let the
	// source/path render as empty.
	marketplaces, _ := readMarketplaces(marketplacesPath)

	entries, found := installed.Plugins[cfgClaude.PluginID]
	if !found || len(entries) == 0 {
		return PluginDetails{}, false
	}

	entry := entries[0]

	d := PluginDetails{
		Scope:           entry.Scope,
		Version:         entry.Version,
		GitCommit:       shortSha(entry.GitCommitSha),
		EnabledGlobally: plugin.EnabledGlobally(),
		EnabledLocally:  plugin.EnabledLocally(),
	}

	if mp, mpFound := marketplaces[cfgClaude.PluginMarketplaceID]; mpFound {
		d.Source = mp.Source.Type
		d.SourcePath = mp.Source.Path
	}

	return d, true
}
