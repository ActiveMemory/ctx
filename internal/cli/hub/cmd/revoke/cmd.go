//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package revoke

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreAdmin "github.com/ActiveMemory/ctx/internal/cli/hub/core/admin"
	coreRevoke "github.com/ActiveMemory/ctx/internal/cli/hub/core/revoke"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the hub revoke subcommand.
//
// Returns:
//   - *cobra.Command: The revoke subcommand
func Cmd() *cobra.Command {
	var adminToken string

	short, long := desc.Command(cmd.DescKeyHubRevoke)

	c := &cobra.Command{
		Use:     cmd.UseHubRevoke,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyHubRevoke),
		Args:    cobra.ExactArgs(1),
		// Hub stores at ~/.ctx/hub-data/, not .context/.
		// Spec: specs/single-source-context-anchor.md.
		Annotations: map[string]string{cli.AnnotationSkipInit: cli.AnnotationTrue},
		RunE: func(
			cobraCmd *cobra.Command, args []string,
		) error {
			token, tokenErr := coreAdmin.Token(adminToken)
			if tokenErr != nil {
				cobraCmd.SilenceUsage = true
				return tokenErr
			}
			return coreRevoke.Run(cobraCmd, args[0], token)
		},
	}

	flagbind.StringFlag(
		c, &adminToken,
		cFlag.Token, flag.DescKeyHubRevokeAuth,
	)

	return c
}
