//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core is the umbrella for the drift command's
// detection output, auto-fix, and sanitisation logic.
//
// # Overview
//
// The drift command detects stale paths, missing files,
// constitution violations, and other quality issues in
// the user's .context/ directory. This package groups
// three sub-packages that handle output formatting,
// automated fixes, and display sanitisation.
//
// # Sub-packages
//
//   - out: renders drift reports as human-readable
//     text (with icons and grouping) or as structured
//     JSON for machine consumption. Exports [out.DriftText]
//     and [out.DriftJSON].
//   - fix: applies automated corrections for fixable
//     drift issues. [fix.Apply] iterates the report,
//     archiving completed tasks for staleness issues
//     and creating files for missing-file issues.
//     Tracks results in [fix.Result].
//   - sanitize: converts internal check identifiers
//     to human-readable labels via [sanitize.FormatCheckName].
//
// # Data Flow
//
// The cmd layer loads context, calls drift.Detect to
// produce a Report, then delegates to this package:
//
//  1. out.DriftText or out.DriftJSON renders the
//     report for display.
//  2. If --fix is passed, fix.Apply walks the report
//     and attempts auto-remediation.
//  3. sanitize helpers are used by the output layer
//     to translate check names for display.
package core
