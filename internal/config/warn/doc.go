//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package warn provides Printf-style format string constants for
// best-effort warning messages routed through log.Warn.
//
// Key exports: [Close], [Write], [Remove], [Mkdir], [Rename], [Walk].
// Using constants prevents typo drift across 40+ call sites.
// Import as config/warn.
package warn
