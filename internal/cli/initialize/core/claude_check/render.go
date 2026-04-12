//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude_check

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	cfgClaude "github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
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
	scope = orDefault(
		d.Scope,
		desc.Text(text.DescKeyWriteClaudecheckUnknown),
	)
	version = formatVersion(d.Version, d.GitCommit)
	source = formatSource(d.Source)
	clonePath = orDefault(
		d.SourcePath,
		desc.Text(text.DescKeyWriteClaudecheckNone),
	)
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
		return desc.Text(text.DescKeyWriteClaudecheckUnknown)
	}
	if sha == "" {
		return version
	}
	return version +
		desc.Text(text.DescKeyWriteClaudecheckVersionOpen) +
		sha +
		desc.Text(text.DescKeyWriteClaudecheckVersionClose)
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
		return desc.Text(text.DescKeyWriteClaudecheckSourceDir)
	case cfgClaude.PluginSourceGitHub:
		return desc.Text(text.DescKeyWriteClaudecheckSourceGitHub)
	case "":
		return desc.Text(text.DescKeyWriteClaudecheckUnknown)
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
		return desc.Text(text.DescKeyWriteClaudecheckEnabledBoth)
	case globally:
		return desc.Text(text.DescKeyWriteClaudecheckEnabledGlobal)
	case locally:
		return desc.Text(text.DescKeyWriteClaudecheckEnabledLocal)
	default:
		return desc.Text(text.DescKeyWriteClaudecheckNone)
	}
}
