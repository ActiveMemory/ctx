//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package command lists available CLI commands for the
// ctx guide output.
//
// The [List] function walks the root cobra command tree and
// prints every non-hidden subcommand with its short
// description. Hidden commands (system hooks, internal
// plumbing) are filtered out so the guide shows only
// user-facing surface area.
//
// Output is rendered through [internal/write/guide] which
// handles the terminal formatting: a header line followed
// by one indented line per command showing the command name
// and its one-line description.
//
// # Design Choice
//
// The guide command serves as an interactive onboarding
// tool. By listing commands from the live cobra tree rather
// than a static list, the output stays in sync with the
// actual binary. Adding or removing a subcommand
// automatically updates the guide.
package command
