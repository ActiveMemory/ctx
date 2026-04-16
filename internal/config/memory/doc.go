//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package memory defines constants, paths, and types
// for the memory bridge subsystem, which imports
// Claude Code's auto-generated MEMORY.md into the
// structured context files (CONVENTIONS.md,
// DECISIONS.md, LEARNINGS.md, TASKS.md).
//
// # Bridge Files
//
// The bridge operates on three files inside
// .context/memory/:
//
//   - Source ("MEMORY.md") -- the original Claude Code
//     auto-memory file.
//   - Mirror ("mirror.md") -- a raw copy used for
//     diff-based change detection.
//   - State ("memory-import.json") -- tracks which
//     entries have been classified and imported.
//
// PathMemoryMirror provides the pre-joined relative
// path from the project root to the mirror file.
//
// # Classification Rules
//
// Each memory entry is classified into a target
// context file using keyword heuristics. The
// DefaultClassifyRules slice maps keyword patterns
// (e.g., "always use", "gotcha", "decided") to
// targets (convention, decision, learning, task).
// Entries that match no rule receive TargetSkip.
//
// The ClassifyRule type holds a Target string and
// a Keywords slice, and is user-overridable via the
// classify_rules key in .ctxrc.
//
// # Publish Budgets
//
// When publishing context back to Claude Code the
// bridge enforces line and entry budgets:
//
//   - DefaultPublishBudget (80 lines total)
//   - PublishMaxTasks (10), PublishMaxDecisions (5),
//     PublishMaxConventions (10), PublishMaxLearnings (5)
//   - PublishRecentDays (7) -- lookback window
package memory
