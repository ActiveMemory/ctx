//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx doctor" command.
//
// Flags:
//   - --json, -j: Machine-readable JSON output
//
// Returns:
//   - *cobra.Command: Configured doctor command with flags registered
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(assets.CmdDescKeyDoctor)
	cmd := &cobra.Command{
		Use:         "doctor",
		Short:       short,
		Annotations: map[string]string{cli.AnnotationSkipInit: cli.AnnotationTrue},
		Long:        long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			jsonOut, _ := cmd.Flags().GetBool("json")
			return Run(cmd, jsonOut)
		},
	}
	cmd.Flags().BoolP("json", "j", false, assets.FlagDesc(assets.FlagDescKeyDoctorJson))
	return cmd
}
