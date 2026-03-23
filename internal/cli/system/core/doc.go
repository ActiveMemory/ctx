//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core provides shared helpers used by all system subcommand
// packages. It owns state management, hook I/O, message templating,
// SMB helpers, and session token utilities; subcommands import core,
// never each other.
package core
