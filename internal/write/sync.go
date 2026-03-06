//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"github.com/spf13/cobra"
)

// DryRun prints the dry-run header to stdout.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func DryRun(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(tplDryRun)
}

// Source prints an indented source path line to stdout.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - path: the source file path to display.
func Source(cmd *cobra.Command, path string) {
	if cmd == nil {
		return
	}
	sprintf(cmd, tplSource, path)
}

// Mirror prints an indented mirror path line to stdout.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - relativePath: the mirror path relative to the project root.
func Mirror(cmd *cobra.Command, relativePath string) {
	if cmd == nil {
		return
	}
	sprintf(cmd, tplMirror, relativePath)
}

// StatusDrift prints that drift was detected.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func StatusDrift(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(tplStatusDrift)
}

// StatusNoDrift prints that no drift was detected.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func StatusNoDrift(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(tplStatusNoDrift)
}

// Archived prints that a previous file was archived.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - filename: the archive filename (basename, not full path).
func Archived(cmd *cobra.Command, filename string) {
	if cmd == nil {
		return
	}
	sprintf(cmd, tplArchived, filename)
}

// Synced prints that a sync completed successfully.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - source: label for the source (e.g. "MEMORY.md").
//   - destination: relative path to the destination.
func Synced(cmd *cobra.Command, source, destination string) {
	if cmd == nil {
		return
	}
	sprintf(cmd, tplSynced, source, destination)
}

// Lines prints a line count, optionally including the previous count.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - count: current line count.
//   - previous: previous line count. Zero omits the "(was N)" suffix.
func Lines(cmd *cobra.Command, count, previous int) {
	if cmd == nil {
		return
	}
	line := tplLines
	if previous > 0 {
		line += tplLinesPrevious
		sprintf(cmd, line, count, previous)
		return
	}
	sprintf(cmd, line, count)
}

// NewContent prints how many new lines appeared since the last sync.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - count: number of new lines.
func NewContent(cmd *cobra.Command, count int) {
	if cmd == nil {
		return
	}
	sprintf(cmd, tplNewContent, count)
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
