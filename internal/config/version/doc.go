//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package version defines constants for ctx's version
// checking, update detection, and release management
// subsystem.
//
// ctx can compare its binary version against the
// VERSION file, the Claude plugin manifest, and the
// VS Code marketplace manifest to detect version
// drift. This package provides the file paths,
// template variable keys, and throttle identifiers
// the checker needs.
//
// # Template Variables
//
// Hook templates can interpolate version data using
// these keys:
//
//   - [VarBinary]      — the running binary version.
//   - [VarFile]        — version from the VERSION
//     file at project root.
//   - [VarPlugin]      — version from the Claude
//     plugin manifest.
//   - [VarMarketplace] — version from the VS Code
//     marketplace manifest.
//   - [VarKeyAgeDays]  — API key age in days.
//
// # Project Paths
//
//   - [FileVersion] ("VERSION") — the project-root
//     version file.
//   - [DirClaudePlugin] (".claude-plugin") — the
//     Claude plugin directory.
//   - [FileMarketplace] ("marketplace.json") — the
//     marketplace manifest filename.
//
// # Throttling
//
//   - [ThrottleID] ("version-checked") — state file
//     name for daily check throttling.
//   - [DevBuild] ("dev") — version string for
//     development builds (skips update checks).
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
package version
