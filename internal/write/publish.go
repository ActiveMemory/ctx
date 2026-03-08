//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"github.com/ActiveMemory/ctx/internal/write/config"
	"github.com/ActiveMemory/ctx/internal/write/io"
	"github.com/spf13/cobra"
)

// UnpublishNotFound prints that no published block was found.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - filename: source file name (e.g. "MEMORY.md").
func UnpublishNotFound(cmd *cobra.Command, filename string) {
	if cmd == nil {
		return
	}
	io.sprintf(cmd, config.tplUnpublishNotFound, filename)
}

// UnpublishDone prints that the published block was removed.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - filename: source file name (e.g. "MEMORY.md").
func UnpublishDone(cmd *cobra.Command, filename string) {
	if cmd == nil {
		return
	}
	io.sprintf(cmd, config.tplUnpublishDone, filename)
}

// PublishPlan prints the full publish plan: header, source files,
// budget, per-file counts, and total.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - budget: maximum line count for the published block.
//   - tasks: number of pending tasks selected.
//   - decisions: number of recent decisions selected.
//   - conventions: number of key conventions selected.
//   - learnings: number of recent learnings selected.
//   - totalLines: total lines in the published block.
func PublishPlan(
	cmd *cobra.Command,
	budget, tasks, decisions, conventions, learnings, totalLines int,
) {
	if cmd == nil {
		return
	}
	cmd.Println(config.tplPublishHeader)
	cmd.Println()
	cmd.Println(config.tplPublishSourceFiles)
	io.sprintf(cmd, config.tplPublishBudget, budget)
	cmd.Println()
	cmd.Println(config.tplPublishBlock)
	if tasks > 0 {
		io.sprintf(cmd, config.tplPublishTasks, tasks)
	}
	if decisions > 0 {
		io.sprintf(cmd, config.tplPublishDecisions, decisions)
	}
	if conventions > 0 {
		io.sprintf(cmd, config.tplPublishConventions, conventions)
	}
	if learnings > 0 {
		io.sprintf(cmd, config.tplPublishLearnings, learnings)
	}
	cmd.Println()
	io.sprintf(cmd, config.tplPublishTotal, totalLines, budget)
}

// PublishDryRun prints the dry-run notice.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func PublishDryRun(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println()
	cmd.Println(config.tplPublishDryRun)
}

// PublishDone prints the success message with marker info.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func PublishDone(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println()
	cmd.Println(config.tplPublishDone)
}
