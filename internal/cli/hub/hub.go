//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/hub/cmd/peer"
	hubStatus "github.com/ActiveMemory/ctx/internal/cli/hub/cmd/status"
	"github.com/ActiveMemory/ctx/internal/cli/hub/cmd/stepdown"
	"github.com/ActiveMemory/ctx/internal/cli/parent"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the hub command with subcommands.
//
// Returns:
//   - *cobra.Command: hub with status, peer, stepdown
func Cmd() *cobra.Command {
	return parent.Cmd(
		cmd.DescKeyHub, cmd.UseHub,
		hubStatus.Cmd(),
		peer.Cmd(),
		stepdown.Cmd(),
	)
}
