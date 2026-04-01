//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package journal scans journal directories for the newest file.
//
// Key exports: [NewestMtime], [CountNewerFiles], [CountUnenriched],
// [CheckStage], [MarkStage].
// Shared helpers used by sibling cmd/ packages.
// Used by core cmd/ packages.
package journal
