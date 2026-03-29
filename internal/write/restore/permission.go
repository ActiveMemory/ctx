//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package restore

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// NoLocal prints the message when golden is restored with no local file.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func NoLocal(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWriteRestoreNoLocal))
}

// Match prints the message when settings already match golden.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func Match(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWriteRestoreMatch))
}

// Diff prints the permission diff block: dropped/restored
// allow and deny entries, or a note that only non-permission settings differ.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - dropped: allow permissions removed.
//   - restored: allow permissions added back.
//   - denyDropped: deny rules removed.
//   - denyRestored: deny rules added back.
func Diff(
	cmd *cobra.Command,
	dropped, restored, denyDropped, denyRestored []string,
) {
	if cmd == nil {
		return
	}
	printSection(cmd,
		desc.Text(text.DescKeyWriteRestoreDroppedHeader),
		desc.Text(text.DescKeyWriteRestoreRemoved),
		dropped)
	printSection(cmd,
		desc.Text(text.DescKeyWriteRestoreRestoredHeader),
		desc.Text(text.DescKeyWriteRestoreAdded),
		restored)
	printSection(cmd,
		desc.Text(text.DescKeyWriteRestoreDenyDroppedHeader),
		desc.Text(text.DescKeyWriteRestoreRemoved),
		denyDropped)
	printSection(cmd,
		desc.Text(text.DescKeyWriteRestoreDenyRestoredHeader),
		desc.Text(text.DescKeyWriteRestoreAdded),
		denyRestored)

	if len(dropped) == 0 && len(restored) == 0 &&
		len(denyDropped) == 0 && len(denyRestored) == 0 {
		cmd.Println(desc.Text(text.DescKeyWriteRestorePermMatch))
	}
}

// Done prints the success message after restore.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func Done(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWriteRestoreDone))
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
		cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteSnapshotUpdated), path))
	} else {
		cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteSnapshotSaved), path))
	}
}
