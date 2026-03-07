//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

// InfoPathConversionExists reports that a path conversion target already
// exists at the destination. Used during init to show which template files
// were skipped.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - rootDir: project root directory for path resolution.
//   - oldPath: original template-relative path.
//   - newPath: destination-relative path joined with rootDir.
func InfoPathConversionExists(
	cmd *cobra.Command, rootDir, oldPath, newPath string,
) {
	if cmd == nil {
		return
	}
	sprintf(cmd, tplPathExists, oldPath, filepath.Join(rootDir, newPath))
}

// InfoAddedTo confirms an entry was added to a context file.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: Name of the file the entry was added to
func InfoAddedTo(cmd *cobra.Command, filename string) {
	sprintf(cmd, tplAddedTo, filename)
}

// InfoMovingTask reports a completed task being moved.
//
// Parameters:
//   - cmd: Cobra command for output
//   - taskText: Truncated task description
func InfoMovingTask(cmd *cobra.Command, taskText string) {
	sprintf(cmd, tplMovingTask, taskText)
}

// InfoSkippingTask reports a task skipped due to incomplete children.
//
// Parameters:
//   - cmd: Cobra command for output
//   - taskText: Truncated task description
func InfoSkippingTask(cmd *cobra.Command, taskText string) {
	sprintf(cmd, tplSkippingTask, taskText)
}

// InfoArchivedTasks reports the number of tasks archived.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: Number of tasks archived
//   - archiveFile: Path to the archive file
//   - days: Age threshold in days
func InfoArchivedTasks(cmd *cobra.Command, count int, archiveFile string, days int) {
	sprintf(cmd, tplArchivedTasks, count, archiveFile, days)
}

// InfoCompletedTask reports a task marked complete.
//
// Parameters:
//   - cmd: Cobra command for output
//   - taskText: The completed task description
func InfoCompletedTask(cmd *cobra.Command, taskText string) {
	sprintf(cmd, tplCompletedTask, taskText)
}

// InfoConfigProfileDev reports that the dev profile is active.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoConfigProfileDev(cmd *cobra.Command) {
	cmd.Println(tplConfigProfileDev)
}

// InfoConfigProfileBase reports that the base profile is active.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoConfigProfileBase(cmd *cobra.Command) {
	cmd.Println(tplConfigProfileBase)
}

// InfoConfigProfileNone reports that no profile exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: The .ctxrc filename
func InfoConfigProfileNone(cmd *cobra.Command, filename string) {
	sprintf(cmd, tplConfigProfileNone, filename)
}

// InfoExistsWritingAsAlternative reports that a file already exists and the
// content is being written to an alternative filename instead.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - path: the original target path that already exists.
//   - alternative: the fallback path where content was written.
func InfoExistsWritingAsAlternative(
	cmd *cobra.Command, path, alternative string,
) {
	if cmd == nil {
		return
	}
	sprintf(cmd, tplExistsWritingAsAlternative, path, alternative)
}
