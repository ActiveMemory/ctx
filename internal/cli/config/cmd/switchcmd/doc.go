//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package switchcmd implements the ctx config switch subcommand.
//
// [Cmd] builds the cobra.Command. [Run] switches the active
// .ctxrc profile by copying the named profile file over .ctxrc.
// Profiles are stored as .ctxrc.<name> files in the project root.
package switchcmd
