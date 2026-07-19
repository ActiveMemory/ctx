//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package apply

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	embedCmd "github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	embedFlag "github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the "ctx disclosure apply" command.
//
// It takes one positional argument (a canonical knowledge file) and a
// digest plan (--plan <path>, or - for stdin), moves the plan's staged
// entries into theme files, and folds their gists into the root. --json
// switches the result summary to machine-readable output.
//
// Returns:
//   - *cobra.Command: configured apply command with --plan and --json
func Cmd() *cobra.Command {
	var planPath string
	var jsonOutput bool

	short, long := desc.Command(embedCmd.DescKeyDisclosureApply)
	c := &cobra.Command{
		Use:   embedCmd.UseDisclosureApply,
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args[0], planPath, jsonOutput)
		},
	}

	flagbind.StringFlagPDefault(
		c, &planPath, cFlag.Plan, "", token.Dash,
		embedFlag.DescKeyDisclosureApplyPlan,
	)
	flagbind.BoolFlag(
		c, &jsonOutput, cFlag.JSON, embedFlag.DescKeyDisclosureApplyJSON,
	)
	return c
}
