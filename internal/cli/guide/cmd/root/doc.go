//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the "ctx guide" command that
// lists available skills and CLI commands as a quick
// reference.
//
// # What It Does
//
// Displays help information about the ctx ecosystem.
// Without flags it shows a default overview. With
// --skills it lists all available slash-command skills
// from the embedded plugin. With --commands it lists
// every registered CLI command.
//
// # Flags
//
//   - --skills: List all available skills with
//     their trigger descriptions.
//   - --commands: List all CLI commands with their
//     short descriptions.
//
// # Output
//
// Default mode prints a concise getting-started
// guide. Skills mode prints a table of skill names
// and descriptions. Commands mode prints a table
// of command paths and short descriptions.
//
// # Delegation
//
// [Cmd] builds the cobra.Command with
// AnnotationSkipInit so it works before context
// initialization. [Run] dispatches to
// [skill.List], [command.List], or
// [guide.Default] based on flags.
package root
