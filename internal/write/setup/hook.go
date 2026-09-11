//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package setup

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
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

// Context prints a JSON hook response line. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - response: JSON-encoded hook response.
func Context(cmd *cobra.Command, response string) {
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
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHookCopilotSkipped),
		targetFile))
	cmd.Println(desc.Text(text.DescKeyWriteHookCopilotForceHint))
}

// InfoCopilotMerged reports that copilot instructions were merged
// into an existing file.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the merged file
func InfoCopilotMerged(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHookCopilotMerged),
		targetFile))
}

// InfoCopilotCreated reports that copilot instructions were created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the created file
func InfoCopilotCreated(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHookCopilotCreated),
		targetFile))
}

// InfoCopilotSessionsDir reports that the sessions directory was created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - sessionsDir: Path to the sessions directory
func InfoCopilotSessionsDir(cmd *cobra.Command, sessionsDir string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHookCopilotSessionsDir),
		sessionsDir))
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
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHookCopilotCLICreated),
		targetFile))
}

// InfoAgentsCreated reports that AGENTS.md was created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the created file
func InfoAgentsCreated(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHookAgentsCreated),
		targetFile))
}

// InfoAgentsMerged reports that ctx content was merged into AGENTS.md.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the merged file
func InfoAgentsMerged(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHookAgentsMerged),
		targetFile))
}

// InfoAgentsSkipped reports that AGENTS.md was skipped because
// ctx markers already exist.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the existing file
func InfoAgentsSkipped(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHookAgentsSkipped),
		targetFile))
}

// InfoAgentsSummary prints the post-write summary for AGENTS.md.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoAgentsSummary(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(desc.Text(text.DescKeyWriteHookAgentsSummary))
}

// InfoOpenCodeCreated reports that an OpenCode integration file was
// created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the created file
func InfoOpenCodeCreated(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHookOpenCodeCreated),
		targetFile))
}

// InfoOpenCodeSkipped reports that an OpenCode integration file was
// skipped because it already exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the existing file
func InfoOpenCodeSkipped(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHookOpenCodeSkipped),
		targetFile))
}

// InfoOpenCodeSummary prints the post-write summary for OpenCode.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoOpenCodeSummary(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(desc.Text(text.DescKeyWriteHookOpenCodeSummary))
}

// InfoCodexCreated reports that a Codex integration file was created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the created file
func InfoCodexCreated(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHookCodexCreated),
		targetFile))
}

// InfoCodexMerged reports that ctx content was merged into an
// existing Codex integration file.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the merged file
func InfoCodexMerged(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHookCodexMerged),
		targetFile))
}

// InfoCodexSkipped reports that a Codex integration file was
// skipped because it is already up to date.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the existing file
func InfoCodexSkipped(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHookCodexSkipped),
		targetFile))
}

// InfoCodexRejected reports that a foreign (not ctx-managed) file
// occupies a Codex integration path and was left untouched.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the foreign file
func InfoCodexRejected(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHookCodexRejected),
		targetFile))
}

// InfoCodexPluginActive reports that the ctx Codex plugin is
// enabled, so hooks, MCP, and skills are not deployed per project.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoCodexPluginActive(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteHookCodexPluginActive))
}

// InfoCodexPluginWrongVariant warns that the enabled ctx Codex
// plugin is the legacy Claude Code variant and that the
// project-local route is being deployed instead.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoCodexPluginWrongVariant(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteHookCodexPluginWrongVariant))
	cmd.Println()
}

// InfoCodexProjectAlso warns that a project-local .codex/hooks.json
// coexists with the enabled plugin, which would run every hook
// twice; names the file to remove.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoCodexProjectAlso(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteHookCodexProjectAlso))
}

// InfoCodexSummary prints the post-write summary for Codex,
// ending with the /hooks trust reminder.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoCodexSummary(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(desc.Text(text.DescKeyWriteHookCodexSummary))
}

// InfoCodexSummaryPlugin prints the post-write summary for the
// plugin-enabled short-circuit: hooks, MCP, and skills come from
// the installed plugin, so only AGENTS.md was deployed and no
// project-local trust step applies.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoCodexSummaryPlugin(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(desc.Text(text.DescKeyWriteHookCodexSummaryPlugin))
}

// InfoCodexState prints the detection state line for
// `ctx setup codex`.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stateLabel: human-readable detection state
func InfoCodexState(cmd *cobra.Command, stateLabel string) {
	cmd.Println()
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHookCodexState), stateLabel))
}

// InfoPiCreated reports that a Pi integration file was created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the created file
func InfoPiCreated(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHookPiCreated),
		targetFile))
}

// InfoPiSkipped reports that a Pi integration file was skipped
// because it already exists.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the existing file
func InfoPiSkipped(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHookPiSkipped),
		targetFile))
}

// InfoPiSummary prints the post-write summary for Pi.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoPiSummary(cmd *cobra.Command) {
	cmd.Println()
	cmd.Println(desc.Text(text.DescKeyWriteHookPiSummary))
}

// InfoCopilotCLISkipped reports that copilot-cli hooks were skipped
// because they already exist.
//
// Parameters:
//   - cmd: Cobra command for output
//   - targetFile: Path to the existing file
func InfoCopilotCLISkipped(cmd *cobra.Command, targetFile string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHookCopilotCLISkipped),
		targetFile))
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

// InfoCursorIntegration prints Cursor integration instructions.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoCursorIntegration(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteSetupCursorHead))
	cmd.Println(desc.Text(text.DescKeyWriteSetupCursorRun))
	cmd.Println(desc.Text(text.DescKeyWriteSetupCursorMCP))
	cmd.Println(desc.Text(text.DescKeyWriteSetupCursorSync))
}

// InfoKiroIntegration prints Kiro integration instructions.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoKiroIntegration(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteSetupKiroHead))
	cmd.Println(desc.Text(text.DescKeyWriteSetupKiroRun))
	cmd.Println(desc.Text(text.DescKeyWriteSetupKiroMCP))
	cmd.Println(desc.Text(text.DescKeyWriteSetupKiroSync))
}

// InfoClineIntegration prints Cline integration instructions.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoClineIntegration(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteSetupClineHead))
	cmd.Println(desc.Text(text.DescKeyWriteSetupClineRun))
	cmd.Println(desc.Text(text.DescKeyWriteSetupClineMCP))
	cmd.Println(desc.Text(text.DescKeyWriteSetupClineSync))
}

// DeployComplete prints the completion message for a tool setup.
//
// Parameters:
//   - cmd: Cobra command for output
//   - tool: Tool name (e.g., "Cursor", "Kiro", "Cline")
//   - mcpPath: Path to the MCP config file
//   - steeringPath: Path to the steering directory
func DeployComplete(cmd *cobra.Command, tool, mcpPath, steeringPath string) {
	cmd.Println()
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteSetupDeployComplete), tool))
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteSetupDeployMCP), mcpPath))
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteSetupDeploySteering),
		steeringPath))
}

// DeployFileExists prints that a file already exists and was skipped.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: Path to the existing file
func DeployFileExists(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteSetupDeployExists), path))
}

// DeployFileCreated prints that a file was created.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: Path to the created file
func DeployFileCreated(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteSetupDeployCreated), path))
}

// DeploySteeringSynced prints that a steering file was synced.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Name of the synced file
func DeploySteeringSynced(cmd *cobra.Command, name string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteSetupDeploySynced), name))
}

// DeploySteeringSkipped prints that a steering file was skipped.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Name of the skipped file
func DeploySteeringSkipped(cmd *cobra.Command, name string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteSetupDeploySkipSteer),
		name))
}

// DeployNoSteering prints that no steering files are
// available to sync.
//
// Parameters:
//   - cmd: Cobra command for output
func DeployNoSteering(cmd *cobra.Command) {
	cmd.Println(desc.Text(
		text.DescKeyWriteSetupNoSteeringToSync))
}
