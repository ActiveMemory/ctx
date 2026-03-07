//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write"
)

// runSync discovers MEMORY.md, mirrors it into .context/memory/, and
// updates the sync state. In dry-run mode it reports what would happen
// without writing any files.
//
// Parameters:
//   - cmd: Cobra command for output routing.
//   - dryRun: when true, report the plan without writing.
//
// Returns:
//   - error: on discovery failure, sync failure, or state persistence failure.
func runSync(cmd *cobra.Command, dryRun bool) error {
	contextDir := rc.ContextDir()
	projectRoot := filepath.Dir(contextDir)

	sourcePath, discoverErr := memory.DiscoverMemoryPath(projectRoot)
	if discoverErr != nil {
		write.ErrAutoMemoryNotActive(cmd, discoverErr)
		return ctxerr.MemoryNotFound()
	}

	if dryRun {
		write.DryRun(cmd)
		write.Source(cmd, sourcePath)
		write.Mirror(cmd, config.PathMemoryMirror)
		if memory.HasDrift(contextDir, sourcePath) {
			write.StatusDrift(cmd)
		} else {
			write.StatusNoDrift(cmd)
		}
		return nil
	}

	result, syncErr := memory.Sync(contextDir, sourcePath)
	if syncErr != nil {
		return ctxerr.SyncFailed(syncErr)
	}

	if result.ArchivedTo != "" {
		write.Archived(cmd, filepath.Base(result.ArchivedTo))
	}

	write.Synced(cmd, config.FileMemorySource, config.PathMemoryMirror)
	write.Source(cmd, result.SourcePath)
	write.Lines(cmd, result.SourceLines, result.MirrorLines)

	if result.SourceLines > result.MirrorLines {
		write.NewContent(cmd, result.SourceLines-result.MirrorLines)
	}

	// Update sync state
	state, loadErr := memory.LoadState(contextDir)
	if loadErr != nil {
		return ctxerr.LoadState(loadErr)
	}
	state.MarkSynced()
	if saveErr := memory.SaveState(contextDir, state); saveErr != nil {
		return ctxerr.SaveState(saveErr)
	}

	return nil
}
