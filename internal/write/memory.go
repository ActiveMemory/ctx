//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"fmt"
	"github.com/ActiveMemory/ctx/internal/write/config"
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
	cmd.Println(config.TplMemoryNoChanges)
}

// MemoryBridgeHeader prints the "Memory Bridge Status" heading.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func MemoryBridgeHeader(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(config.TplMemoryBridgeHeader)
}

// MemorySourceNotActive prints that auto memory is not active.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func MemorySourceNotActive(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(config.TplMemorySourceNotActive)
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
	cmd.Println(fmt.Sprintf(config.TplMemorySource, path))
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
	cmd.Println(fmt.Sprintf(config.TplMemoryMirror, relativePath))
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
	cmd.Println(fmt.Sprintf(config.TplMemoryLastSync, formatted, ago))
}

// MemoryLastSyncNever prints that no sync has occurred.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func MemoryLastSyncNever(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(config.TplMemoryLastSyncNever)
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
		cmd.Println(fmt.Sprintf(config.TplMemorySourceLinesDrift, count))
		return
	}
	cmd.Println(fmt.Sprintf(config.TplMemorySourceLines, count))
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
	cmd.Println(fmt.Sprintf(config.TplMemoryMirrorLines, count))
}

// MemoryMirrorNotSynced prints that the mirror has not been synced yet.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func MemoryMirrorNotSynced(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(config.TplMemoryMirrorNotSynced)
}

// MemoryDriftDetected prints that drift was detected.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func MemoryDriftDetected(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(config.TplMemoryDriftDetected)
}

// MemoryDriftNone prints that no drift was detected.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func MemoryDriftNone(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(config.TplMemoryDriftNone)
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
	cmd.Println(fmt.Sprintf(config.TplMemoryArchives, count, dir))
}
