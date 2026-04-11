//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claudecheck

import (
	cfgClaude "github.com/ActiveMemory/ctx/internal/config/claude"
)

// Placeholder strings shown in the Ready-state detail
// block when a given field is not available.
const (
	// fieldUnknown is the placeholder used when a string
	// field can't be determined from plugin metadata.
	fieldUnknown = "unknown"
	// fieldNone is the placeholder used when a field is
	// legitimately empty (e.g. clone path for a GitHub
	// install).
	fieldNone = "(none)"
	// sourceDirHot is the label for directory-sourced
	// installs that reflect edits immediately.
	sourceDirHot = "directory (hot-reload)"
	// sourceGitHub is the label for marketplace installs
	// pulled from the Anthropic-hosted GitHub repo.
	sourceGitHub = "github (marketplace)"
	// enabledBoth is the summary when the plugin is
	// enabled in both the global and project-local
	// settings files.
	enabledBoth = "user scope + project"
	// enabledGlobal is the summary when the plugin is
	// enabled only in ~/.claude/settings.json.
	enabledGlobal = "user scope (all projects)"
	// enabledLocal is the summary when the plugin is
	// enabled only in the project's
	// .claude/settings.local.json.
	enabledLocal = "project only"
	// enabledNeither is the summary when the plugin is
	// installed but not enabled anywhere. Ready state
	// excludes this case, but the constant is here so the
	// switch in [formatEnabled] is exhaustive.
	enabledNeither = "(none)"
	// versionOpen is the opening bracket placed before
	// the short SHA in the formatted version string.
	versionOpen = " ("
	// versionClose is the closing bracket placed after
	// the short SHA in the formatted version string.
	versionClose = ")"
)

// renderDetails converts a [PluginDetails] struct into the
// five human-readable strings expected by the Ready-state
// template placeholders. Each field falls back to a
// sentinel when the underlying data is empty so the output
// never shows a bare empty column.
//
// Parameters:
//   - d: populated plugin details
//
// Returns:
//   - scope: installation scope
//   - version: version + short git SHA
//   - source: marketplace source type label
//   - clonePath: filesystem clone path or sentinel
//   - enabled: enablement summary
func renderDetails(
	d PluginDetails,
) (scope, version, source, clonePath, enabled string) {
	scope = orDefault(d.Scope, fieldUnknown)
	version = formatVersion(d.Version, d.GitCommit)
	source = formatSource(d.Source)
	clonePath = orDefault(d.SourcePath, fieldNone)
	enabled = formatEnabled(d.EnabledGlobally, d.EnabledLocally)
	return
}

// orDefault returns s when non-empty, otherwise fallback.
//
// Parameters:
//   - s: candidate string
//   - fallback: default when s is empty
//
// Returns:
//   - string: s or fallback
func orDefault(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

// formatVersion renders "0.8.1 (b4cdb428)" when both
// fields are present, or just the version if the SHA is
// missing.
//
// Parameters:
//   - version: plugin version string
//   - sha: short git commit SHA
//
// Returns:
//   - string: formatted version line
func formatVersion(version, sha string) string {
	if version == "" {
		return fieldUnknown
	}
	if sha == "" {
		return version
	}
	return version + versionOpen + sha + versionClose
}

// formatSource maps the marketplace source type into a
// human-readable label.
//
// Parameters:
//   - source: raw source type from
//     known_marketplaces.json
//
// Returns:
//   - string: display label
func formatSource(source string) string {
	switch source {
	case cfgClaude.PluginSourceDirectory:
		return sourceDirHot
	case cfgClaude.PluginSourceGitHub:
		return sourceGitHub
	case "":
		return fieldUnknown
	default:
		return source
	}
}

// formatEnabled renders an enablement summary from the two
// boolean flags.
//
// Parameters:
//   - globally: plugin enabled in ~/.claude/settings.json
//   - locally: plugin enabled in this project's
//     .claude/settings.local.json
//
// Returns:
//   - string: enablement summary
func formatEnabled(globally, locally bool) string {
	switch {
	case globally && locally:
		return enabledBoth
	case globally:
		return enabledGlobal
	case locally:
		return enabledLocal
	default:
		return enabledNeither
	}
}
