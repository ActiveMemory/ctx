//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package switchcmd

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/config/core/profile"
	cfgCli "github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx config switch" subcommand.
//
// Returns:
//   - *cobra.Command: Configured switch subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyConfigSwitch)

	return &cobra.Command{
		Use:         cmd.UseConfigSwitch,
		Short:       short,
		Annotations: map[string]string{cfgCli.AnnotationSkipInit: ""},
		Long:        long,
		Args:        cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root, rootErr := profile.GitRoot()
			if rootErr != nil {
				return rootErr
			}
			return Run(cmd, root, args)
		},
	}
}
