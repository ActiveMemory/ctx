//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

import (
	"github.com/spf13/cobra"

	agentRoot "github.com/ActiveMemory/ctx/internal/cli/agent/cmd/root"
)

// Cmd returns the "ctx agent" command for generating AI-ready context packets.
//
// Returns:
//   - *cobra.Command: The agent command with subcommands registered
func Cmd() *cobra.Command {
	return agentRoot.Cmd()
}
