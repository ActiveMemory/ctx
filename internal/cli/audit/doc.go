//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package audit implements the "ctx audit" command for the
// out-of-band audit channel. Audit reports live as
// `.context/audit/<id>.md` files, dropped there by a
// separate Claude Code session running an audit skill
// (e.g. /ctx-surface-audit). The
// `ctx system checkaudit` UserPromptSubmit hook
// verbatim-relays each not-yet-dismissed report on the
// next interactive turn.
//
// # Subcommands
//
//   - list: show all reports with status and age (default)
//   - show: print one report's body verbatim
//   - dismiss: dismiss one or all reports
//
// # Subpackages
//
//	cmd/list, cmd/show, cmd/dismiss: subcommand impls
//	core/store: report I/O + dismissal ledger
//	core/parse: YAML frontmatter + body parser
//
// See specs/audit-channel.md for the design rationale.
package audit
