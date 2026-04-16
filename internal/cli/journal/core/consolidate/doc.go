//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package consolidate merges consecutive identical tool
// runs in journal markdown for cleaner reading.
//
// # Overview
//
// AI session transcripts often contain long sequences
// of identical tool calls (e.g. repeated file reads or
// test runs). This package detects those repetitions
// and collapses them into a single entry with a count
// annotation.
//
// # Public Surface
//
//   - [ToolRuns] -- collapses consecutive turns with
//     identical body content.
//
// # Algorithm
//
// ToolRuns processes the content line by line:
//
//  1. Scans for turn headers matching the turn header
//     regex pattern.
//  2. For each turn header, extracts the body content
//     up to the next header or end of file using the
//     turn.Body helper.
//  3. Counts consecutive turns that share the same
//     body content.
//  4. If the count exceeds one, emits a single copy
//     of the header and body followed by a count
//     annotation (e.g. "repeated 5 times").
//  5. If the count is one, preserves the original
//     lines unchanged.
//
// Non-turn lines (narrative text, headings, etc.) pass
// through unchanged.
package consolidate
