//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// CtxSyncInSync prints the all-clear message when context is in sync.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func CtxSyncInSync(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWriteSyncInSync))
}

// CtxSyncHeader prints the sync analysis heading and optional dry-run notice.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - dryRun: If true, includes the dry-run notice.
func CtxSyncHeader(cmd *cobra.Command, dryRun bool) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWriteSyncHeader))
	cmd.Println(desc.Text(text.DescKeyWriteSyncSeparator))
	cmd.Println()
	if dryRun {
		cmd.Println(desc.Text(text.DescKeyWriteSyncDryRun))
		cmd.Println()
	}
}

// CtxSyncAction prints a single sync action item with optional suggestion.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - index: 1-based action number.
//   - actionType: Action type label (e.g. "DEPS", "CONFIG").
//   - description: Action description.
//   - suggestion: Optional suggestion text (empty string skips).
func CtxSyncAction(cmd *cobra.Command, index int, actionType, description, suggestion string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteSyncAction), index, actionType, description))
	if suggestion != "" {
		cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteSyncSuggestion), suggestion))
	}
	cmd.Println()
}

// CtxSyncSummary prints the sync summary with item count.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - count: Number of sync items found.
//   - dryRun: If true, shows the dry-run variant.
func CtxSyncSummary(cmd *cobra.Command, count int, dryRun bool) {
	if cmd == nil {
		return
	}
	if dryRun {
		cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteSyncDryRunSummary), count))
	} else {
		cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteSyncSummary), count))
	}
}
