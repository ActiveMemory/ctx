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
// from embedded assets. Skips if ctx-hooks.json already exists.
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

	hook.InfoCopilotCLISummary(cmd)
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
