//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package cli defines the typed error constructors and
// sentinels used by the `ctx kb` CLI subcommands. The package
// is the single source for refusal messages (no question, no
// sources, missing managed block, etc.) and the wrappers
// surfaced around `io.Safe*` failures.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/config/kb/cli]
//     supplies the message and format-string constants.
//   - [github.com/ActiveMemory/ctx/internal/cli/kb/cmd/...]
//     are the primary callers.
package cli
