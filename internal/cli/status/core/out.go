//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"encoding/json"
	"time"

	ctxtime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/write"
)

// OutputStatusJSON writes context status as JSON to the command output.
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
func OutputStatusJSON(
	cmd *cobra.Command, ctx *context.Context, verbose bool,
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
			fs.Preview = ContentPreview(string(f.Content), 5)
		}
		output.Files = append(output.Files, fs)
	}

	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	return enc.Encode(output)
}

// OutputStatusText writes context status as formatted text to the command output.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - ctx: Loaded context to display
//   - verbose: If true, include detailed info and content previews
//
// Returns:
//   - error: Always nil (included for interface consistency)
func OutputStatusText(
	cmd *cobra.Command, ctx *context.Context, verbose bool,
) error {
	write.StatusHeader(cmd, ctx.Dir, len(ctx.Files), ctx.TotalTokens)

	sortedFiles := make([]context.FileInfo, len(ctx.Files))
	copy(sortedFiles, ctx.Files)
	SortFilesByPriority(sortedFiles)

	for _, f := range sortedFiles {
		fi := write.StatusFileInfo{
			Name:   f.Name,
			Tokens: f.Tokens,
			Size:   f.Size,
		}
		if f.IsEmpty {
			fi.Indicator = "\u25cb"
			fi.Status = "empty"
		} else {
			fi.Indicator = "\u2713"
			fi.Status = f.Summary
		}
		if verbose && !f.IsEmpty {
			fi.Preview = ContentPreview(string(f.Content), 3)
		}
		write.StatusFileItem(cmd, fi, verbose)
	}

	recentFiles := GetRecentFiles(ctx.Files, 3)
	entries := make([]write.StatusActivityInfo, len(recentFiles))
	for i, f := range recentFiles {
		d := time.Since(f.ModTime)
		entries[i] = write.StatusActivityInfo{
			Name: f.Name,
			Ago:  write.FormatTimeAgo(d.Hours(), int(d.Minutes()), f.ModTime.Format(ctxtime.OlderFormat)),
		}
	}
	write.StatusActivity(cmd, entries)

	return nil
}
