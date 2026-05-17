//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package note

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// Cmd returns the `ctx kb note` command.
//
// Returns:
//   - *cobra.Command: configured command.
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyKBNote)
	return &cobra.Command{
		Use:   cmd.UseKBNote,
		Short: short,
		Long:  long,
		Args:  cobra.ArbitraryArgs,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			text := strings.TrimSpace(strings.Join(args, token.Space))
			return Run(cobraCmd, text)
		},
	}
}
