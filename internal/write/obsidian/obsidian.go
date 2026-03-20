//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package obsidian

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
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
		fmt.Sprintf(desc.TextDesc(text.DescKeyWriteObsidianGenerated),
			count, output,
		),
	)
	cmd.Println()
	cmd.Println(desc.TextDesc(text.DescKeyWriteObsidianNextStepsHeading))
	cmd.Println(
		fmt.Sprintf(
			desc.TextDesc(text.DescKeyWriteObsidianNextSteps),
			output,
		),
	)
}
