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

// NudgeBlock prints a nudge box followed by an empty line.
// Empty box or nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - nudgeBox: fully formatted nudge box string.
func NudgeBlock(cmd *cobra.Command, nudgeBox string) {
	if cmd == nil || nudgeBox == "" {
		return
	}
	cmd.Println(nudgeBox)
	cmd.Println()
}

// HookContext prints a JSON hook response line. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - response: JSON-encoded hook response.
func HookContext(cmd *cobra.Command, response string) {
	if cmd == nil {
		return
	}
	cmd.Println(response)
}

// BlockResponse prints a JSON block response line. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - response: JSON-encoded block response.
func BlockResponse(cmd *cobra.Command, response string) {
	if cmd == nil {
		return
	}
	cmd.Println(response)
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
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteHookCopilotSkipped), targetFile))
	cmd.Println(desc.Text(text.DescKeyWriteHookCopilotForceHint))
}

// InfoCopilotMerged reports that copilot instructions were merged
// into an existing file.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the merged file
func InfoCopilotMerged(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteHookCopilotMerged), targetFile))
}

// InfoCopilotCreated reports that copilot instructions were created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the created file
func InfoCopilotCreated(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteHookCopilotCreated), targetFile))
}

// InfoCopilotSessionsDir reports that the sessions directory was created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - sessionsDir: Path to the sessions directory
func InfoCopilotSessionsDir(cmd *cobra.Command, sessionsDir string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteHookCopilotSessionsDir), sessionsDir))
}

// InfoCopilotSummary prints the post-write summary for copilot.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoCopilotSummary(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(desc.Text(text.DescKeyWriteHookCopilotSummary))
}

// InfoCopilotCLICreated reports that copilot-cli hook files were created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the created file
func InfoCopilotCLICreated(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteHookCopilotCLICreated), targetFile))
}

// InfoCopilotCLISkipped reports that copilot-cli hooks were skipped
// because they already exist.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the existing file
func InfoCopilotCLISkipped(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteHookCopilotCLISkipped), targetFile))
}

// InfoCopilotCLISummary prints the post-write summary for copilot-cli.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoCopilotCLISummary(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(desc.Text(text.DescKeyWriteHookCopilotCLISummary))
}

// InfoUnknownTool prints the unknown tool message.
//
// Parameters:
//   - cmd: Cobra command for output
//   - tool: The unrecognized tool name
func InfoUnknownTool(cmd *cobra.Command, tool string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteHookUnknownTool), tool))
}

// Separator prints a blank line between hook output sections.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func Separator(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println()
}

// Content prints raw hook content to stdout.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - content: Pre-rendered content string.
func Content(cmd *cobra.Command, content string) {
	if cmd == nil {
		return
	}
	cmd.Print(content)
}
