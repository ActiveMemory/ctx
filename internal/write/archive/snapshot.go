//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package archive

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/spf13/cobra"
)

// SnapshotSaved prints the result of a successful task snapshot.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - snapshotPath: path to the created snapshot file.
func SnapshotSaved(cmd *cobra.Command, snapshotPath string) {
	if cmd == nil {
		return
	}
	cmd.Println(
		fmt.Sprintf(
			desc.Text(text.DescKeyTaskSnapshotSaved), snapshotPath))
}

// SnapshotContent builds the snapshot file content with header and body.
//
// Parameters:
//   - name: snapshot name.
//   - created: RFC3339 formatted creation timestamp.
//   - separator: the separator string.
//   - nl: newline string.
//   - body: the original TASKS.md content.
//
// Returns:
//   - string: formatted snapshot content.
func SnapshotContent(name, created, separator, nl, body string) string {
	return fmt.Sprintf(
		desc.Text(text.DescKeyTaskSnapshotHeaderFormat)+
			nl+nl+
			desc.Text(text.DescKeyTaskSnapshotCreatedFormat)+
			nl+nl+separator+nl+nl+"%s",
		name, created, body,
	)
}
