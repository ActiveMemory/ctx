//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package publish

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
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
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteUnpublishNotFound), filename))
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
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteUnpublishDone), filename))
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
	cmd.Println(desc.Text(text.DescKeyWritePublishHeader))
	cmd.Println()
	cmd.Println(desc.Text(text.DescKeyWritePublishSourceFiles))
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePublishBudget), budget))
	cmd.Println()
	cmd.Println(desc.Text(text.DescKeyWritePublishBlock))
	if tasks > 0 {
		cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePublishTasks), tasks))
	}
	if decisions > 0 {
		cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePublishDecisions), decisions))
	}
	if conventions > 0 {
		cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePublishConventions), conventions))
	}
	if learnings > 0 {
		cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePublishLearnings), learnings))
	}
	cmd.Println()
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePublishTotal), totalLines, budget))
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
	cmd.Println(desc.Text(text.DescKeyWritePublishDryRun))
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
	cmd.Println(desc.Text(text.DescKeyWritePublishDone))
}
