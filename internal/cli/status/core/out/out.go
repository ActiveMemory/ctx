//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package out

import (
	"encoding/json"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/status/core"
	"github.com/ActiveMemory/ctx/internal/cli/status/core/preview"
	"github.com/ActiveMemory/ctx/internal/cli/status/core/sort"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/format"
	"github.com/ActiveMemory/ctx/internal/write/status"
)

// PersistStatusJSON writes context status as JSON to the command output.
//
// When verbose is true, includes content previews for each file.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - ctx: Loaded context to display
//   - verbose: If true, include file content previews
//
// Returns:
//   - error: Non-nil if JSON encoding fails
func PersistStatusJSON(
	cmd *cobra.Command, ctx *entity.Context, verbose bool,
) error {
	output := core.Output{
		ContextDir:  ctx.Dir,
		TotalFiles:  len(ctx.Files),
		TotalTokens: ctx.TotalTokens,
		TotalSize:   ctx.TotalSize,
		Files:       make([]core.FileStatus, 0, len(ctx.Files)),
	}

	for _, f := range ctx.Files {
		fs := core.FileStatus{
			Name:    f.Name,
			Tokens:  f.Tokens,
			Size:    f.Size,
			IsEmpty: f.IsEmpty,
			Summary: f.Summary,
			ModTime: f.ModTime.Format(time.RFC3339),
		}
		if verbose && !f.IsEmpty {
			fs.Preview = preview.ContentPreview(string(f.Content), 5)
		}
		output.Files = append(output.Files, fs)
	}

	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	return enc.Encode(output)
}

// PersistStatusText writes context status as formatted text to the
// command output.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - ctx: Loaded context to display
//   - verbose: If true, include detailed info and content previews
//
// Returns:
//   - error: Always nil (included for interface consistency)
func PersistStatusText(
	cmd *cobra.Command, ctx *entity.Context, verbose bool,
) error {
	status.StatusHeader(cmd, ctx.Dir, len(ctx.Files), ctx.TotalTokens)

	sortedFiles := make([]entity.FileInfo, len(ctx.Files))
	copy(sortedFiles, ctx.Files)
	sort.SortFilesByPriority(sortedFiles)

	for _, f := range sortedFiles {
		fi := status.StatusFileInfo{
			Name:   f.Name,
			Tokens: f.Tokens,
			Size:   f.Size,
		}
		if f.IsEmpty {
			fi.Indicator = token.IconEmpty
			fi.Status = desc.Text(text.DescKeyWriteStatusEmpty)
		} else {
			fi.Indicator = token.IconOK
			fi.Status = f.Summary
		}
		if verbose && !f.IsEmpty {
			fi.Preview = preview.ContentPreview(string(f.Content), 3)
		}
		status.StatusFileItem(cmd, fi, verbose)
	}

	recentFiles := sort.RecentFilesSorted(ctx.Files, 3)
	entries := make([]status.StatusActivityInfo, len(recentFiles))
	for i, f := range recentFiles {
		d := time.Since(f.ModTime)
		entries[i] = status.StatusActivityInfo{
			Name: f.Name,
			Ago: format.TimeAgo(
				d.Hours(),
				int(d.Minutes()),
				f.ModTime.Format(cfgTime.OlderFormat),
			),
		}
	}
	status.StatusActivity(cmd, entries)

	return nil
}
