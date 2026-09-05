//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package codex

import (
	"github.com/spf13/cobra"

	coreAgents "github.com/ActiveMemory/ctx/internal/cli/setup/core/agents"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	writeErr "github.com/ActiveMemory/ctx/internal/write/err"
)

// deployAgents deploys AGENTS.md, warning (not failing) on error.
//
// Parameters:
//   - cmd: Cobra command for output messages
func deployAgents(cmd *cobra.Command) {
	if agentsErr := coreAgents.Deploy(cmd); agentsErr != nil {
		writeErr.WarnFile(cmd, cfgHook.FileAgentsMd, agentsErr)
	}
}
