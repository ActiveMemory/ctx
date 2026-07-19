//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/disclosure/cmd/apply"
	"github.com/ActiveMemory/ctx/internal/cli/disclosure/cmd/inspect"
	"github.com/ActiveMemory/ctx/internal/cli/parent"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the `ctx disclosure` parent command, which inspects and
// digests the canonical knowledge files under progressive disclosure.
//
// Returns:
//   - *cobra.Command: disclosure parent with the inspect and apply
//     subcommands
func Cmd() *cobra.Command {
	return parent.Cmd(cmd.DescKeyDisclosure, cmd.UseDisclosure,
		inspect.Cmd(),
		apply.Cmd(),
	)
}
