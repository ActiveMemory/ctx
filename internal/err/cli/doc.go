//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package cli defines the typed error constructors
// for generic CLI input validation: flags,
// arguments, interactive selections, and tool
// configuration. These errors are not tied to any
// single command; they appear wherever the CLI
// framework validates user input.
//
// # Domain
//
// Errors fall into two categories:
//
//   - **Input validation**: a required flag or
//     argument is missing, or an interactive
//     selection is out of range. Constructors:
//     [FlagRequired], [ArgRequired],
//     [InvalidSelection], [UnknownDocument].
//   - **Tool configuration**: no AI tool was
//     specified via flag or .ctxrc. Constructor:
//     [NoToolSpecified].
//
// # Wrapping Strategy
//
// These constructors return plain errors (no cause
// wrapping) because the failures are pure
// validation; there is no underlying system
// error to chain. All user-facing text is resolved
// through [internal/assets/read/desc].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package cli
