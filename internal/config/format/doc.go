//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package format provides numeric and display constants
// used when rendering human-readable output in the ctx
// CLI.
//
// Whenever the CLI needs to format byte counts, token
// numbers, hash prefixes, or truncated text, it uses
// the thresholds and widths defined here instead of
// inline magic numbers. This keeps display behavior
// consistent across commands.
//
// # SI and IEC Thresholds
//
//   - SIThreshold (1000) -- boundary between raw
//     numbers and abbreviated display (e.g. 999 vs
//     1.0K)
//   - SIThresholdM (1,000,000) -- boundary between
//     K and M abbreviations
//   - IECUnit (1024) -- binary base for byte
//     formatting (KiB, MiB, etc.)
//   - IECPrefixes ("KMGTPE") -- unit prefix letters
//     for binary byte formatting
//
// # Hash Display
//
//   - HashPrefixLen (8) -- number of hex characters
//     shown for truncated commit or content hashes
//
// # Text Truncation Widths
//
// These constants control how long strings are clipped
// in compact displays:
//
//   - TruncateDetail (120) -- max width for detail
//     strings in governance violation reports
//   - TruncateTitle (60) -- max width for titles in
//     import previews and list views
//   - TruncateDescription (70) -- max width for
//     descriptions in skill listings
//
// # Preview Line Counts
//
//   - PreviewLines (5) -- content lines shown in
//     status previews
//   - StatusPreviewLines (3) -- content lines shown
//     in verbose status file previews
//   - StatusRecentFiles (3) -- recently modified
//     files shown in the status activity section
//
// # Why Centralized
//
// Display constants are consumed by multiple rendering
// paths: status output, doctor reports, skill
// listings, governance summaries, and import previews.
// Defining them in one place ensures visual consistency
// and makes it easy to adjust column widths or
// thresholds globally.
package format
