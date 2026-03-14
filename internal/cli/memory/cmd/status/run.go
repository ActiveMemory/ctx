//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/memory"
	time2 "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/memory/core"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	mem "github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write"
)

// Run prints memory bridge status including source location,
// last sync time, line counts, drift indicator, and archive count.
//
// Parameters:
//   - cmd: Cobra command for output routing.
//
// Returns:
//   - error: on discovery failure.
func Run(cmd *cobra.Command) error {
	contextDir := rc.ContextDir()
	projectRoot := filepath.Dir(contextDir)

	sourcePath, discoverErr := mem.DiscoverMemoryPath(projectRoot)
	if discoverErr != nil {
		write.MemoryBridgeHeader(cmd)
		write.MemorySourceNotActive(cmd)
		return ctxerr.MemoryNotFound()
	}

	write.MemoryBridgeHeader(cmd)
	write.MemorySource(cmd, sourcePath)
	write.MemoryMirror(cmd, memory.PathMemoryMirror)

	// Last sync time
	state, _ := mem.LoadState(contextDir)
	if state.LastSync != nil {
		ago := time.Since(*state.LastSync).Truncate(time.Minute)
		write.MemoryLastSync(cmd,
			state.LastSync.Local().Format(time2.DateTimeFormat),
			core.FormatDuration(ago))
	} else {
		write.MemoryLastSyncNever(cmd)
	}

	cmd.Println()

	// Source line count
	hasDrift := mem.HasDrift(contextDir, sourcePath)
	if sourceData, readErr := io.SafeReadFile(
		filepath.Dir(sourcePath), filepath.Base(sourcePath),
	); readErr == nil {
		write.MemorySourceLines(cmd, core.CountFileLines(sourceData), hasDrift)
	}

	// Mirror line count
	memoryDir := filepath.Join(contextDir, dir.Memory)
	if mirrorData, readErr := io.SafeReadFile(
		memoryDir, memory.MemoryMirror,
	); readErr == nil {
		write.MemoryMirrorLines(cmd, core.CountFileLines(mirrorData))
	} else {
		write.MemoryMirrorNotSynced(cmd)
	}

	// Drift
	if hasDrift {
		write.MemoryDriftDetected(cmd)
	} else {
		write.MemoryDriftNone(cmd)
	}

	// Archives
	count := mem.ArchiveCount(contextDir)
	write.MemoryArchives(cmd, count, dir.MemoryArchive)

	if hasDrift {
		// Exit code 2 for drift
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true
		os.Exit(2) //nolint:revive // spec-defined exit code
	}

	return nil
}
