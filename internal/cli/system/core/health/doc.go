//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package health monitors context health by detecting stale.
//
// Key exports: [ReadMapTracking], [CountModuleCommits],
// [EmitMapStalenessWarning], [UUIDPattern], [AutoPrune].
// Shared helpers used by sibling cmd/ packages.
// Used by core cmd/ packages.
package health
