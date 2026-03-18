//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package loop

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/embed"
	"github.com/spf13/cobra"
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
	iterLine := assets.TextDesc(embed.TextDescKeyWriteLoopUnlimited)
	if maxIterations > 0 {
		iterLine = fmt.Sprintf(assets.TextDesc(embed.TextDescKeyWriteLoopMaxIterations), maxIterations)
	}
	cmd.Println(fmt.Sprintf(
		assets.TextDesc(embed.TextDescKeyWriteLoopGeneratedBlock),
		outputFile, heading, outputFile, tool, promptFile, iterLine, completionMsg,
	))
}
