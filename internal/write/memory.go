//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"github.com/ActiveMemory/ctx/internal/write/config"
	"github.com/ActiveMemory/ctx/internal/write/io"
	"github.com/spf13/cobra"
)

// MemoryNoChanges prints that no changes exist since last sync.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func MemoryNoChanges(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(config.tplMemoryNoChanges)
}

// MemoryBridgeHeader prints the "Memory Bridge Status" heading.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func MemoryBridgeHeader(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(config.tplMemoryBridgeHeader)
}

// MemorySourceNotActive prints that auto memory is not active.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func MemorySourceNotActive(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(config.tplMemorySourceNotActive)
}

// MemorySource prints the source path.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - path: absolute path to MEMORY.md.
func MemorySource(cmd *cobra.Command, path string) {
	if cmd == nil {
		return
	}
	io.sprintf(cmd, config.tplMemorySource, path)
}

// MemoryMirror prints the mirror relative path.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - relativePath: mirror path relative to project root.
func MemoryMirror(cmd *cobra.Command, relativePath string) {
	if cmd == nil {
		return
	}
	io.sprintf(cmd, config.tplMemoryMirror, relativePath)
}

// MemoryLastSync prints the last sync timestamp with age.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - formatted: formatted timestamp string.
//   - ago: human-readable duration since sync.
func MemoryLastSync(cmd *cobra.Command, formatted, ago string) {
	if cmd == nil {
		return
	}
	io.sprintf(cmd, config.tplMemoryLastSync, formatted, ago)
}

// MemoryLastSyncNever prints that no sync has occurred.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func MemoryLastSyncNever(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(config.tplMemoryLastSyncNever)
}

// MemorySourceLines prints the MEMORY.md line count.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - count: number of lines.
//   - drifted: whether the source has changed since last sync.
func MemorySourceLines(cmd *cobra.Command, count int, drifted bool) {
	if cmd == nil {
		return
	}
	if drifted {
		io.sprintf(cmd, config.tplMemorySourceLinesDrift, count)
		return
	}
	io.sprintf(cmd, config.tplMemorySourceLines, count)
}

// MemoryMirrorLines prints the mirror line count.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - count: number of lines.
func MemoryMirrorLines(cmd *cobra.Command, count int) {
	if cmd == nil {
		return
	}
	io.sprintf(cmd, config.tplMemoryMirrorLines, count)
}

// MemoryMirrorNotSynced prints that the mirror has not been synced yet.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func MemoryMirrorNotSynced(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(config.tplMemoryMirrorNotSynced)
}

// MemoryDriftDetected prints that drift was detected.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func MemoryDriftDetected(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(config.tplMemoryDriftDetected)
}

// MemoryDriftNone prints that no drift was detected.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func MemoryDriftNone(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(config.tplMemoryDriftNone)
}

// MemoryArchives prints the archive snapshot count.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - count: number of archived snapshots.
//   - dir: archive directory name relative to .context/.
func MemoryArchives(cmd *cobra.Command, count int, dir string) {
	if cmd == nil {
		return
	}
	io.sprintf(cmd, config.tplMemoryArchives, count, dir)
}
