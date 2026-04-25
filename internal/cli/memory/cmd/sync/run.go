//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/memory/core/resolve"
	cfgMem "github.com/ActiveMemory/ctx/internal/config/memory"
	errMem "github.com/ActiveMemory/ctx/internal/err/memory"
	errState "github.com/ActiveMemory/ctx/internal/err/state"
	"github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/write/sync"
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
	contextDir, projectRoot, err := resolve.ContextAndRoot(cmd)
	if err != nil {
		return err
	}

	sourcePath, discoverErr := resolve.DiscoverSource(cmd, projectRoot)
	if discoverErr != nil {
		return discoverErr
	}

	if dryRun {
		sync.DryRun(cmd, sourcePath, cfgMem.PathMemoryMirror,
			memory.HasDrift(contextDir, sourcePath))
		return nil
	}

	result, syncErr := memory.Sync(contextDir, sourcePath)
	if syncErr != nil {
		return errMem.Sync(syncErr)
	}

	sync.Result(cmd,
		cfgMem.Source, cfgMem.PathMemoryMirror,
		result.SourcePath, filepath.Base(result.ArchivedTo),
		result.SourceLines, result.MirrorLines,
	)

	// Update sync state
	state, loadErr := memory.LoadState(contextDir)
	if loadErr != nil {
		return errState.Load(loadErr)
	}
	state.MarkSynced()
	if saveErr := memory.SaveState(contextDir, state); saveErr != nil {
		return errState.Save(saveErr)
	}

	return nil
}
