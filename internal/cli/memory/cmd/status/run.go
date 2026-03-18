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
	ctxerr "github.com/ActiveMemory/ctx/internal/err/memory"
	"github.com/ActiveMemory/ctx/internal/io"
	memory2 "github.com/ActiveMemory/ctx/internal/write/memory"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/memory/core"
	mem "github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
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
		memory2.BridgeHeader(cmd)
		memory2.SourceNotActive(cmd)
		return ctxerr.NotFound()
	}

	memory2.BridgeHeader(cmd)
	memory2.Source(cmd, sourcePath)
	memory2.Mirror(cmd, memory.PathMemoryMirror)

	// Last sync time
	state, _ := mem.LoadState(contextDir)
	if state.LastSync != nil {
		ago := time.Since(*state.LastSync).Truncate(time.Minute)
		memory2.LastSync(cmd,
			state.LastSync.Local().Format(time2.DateTimeFormat),
			core.FormatDuration(ago))
	} else {
		memory2.LastSyncNever(cmd)
	}

	cmd.Println()

	// Source line count
	hasDrift := mem.HasDrift(contextDir, sourcePath)
	if sourceData, readErr := io.SafeReadFile(
		filepath.Dir(sourcePath), filepath.Base(sourcePath),
	); readErr == nil {
		memory2.SourceLines(cmd, core.CountFileLines(sourceData), hasDrift)
	}

	// Mirror line count
	memoryDir := filepath.Join(contextDir, dir.Memory)
	if mirrorData, readErr := io.SafeReadFile(
		memoryDir, memory.MemoryMirror,
	); readErr == nil {
		memory2.MirrorLines(cmd, core.CountFileLines(mirrorData))
	} else {
		memory2.MirrorNotSynced(cmd)
	}

	// Drift
	if hasDrift {
		memory2.DriftDetected(cmd)
	} else {
		memory2.DriftNone(cmd)
	}

	// Archives
	count := mem.ArchiveCount(contextDir)
	memory2.Archives(cmd, count, dir.MemoryArchive)

	if hasDrift {
		// Exit code 2 for drift
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true
		os.Exit(2) //nolint:revive // spec-defined exit code
	}

	return nil
}
