//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/spf13/cobra"
)

// CtxSyncInSync prints the all-clear message when context is in sync.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func CtxSyncInSync(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(assets.TextDesc(assets.TextDescKeyWriteSyncInSync))
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
	cmd.Println(assets.TextDesc(assets.TextDescKeyWriteSyncHeader))
	cmd.Println(assets.TextDesc(assets.TextDescKeyWriteSyncSeparator))
	cmd.Println()
	if dryRun {
		cmd.Println(assets.TextDesc(assets.TextDescKeyWriteSyncDryRun))
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
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteSyncAction), index, actionType, description))
	if suggestion != "" {
		cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteSyncSuggestion), suggestion))
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
		cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteSyncDryRunSummary), count))
	} else {
		cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteSyncSummary), count))
	}
}
