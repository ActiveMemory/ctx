//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the "ctx fmt" command logic.
//
// [Cmd] registers the command and its flags (--width, --check).
// [Run] iterates over the four context files, applies wrapping,
// and writes back only files that changed.
package root
