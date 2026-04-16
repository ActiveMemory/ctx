//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the "ctx compact" command for
// cleaning up context files by archiving completed tasks
// and removing empty sections.
//
// # What It Does
//
// The command performs two housekeeping passes over
// the .context/ directory:
//
//  1. Task compaction -- moves completed tasks in
//     TASKS.md to a "Completed (Recent)" section.
//     When --archive or .ctxrc auto_archive is
//     enabled, older completed tasks are moved to
//     .context/archive/ files.
//
//  2. Section cleanup -- removes empty markdown
//     sections from all context files (TASKS.md,
//     DECISIONS.md, LEARNINGS.md, CONVENTIONS.md).
//
// # Flags
//
//   - --archive: Create .context/archive/ for old
//     completed tasks. Also enabled automatically
//     when the auto_archive option is set in .ctxrc.
//
// # Output
//
// Prints a heading, per-file change counts (tasks
// moved, sections removed), and a final summary
// line. When nothing changed it prints "all clean".
//
// # Delegation
//
// [Cmd] builds the cobra.Command with the --archive
// flag. [Run] loads context via [context/load],
// delegates task compaction to [core/task], reloads
// context, then runs [tidy.CompactContext] for
// section cleanup.
package root
