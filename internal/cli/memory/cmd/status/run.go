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

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/memory/core"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/memory"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	errMemory "github.com/ActiveMemory/ctx/internal/err/memory"
	"github.com/ActiveMemory/ctx/internal/format"
	"github.com/ActiveMemory/ctx/internal/io"
	mem "github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
	writeMem "github.com/ActiveMemory/ctx/internal/write/memory"
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

	sourcePath, discoverErr := mem.DiscoverPath(projectRoot)
	if discoverErr != nil {
		writeMem.BridgeHeader(cmd)
		writeMem.SourceNotActive(cmd)
		return errMemory.NotFound()
	}

	writeMem.BridgeHeader(cmd)
	writeMem.Source(cmd, sourcePath)
	writeMem.Mirror(cmd, memory.PathMemoryMirror)

	// Last sync time
	state, _ := mem.LoadState(contextDir)
	if state.LastSync != nil {
		ago := time.Since(*state.LastSync).Truncate(time.Minute)
		writeMem.LastSync(cmd,
			state.LastSync.Local().Format(cfgTime.DateTimeFmt),
			format.Duration(ago))
	} else {
		writeMem.LastSyncNever(cmd)
	}

	writeMem.StatusSeparator(cmd)

	// Source line count
	hasDrift := mem.HasDrift(contextDir, sourcePath)
	if sourceData, readErr := io.SafeReadFile(
		filepath.Dir(sourcePath), filepath.Base(sourcePath),
	); readErr == nil {
		writeMem.SourceLines(cmd, core.CountFileLines(sourceData), hasDrift)
	}

	// Mirror line count
	memoryDir := filepath.Join(contextDir, dir.Memory)
	if mirrorData, readErr := io.SafeReadFile(
		memoryDir, memory.Mirror,
	); readErr == nil {
		writeMem.MirrorLines(cmd, core.CountFileLines(mirrorData))
	} else {
		writeMem.MirrorNotSynced(cmd)
	}

	// Drift
	if hasDrift {
		writeMem.DriftDetected(cmd)
	} else {
		writeMem.DriftNone(cmd)
	}

	// Archives
	count := mem.ArchiveCount(contextDir)
	writeMem.Archives(cmd, count, dir.MemoryArchive)

	if hasDrift {
		// Exit code 2 for drift
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true
		os.Exit(2) //nolint:revive // spec-defined exit code
	}

	return nil
}
