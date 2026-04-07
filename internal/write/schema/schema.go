//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

import (
	"fmt"

	"github.com/spf13/cobra"

	cfgSchema "github.com/ActiveMemory/ctx/internal/config/schema"
)

// NoDirs prints that no session directories were found.
//
// Parameters:
//   - cmd: Cobra command for output
func NoDirs(cmd *cobra.Command) {
	cmd.Println(cfgSchema.MsgNoDirs)
}

// NoFiles prints that no session files were found.
//
// Parameters:
//   - cmd: Cobra command for output
func NoFiles(cmd *cobra.Command) {
	cmd.Println(cfgSchema.MsgNoFiles)
}

// Clean prints that no drift was detected.
//
// Parameters:
//   - cmd: Cobra command for output
//   - files: number of files scanned
//   - lines: number of lines scanned
func Clean(cmd *cobra.Command, files, lines int) {
	cmd.Println(fmt.Sprintf(
		cfgSchema.FmtClean, files, lines))
}

// DriftSummary prints a drift summary to stderr.
//
// Parameters:
//   - cmd: Cobra command for output
//   - summary: pre-formatted drift summary string
func DriftSummary(cmd *cobra.Command, summary string) {
	cmd.PrintErrln()
	cmd.PrintErrln(summary)
}

// DumpLine prints a single line of schema dump output.
//
// Parameters:
//   - cmd: Cobra command for output
//   - line: text to print
func DumpLine(cmd *cobra.Command, line string) {
	cmd.Println(line)
}

// DumpBlank prints a blank line in schema dump output.
//
// Parameters:
//   - cmd: Cobra command for output
func DumpBlank(cmd *cobra.Command) {
	cmd.Println()
}
