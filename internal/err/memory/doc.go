//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package memory provides error constructors for memory bridge operations.
//
// Error constructors return structured errors with context for
// user-facing messages routed through internal/assets text lookups.
// Exports: [NotFound], [DiscoverFailed], [DiffFailed], [SelectContentFailed], [PublishFailed], [Read].
package memory
