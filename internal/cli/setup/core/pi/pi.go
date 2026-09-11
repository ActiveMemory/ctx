//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pi

import (
	"path/filepath"

	"github.com/spf13/cobra"

	coreAgents "github.com/ActiveMemory/ctx/internal/cli/setup/core/agents"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	writeErr "github.com/ActiveMemory/ctx/internal/write/err"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Deploy generates all Pi integration files.
//
// Creates the .pi/extensions/ctx.ts extension file, deploys
// AGENTS.md with shared instructions, and copies ctx skills to
// .pi/skills/.
//
// Refreshes stale managed files and skips ones that already match.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if extension deployment fails (other errors are
//     warned but do not halt deployment)
func Deploy(cmd *cobra.Command) error {
	if extErr := deployExtension(cmd); extErr != nil {
		return extErr
	}

	if agentsErr := coreAgents.Deploy(cmd); agentsErr != nil {
		writeErr.WarnFile(
			cmd, cfgHook.FileAgentsMd, agentsErr,
		)
	}

	if skillErr := deploySkills(cmd); skillErr != nil {
		skillsBase := filepath.Join(
			cfgHook.DirPi, cfgHook.DirPiSkills,
		)
		writeErr.WarnFile(cmd, skillsBase, skillErr)
	}

	writeSetup.InfoPiSummary(cmd)
	return nil
}
