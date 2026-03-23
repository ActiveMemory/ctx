//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core provides shared scratchpad operations used by all pad
// subcommand packages. It owns storage, encryption, blob handling,
// parsing, and validation; subcommands import core, never each other.
package core
