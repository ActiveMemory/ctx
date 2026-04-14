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
// The package is a typed constants registry — no logic.
//
// # Public Surface
//
//   - **[StatusIcon]** — map from status level to
//     the rendered glyph; consumed by `ctx status`,
//     `ctx doctor`, and the per-section size lines
//     in the agent context packet.
//   - **Threshold percentage constants** — `ok` /
//     `warn` / `danger` boundaries the renderers
//     use to pick the matching icon.
//   - **Format-string constants** — the per-line
//     templates used to render "FILE:  N tokens
//     (PCT%)" rows.
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
//
// # Related Packages
//
//   - [internal/cli/status]              — chief
//     consumer for the status one-liner.
//   - [internal/cli/agent]               — uses
//     [StatusIcon] in the per-section budget
//     summary at the head of the context packet.
//   - [internal/cli/doctor]              — uses
//     the same icons in its checklist output.
package stats
