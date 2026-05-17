//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ask

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// Cmd returns the `ctx kb ask` command.
//
// Returns:
//   - *cobra.Command: configured command.
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyKBAsk)
	return &cobra.Command{
		Use:   cmd.UseKBAsk,
		Short: short,
		Long:  long,
		Args:  cobra.ArbitraryArgs,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			q := strings.TrimSpace(strings.Join(args, token.Space))
			return Run(cobraCmd, q, args)
		},
	}
}
