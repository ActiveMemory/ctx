//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package parser provides error constructors for session transcript parsing.
//
// Error constructors return structured errors with context for
// user-facing messages routed through internal/assets text lookups.
// Exports: [ReadFile], [OpenFile], [NoMatch], [WalkDir], [FileError], [ScanFile].
package parser
