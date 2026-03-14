//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/write/config"
	"github.com/spf13/cobra"
)

// HookNudge prints a pre-built nudge box to stdout.
//
// Used by system hooks to emit nudge messages through the write layer
// rather than calling cmd.Println directly.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - nudgeBox: fully formatted nudge box string.
func HookNudge(cmd *cobra.Command, nudgeBox string) {
	if cmd == nil {
		return
	}
	cmd.Println(nudgeBox)
}

// InfoHookTool prints a tool integration section to stdout.
//
// The content is a pre-formatted multi-line text block loaded from
// commands.yaml. A trailing newline is not added: the content is
// expected to include its own formatting.
//
// Parameters:
//   - cmd: Cobra command for output
//   - content: Pre-formatted text block
func InfoHookTool(cmd *cobra.Command, content string) {
	cmd.Print(content)
}

// InfoHookCopilotSkipped reports that copilot instructions were skipped
// because the ctx marker already exists in the target file.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the existing file
func InfoHookCopilotSkipped(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(config.TplHookCopilotSkipped, targetFile))
	cmd.Println(config.TplHookCopilotForceHint)
}

// InfoHookCopilotMerged reports that copilot instructions were merged
// into an existing file.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the merged file
func InfoHookCopilotMerged(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(config.TplHookCopilotMerged, targetFile))
}

// InfoHookCopilotCreated reports that copilot instructions were created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the created file
func InfoHookCopilotCreated(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(config.TplHookCopilotCreated, targetFile))
}

// InfoHookCopilotSessionsDir reports that the sessions directory was created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - sessionsDir: Path to the sessions directory
func InfoHookCopilotSessionsDir(cmd *cobra.Command, sessionsDir string) {
	cmd.Println(fmt.Sprintf(config.TplHookCopilotSessionsDir, sessionsDir))
}

// InfoHookCopilotSummary prints the post-write summary for copilot.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoHookCopilotSummary(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(config.TplHookCopilotSummary)
}

// InfoHookUnknownTool prints the unknown tool message.
//
// Parameters:
//   - cmd: Cobra command for output
//   - tool: The unrecognized tool name
func InfoHookUnknownTool(cmd *cobra.Command, tool string) {
	cmd.Println(fmt.Sprintf(config.TplHookUnknownTool, tool))
}
