//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx doctor" command.
//
// Flags:
//   - --json, -j: Machine-readable JSON output
//
// Returns:
//   - *cobra.Command: Configured doctor command with flags registered
func Cmd() *cobra.Command {
	short, long := desc.CommandDesc(cmd.DescKeyDoctor)
	cmd := &cobra.Command{
		Use:         cmd.DescKeyDoctor,
		Short:       short,
		Annotations: map[string]string{cli.AnnotationSkipInit: cli.AnnotationTrue},
		Long:        long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			jsonOut, _ := cmd.Flags().GetBool("json")
			return Run(cmd, jsonOut)
		},
	}
	cmd.Flags().BoolP("json", "j", false, desc.FlagDesc(flag.DescKeyDoctorJson))
	return cmd
}
