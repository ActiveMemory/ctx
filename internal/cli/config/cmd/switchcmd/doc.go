//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package switchcmd implements the "ctx config switch"
// subcommand for switching between .ctxrc profiles.
//
// # What It Does
//
// Switches the active configuration by copying a
// named profile file over .ctxrc. Profiles are stored
// as .ctxrc.<name> files in the project root (e.g.
// .ctxrc.dev, .ctxrc.base).
//
// # Arguments
//
// An optional positional argument specifying the
// target profile:
//
//   - dev: switch to the development profile
//   - base: switch to the base (production) profile
//   - prod: alias for base
//   - (none): toggle between dev and base
//
// # Flags
//
// None.
//
// # Output
//
// A confirmation line showing which profile is now
// active, e.g. "Switched to dev profile."
//
// # Delegation
//
// [Cmd] builds the cobra.Command with the
// AnnotationSkipInit annotation so it works before
// full context initialization. [Run] normalizes the
// profile name (e.g. "prod" -> "base"), calls
// [profile.SwitchTo] to copy the profile file, and
// writes confirmation via [writeConfig.SwitchConfirm].
package switchcmd
