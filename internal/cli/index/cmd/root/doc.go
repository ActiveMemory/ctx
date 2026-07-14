//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root wires the `ctx index` cobra command.
//
// Cmd builds the command (positional FILE argument, --depth and --json
// flags); Run reads the file, delegates heading extraction to
// [internal/heading], and rendering to [internal/write/index]. Keeping the
// recognizer and renderer in their own packages leaves this package as thin
// argument-and-flag plumbing.
package root
