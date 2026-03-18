//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package loop

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
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
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteLoopGenerated), outputFile))
	cmd.Println()
	cmd.Println(heading)
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteLoopRunCmd), outputFile))
	cmd.Println()
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteLoopTool), tool))
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteLoopPrompt), promptFile))
	if maxIterations > 0 {
		cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteLoopMaxIterations), maxIterations))
	} else {
		cmd.Println(assets.TextDesc(assets.TextDescKeyWriteLoopUnlimited))
	}
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteLoopCompletion), completionMsg))
}
