//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package merge handles create-or-merge file operations during init.
//
// [OrCreate] creates a file from template, or merges the template's
// marked section into an existing file. [UpdateMarkedSection]
// replaces content between start/end markers. [SettingsPermissions]
// merges Claude Code permission settings. [Permissions] deduplicates
// and merges allow/deny permission lists.
package merge
