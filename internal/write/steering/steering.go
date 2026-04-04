//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

import (
	"fmt"

	"github.com/spf13/cobra"
)

// User-facing messages for steering commands.
const (
	// msgNoFiles is shown when no steering files exist.
	msgNoFiles = "No steering files found."
	// msgNoMatch is shown when no files match the prompt.
	msgNoMatch = "No steering files match the given prompt."
)

// Created prints confirmation that a steering file was created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: Path to the created file
func Created(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf("Created %s", path))
}

// Skipped prints that a steering file was skipped because it exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: Path to the existing file
func Skipped(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf("Skipped %s (already exists)", path))
}

// InitSummary prints the summary after steering init.
//
// Parameters:
//   - cmd: Cobra command for output
//   - created: Number of files created
//   - skipped: Number of files skipped
func InitSummary(cmd *cobra.Command, created, skipped int) {
	cmd.Println(fmt.Sprintf("\n%d created, %d skipped", created, skipped))
}

// NoFilesFound prints a message indicating no steering files exist.
//
// Parameters:
//   - cmd: Cobra command for output
func NoFilesFound(cmd *cobra.Command) {
	cmd.Println(msgNoFiles)
}

// FileEntry prints a single steering file entry with metadata.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: File name
//   - inclusion: Inclusion mode
//   - priority: Priority value
//   - tools: Comma-separated tool list or "all"
func FileEntry(
	cmd *cobra.Command, name, inclusion string,
	priority int, tools string,
) {
	cmd.Println(fmt.Sprintf("%-20s  inclusion=%-7s  priority=%-3d  tools=%s",
		name, inclusion, priority, tools))
}

// FileCount prints the total steering file count.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: Number of steering files
func FileCount(cmd *cobra.Command, count int) {
	cmd.Println(fmt.Sprintf("\n%d steering file(s)", count))
}

// NoFilesMatch prints a message indicating no files match the prompt.
//
// Parameters:
//   - cmd: Cobra command for output
func NoFilesMatch(cmd *cobra.Command) {
	cmd.Println(msgNoMatch)
}

// PreviewHeader prints the header for steering preview output.
//
// Parameters:
//   - cmd: Cobra command for output
//   - prompt: The prompt being matched against
func PreviewHeader(cmd *cobra.Command, prompt string) {
	cmd.Println(fmt.Sprintf("Steering files matching prompt %q:", prompt))
	cmd.Println()
}

// PreviewEntry prints a single preview match entry.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: File name
//   - inclusion: Inclusion mode
//   - priority: Priority value
//   - tools: Comma-separated tool list or "all"
func PreviewEntry(
	cmd *cobra.Command, name, inclusion string,
	priority int, tools string,
) {
	cmd.Println(fmt.Sprintf("  %-20s  inclusion=%-7s  priority=%-3d  tools=%s",
		name, inclusion, priority, tools))
}

// PreviewCount prints the count of files that would be included.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: Number of matched files
func PreviewCount(cmd *cobra.Command, count int) {
	cmd.Println(fmt.Sprintf("\n%d file(s) would be included", count))
}

// SyncWritten prints that a file was written during sync.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Name of the written file
func SyncWritten(cmd *cobra.Command, name string) {
	cmd.Println(fmt.Sprintf("Written: %s", name))
}

// SyncSkipped prints that a file was skipped during sync.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Name of the skipped file
func SyncSkipped(cmd *cobra.Command, name string) {
	cmd.Println(fmt.Sprintf("Skipped: %s", name))
}

// SyncError prints a sync error.
//
// Parameters:
//   - cmd: Cobra command for output
//   - errMsg: The error message
func SyncError(cmd *cobra.Command, errMsg string) {
	cmd.Println(fmt.Sprintf("Error: %s", errMsg))
}

// SyncSummary prints the sync summary with counts.
//
// Parameters:
//   - cmd: Cobra command for output
//   - written: Number of files written
//   - skipped: Number of files skipped
//   - errors: Number of errors
func SyncSummary(cmd *cobra.Command, written, skipped, errors int) {
	cmd.Println(fmt.Sprintf("\n%d written, %d skipped, %d errors",
		written, skipped, errors))
}
