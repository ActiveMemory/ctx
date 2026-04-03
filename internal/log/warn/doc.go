//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package warn provides a centralized stderr warning sink for best-effort
// operations whose errors would otherwise be silently discarded.
//
// [Warn] formats and writes a message to stderr. The sink is
// replaceable for testing. Callers use this instead of log.Println
// to keep warning output consistent across the codebase.
//
// Key exports: [Warn].
// Used throughout ctx for non-fatal error reporting.
package warn
