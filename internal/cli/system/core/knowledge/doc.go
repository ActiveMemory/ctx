//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package knowledge is the **measurement helper** the
// `check_knowledge` system hook uses to evaluate whether
// a project's knowledge files (DECISIONS.md, LEARNINGS.md,
// CONVENTIONS.md) have outgrown the configured per-file
// thresholds and warrant a consolidation nudge.
//
// # Public Surface
//
//   - **[ScanFiles](contextDir)**: counts entries
//     in DECISIONS.md and LEARNINGS.md, and lines
//     in CONVENTIONS.md, and returns the result.
//   - **[FormatWarnings](report, thresholds)**:
//     turns the scan into the human-readable
//     warning text emitted via the VERBATIM relay.
//   - **[EmitWarning](text)**: writes the warning
//     through the standard nudge path.
//   - **[CheckHealth](contextDir)**: convenience;
//     scan + threshold compare in one call;
//     returns the warning text or empty.
//
// # Thresholds
//
// Per-file thresholds come from `.ctxrc`:
//
//   - `entry_count_decisions`     (default 20;
//     0 disables)
//   - `entry_count_learnings`     (default 30;
//     0 disables)
//   - `convention_line_count`     (default 200;
//     0 disables)
//
// Crossing a threshold means "consider running
// `/ctx-consolidate`", not "stop adding entries".
//
// # Concurrency
//
// Filesystem-bound. Stateless.
package knowledge
