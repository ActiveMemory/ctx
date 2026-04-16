//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check_version implements the
// **`ctx system check-version`** hidden hook, which
// warns when the ctx binary and the installed plugin
// are on different versions.
//
// # What It Does
//
// The hook compares the running binary version (from
// cobra's root command) against the embedded plugin
// version (from the Claude Code plugin manifest). If
// the major.minor components differ, it emits a nudge
// warning about the version mismatch and suggesting
// an upgrade.
//
// Dev builds (version string "dev") are silently
// skipped. If the plugin version cannot be read, an
// error nudge is emitted instead.
//
// As a secondary check, the hook piggybacks a key
// rotation age check on the daily version check cycle.
//
// # Input
//
// A JSON hook envelope on stdin with session metadata.
//
// # Output
//
// On version mismatch: a nudge box showing the binary
// and plugin versions with upgrade instructions.
// On match, dev build, or throttled: no output.
// On key age issue: an additional nudge from the
// key-age checker.
//
// # Throttling
//
// The hook is throttled to fire at most once per day
// using a marker file in the state directory.
//
// # Delegation
//
// [Cmd] builds the hidden cobra command. [Run] reads
// stdin via [core/check.Preamble], reads the plugin
// version from [assets/read/claude.PluginVersion],
// parses major.minor via [core/version.ParseMajorMinor],
// loads the mismatch message via [core/message.Load],
// and emits the nudge through [write/setup.Nudge].
// Key age is checked via [core/version.CheckKeyAge].
package check_version
