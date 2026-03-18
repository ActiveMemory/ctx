//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package obsidian

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/embed"
	"github.com/spf13/cobra"
)

// InfoGenerated reports successful Obsidian vault generation.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: Number of entries generated
//   - output: Output directory path
func InfoGenerated(cmd *cobra.Command, count int, output string) {
	cmd.Println(
		fmt.Sprintf(assets.TextDesc(embed.TextDescKeyWriteObsidianGenerated),
			count, output,
		),
	)
	cmd.Println()
	cmd.Println(assets.TextDesc(embed.TextDescKeyWriteObsidianNextStepsHeading))
	cmd.Println(
		fmt.Sprintf(
			assets.TextDesc(embed.TextDescKeyWriteObsidianNextSteps),
			output,
		),
	)
}
