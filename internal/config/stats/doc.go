//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package stats holds the **constants** used by ctx's
// context-size monitoring and status reporting: the
// status-icon glyphs (`✓`, `⚠`, `✗`), the threshold
// percentages that decide which icon to use, and the
// per-line format strings the renderers consume.
//
// The package is a typed constants registry, not logic.
//
// # Public Surface
//
//   - **[StatusIcon]**: map from status level to
//     the rendered glyph; consumed by `ctx status`,
//     `ctx doctor`, and the per-section size lines
//     in the agent context packet.
//   - **Threshold percentage constants**: `ok` /
//     `warn` / `danger` boundaries the renderers
//     use to pick the matching icon.
//   - **Format-string constants**: the per-line
//     templates used to render "FILE:  N tokens
//     (PCT%)" rows.
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
package stats
