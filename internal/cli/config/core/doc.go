//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core contains business logic for the config
// command, which manages runtime configuration profiles
// (.ctxrc files).
//
// This package delegates its work to the profile
// subpackage. It does not export functions directly.
//
// # Profile Management (profile/)
//
// The profile subpackage provides three exported
// functions:
//
//   - [profile.Detect] reads the active profile name
//     from the parsed .ctxrc via rc.RC(). Returns an
//     empty string when no profile is set.
//   - [profile.SwitchTo] copies the requested profile
//     file over .ctxrc. It handles three cases: the
//     profile is already active (no-op message), .ctxrc
//     did not exist (created message), or a switch
//     occurred (switched message). The source file is
//     .ctxrc.base for the base profile and .ctxrc.dev
//     for the dev profile.
//   - [profile.Copy] performs the low-level file copy
//     from a source profile file to .ctxrc using safe
//     I/O helpers.
//
// [profile.GitRoot] resolves the git repository root
// directory, used by the cmd/ layer to locate project-
// level config files. It returns an error when git is
// unavailable or the working directory is outside a
// repository.
//
// # Data Flow
//
// The cmd/ layer calls profile.GitRoot to find the
// project root, then profile.Detect to check the
// current profile, and profile.SwitchTo to apply the
// requested change. Status messages are returned as
// strings for the write/ layer to display.
package core
