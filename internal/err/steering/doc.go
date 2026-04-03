//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package steering provides error constructors for steering operations.
//
// Error constructors return structured errors with context for
// user-facing messages routed through internal/assets text lookups.
// Exports: [ComputeRelPath], [ContextDirMissing], [CreateDir],
// [FileExists], [InvalidYAML], [MissingClosingDelimiter],
// [MissingOpeningDelimiter], [NoTool], [OutputEscapesRoot],
// [Parse], [ReadDir], [ReadFile], [ResolveOutput],
// [ResolveRoot], [SyncAll], [SyncName], [UnsupportedTool],
// [WriteFile], [WriteSteeringFile], [WriteInitFile].
package steering
