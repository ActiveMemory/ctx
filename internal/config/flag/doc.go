//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package flag centralizes CLI flag name constants and
// their shorthand letters for every ctx command.
//
// The ctx CLI uses cobra for command registration.
// Rather than embedding flag name strings directly in
// each command's init function, every flag name and
// shorthand is defined here. This prevents naming
// collisions, keeps flag names consistent across
// commands, and makes it easy to find every flag the
// CLI accepts.
//
// # Global Flags
//
//   - Tool ("tool"): override the active AI tool
//     identifier (e.g. claude, cursor, kiro).
//
// PrefixLong ("--") is the long-flag prefix used in
// error messages and help text formatting.
//
// # Add Command Flags
//
// Flags specific to the add command for creating
// context entries:
//
//   - Application, Branch, Commit, Consequence,
//     Context, File, Lesson, Priority, Rationale,
//     Section
//   - Corresponding ShortApplication ("a"),
//     ShortContext ("c"), ShortFile ("f"), etc.
//
// # Agent Command Flags
//
//   - Budget, Cooldown, Follow, Format, Session,
//     Skill: control agent behavior and output
//
// # Shared Flags
//
// Flags used across multiple commands (a large set
// including After, All, DryRun, Force, JSON, Limit,
// Output, Quiet, Verbose, and many more). Each has a
// corresponding Short* constant where a shorthand
// letter is assigned.
//
// # Time-Range Flags
//
//   - Log, Since, Until: used by commands that
//     filter by date range
//
// # Why Centralized
//
// Flag names appear in command registration, flag
// value retrieval, error messages, and shell completion
// generators. Scattering these strings would cause
// silent breakage when a flag is renamed in one place
// but not another. This package makes renames a
// single-point edit with compile-time verification.
package flag
