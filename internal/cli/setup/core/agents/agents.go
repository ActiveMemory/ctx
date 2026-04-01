//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agents

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/agent"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Deploy generates AGENTS.md in the project root.
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
func Deploy(cmd *cobra.Command) error {
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
		if strings.Contains(existingStr, marker.AgentsStart) {
			writeSetup.InfoAgentsSkipped(cmd, targetFile)
			return nil
		}

		// File exists without ctx markers: append ctx content
		merged := existingStr + token.NewlineLF + string(agentsContent)
		if wErr := os.WriteFile(targetFile, []byte(merged), fs.PermFile); wErr != nil { //nolint:gosec // targetFile from known tool config path
			return errFs.FileWrite(targetFile, wErr)
		}
		writeSetup.InfoAgentsMerged(cmd, targetFile)
		return nil
	}

	// File doesn't exist: create it
	if wErr := os.WriteFile(targetFile, agentsContent, fs.PermFile); wErr != nil {
		return errFs.FileWrite(targetFile, wErr)
	}
	writeSetup.InfoAgentsCreated(cmd, targetFile)

	writeSetup.InfoAgentsSummary(cmd)
	return nil
}
