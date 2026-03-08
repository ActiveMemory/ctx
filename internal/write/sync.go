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

// SyncDryRun prints the full dry-run plan block: header, source path,
// mirror path, and drift status.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - sourcePath: absolute path to MEMORY.md.
//   - mirrorPath: relative mirror path.
//   - hasDrift: whether the source has changed since last sync.
func SyncDryRun(cmd *cobra.Command, sourcePath, mirrorPath string, hasDrift bool) {
	if cmd == nil {
		return
	}
	cmd.Println(config.tplDryRun)
	io.sprintf(cmd, config.tplSource, sourcePath)
	io.sprintf(cmd, config.tplMirror, mirrorPath)
	if hasDrift {
		cmd.Println(config.tplStatusDrift)
	} else {
		cmd.Println(config.tplStatusNoDrift)
	}
}

// SyncResult prints the full sync result block: optional archive notice,
// synced confirmation, source path, line counts, and optional new content.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - sourceLabel: source label (e.g. "MEMORY.md").
//   - mirrorPath: relative mirror path.
//   - sourcePath: absolute source path for display.
//   - archivedTo: archive basename, or empty if no archive was created.
//   - sourceLines: current source line count.
//   - mirrorLines: previous mirror line count.
func SyncResult(
	cmd *cobra.Command,
	sourceLabel, mirrorPath, sourcePath, archivedTo string,
	sourceLines, mirrorLines int,
) {
	if cmd == nil {
		return
	}
	if archivedTo != "" {
		io.sprintf(cmd, config.tplArchived, archivedTo)
	}
	io.sprintf(cmd, config.tplSynced, sourceLabel, mirrorPath)
	io.sprintf(cmd, config.tplSource, sourcePath)

	line := config.tplLines
	if mirrorLines > 0 {
		line += config.tplLinesPrevious
		io.sprintf(cmd, line, sourceLines, mirrorLines)
	} else {
		io.sprintf(cmd, line, sourceLines)
	}

	if sourceLines > mirrorLines {
		io.sprintf(cmd, config.tplNewContent, sourceLines-mirrorLines)
	}
}

// ErrAutoMemoryNotActive prints an informational stderr message when
// auto memory discovery fails.
//
// Parameters:
//   - cmd: Cobra command whose stderr stream receives the message. Nil is a no-op.
//   - cause: the discovery error to display.
func ErrAutoMemoryNotActive(cmd *cobra.Command, cause error) {
	if cmd == nil {
		return
	}
	cmd.PrintErrln("Auto memory not active:", cause)
}
