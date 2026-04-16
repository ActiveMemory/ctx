//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package text holds the **lookup keys** for every piece of
// user-facing display text emitted anywhere in ctx: error
// messages, status banners, MCP responses, hook nudges,
// admonition templates, you name it.
//
// The package is one half of a deliberate two-step indirection:
//
//  1. **Here**: typed `DescKeyXxx` Go constants. Compile-time
//     guarantees that every reference is a real key.
//  2. **In** [internal/assets/commands/text/*.yaml]: the
//     actual strings, embedded into the binary at build time.
//     Reachable via [internal/assets/read/desc.Text](key).
//
// The split exists for three reasons:
//
//   - **Editing copy stops touching Go code**: copywriters,
//     translators, and product owners can edit YAML without
//     a Go toolchain.
//   - **i18n is structurally possible**: adding a locale is
//     a parallel YAML tree, not a fork of every package.
//   - **One sentence cannot quietly drift between two
//     callers**: both grab the same key, both render the
//     same text.
//
// # File Layout
//
// One Go file per subsystem (`agent.go`, `bootstrap.go`,
// `mcp.go`, `journal.go`, `steering.go`, …). Each file groups
// the `DescKeyXxx` constants for that subsystem so adding a
// new message to the agent flow only edits `agent.go` and the
// corresponding YAML file.
//
// # Naming Convention
//
// Constants follow `DescKey<Subsystem><Specifier>`; the
// underlying YAML key follows the dotted form
// `<subsystem>.<specifier>`. Both halves are validated by the
// `desckey_namespace_test` audit so a typo in either side
// fails CI.
//
// # Consumers
//
// Pretty much every CLI subcommand, every MCP handler, and
// every write-side terminal-output package imports this
// package. Common usage:
//
//	import (
//	    "github.com/ActiveMemory/ctx/internal/assets/read/desc"
//	    "github.com/ActiveMemory/ctx/internal/config/embed/text"
//	)
//	msg := desc.Text(text.DescKeyAgentInstruction)
package text
