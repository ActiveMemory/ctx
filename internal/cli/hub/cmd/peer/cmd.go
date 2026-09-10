//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package peer

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreAdmin "github.com/ActiveMemory/ctx/internal/cli/hub/core/admin"
	corePeer "github.com/ActiveMemory/ctx/internal/cli/hub/core/peer"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the hub peer subcommand.
//
// Returns:
//   - *cobra.Command: The peer subcommand
func Cmd() *cobra.Command {
	var adminToken string

	short, long := desc.Command(cmd.DescKeyHubPeer)

	c := &cobra.Command{
		Use:     cmd.UseHubPeer,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyHubPeer),
		Args:    cobra.ExactArgs(2),
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
			return corePeer.Run(
				cobraCmd, args[0], args[1], token,
			)
		},
	}

	flagbind.StringFlag(
		c, &adminToken,
		cFlag.Token, flag.DescKeyHubPeerAuth,
	)

	return c
}
