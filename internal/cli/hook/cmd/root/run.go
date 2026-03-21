//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/agent"
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	hookCfg "github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/err/config"
	ctxErr "github.com/ActiveMemory/ctx/internal/err/fs"
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
	case hookCfg.ToolClaudeCode, hookCfg.ToolClaude:
		hook.InfoTool(cmd, desc.TextDesc(text.DescKeyHookClaude))

	case hookCfg.ToolCursor:
		hook.InfoTool(cmd, desc.TextDesc(text.DescKeyHookCursor))

	case hookCfg.ToolAider:
		hook.InfoTool(cmd, desc.TextDesc(text.DescKeyHookAider))

	case hookCfg.ToolCopilot:
		if writeFile {
			return WriteCopilotInstructions(cmd)
		}
		hook.InfoTool(cmd, desc.TextDesc(text.DescKeyHookCopilot))
		cmd.Println()
		content, readErr := agent.CopilotInstructions()
		if readErr != nil {
			return readErr
		}
		cmd.Print(string(content))

	case hookCfg.ToolWindsurf:
		hook.InfoTool(cmd, desc.TextDesc(text.DescKeyHookWindsurf))

	default:
		hook.InfoUnknownTool(cmd, tool)
		hook.InfoTool(cmd, desc.TextDesc(text.DescKeyHookSupportedTools))
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
	targetFile := filepath.Join(hookCfg.DirGitHub, hookCfg.FileCopilotInstructions)

	// Create .github/ directory if needed
	if err := os.MkdirAll(hookCfg.DirGitHub, fs.PermExec); err != nil {
		return ctxErr.Mkdir(hookCfg.DirGitHub, err)
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
			return ctxErr.FileWrite(targetFile, wErr)
		}
		hook.InfoCopilotMerged(cmd, targetFile)
		return nil
	}

	// File doesn't exist: create it
	if wErr := os.WriteFile(
		targetFile, instructions, fs.PermFile,
	); wErr != nil {
		return ctxErr.FileWrite(targetFile, wErr)
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

	return nil
}
