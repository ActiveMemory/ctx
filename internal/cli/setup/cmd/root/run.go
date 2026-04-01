//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/agent"
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreAgents "github.com/ActiveMemory/ctx/internal/cli/setup/core/agents"
	coreCopilot "github.com/ActiveMemory/ctx/internal/cli/setup/core/copilot"
	coreCopilotCLI "github.com/ActiveMemory/ctx/internal/cli/setup/core/copilot_cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/err/config"
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
func Run(cmd *cobra.Command, args []string, writeFile bool) error {
	tool := strings.ToLower(args[0])

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
		writeSetup.InfoTool(cmd, desc.Text(text.DescKeyHookClaude))

	case cfgHook.ToolCursor:
		writeSetup.InfoTool(cmd, desc.Text(text.DescKeyHookCursor))

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
			return coreCopilotCLI.Deploy(cmd)
		}
		writeSetup.InfoTool(cmd, desc.Text(text.DescKeyHookCopilotCLI))

	case cfgHook.ToolWindsurf:
		writeSetup.InfoTool(cmd, desc.Text(text.DescKeyHookWindsurf))

	default:
		writeSetup.InfoUnknownTool(cmd, tool)
		writeSetup.InfoTool(cmd, desc.Text(text.DescKeyHookSupportedTools))
		return config.UnsupportedTool(tool)
	}

	return nil
}
