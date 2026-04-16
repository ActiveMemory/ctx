//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package mark_journal provides terminal output for the
// journal processing stage marker commands (ctx journal
// mark).
//
// # Exported Functions
//
// [StageChecked] prints the result of a --check query,
// showing the journal filename, the processing stage
// name, and the current stage value. This lets users
// inspect which stages have been completed for a given
// journal entry.
//
// [StageMarked] prints a confirmation after a stage is
// marked as complete, showing the journal filename and
// the stage name that was just recorded.
//
// # Message Categories
//
//   - Info: stage check results and mark confirmations
//
// # Nil Safety
//
// Both functions treat a nil *cobra.Command as a no-op.
//
// # Usage
//
//	if checkOnly {
//	    mark_journal.StageChecked(cmd, file, stage, val)
//	} else {
//	    mark_journal.StageMarked(cmd, file, stage)
//	}
package mark_journal
