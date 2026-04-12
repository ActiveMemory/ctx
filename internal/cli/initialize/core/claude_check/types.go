//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude_check

// PluginDetails describes a fully-installed ctx plugin for
// the Ready-state output. All string fields are
// human-readable and safe to print verbatim; none are
// required.
//
// Fields:
//   - Scope: installation scope recorded by Claude Code
//     ("user", "project", "local"). Empty if unknown.
//   - Version: plugin version string (e.g. "0.8.1").
//   - GitCommit: short git SHA the plugin was built from,
//     or empty when unavailable.
//   - Source: marketplace source type ("directory" for
//     local clone installs, "github" for remote
//     marketplace installs).
//   - SourcePath: filesystem path to the source clone when
//     Source is "directory". Empty for github-sourced
//     installs.
//   - EnabledGlobally: plugin is enabled in
//     ~/.claude/settings.json.
//   - EnabledLocally: plugin is enabled in the current
//     project's .claude/settings.local.json.
type PluginDetails struct {
	Scope           string
	Version         string
	GitCommit       string
	Source          string
	SourcePath      string
	EnabledGlobally bool
	EnabledLocally  bool
}

// installedPluginsFile mirrors the shape of
// ~/.claude/plugins/installed_plugins.json used only for
// dev-loop metadata extraction. The full schema has more
// fields than we care about; what's here is what we parse.
type installedPluginsFile struct {
	Version int                            `json:"version"`
	Plugins map[string][]installedPluginV2 `json:"plugins"`
}

// installedPluginV2 is a single entry under a plugin ID in
// installed_plugins.json (the format uses an array to allow
// multiple installations of the same plugin at different
// scopes).
type installedPluginV2 struct {
	Scope        string `json:"scope"`
	InstallPath  string `json:"installPath"`
	Version      string `json:"version"`
	GitCommitSha string `json:"gitCommitSha"`
}

// knownMarketplacesFile mirrors the shape of
// ~/.claude/plugins/known_marketplaces.json. Each top-level
// key is a marketplace ID; we only read the source block to
// discover the clone path for directory-sourced installs.
type knownMarketplacesFile map[string]marketplaceEntry

// marketplaceEntry is a single marketplace registration in
// known_marketplaces.json.
type marketplaceEntry struct {
	Source marketplaceSource `json:"source"`
}

// marketplaceSource captures the origin of a marketplace —
// either a github repo or a local directory (the dev
// flow).
type marketplaceSource struct {
	Type string `json:"source"`
	Path string `json:"path"`
	Repo string `json:"repo"`
}

// State represents the current setup state of Claude Code
// and the ctx plugin together.
type State int

const (
	// StateClaudeAbsent means the `claude` binary is not on
	// PATH. The user hasn't installed Claude Code yet, or
	// it is not in a directory the shell can find.
	StateClaudeAbsent State = iota
	// StatePluginNotInstalled means `claude` is present but
	// the ctx plugin is not registered in
	// ~/.claude/plugins/installed_plugins.json.
	StatePluginNotInstalled
	// StatePluginInstalledNotEnabled means the plugin is
	// registered but is not enabled in either the global
	// settings or the project-local settings.
	StatePluginInstalledNotEnabled
	// StatePluginReady means `claude` is present, the
	// plugin is registered, and the plugin is enabled
	// globally or locally (or both). The user is fully
	// wired up.
	StatePluginReady
)
