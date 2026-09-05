//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package text

// DescKeys for hook messages.
const (
	// DescKeyHookAider is the text key for hook aider messages.
	DescKeyHookAider = "hook.aider"
	// DescKeyHookAgents is the text key for hook agents messages.
	DescKeyHookAgents = "hook.agents"
	// DescKeyHookCodex is the text key for hook codex messages.
	DescKeyHookCodex = "hook.codex"
	// DescKeyHookCopilot is the text key for hook copilot messages.
	DescKeyHookCopilot = "hook.copilot"
	// DescKeyHookCopilotCLI is the text key for hook copilot cli messages.
	DescKeyHookCopilotCLI = "hook.copilot-cli"
	// DescKeyHookOpenCode is the text key for hook opencode messages.
	DescKeyHookOpenCode = "hook.opencode"
	// DescKeyHookPi is the text key for hook pi messages.
	DescKeyHookPi = "hook.pi"
	// DescKeyHookSupportedTools is the text key for hook supported tools messages.
	DescKeyHookSupportedTools = "hook.supported-tools"
	// DescKeyHookWindsurf is the text key for hook windsurf messages.
	DescKeyHookWindsurf = "hook.windsurf"
)

// DescKeys for hook write output.
const (
	// DescKeyWriteHookAgentsCreated is the text key for write hook agents created
	// messages.
	DescKeyWriteHookAgentsCreated = "write.hook-agents-created"
	// DescKeyWriteHookAgentsMerged is the text key for write hook agents merged
	// messages.
	DescKeyWriteHookAgentsMerged = "write.hook-agents-merged"
	// DescKeyWriteHookAgentsSkipped is the text key for write hook agents skipped
	// messages.
	DescKeyWriteHookAgentsSkipped = "write.hook-agents-skipped"
	// DescKeyWriteHookAgentsSummary is the text key for write hook agents summary
	// messages.
	DescKeyWriteHookAgentsSummary = "write.hook-agents-summary"
	// DescKeyWriteHookCodexCreated is the text key for write hook codex
	// created messages.
	DescKeyWriteHookCodexCreated = "write.hook-codex-created"
	// DescKeyWriteHookCodexMerged is the text key for write hook codex
	// merged messages.
	DescKeyWriteHookCodexMerged = "write.hook-codex-merged"
	// DescKeyWriteHookCodexSkipped is the text key for write hook codex
	// skipped messages.
	DescKeyWriteHookCodexSkipped = "write.hook-codex-skipped"
	// DescKeyWriteHookCodexRejected is the text key for write hook codex
	// rejected (foreign file) messages.
	DescKeyWriteHookCodexRejected = "write.hook-codex-rejected"
	// DescKeyWriteHookCodexPluginActive is the text key for the notice
	// printed when the ctx Codex plugin already provides hooks/skills/MCP.
	DescKeyWriteHookCodexPluginActive = "write.hook-codex-plugin-active"
	// DescKeyWriteHookCodexProjectAlso is the text key for the warning
	// that project-local hooks coexist with the enabled plugin.
	DescKeyWriteHookCodexProjectAlso = "write.hook-codex-project-also"
	// DescKeyWriteHookCodexPluginWrongVariant is the text key for the
	// warning printed when the enabled plugin is the legacy Claude
	// Code variant and the project-local route is deployed instead.
	DescKeyWriteHookCodexPluginWrongVariant = "write.hook-codex-plugin-wrong-variant"
	// DescKeyWriteHookCodexSummary is the text key for write hook codex
	// summary messages.
	DescKeyWriteHookCodexSummary = "write.hook-codex-summary"
	// DescKeyWriteHookCodexSummaryPlugin is the text key for the
	// summary printed when the ctx Codex plugin provides
	// hooks/MCP/skills and only AGENTS.md was deployed.
	DescKeyWriteHookCodexSummaryPlugin = "write.hook-codex-summary-plugin"
	// DescKeyWriteHookCodexState is the text key for the detection state
	// line printed by `ctx setup codex`.
	DescKeyWriteHookCodexState = "write.hook-codex-state"
	// DescKeyWriteHookCodexStateConfigured is the state label used when
	// the project-local Codex integration is already deployed.
	DescKeyWriteHookCodexStateConfigured = "write.hook-codex-state-configured"
	// DescKeyWriteHookCodexStateAbsent labels the state where the codex
	// binary is not on PATH.
	DescKeyWriteHookCodexStateAbsent = "write.hook-codex-state-absent"
	// DescKeyWriteHookCodexStateNotInstalled labels the state where codex
	// is present but the ctx plugin is not in its plugin cache.
	DescKeyWriteHookCodexStateNotInstalled = "write.hook-codex-state-not-installed"
	// DescKeyWriteHookCodexStateNotEnabled labels the state where the ctx
	// plugin is cached but not enabled in config.toml.
	DescKeyWriteHookCodexStateNotEnabled = "write.hook-codex-state-not-enabled"
	// DescKeyWriteHookCodexStateReady labels the state where codex and the
	// ctx plugin are both present and enabled.
	DescKeyWriteHookCodexStateReady = "write.hook-codex-state-ready"
	// DescKeyWriteHookCopilotCLICreated is the text key for write hook copilot
	// cli created messages.
	DescKeyWriteHookCopilotCLICreated = "write.hook-copilot-cli-created"
	// DescKeyWriteHookCopilotCLISkipped is the text key for write hook copilot
	// cli skipped messages.
	DescKeyWriteHookCopilotCLISkipped = "write.hook-copilot-cli-skipped"
	// DescKeyWriteHookCopilotCLISummary is the text key for write hook copilot
	// cli summary messages.
	DescKeyWriteHookCopilotCLISummary = "write.hook-copilot-cli-summary"
	// DescKeyWriteHookCopilotCreated is the text key for write hook copilot
	// created messages.
	DescKeyWriteHookCopilotCreated = "write.hook-copilot-created"
	// DescKeyWriteHookCopilotForceHint is the text key for write hook copilot
	// force hint messages.
	DescKeyWriteHookCopilotForceHint = "write.hook-copilot-force-hint"
	// DescKeyWriteHookCopilotMerged is the text key for write hook copilot merged
	// messages.
	DescKeyWriteHookCopilotMerged = "write.hook-copilot-merged"
	// DescKeyWriteHookCopilotSessionsDir is the text key for write hook copilot
	// sessions dir messages.
	DescKeyWriteHookCopilotSessionsDir = "write.hook-copilot-sessions-dir"
	// DescKeyWriteHookCopilotSkipped is the text key for write hook copilot
	// skipped messages.
	DescKeyWriteHookCopilotSkipped = "write.hook-copilot-skipped"
	// DescKeyWriteHookCopilotSummary is the text key for write hook copilot
	// summary messages.
	DescKeyWriteHookCopilotSummary = "write.hook-copilot-summary"
	// DescKeyWriteHookOpenCodeCreated is the text key for write hook opencode
	// created messages.
	DescKeyWriteHookOpenCodeCreated = "write.hook-opencode-created"
	// DescKeyWriteHookOpenCodeSkipped is the text key for write hook opencode
	// skipped messages.
	DescKeyWriteHookOpenCodeSkipped = "write.hook-opencode-skipped"
	// DescKeyWriteHookOpenCodeSummary is the text key for write hook opencode
	// summary messages.
	DescKeyWriteHookOpenCodeSummary = "write.hook-opencode-summary"
	// DescKeyWriteHookPiCreated is the text key for write hook pi created
	// messages.
	DescKeyWriteHookPiCreated = "write.hook-pi-created"
	// DescKeyWriteHookPiSkipped is the text key for write hook pi skipped
	// messages.
	DescKeyWriteHookPiSkipped = "write.hook-pi-skipped"
	// DescKeyWriteHookPiSummary is the text key for write hook pi summary
	// messages.
	DescKeyWriteHookPiSummary = "write.hook-pi-summary"
	// DescKeyWriteHookUnknownTool is the text key for write hook unknown tool
	// messages.
	DescKeyWriteHookUnknownTool = "write.hook-unknown-tool"
)
