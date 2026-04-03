//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for hook messages.
const (
	DescKeyHookAider          = "hook.aider"
	DescKeyHookAgents         = "hook.agents"
	DescKeyHookClaude         = "hook.claude"
	DescKeyHookCopilot        = "hook.copilot"
	DescKeyHookCopilotCLI     = "hook.copilot-cli"
	DescKeyHookSupportedTools = "hook.supported-tools"
	DescKeyHookWindsurf       = "hook.windsurf"
)

// DescKeys for hook write output.
const (
	DescKeyWriteHookAgentsCreated      = "write.hook-agents-created"
	DescKeyWriteHookAgentsMerged       = "write.hook-agents-merged"
	DescKeyWriteHookAgentsSkipped      = "write.hook-agents-skipped"
	DescKeyWriteHookAgentsSummary      = "write.hook-agents-summary"
	DescKeyWriteHookCopilotCLICreated  = "write.hook-copilot-cli-created"
	DescKeyWriteHookCopilotCLISkipped  = "write.hook-copilot-cli-skipped"
	DescKeyWriteHookCopilotCLISummary  = "write.hook-copilot-cli-summary"
	DescKeyWriteHookCopilotCreated     = "write.hook-copilot-created"
	DescKeyWriteHookCopilotForceHint   = "write.hook-copilot-force-hint"
	DescKeyWriteHookCopilotMerged      = "write.hook-copilot-merged"
	DescKeyWriteHookCopilotSessionsDir = "write.hook-copilot-sessions-dir"
	DescKeyWriteHookCopilotSkipped     = "write.hook-copilot-skipped"
	DescKeyWriteHookCopilotSummary     = "write.hook-copilot-summary"
	DescKeyWriteHookUnknownTool        = "write.hook-unknown-tool"
)
