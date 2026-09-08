//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stepdown

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreAdmin "github.com/ActiveMemory/ctx/internal/cli/hub/core/admin"
	coreStep "github.com/ActiveMemory/ctx/internal/cli/hub/core/stepdown"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the hub stepdown subcommand.
//
// Returns:
//   - *cobra.Command: The stepdown subcommand
func Cmd() *cobra.Command {
	var adminToken string

	short, long := desc.Command(cmd.DescKeyHubStepdown)

	c := &cobra.Command{
		Use:     cmd.UseHubStepdown,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyHubStepdown),
		Args:    cobra.NoArgs,
		// Hub stores at ~/.ctx/hub-data/, not .context/.
		// Spec: specs/single-source-context-anchor.md.
		Annotations: map[string]string{cli.AnnotationSkipInit: cli.AnnotationTrue},
		RunE: func(
			cobraCmd *cobra.Command, _ []string,
		) error {
			token, tokenErr := coreAdmin.Token(adminToken)
			if tokenErr != nil {
				cobraCmd.SilenceUsage = true
				return tokenErr
			}
			return coreStep.Run(cobraCmd, token)
		},
	}

	flagbind.StringFlag(
		c, &adminToken,
		cFlag.Token, flag.DescKeyHubStepdownAuth,
	)

	return c
}
