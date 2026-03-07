//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

// prefixError is prepended to all error messages written to stderr.
const prefixError = "Error: "

// tplPathExists is a format template for reporting that a destination path
// already exists. Arguments: original path, resolved destination path.
const tplPathExists = "  %s -> %s (exists)"

// tplExistsWritingAsAlternative is a format template for reporting that a
// file exists and content was written to an alternative filename instead.
// Arguments: original path, alternative path.
const tplExistsWritingAsAlternative = "  ! %s exists, writing as %s"

// tplDryRun is printed when a command runs in dry-run mode.
const tplDryRun = "Dry run — no files will be written."

// tplSource is a format template for reporting a source path.
// Arguments: path.
const tplSource = "  Source: %s"

// tplMirror is a format template for reporting a mirror path.
// Arguments: relative mirror path.
const tplMirror = "  Mirror: %s"

// tplStatusDrift is printed when drift is detected.
const tplStatusDrift = "  Status: drift detected (source is newer)"

// tplStatusNoDrift is printed when no drift is detected.
const tplStatusNoDrift = "  Status: no drift"

// tplArchived is a format template for reporting an archived file.
// Arguments: archive filename.
const tplArchived = "Archived previous mirror to %s"

// tplSynced is a format template for reporting a successful sync.
// Arguments: source label, destination relative path.
const tplSynced = "Synced %s -> %s"

// tplLines is a format template for reporting line counts.
// Arguments: line count.
const tplLines = "  Lines: %d"

// tplLinesPrevious is a format template appended to line counts when a
// previous count is available. Arguments: previous line count.
const tplLinesPrevious = " (was %d)"

// tplNewContent is a format template for reporting new content since last sync.
// Arguments: line count.
const tplNewContent = "  New content: %d lines since last sync"

// tplAddedTo is a format template for confirming an entry was added.
// Arguments: filename.
const tplAddedTo = "✓ Added to %s"

// tplMovingTask is a format template for a task being moved to completed.
// Arguments: truncated task text.
const tplMovingTask = "✓ Moving completed task: %s"

// tplSkippingTask is a format template for a task skipped due to
// incomplete children. Arguments: truncated task text.
const tplSkippingTask = "! Skipping (has incomplete children): %s"

// tplArchivedTasks is a format template for archived tasks summary.
// Arguments: count, archive file path, days threshold.
const tplArchivedTasks = "✓ Archived %d tasks to %s (older than %d days)"

// tplCompletedTask is a format template for a task marked complete.
// Arguments: task text.
const tplCompletedTask = "✓ Completed: %s"

// tplConfigProfileDev is the status output for dev profile.
const tplConfigProfileDev = "active: dev (verbose logging enabled)"

// tplConfigProfileBase is the status output for base profile.
const tplConfigProfileBase = "active: base (defaults)"

// tplConfigProfileNone is the status output when no profile exists.
// Arguments: ctxrc filename.
const tplConfigProfileNone = "active: none (%s does not exist)"
