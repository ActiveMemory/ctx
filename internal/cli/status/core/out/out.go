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
	"github.com/ActiveMemory/ctx/internal/cli/status/core/preview"
	"github.com/ActiveMemory/ctx/internal/cli/status/core/sort"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgFmt "github.com/ActiveMemory/ctx/internal/config/format"
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
	output := Output{
		ContextDir:  ctx.Dir,
		TotalFiles:  len(ctx.Files),
		TotalTokens: ctx.TotalTokens,
		TotalSize:   ctx.TotalSize,
		Files:       make([]FileStatus, 0, len(ctx.Files)),
	}

	for _, f := range ctx.Files {
		fs := FileStatus{
			Name:    f.Name,
			Tokens:  f.Tokens,
			Size:    f.Size,
			IsEmpty: f.IsEmpty,
			Summary: f.Summary,
			ModTime: f.ModTime.Format(time.RFC3339),
		}
		if verbose && !f.IsEmpty {
			fs.Preview = preview.Content(string(f.Content), cfgFmt.PreviewLines)
		}
		output.Files = append(output.Files, fs)
	}

	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", token.Indent2)
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
	status.Header(cmd, ctx.Dir, len(ctx.Files), ctx.TotalTokens)

	sortedFiles := make([]entity.FileInfo, len(ctx.Files))
	copy(sortedFiles, ctx.Files)
	sort.FilesByPriority(sortedFiles)

	for _, f := range sortedFiles {
		fi := status.FileInfo{
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
			fi.Preview = preview.Content(string(f.Content), cfgFmt.StatusPreviewLines)
		}
		status.FileItem(cmd, fi, verbose)
	}

	recentFiles := sort.RecentFiles(ctx.Files, cfgFmt.StatusRecentFiles)
	entries := make([]status.ActivityInfo, len(recentFiles))
	for i, f := range recentFiles {
		d := time.Since(f.ModTime)
		entries[i] = status.ActivityInfo{
			Name: f.Name,
			Ago: format.TimeAgo(
				d.Hours(),
				int(d.Minutes()),
				f.ModTime.Format(cfgTime.OlderFormat),
			),
		}
	}
	status.Activity(cmd, entries)

	return nil
}
