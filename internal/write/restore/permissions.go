//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package restore

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/spf13/cobra"
)

// RestoreNoLocal prints the message when golden is restored with no local file.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func RestoreNoLocal(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.TextDesc(text.TextDescKeyWriteRestoreNoLocal))
}

// RestoreMatch prints the message when settings already match golden.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func RestoreMatch(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.TextDesc(text.TextDescKeyWriteRestoreMatch))
}

// RestoreDiff prints the permission diff block: dropped/restored
// allow and deny entries, or a note that only non-permission settings differ.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - dropped: allow permissions removed.
//   - restored: allow permissions added back.
//   - denyDropped: deny rules removed.
//   - denyRestored: deny rules added back.
func RestoreDiff(
	cmd *cobra.Command,
	dropped, restored, denyDropped, denyRestored []string,
) {
	if cmd == nil {
		return
	}
	printSection(cmd, desc.TextDesc(text.TextDescKeyWriteRestoreDroppedHeader), desc.TextDesc(text.TextDescKeyWriteRestoreRemoved), dropped)
	printSection(cmd, desc.TextDesc(text.TextDescKeyWriteRestoreRestoredHeader), desc.TextDesc(text.TextDescKeyWriteRestoreAdded), restored)
	printSection(cmd, desc.TextDesc(text.TextDescKeyWriteRestoreDenyDroppedHeader), desc.TextDesc(text.TextDescKeyWriteRestoreRemoved), denyDropped)
	printSection(cmd, desc.TextDesc(text.TextDescKeyWriteRestoreDenyRestoredHeader), desc.TextDesc(text.TextDescKeyWriteRestoreAdded), denyRestored)

	if len(dropped) == 0 && len(restored) == 0 &&
		len(denyDropped) == 0 && len(denyRestored) == 0 {
		cmd.Println(desc.TextDesc(text.TextDescKeyWriteRestorePermMatch))
	}
}

// RestoreDone prints the success message after restore.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func RestoreDone(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.TextDesc(text.TextDescKeyWriteRestoreDone))
}

// SnapshotDone prints the golden image save/update confirmation.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - updated: true if golden already existed (update vs save).
//   - path: golden file path.
func SnapshotDone(cmd *cobra.Command, updated bool, path string) {
	if cmd == nil {
		return
	}
	if updated {
		cmd.Println(fmt.Sprintf(desc.TextDesc(text.TextDescKeyWriteSnapshotUpdated), path))
	} else {
		cmd.Println(fmt.Sprintf(desc.TextDesc(text.TextDescKeyWriteSnapshotSaved), path))
	}
}

// printSection prints a header and list items if the list is non-empty.
func printSection(cmd *cobra.Command, headerTpl, itemTpl string, items []string) {
	if len(items) == 0 {
		return
	}
	cmd.Println(fmt.Sprintf(headerTpl, len(items)))
	for _, item := range items {
		cmd.Println(fmt.Sprintf(itemTpl, item))
	}
}
