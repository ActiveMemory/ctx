//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the "ctx drift" command that
// detects stale or invalid context and optionally
// auto-fixes supported issues.
//
// # What It Does
//
// Loads the parsed context from .context/, runs the
// drift detection engine, and reports any issues
// found. Detected issue categories include:
//
//   - Broken path references in context files
//   - Staleness indicators (old timestamps, etc.)
//   - Constitution violations
//   - Missing required files
//
// When --fix is set, the command attempts to resolve
// supported issue types automatically, then re-runs
// detection to show the updated state.
//
// # Flags
//
//   - --json: Output results as machine-readable
//     JSON instead of human-readable text.
//   - --fix: Attempt to auto-fix supported issues
//     (staleness, missing_file). Prints a summary
//     of fixed, skipped, and errored items.
//
// # Output
//
// In text mode, prints warnings and violations with
// file paths and descriptions. In JSON mode, outputs
// a structured report. When --fix is active, prints
// a fix header, per-item results, counts, and then
// the re-checked state.
//
// # Delegation
//
// [Cmd] builds the cobra.Command with --json and
// --fix flags. [Run] loads context via
// [context/load], runs [drift.Detect], optionally
// applies fixes via [core/fix], and renders output
// through [core/out].
package root
