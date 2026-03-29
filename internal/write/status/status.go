//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/format"
)

// Header prints the status heading and summary block.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - dir: Context directory path.
//   - fileCount: Number of context files.
//   - totalTokens: Estimated total token count.
func Header(cmd *cobra.Command, dir string, fileCount, totalTokens int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteStatusHeaderBlock),
		dir, fileCount, format.Number(totalTokens),
	))
}

// FileItem prints a single file entry in the status list.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - f: Prepared file info.
//   - verbose: If true, include tokens, size, and preview.
func FileItem(cmd *cobra.Command, f FileInfo, verbose bool) {
	if cmd == nil {
		return
	}
	if verbose {
		cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteStatusFileVerbose),
			f.Indicator, f.Name, f.Status,
			format.Number(f.Tokens), format.Bytes(f.Size)))
		for _, line := range f.Preview {
			cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteStatusPreviewLine), line))
		}
	} else {
		cmd.Println(fmt.Sprintf(
			desc.Text(text.DescKeyWriteStatusFileCompact),
			f.Indicator, f.Name, f.Status))
	}
}

// Activity prints the recent activity section.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - entries: Recent activity entries.
func Activity(cmd *cobra.Command, entries []ActivityInfo) {
	if cmd == nil {
		return
	}
	cmd.Println()
	cmd.Println(desc.Text(text.DescKeyWriteStatusActivityHeader))
	for _, e := range entries {
		cmd.Println(fmt.Sprintf(
			desc.Text(text.DescKeyWriteStatusActivityItem),
			e.Name, e.Ago))
	}
}
