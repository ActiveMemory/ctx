//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package loop

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// InfoGenerated reports successful loop script generation with details.
//
// Parameters:
//   - cmd: Cobra command for output
//   - outputFile: Generated script path
//   - heading: Start heading text
//   - tool: Selected AI tool
//   - promptFile: Prompt file path
//   - maxIterations: Max iterations (0 = unlimited)
//   - completionMsg: Completion signal string
func InfoGenerated(
	cmd *cobra.Command,
	outputFile, heading, tool, promptFile string,
	maxIterations int,
	completionMsg string,
) {
	iterLine := desc.Text(text.DescKeyWriteLoopUnlimited)
	if maxIterations > 0 {
		iterLine = fmt.Sprintf(desc.Text(text.DescKeyWriteLoopMaxIterations), maxIterations)
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteLoopGeneratedBlock),
		outputFile, heading, outputFile, tool, promptFile, iterLine, completionMsg,
	))
}
