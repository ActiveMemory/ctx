//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package box defines box-drawing characters and layout
// constants used to render nudge boxes in the terminal.
//
// When ctx needs to surface a nudge (a contextual hint or
// reminder) to the user, it wraps the message in a Unicode
// box with consistent borders. This package provides the
// border characters, width constraints, and separators
// that make those boxes visually uniform.
//
// # Box Structure
//
// A nudge box has three visual regions:
//
//   - [Top] renders the top-left corner and opening border
//     (e.g. "--- ctx: backup stale").
//   - [LinePrefix] prefixes every content line with a left
//     border character for visual containment.
//   - [Bottom] closes the box with a horizontal rule.
//
// # Layout Constants
//
//   - [NudgeBoxWidth] sets the inner width to 51 characters
//     so boxes fit in standard 80-column terminals with
//     room for indentation.
//   - [BorderFill] is the repeating dash character used to
//     pad the top border to the target width.
//
// # Navigation Links
//
// Nudge boxes sometimes include navigation links at the
// bottom (e.g. "docs | settings | help"). [PipeSeparator]
// provides the " | " delimiter between those links.
//
// # Why Centralized
//
// Multiple nudge producers (backup staleness, ceremony
// reminders, architecture drift) share the same box
// layout. Centralizing the border constants ensures
// visual consistency and makes width changes atomic.
package box
