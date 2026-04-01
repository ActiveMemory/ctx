//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package journal provides error constructors for journal pipeline operations.
//
// Error constructors return structured errors with context for
// user-facing messages routed through internal/assets text lookups.
// Exports: [LoadState], [SaveState], [LoadStateErr], [LoadStateFailed], [SaveStateFailed], [NoDir].
package journal
