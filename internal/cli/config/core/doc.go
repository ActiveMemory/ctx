//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core contains configuration management helpers for the
// config command.
//
// [DetectProfile] reads the active profile name from .ctxrc.
// [SwitchTo] copies a named profile over .ctxrc. [CopyProfile]
// performs the file copy. [GitRoot] resolves the repository root
// for locating project-level config files.
package core
