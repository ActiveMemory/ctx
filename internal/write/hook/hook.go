//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/spf13/cobra"
)

// Nudge prints a pre-built nudge box to stdout.
//
// Used by system hooks to emit nudge messages through the write layer
// rather than calling cmd.Println directly.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - nudgeBox: fully formatted nudge box string.
func Nudge(cmd *cobra.Command, nudgeBox string) {
	if cmd == nil {
		return
	}
	cmd.Println(nudgeBox)
}

// InfoTool prints a tool integration section to stdout.
//
// The content is a pre-formatted multi-line text block loaded from
// commands.yaml. A trailing newline is not added: the content is
// expected to include its own formatting.
//
// Parameters:
//   - cmd: Cobra command for output
//   - content: Pre-formatted text block
func InfoTool(cmd *cobra.Command, content string) {
	cmd.Print(content)
}

// InfoCopilotSkipped reports that copilot instructions were skipped
// because the ctx marker already exists in the target file.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the existing file
func InfoCopilotSkipped(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyWriteHookCopilotSkipped), targetFile))
	cmd.Println(desc.TextDesc(text.DescKeyWriteHookCopilotForceHint))
}

// InfoCopilotMerged reports that copilot instructions were merged
// into an existing file.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the merged file
func InfoCopilotMerged(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyWriteHookCopilotMerged), targetFile))
}

// InfoCopilotCreated reports that copilot instructions were created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the created file
func InfoCopilotCreated(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyWriteHookCopilotCreated), targetFile))
}

// InfoCopilotSessionsDir reports that the sessions directory was created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - sessionsDir: Path to the sessions directory
func InfoCopilotSessionsDir(cmd *cobra.Command, sessionsDir string) {
	cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyWriteHookCopilotSessionsDir), sessionsDir))
}

// InfoCopilotSummary prints the post-write summary for copilot.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoCopilotSummary(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(desc.TextDesc(text.DescKeyWriteHookCopilotSummary))
}

// InfoUnknownTool prints the unknown tool message.
//
// Parameters:
//   - cmd: Cobra command for output
//   - tool: The unrecognized tool name
func InfoUnknownTool(cmd *cobra.Command, tool string) {
	cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyWriteHookUnknownTool), tool))
}
