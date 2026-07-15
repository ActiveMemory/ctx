//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package revoke

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreRevoke "github.com/ActiveMemory/ctx/internal/cli/hub/core/revoke"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/ActiveMemory/ctx/internal/config/env"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	errHub "github.com/ActiveMemory/ctx/internal/err/hub"
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
			// Admin token: --token flag takes precedence, then
			// the CTX_HUB_ADMIN_TOKEN environment variable.
			token := adminToken
			if token == "" {
				token = os.Getenv(env.HubAdmin)
			}
			if token == "" {
				cobraCmd.SilenceUsage = true
				return errHub.AdminTokenRequired()
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
