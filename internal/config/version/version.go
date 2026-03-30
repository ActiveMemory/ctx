//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package version

// Check-version configuration.
const (
	// ThrottleID is the state file name for daily throttle of version checks.
	ThrottleID = "version-checked"
	// DevBuild is the version string used for development builds.
	DevBuild = "dev"
)

// Project-root paths for version checking.
const (
	// DirClaudePlugin is the Claude plugin directory at project root.
	DirClaudePlugin = ".claude-plugin"
	// FileMarketplace is the marketplace manifest filename.
	FileMarketplace = "marketplace.json"
)
