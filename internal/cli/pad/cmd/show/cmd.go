//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
)

// Cmd returns the pad show subcommand.
//
// Outputs the raw text of entry N (1-based) with no numbering prefix.
// Designed for pipe composability:
//
//	ctx pad edit 1 --append "$(ctx pad show 3)"
//
// Returns:
//   - *cobra.Command: Configured show subcommand
func Cmd() *cobra.Command {
	var outPath string

	short, long := assets.CommandDesc(assets.CmdDescKeyPadShow)
	cmd := &cobra.Command{
		Use:   "show N",
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				return ctxerr.InvalidIndex(args[0])
			}
			return Run(cmd, n, outPath)
		},
	}

	cmd.Flags().StringVar(&outPath,
		"out", "", assets.FlagDesc(assets.FlagDescKeyPadShowOut),
	)

	return cmd
}
