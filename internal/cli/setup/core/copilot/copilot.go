//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package copilot

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/agent"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/token"
	cfgVscode "github.com/ActiveMemory/ctx/internal/config/vscode"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	"github.com/ActiveMemory/ctx/internal/io"
	writeErr "github.com/ActiveMemory/ctx/internal/write/err"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// DeployInstructions generates .github/copilot-instructions.md.
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
func DeployInstructions(cmd *cobra.Command) error {
	targetFile := filepath.Join(cfgHook.DirGitHub, cfgHook.FileCopilotInstructions)

	// Create .github/ directory if needed
	if mkdirErr := os.MkdirAll(cfgHook.DirGitHub, fs.PermExec); mkdirErr != nil {
		return errFs.Mkdir(cfgHook.DirGitHub, mkdirErr)
	}

	// Load the copilot instructions from embedded assets
	instructions, readErr := agent.CopilotInstructions()
	if readErr != nil {
		return readErr
	}

	// Check if the file exists
	existingContent, existErr := os.ReadFile(filepath.Clean(targetFile))
	fileExists := existErr == nil

	if fileExists {
		existingStr := string(existingContent)
		if strings.Contains(existingStr, marker.CopilotStart) {
			writeSetup.InfoCopilotSkipped(cmd, targetFile)
			return nil
		}

		// File exists without ctx markers: append ctx content
		merged := existingStr + token.NewlineLF + string(instructions)
		if wErr := io.SafeWriteFile(
			targetFile, []byte(merged), fs.PermFile,
		); wErr != nil {
			return errFs.FileWrite(targetFile, wErr)
		}
		writeSetup.InfoCopilotMerged(cmd, targetFile)
		return nil
	}

	// File doesn't exist: create it
	if wErr := io.SafeWriteFile(
		targetFile, instructions, fs.PermFile,
	); wErr != nil {
		return errFs.FileWrite(targetFile, wErr)
	}
	writeSetup.InfoCopilotCreated(cmd, targetFile)

	// Also create .context/sessions/ if it doesn't exist
	sessionsDir := filepath.Join(dir.Context, dir.Sessions)
	if mkErr := os.MkdirAll(sessionsDir, fs.PermExec); mkErr != nil {
		writeErr.WarnFile(cmd, sessionsDir, mkErr)
	} else {
		writeSetup.InfoCopilotSessionsDir(cmd, sessionsDir)
	}

	writeSetup.InfoCopilotSummary(cmd)

	// Also create .vscode/mcp.json if it doesn't exist
	if err := ensureVSCodeMCP(cmd); err != nil {
		writeErr.WarnFile(cmd, cfgVscode.FileMCPJSON, err)
	}

	return nil
}
