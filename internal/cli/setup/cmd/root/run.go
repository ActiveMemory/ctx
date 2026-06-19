//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/agent"
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreCC "github.com/ActiveMemory/ctx/internal/cli/initialize/core/claudecheck"
	coreAgents "github.com/ActiveMemory/ctx/internal/cli/setup/core/agents"
	coreBackend "github.com/ActiveMemory/ctx/internal/cli/setup/core/backend"
	coreCline "github.com/ActiveMemory/ctx/internal/cli/setup/core/cline"
	coreCopilot "github.com/ActiveMemory/ctx/internal/cli/setup/core/copilot"
	coreCopCLI "github.com/ActiveMemory/ctx/internal/cli/setup/core/copilotcli"
	coreCursor "github.com/ActiveMemory/ctx/internal/cli/setup/core/cursor"
	coreKiro "github.com/ActiveMemory/ctx/internal/cli/setup/core/kiro"
	coreOpenCode "github.com/ActiveMemory/ctx/internal/cli/setup/core/opencode"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/err/config"
	"github.com/ActiveMemory/ctx/internal/i18n"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Run executes the setup command logic.
//
// Outputs integration instructions and configuration snippets for the
// specified AI tool. With --write, generates the configuration file
// directly.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - args: Command arguments; args[0] is the tool name
//   - writeFile: If true, write the configuration file instead of printing
//
// Returns:
//   - error: Non-nil if the tool is not supported or file write fails
func Run(
	cmd *cobra.Command,
	args []string,
	writeFile bool,
	backend string,
	endpoint string,
	apiKeyEnv string,
	model string,
	timeout string,
) error {
	if backend != "" {
		return coreBackend.Run(cmd.OutOrStdout(), coreBackend.Options{
			Name:      i18n.Fold(backend),
			Endpoint:  endpoint,
			APIKeyEnv: apiKeyEnv,
			Model:     model,
			Timeout:   timeout,
			Write:     writeFile,
		})
	}
	tool := i18n.Fold(args[0])

	switch tool {
	case cfgHook.ToolAgents:
		if writeFile {
			return coreAgents.Deploy(cmd)
		}
		writeSetup.InfoTool(cmd, desc.Text(text.DescKeyHookAgents))
		writeSetup.Separator(cmd)
		content, readErr := agent.AgentsMd()
		if readErr != nil {
			return readErr
		}
		writeSetup.Content(cmd, string(content))

	case cfgHook.ToolClaudeCode, cfgHook.ToolClaude:
		coreCC.SetupHint(cmd)

	case cfgHook.ToolCursor:
		if writeFile {
			return coreCursor.Deploy(cmd)
		}
		writeSetup.InfoCursorIntegration(cmd)

	case cfgHook.ToolKiro:
		if writeFile {
			return coreKiro.Deploy(cmd)
		}
		writeSetup.InfoKiroIntegration(cmd)

	case cfgHook.ToolCline:
		if writeFile {
			return coreCline.Deploy(cmd)
		}
		writeSetup.InfoClineIntegration(cmd)

	case cfgHook.ToolAider:
		writeSetup.InfoTool(cmd, desc.Text(text.DescKeyHookAider))

	case cfgHook.ToolCopilot:
		if writeFile {
			return coreCopilot.DeployInstructions(cmd)
		}
		writeSetup.InfoTool(cmd, desc.Text(text.DescKeyHookCopilot))
		writeSetup.Separator(cmd)
		content, readErr := agent.CopilotInstructions()
		if readErr != nil {
			return readErr
		}
		writeSetup.Content(cmd, string(content))

	case cfgHook.ToolCopilotCLI:
		if writeFile {
			return coreCopCLI.Deploy(cmd)
		}
		writeSetup.InfoTool(cmd, desc.Text(text.DescKeyHookCopilotCLI))

	case cfgHook.ToolOpenCode:
		if writeFile {
			return coreOpenCode.Deploy(cmd)
		}
		writeSetup.InfoTool(cmd, desc.Text(text.DescKeyHookOpenCode))

	case cfgHook.ToolWindsurf:
		writeSetup.InfoTool(cmd, desc.Text(text.DescKeyHookWindsurf))

	default:
		writeSetup.InfoUnknownTool(cmd, tool)
		writeSetup.InfoTool(cmd, desc.Text(text.DescKeyHookSupportedTools))
		return config.UnsupportedTool(tool)
	}

	return nil
}
