//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package status implements the "ctx config status"
// subcommand that displays the active .ctxrc profile.
//
// # What It Does
//
// Detects which configuration profile is currently
// active (base or dev) and prints it to stdout. The
// profile is determined by comparing the contents of
// .ctxrc against known profile files (.ctxrc.dev,
// .ctxrc.base).
//
// # Flags
//
// None. The command accepts no arguments.
//
// # Output
//
// A single line identifying the active profile name
// (e.g. "base" or "dev"). This is useful for scripts
// and shell prompts that need to know which config
// is active.
//
// # Delegation
//
// [Cmd] builds the cobra.Command with the
// AnnotationSkipInit annotation so it works before
// full context initialization. [Run] calls
// [profile.Detect] to identify the active profile
// and writes the result through
// [writeConfig.ProfileStatus].
package status
