//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root provides the "ctx hook message" parent
// command.
//
// # Overview
//
// This package groups message template management
// subcommands under a single namespace. It does not
// contain business logic itself; it delegates to its
// four children:
//
//   - list: displays all registered hook message
//     templates with their override status.
//   - show: prints the content of a specific template,
//     using the override if one exists.
//   - edit: creates a local override file so the user
//     can customize a template.
//   - reset: removes a local override to restore the
//     embedded default.
//
// # Usage
//
//	ctx hook message list [--json]
//	ctx hook message show <hook> <variant>
//	ctx hook message edit <hook> <variant>
//	ctx hook message reset <hook> <variant>
//
// # Behavior
//
// [Cmd] uses the parent.Cmd helper to build a
// cobra.Command with descriptions loaded from embedded
// assets. It attaches the list, show, edit, and reset
// subcommands as children. Running the parent without a
// subcommand prints the help text.
package root
