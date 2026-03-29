//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package publish

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	writeIo "github.com/ActiveMemory/ctx/internal/write/line"
)

// NotFound prints that no published block was found.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - filename: source file name (e.g. "MEMORY.md").
func NotFound(cmd *cobra.Command, filename string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteUnpublishNotFound),
		filename))
}

// Unpublished prints that the published block was removed.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - filename: source file name (e.g. "MEMORY.md").
func Unpublished(cmd *cobra.Command, filename string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteUnpublishDone), filename))
}

// Plan prints the full publish plan: header, source files,
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
func Plan(
	cmd *cobra.Command,
	budget, tasks, decisions, conventions, learnings, totalLines int,
) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWritePublishHeader))
	cmd.Println()
	cmd.Println(desc.Text(text.DescKeyWritePublishSourceFiles))
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePublishBudget), budget))
	cmd.Println()
	cmd.Println(desc.Text(text.DescKeyWritePublishBlock))
	writeIo.Count(cmd, text.DescKeyWritePublishTasks, tasks)
	writeIo.Count(cmd, text.DescKeyWritePublishDecisions, decisions)
	writeIo.Count(cmd, text.DescKeyWritePublishConventions, conventions)
	writeIo.Count(cmd, text.DescKeyWritePublishLearnings, learnings)
	cmd.Println()
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWritePublishTotal),
		totalLines, budget))
}

// DryRun prints the dry-run notice.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func DryRun(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println()
	cmd.Println(desc.Text(text.DescKeyWritePublishDryRun))
}

// Done prints the success message with marker info.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func Done(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println()
	cmd.Println(desc.Text(text.DescKeyWritePublishDone))
}
