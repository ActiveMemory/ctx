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

// RestoreNoLocal prints the message when golden is restored with no local file.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func RestoreNoLocal(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(config.tplRestoreNoLocal)
}

// RestoreMatch prints the message when settings already match golden.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func RestoreMatch(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(config.tplRestoreMatch)
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
	printSection(cmd, config.tplRestoreDroppedHeader, config.tplRestoreRemoved, dropped)
	printSection(cmd, config.tplRestoreRestoredHeader, config.tplRestoreAdded, restored)
	printSection(cmd, config.tplRestoreDenyDroppedHeader, config.tplRestoreRemoved, denyDropped)
	printSection(cmd, config.tplRestoreDenyRestoredHeader, config.tplRestoreAdded, denyRestored)

	if len(dropped) == 0 && len(restored) == 0 &&
		len(denyDropped) == 0 && len(denyRestored) == 0 {
		cmd.Println(config.tplRestorePermMatch)
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
	cmd.Println(config.tplRestoreDone)
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
		io.sprintf(cmd, config.tplSnapshotUpdated, path)
	} else {
		io.sprintf(cmd, config.tplSnapshotSaved, path)
	}
}

// printSection prints a header and list items if the list is non-empty.
func printSection(cmd *cobra.Command, headerTpl, itemTpl string, items []string) {
	if len(items) == 0 {
		return
	}
	io.sprintf(cmd, headerTpl, len(items))
	for _, item := range items {
		io.sprintf(cmd, itemTpl, item)
	}
}
