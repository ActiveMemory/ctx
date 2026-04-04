//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cmd

// Use strings for task subcommands.
const (
	// UseTaskArchive is the cobra Use string for the task archive command.
	UseTaskArchive = "archive"
	// UseTaskSnapshot is the cobra Use string for the task snapshot command.
	UseTaskSnapshot = "snapshot [name]"
)

// DescKeys for task subcommands.
const (
	// DescKeyTask is the description key for the task command.
	DescKeyTask = "task"
	// DescKeyTaskArchive is the description key for the task archive command.
	DescKeyTaskArchive = "task.archive"
	// DescKeyTaskSnapshot is the description key for the task snapshot command.
	DescKeyTaskSnapshot = "task.snapshot"
)
