//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package knowledge is the **measurement helper** the
// `checkknowledge` system hook uses to evaluate whether
// a project's knowledge files (DECISIONS.md, LEARNINGS.md,
// CONVENTIONS.md) have outgrown the configured per-file
// thresholds and warrant a consolidation nudge.
//
// # Public Surface
//
//   - **[Health](contextDir, thresholds)**: emits the
//     two M5 signals — foldable roots (staging count
//     via disclosure.StagedEntries → /ctx-digest) and
//     heavy pages (bytes over the root and every theme
//     file → split / extract-to-tooling), foldable
//     first.
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
//   - `convention_section_count`  (default 12;
//     0 disables)
//   - `theme_page_byte_ceiling`   (default 65536;
//     0 disables)
//
// Crossing a foldability threshold means "consider
// running `/ctx-digest`"; crossing the weight ceiling
// means "split the theme or extract to tooling" —
// never "stop adding entries".
//
// # Concurrency
//
// Filesystem-bound. Stateless.
package knowledge
