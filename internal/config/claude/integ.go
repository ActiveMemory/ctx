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

	// KeyEnabledPlugins is the JSON key for enabled plugins
	// in Claude Code settings.
	KeyEnabledPlugins = "enabledPlugins"

	// PluginScope is the permission scope prefix for plugin-scoped skills.
	PluginScope = "ctx:"
	// PluginScopeWildcard is the wildcard suffix for plugin-scoped permissions.
	PluginScopeWildcard = ":*"

	// PermSkillPrefix is the opening token of a Claude Code skill permission.
	PermSkillPrefix = "Skill("
	// PermSkillSuffix is the closing token of a Claude Code skill permission.
	PermSkillSuffix = ")"
)
