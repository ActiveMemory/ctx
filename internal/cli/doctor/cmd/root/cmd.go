//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the "ctx doctor" command.
//
// Flags:
//   - --json, -j: Machine-readable JSON output
//
// Returns:
//   - *cobra.Command: Configured doctor command with flags registered
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyDoctor)
	c := &cobra.Command{
		Use:         cmd.UseDoctor,
		Short:       short,
		Annotations: map[string]string{cli.AnnotationSkipInit: cli.AnnotationTrue},
		Long:        long,
		Example:     desc.Example(cmd.DescKeyDoctor),
		RunE: func(cmd *cobra.Command, _ []string) error {
			jsonOut, _ := cmd.Flags().GetBool(cFlag.JSON)
			return Run(cmd, jsonOut)
		},
	}
	flagbind.BoolFlagShort(
		c, cFlag.JSON, cFlag.ShortJSON,
		flag.DescKeyDoctorJson,
	)
	return c
}
