//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package hook provides error constructors for hook template operations.
//
// Error constructors return structured errors with context for
// user-facing messages routed through internal/assets text lookups.
// Exports: [EmbeddedTemplateNotFound], [OverrideExists], [WriteOverride], [RemoveOverride], [Unknown], [UnknownVariant].
package hook
