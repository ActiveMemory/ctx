//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package subscribe

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreSub "github.com/ActiveMemory/ctx/internal/cli/connect/core/subscribe"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the connect subscribe subcommand.
//
// Returns:
//   - *cobra.Command: The subscribe subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyConnectSubscribe)

	return &cobra.Command{
		Use:     cmd.UseConnectSubscribe,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyConnectSubscribe),
		Args:    cobra.MinimumNArgs(1),
		RunE:    coreSub.Run,
	}
}
