//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package trigger provides error constructors for trigger operations.
//
// Error constructors return structured errors with context for
// user-facing messages routed through internal/assets text lookups.
// Exports: [Chmod], [CreateDir], [DiscoverFailed],
// [EmbeddedTemplateNotFound], [Exit], [InvalidJSONOutput],
// [InvalidType], [MarshalInput], [NotFound], [OverrideExists],
// [RemoveOverride], [ResolveHooksDir], [ResolvePath],
// [ScriptExists], [Stat], [StatPath], [Timeout], [Unknown],
// [UnknownVariant], [Validate], [WriteScript], [WriteOverride].
package trigger
