//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package claude provides access to Claude Code
// integration files embedded in the assets filesystem.
//
// # CLAUDE.md Template
//
// Md returns the embedded CLAUDE.md template for
// project-level Claude Code instructions. This file
// is deployed to the project root during ctx init,
// separate from the .context/ template files.
//
//	content, err := claude.Md()
//
// # Plugin Version
//
// PluginVersion extracts the semver version string
// from the embedded plugin.json manifest. This is
// used to report the ctx plugin version without
// parsing the manifest at the call site.
//
//	version, err := claude.PluginVersion()
//	// => "0.8.1"
//
// # Plugin Manifest
//
// The plugin.json file follows the Claude Code plugin
// specification and declares the ctx plugin identity,
// version, and capabilities. It is embedded alongside
// CLAUDE.md under the claude/ asset directory.
package claude
