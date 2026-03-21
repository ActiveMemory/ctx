//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/spf13/cobra"
)

// NoChanges prints that no changes exist since last sync.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func NoChanges(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWriteMemoryNoChanges))
}

// BridgeHeader prints the "Memory Bridge Status" heading.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func BridgeHeader(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWriteMemoryBridgeHeader))
}

// SourceNotActive prints that auto memory is not active.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func SourceNotActive(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWriteMemorySourceNotActive))
}

// Source prints the source path.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - path: absolute path to MEMORY.md.
func Source(cmd *cobra.Command, path string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteMemorySource), path))
}

// Mirror prints the mirror relative path.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - relativePath: mirror path relative to project root.
func Mirror(cmd *cobra.Command, relativePath string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteMemoryMirror), relativePath))
}

// LastSync prints the last sync timestamp with age.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - formatted: formatted timestamp string.
//   - ago: human-readable duration since sync.
func LastSync(cmd *cobra.Command, formatted, ago string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteMemoryLastSync), formatted, ago))
}

// LastSyncNever prints that no sync has occurred.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func LastSyncNever(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWriteMemoryLastSyncNever))
}

// SourceLines prints the MEMORY.md line count.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - count: number of lines.
//   - drifted: whether the source has changed since last sync.
func SourceLines(cmd *cobra.Command, count int, drifted bool) {
	if cmd == nil {
		return
	}
	if drifted {
		cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteMemorySourceLinesDrift), count))
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteMemorySourceLines), count))
}

// MirrorLines prints the mirror line count.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - count: number of lines.
func MirrorLines(cmd *cobra.Command, count int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteMemoryMirrorLines), count))
}

// MirrorNotSynced prints that the mirror has not been synced yet.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func MirrorNotSynced(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWriteMemoryMirrorNotSynced))
}

// DriftDetected prints that drift was detected.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func DriftDetected(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWriteMemoryDriftDetected))
}

// DriftNone prints that no drift was detected.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func DriftNone(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWriteMemoryDriftNone))
}

// Archives prints the archive snapshot count.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - count: number of archived snapshots.
//   - dir: archive directory name relative to .context/.
func Archives(cmd *cobra.Command, count int, dir string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteMemoryArchives), count, dir))
}
