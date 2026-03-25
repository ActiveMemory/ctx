//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/agent"
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/err/config"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	writeErr "github.com/ActiveMemory/ctx/internal/write/err"
	"github.com/ActiveMemory/ctx/internal/write/hook"
)

// Run executes the hook command logic.
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
			return WriteAgentsMd(cmd)
		}
		hook.InfoTool(cmd, desc.Text(text.DescKeyHookAgents))
		hook.Separator(cmd)
		content, readErr := agent.AgentsMd()
		if readErr != nil {
			return readErr
		}
		hook.Content(cmd, string(content))

	case cfgHook.ToolClaudeCode, cfgHook.ToolClaude:
		hook.InfoTool(cmd, desc.Text(text.DescKeyHookClaude))

	case cfgHook.ToolCursor:
		hook.InfoTool(cmd, desc.Text(text.DescKeyHookCursor))

	case cfgHook.ToolAider:
		hook.InfoTool(cmd, desc.Text(text.DescKeyHookAider))

	case cfgHook.ToolCopilot:
		if writeFile {
			return WriteCopilotInstructions(cmd)
		}
		hook.InfoTool(cmd, desc.Text(text.DescKeyHookCopilot))
		hook.Separator(cmd)
		content, readErr := agent.CopilotInstructions()
		if readErr != nil {
			return readErr
		}
		hook.Content(cmd, string(content))

	case cfgHook.ToolCopilotCLI:
		if writeFile {
			return WriteCopilotCLIHooks(cmd)
		}
		hook.InfoTool(cmd, desc.Text(text.DescKeyHookCopilotCLI))

	case cfgHook.ToolWindsurf:
		hook.InfoTool(cmd, desc.Text(text.DescKeyHookWindsurf))

	default:
		hook.InfoUnknownTool(cmd, tool)
		hook.InfoTool(cmd, desc.Text(text.DescKeyHookSupportedTools))
		return config.UnsupportedTool(tool)
	}

	return nil
}

// WriteCopilotInstructions generates .github/copilot-instructions.md.
//
// Creates the .github/ directory if needed and writes the comprehensive
// Copilot instructions file. Preserves existing non-ctx content by
// checking for ctx markers.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if directory creation or file write fails
func WriteCopilotInstructions(cmd *cobra.Command) error {
	targetFile := filepath.Join(cfgHook.DirGitHub, cfgHook.FileCopilotInstructions)

	// Create .github/ directory if needed
	if err := os.MkdirAll(cfgHook.DirGitHub, fs.PermExec); err != nil {
		return errFs.Mkdir(cfgHook.DirGitHub, err)
	}

	// Load the copilot instructions from embedded assets
	instructions, readErr := agent.CopilotInstructions()
	if readErr != nil {
		return readErr
	}

	// Check if the file exists
	existingContent, err := os.ReadFile(filepath.Clean(targetFile))
	fileExists := err == nil

	if fileExists {
		existingStr := string(existingContent)
		if strings.Contains(existingStr, marker.CopilotMarkerStart) {
			hook.InfoCopilotSkipped(cmd, targetFile)
			return nil
		}

		// File exists without ctx markers: append ctx content
		merged := existingStr + token.NewlineLF + string(instructions)
		if wErr := os.WriteFile(targetFile, []byte(merged), fs.PermFile); wErr != nil {
			return errFs.FileWrite(targetFile, wErr)
		}
		hook.InfoCopilotMerged(cmd, targetFile)
		return nil
	}

	// File doesn't exist: create it
	if wErr := os.WriteFile(
		targetFile, instructions, fs.PermFile,
	); wErr != nil {
		return errFs.FileWrite(targetFile, wErr)
	}
	hook.InfoCopilotCreated(cmd, targetFile)

	// Also create .context/sessions/ if it doesn't exist
	sessionsDir := filepath.Join(dir.Context, dir.Sessions)
	if mkErr := os.MkdirAll(sessionsDir, fs.PermExec); mkErr != nil {
		writeErr.WarnFile(cmd, sessionsDir, mkErr)
	} else {
		hook.InfoCopilotSessionsDir(cmd, sessionsDir)
	}

	hook.InfoCopilotSummary(cmd)

	// Also create .vscode/mcp.json if it doesn't exist
	if err := ensureVSCodeMCPJSON(cmd); err != nil {
		cmd.Println("  ⚠ .vscode/mcp.json: " + err.Error())
	}

	return nil
}

// WriteCopilotCLIHooks generates .github/hooks/ctx-hooks.json and the
// accompanying hook scripts for GitHub Copilot CLI integration.
//
// Creates the .github/hooks/ and .github/hooks/scripts/ directories if
// needed and writes the JSON config plus bash and PowerShell scripts
// from embedded assets. Also writes .github/agents/ctx.md and
// .github/instructions/context.instructions.md for Copilot CLI.
// Skips if ctx-hooks.json already exists.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if directory creation or file write fails
func WriteCopilotCLIHooks(cmd *cobra.Command) error {
	hooksDir := filepath.Join(cfgHook.DirGitHub, cfgHook.DirGitHubHooks)
	scriptsDir := filepath.Join(hooksDir, cfgHook.DirGitHubHooksScripts)
	targetJSON := filepath.Join(hooksDir, cfgHook.FileCopilotCLIHooksJSON)

	// Check if ctx-hooks.json already exists
	if _, err := os.Stat(targetJSON); err == nil {
		hook.InfoCopilotCLISkipped(cmd, targetJSON)
		return nil
	}

	// Create directories
	if err := os.MkdirAll(scriptsDir, fs.PermExec); err != nil {
		return errFs.Mkdir(scriptsDir, err)
	}

	// Write ctx-hooks.json
	jsonContent, readErr := agent.CopilotCLIHooksJSON()
	if readErr != nil {
		return readErr
	}
	if wErr := os.WriteFile(targetJSON, jsonContent, fs.PermFile); wErr != nil {
		return errFs.FileWrite(targetJSON, wErr)
	}
	hook.InfoCopilotCLICreated(cmd, targetJSON)

	// Write all hook scripts
	scripts, scrErr := agent.CopilotCLIScripts()
	if scrErr != nil {
		return scrErr
	}
	for name, content := range scripts {
		target := filepath.Join(scriptsDir, name)
		if wErr := os.WriteFile(target, content, fs.PermExec); wErr != nil {
			return errFs.FileWrite(target, wErr)
		}
		hook.InfoCopilotCLICreated(cmd, target)
	}

	// Write .github/agents/ctx.md
	if err := writeCopilotCLIAgent(cmd); err != nil {
		writeErr.WarnFile(cmd, cfgHook.DirGitHubAgents, err)
	}

	// Write .github/instructions/context.instructions.md
	if err := writeCopilotCLIInstructions(cmd); err != nil {
		writeErr.WarnFile(cmd, cfgHook.DirGitHubInstructions, err)
	}

	// Register ctx MCP server in ~/.copilot/mcp-config.json
	if err := ensureCopilotCLIMCPConfig(cmd); err != nil {
		cmd.Println("  ⚠ mcp-config.json: " + err.Error())
	}

	hook.InfoCopilotCLISummary(cmd)
	return nil
}

// writeCopilotCLIAgent creates .github/agents/ctx.md for Copilot CLI
// custom agent delegation. Skips if the file already exists.
func writeCopilotCLIAgent(cmd *cobra.Command) error {
	agentsDir := filepath.Join(cfgHook.DirGitHub, cfgHook.DirGitHubAgents)
	target := filepath.Join(agentsDir, cfgHook.FileAgentsCtxMd)

	if _, err := os.Stat(target); err == nil {
		hook.InfoCopilotCLICreated(cmd, target+" (exists, skipped)")
		return nil
	}

	if err := os.MkdirAll(agentsDir, fs.PermExec); err != nil {
		return err
	}

	content, readErr := agent.AgentsCtxMd()
	if readErr != nil {
		return readErr
	}
	if wErr := os.WriteFile(target, content, fs.PermFile); wErr != nil {
		return wErr
	}
	hook.InfoCopilotCLICreated(cmd, target)
	return nil
}

// writeCopilotCLIInstructions creates
// .github/instructions/context.instructions.md for path-specific
// context file instructions. Skips if the file already exists.
func writeCopilotCLIInstructions(cmd *cobra.Command) error {
	instrDir := filepath.Join(cfgHook.DirGitHub, cfgHook.DirGitHubInstructions)
	target := filepath.Join(instrDir, cfgHook.FileInstructionsCtxMd)

	if _, err := os.Stat(target); err == nil {
		hook.InfoCopilotCLICreated(cmd, target+" (exists, skipped)")
		return nil
	}

	if err := os.MkdirAll(instrDir, fs.PermExec); err != nil {
		return err
	}

	content, readErr := agent.InstructionsCtxMd()
	if readErr != nil {
		return readErr
	}
	if wErr := os.WriteFile(target, content, fs.PermFile); wErr != nil {
		return wErr
	}
	hook.InfoCopilotCLICreated(cmd, target)
	return nil
}

// WriteAgentsMd generates AGENTS.md in the project root.
//
// Creates AGENTS.md with universal agent instructions. Preserves existing
// non-ctx content by checking for ctx:agents markers. If the file exists
// with markers, skips. If it exists without markers, merges.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if file write fails
func WriteAgentsMd(cmd *cobra.Command) error {
	targetFile := cfgHook.FileAgentsMd

	// Load the AGENTS.md template
	agentsContent, readErr := agent.AgentsMd()
	if readErr != nil {
		return readErr
	}

	// Check if the file exists
	existingContent, err := os.ReadFile(filepath.Clean(targetFile))
	fileExists := err == nil

	if fileExists {
		existingStr := string(existingContent)
		if strings.Contains(existingStr, marker.AgentsMarkerStart) {
			hook.InfoAgentsSkipped(cmd, targetFile)
			return nil
		}

		// File exists without ctx markers: append ctx content
		merged := existingStr + token.NewlineLF + string(agentsContent)
		if wErr := os.WriteFile(targetFile, []byte(merged), fs.PermFile); wErr != nil {
			return errFs.FileWrite(targetFile, wErr)
		}
		hook.InfoAgentsMerged(cmd, targetFile)
		return nil
	}

	// File doesn't exist: create it
	if wErr := os.WriteFile(targetFile, agentsContent, fs.PermFile); wErr != nil {
		return errFs.FileWrite(targetFile, wErr)
	}
	hook.InfoAgentsCreated(cmd, targetFile)

	hook.InfoAgentsSummary(cmd)
	return nil
}

// ensureVSCodeMCPJSON creates .vscode/mcp.json to register the ctx MCP
// server for VS Code Copilot. Skips if the file already exists.
func ensureVSCodeMCPJSON(cmd *cobra.Command) error {
	vsDir := ".vscode"
	target := filepath.Join(vsDir, "mcp.json")

	if _, err := os.Stat(target); err == nil {
		cmd.Println("  ○ " + target + " (exists, skipped)")
		return nil
	}

	if err := os.MkdirAll(vsDir, fs.PermExec); err != nil {
		return err
	}

	mcpCfg := map[string]interface{}{
		"servers": map[string]interface{}{
			"ctx": map[string]interface{}{
				"command": "ctx",
				"args":    []string{"mcp", "serve"},
			},
		},
	}
	data, _ := json.MarshalIndent(mcpCfg, "", "  ")
	data = append(data, '\n')

	if err := os.WriteFile(target, data, fs.PermFile); err != nil {
		return err
	}
	cmd.Println("  ✓ " + target)
	return nil
}

// ensureCopilotCLIMCPConfig registers the ctx MCP server in
// ~/.copilot/mcp-config.json (or $COPILOT_HOME/mcp-config.json).
// Merge-safe: reads existing config, adds ctx server, writes back.
// Skips if ctx server is already registered.
func ensureCopilotCLIMCPConfig(cmd *cobra.Command) error {
	copilotHome := os.Getenv(cfgHook.EnvCopilotHome)
	if copilotHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		copilotHome = filepath.Join(home, cfgHook.DirCopilotHome)
	}

	target := filepath.Join(copilotHome, cfgHook.FileMCPConfigJSON)

	// Read existing config if it exists
	existing := make(map[string]interface{})
	if data, err := os.ReadFile(filepath.Clean(target)); err == nil {
		if jErr := json.Unmarshal(data, &existing); jErr != nil {
			return jErr
		}
	}

	// Get or create mcpServers map
	servers, _ := existing["mcpServers"].(map[string]interface{})
	if servers == nil {
		servers = make(map[string]interface{})
	}

	// Check if ctx is already registered
	if _, ok := servers["ctx"]; ok {
		cmd.Println("  ○ " + target + " (ctx server exists, skipped)")
		return nil
	}

	// Add ctx MCP server
	servers["ctx"] = map[string]interface{}{
		"type":    "local",
		"command": "ctx",
		"args":    []string{"mcp", "serve"},
		"tools":   []string{"*"},
	}
	existing["mcpServers"] = servers

	// Create directory if needed
	if err := os.MkdirAll(copilotHome, fs.PermExec); err != nil {
		return err
	}

	data, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	if wErr := os.WriteFile(target, data, fs.PermFile); wErr != nil {
		return wErr
	}
	cmd.Println("  ✓ " + target)
	return nil
}
