//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ai

import (
	"github.com/spf13/cobra"

	aiExtract "github.com/ActiveMemory/ctx/internal/cli/ai/cmd/extract"
	aiPing "github.com/ActiveMemory/ctx/internal/cli/ai/cmd/ping"
	"github.com/ActiveMemory/ctx/internal/cli/parent"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the `ctx ai` parent command with its
// subcommands wired in.
//
// Returns:
//   - *cobra.Command: parent command with ping (and
//     future extract) registered as subcommands.
func Cmd() *cobra.Command {
	return parent.Cmd(
		cmd.DescKeyAI, cmd.UseAI,
		aiPing.Cmd(),
		aiExtract.Cmd(),
	)
}
