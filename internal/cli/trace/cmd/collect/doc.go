//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package collect provides the hidden "ctx trace collect" CLI subcommand.
//
// It gathers context refs from all sources and outputs them as a git
// commit trailer. When invoked with --record, it records refs from a
// commit trailer into history and truncates pending state.
//
// Key exports: [Cmd], [Run].
// Called by git hooks to inject and persist context refs at commit time.
package collect
