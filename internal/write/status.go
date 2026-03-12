//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"fmt"
	"github.com/ActiveMemory/ctx/internal/write/config"
	"github.com/spf13/cobra"
)

// StatusFileInfo holds prepared data for a single file in status output.
type StatusFileInfo struct {
	Indicator string
	Name      string
	Status    string
	Tokens    int
	Size      int64
	Preview   []string
}

// StatusActivityInfo holds prepared data for a recent activity entry.
type StatusActivityInfo struct {
	Name string
	Ago  string
}

// StatusHeader prints the status heading and summary block.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - dir: Context directory path.
//   - fileCount: Number of context files.
//   - totalTokens: Estimated total token count.
func StatusHeader(cmd *cobra.Command, dir string, fileCount, totalTokens int) {
	if cmd == nil {
		return
	}
	cmd.Println(config.TplStatusTitle)
	cmd.Println(config.TplStatusSeparator)
	cmd.Println()
	cmd.Println(fmt.Sprintf(config.TplStatusDir, dir))
	cmd.Println(fmt.Sprintf(config.TplStatusFiles, fileCount))
	cmd.Println(fmt.Sprintf(config.TplStatusTokens, FormatNumber(totalTokens)))
	cmd.Println()
	cmd.Println(config.TplStatusFilesHeader)
}

// StatusFileItem prints a single file entry in the status list.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - f: Prepared file info.
//   - verbose: If true, include tokens, size, and preview.
func StatusFileItem(cmd *cobra.Command, f StatusFileInfo, verbose bool) {
	if cmd == nil {
		return
	}
	if verbose {
		cmd.Println(fmt.Sprintf(config.TplStatusFileVerbose,
			f.Indicator, f.Name, f.Status,
			FormatNumber(f.Tokens), FormatBytes(f.Size)))
		for _, line := range f.Preview {
			cmd.Println(fmt.Sprintf(config.TplStatusPreviewLine, line))
		}
	} else {
		cmd.Println(fmt.Sprintf(config.TplStatusFileCompact, f.Indicator, f.Name, f.Status))
	}
}

// StatusActivity prints the recent activity section.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - entries: Recent activity entries.
func StatusActivity(cmd *cobra.Command, entries []StatusActivityInfo) {
	if cmd == nil {
		return
	}
	cmd.Println()
	cmd.Println(config.TplStatusActivityHeader)
	for _, e := range entries {
		cmd.Println(fmt.Sprintf(config.TplStatusActivityItem, e.Name, e.Ago))
	}
}
