//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/config/core/profile"
	cfgCli "github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx config status" subcommand.
//
// Returns:
//   - *cobra.Command: Configured status subcommand
func Cmd() *cobra.Command {
	short, _ := desc.Command(cmd.DescKeyConfigStatus)

	return &cobra.Command{
		Use:         cmd.UseConfigStatus,
		Short:       short,
		Annotations: map[string]string{cfgCli.AnnotationSkipInit: ""},
		Args:        cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			root, rootErr := profile.GitRoot()
			if rootErr != nil {
				return rootErr
			}
			return Run(cmd, root)
		},
	}
}
