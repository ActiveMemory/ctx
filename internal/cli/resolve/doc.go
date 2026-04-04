//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package resolve provides shared CLI flag resolution helpers
// used across multiple command packages.
//
// [Tool] resolves the active tool identifier from the --tool
// flag or the .ctxrc tool field, returning an error when
// neither source provides a value.
package resolve
