//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package deps

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// InfoNoProject reports that no supported project was detected.
//
// Parameters:
//   - cmd: Cobra command for output
//   - builderNames: Comma-separated list of supported project types
func InfoNoProject(cmd *cobra.Command, builderNames string) {
	cmd.Println(desc.Text(text.DescKeyWriteDepsNoProject))
	cmd.Println(desc.Text(text.DescKeyWriteDepsLookingFor))
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteDepsUseType), builderNames))
}

// NoDeps reports that no dependencies were found.
//
// Parameters:
//   - cmd: Cobra command for output
func NoDeps(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteDepsNoDeps))
}

// Mermaid prints rendered Mermaid dependency graph output.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - content: Pre-rendered Mermaid graph string.
func Mermaid(cmd *cobra.Command, content string) {
	if cmd == nil {
		return
	}
	cmd.Print(content)
}

// Table prints rendered table dependency output.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - content: Pre-rendered table string.
func Table(cmd *cobra.Command, content string) {
	if cmd == nil {
		return
	}
	cmd.Print(content)
}

// JSON prints rendered JSON dependency output.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - content: Pre-rendered JSON string.
func JSON(cmd *cobra.Command, content string) {
	if cmd == nil {
		return
	}
	cmd.Print(content)
}
