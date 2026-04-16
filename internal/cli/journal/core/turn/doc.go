//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package turn handles **conversation turn parsing and
// merging** in journal markdown: finding where each
// user/assistant turn begins and ends, and merging
// adjacent turns from the same role into one block when
// the original transcript had artificial splits.
//
// The package operates on already-normalized journal
// content (after [normalize] and friends have run); it
// is the per-turn slicer the renderers and the
// per-turn-anchor navigator both rely on.
//
// # Public Surface
//
//   - **[Body](lines, startIdx)**: extracts the body
//     text of a single turn starting at `startIdx`.
//     Reads forward to the next turn header (or EOF)
//     and returns the in-between lines.
//   - **[MergeConsecutive](lines)**: collapses
//     adjacent turns from the same role into a
//     single combined block. Useful when Claude
//     Code split a long assistant response across
//     two consecutive `assistant:` turns due to
//     internal pacing.
//
// # Concurrency
//
// Pure data transformation. Concurrent callers never
// race.
package turn
