//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude

// Claude Code integration file names.
const (
	// Md is the Claude Code configuration file in the project root.
	Md = "CLAUDE.md"

	// Settings is the Claude Code local settings file.
	Settings = ".claude/settings.local.json"
	// SettingsGolden is the golden image of the Claude Code settings.
	SettingsGolden = ".claude/settings.golden.json"

	// GlobalSettings is the Claude Code global settings file.
	// Located at ~/.claude/settings.json (not the project-local one).
	GlobalSettings = "settings.json"
	// InstalledPlugins is the Claude Code installed plugins registry.
	// Located at ~/.claude/plugins/installed_plugins.json.
	InstalledPlugins = "plugins/installed_plugins.json"

	// PluginID is the ctx plugin identifier in Claude Code.
	PluginID = "ctx@activememory-ctx"
)
