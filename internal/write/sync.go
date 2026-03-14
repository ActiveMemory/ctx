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
	cmd.Println(config.TplDryRun)
	cmd.Println(fmt.Sprintf(config.TplSource, sourcePath))
	cmd.Println(fmt.Sprintf(config.TplMirror, mirrorPath))
	if hasDrift {
		cmd.Println(config.TplStatusDrift)
	} else {
		cmd.Println(config.TplStatusNoDrift)
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
		cmd.Println(fmt.Sprintf(config.TplArchived, archivedTo))
	}
	cmd.Println(fmt.Sprintf(config.TplSynced, sourceLabel, mirrorPath))
	cmd.Println(fmt.Sprintf(config.TplSource, sourcePath))

	line := config.TplLines
	if mirrorLines > 0 {
		line += config.TplLinesPrevious
		cmd.Println(fmt.Sprintf(line, sourceLines, mirrorLines))
	} else {
		cmd.Println(fmt.Sprintf(line, sourceLines))
	}

	if sourceLines > mirrorLines {
		cmd.Println(fmt.Sprintf(config.TplNewContent, sourceLines-mirrorLines))
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
