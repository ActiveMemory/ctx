//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package task provides error constructors for task file operations.
//
// Error constructors return structured errors with context for
// user-facing messages routed through internal/assets text lookups.
// Exports: [FileNotFound], [FileRead], [FileWrite], [MultipleMatches], [NotFound], [NoneCompleted].
package task
