//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"path/filepath"

	memory2 "github.com/ActiveMemory/ctx/internal/config/memory"
	"github.com/spf13/cobra"

	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write"
)

// Run discovers MEMORY.md, mirrors it into .context/memory/, and
// updates the sync state. In dry-run mode it reports what would happen
// without writing any files.
//
// Parameters:
//   - cmd: Cobra command for output routing.
//   - dryRun: when true, report the plan without writing.
//
// Returns:
//   - error: on discovery failure, sync failure, or state persistence failure.
func Run(cmd *cobra.Command, dryRun bool) error {
	contextDir := rc.ContextDir()
	projectRoot := filepath.Dir(contextDir)

	sourcePath, discoverErr := memory.DiscoverMemoryPath(projectRoot)
	if discoverErr != nil {
		write.ErrAutoMemoryNotActive(cmd, discoverErr)
		return ctxerr.MemoryNotFound()
	}

	if dryRun {
		write.SyncDryRun(cmd, sourcePath, memory2.PathMemoryMirror,
			memory.HasDrift(contextDir, sourcePath))
		return nil
	}

	result, syncErr := memory.Sync(contextDir, sourcePath)
	if syncErr != nil {
		return ctxerr.SyncFailed(syncErr)
	}

	write.SyncResult(cmd,
		memory2.MemorySource, memory2.PathMemoryMirror,
		result.SourcePath, filepath.Base(result.ArchivedTo),
		result.SourceLines, result.MirrorLines,
	)

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
