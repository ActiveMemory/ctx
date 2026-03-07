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

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/write"
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
	case "claude-code", "claude":
		write.InfoHookTool(cmd, assets.TextDesc("hook.claude"))

	case "cursor":
		write.InfoHookTool(cmd, assets.TextDesc("hook.cursor"))

	case "aider":
		write.InfoHookTool(cmd, assets.TextDesc("hook.aider"))

	case "copilot":
		if writeFile {
			return WriteCopilotInstructions(cmd)
		}
		write.InfoHookTool(cmd, assets.TextDesc("hook.copilot"))
		cmd.Println()
		content, readErr := assets.CopilotInstructions()
		if readErr != nil {
			return readErr
		}
		cmd.Print(string(content))

	case "windsurf":
		write.InfoHookTool(cmd, assets.TextDesc("hook.windsurf"))

	default:
		write.InfoHookUnknownTool(cmd, tool)
		write.InfoHookTool(cmd, assets.TextDesc("hook.supported-tools"))
		return ctxerr.UnsupportedTool(tool)
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
	targetDir := ".github"
	targetFile := filepath.Join(targetDir, "copilot-instructions.md")

	// Create .github/ directory if needed
	if err := os.MkdirAll(targetDir, config.PermExec); err != nil {
		return ctxerr.Mkdir(targetDir, err)
	}

	// Load the copilot instructions from embedded assets
	instructions, readErr := assets.CopilotInstructions()
	if readErr != nil {
		return readErr
	}

	// Check if file exists
	existingContent, err := os.ReadFile(filepath.Clean(targetFile))
	fileExists := err == nil

	if fileExists {
		existingStr := string(existingContent)
		if strings.Contains(existingStr, config.CopilotMarkerStart) {
			write.InfoHookCopilotSkipped(cmd, targetFile)
			return nil
		}

		// File exists without ctx markers: append ctx content
		merged := existingStr + config.NewlineLF + string(instructions)
		if writeErr := os.WriteFile(targetFile, []byte(merged), config.PermFile); writeErr != nil {
			return ctxerr.FileWrite(targetFile, writeErr)
		}
		write.InfoHookCopilotMerged(cmd, targetFile)
		return nil
	}

	// File doesn't exist: create it
	if writeErr := os.WriteFile(
		targetFile, instructions, config.PermFile,
	); writeErr != nil {
		return ctxerr.FileWrite(targetFile, writeErr)
	}
	write.InfoHookCopilotCreated(cmd, targetFile)

	// Also create .context/sessions/ if it doesn't exist
	sessionsDir := filepath.Join(config.DirContext, config.DirSessions)
	if mkErr := os.MkdirAll(sessionsDir, config.PermExec); mkErr != nil {
		write.WarnFileErr(cmd, sessionsDir, mkErr)
	} else {
		write.InfoHookCopilotSessionsDir(cmd, sessionsDir)
	}

	write.InfoHookCopilotSummary(cmd)

	return nil
}
