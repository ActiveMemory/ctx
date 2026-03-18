//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package publish

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/embed"
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
	cmd.Println(fmt.Sprintf(assets.TextDesc(embed.TextDescKeyWriteUnpublishNotFound), filename))
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
	cmd.Println(fmt.Sprintf(assets.TextDesc(embed.TextDescKeyWriteUnpublishDone), filename))
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
	cmd.Println(assets.TextDesc(embed.TextDescKeyWritePublishHeader))
	cmd.Println()
	cmd.Println(assets.TextDesc(embed.TextDescKeyWritePublishSourceFiles))
	cmd.Println(fmt.Sprintf(assets.TextDesc(embed.TextDescKeyWritePublishBudget), budget))
	cmd.Println()
	cmd.Println(assets.TextDesc(embed.TextDescKeyWritePublishBlock))
	if tasks > 0 {
		cmd.Println(fmt.Sprintf(assets.TextDesc(embed.TextDescKeyWritePublishTasks), tasks))
	}
	if decisions > 0 {
		cmd.Println(fmt.Sprintf(assets.TextDesc(embed.TextDescKeyWritePublishDecisions), decisions))
	}
	if conventions > 0 {
		cmd.Println(fmt.Sprintf(assets.TextDesc(embed.TextDescKeyWritePublishConventions), conventions))
	}
	if learnings > 0 {
		cmd.Println(fmt.Sprintf(assets.TextDesc(embed.TextDescKeyWritePublishLearnings), learnings))
	}
	cmd.Println()
	cmd.Println(fmt.Sprintf(assets.TextDesc(embed.TextDescKeyWritePublishTotal), totalLines, budget))
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
	cmd.Println(assets.TextDesc(embed.TextDescKeyWritePublishDryRun))
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
	cmd.Println(assets.TextDesc(embed.TextDescKeyWritePublishDone))
}
