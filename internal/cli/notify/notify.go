//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package notify

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/notify/cmd/setup"
	"github.com/ActiveMemory/ctx/internal/cli/notify/cmd/test"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	errCli "github.com/ActiveMemory/ctx/internal/err/cli"
	"github.com/ActiveMemory/ctx/internal/flagbind"
	iNotify "github.com/ActiveMemory/ctx/internal/notify"
)

// Cmd returns the "ctx notify" parent command.
//
// Returns:
//   - *cobra.Command: Configured notify command with subcommands
func Cmd() *cobra.Command {
	var event string
	var sessionID string
	var hook string
	var variant string

	short, long := desc.Command(cmd.DescKeyNotify)
	c := &cobra.Command{
		Use:   cmd.UseNotify,
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if event == "" {
				return errCli.FlagRequired(cFlag.Event)
			}
			if len(args) == 0 {
				return errCli.ArgRequired(cFlag.Message)
			}
			message := strings.Join(args, token.Space)
			var ref *entity.TemplateRef
			if hook != "" {
				ref = entity.NewTemplateRef(hook, variant, nil)
			}
			return iNotify.Send(event, message, sessionID, ref)
		},
	}

	flagbind.StringFlagP(
		c, &event, cFlag.Event,
		cFlag.ShortEvent, flag.DescKeyNotifyEvent,
	)
	flagbind.StringFlagP(
		c, &sessionID, cFlag.SessionID,
		cFlag.ShortSessionID, flag.DescKeyNotifySessionId,
	)
	flagbind.StringFlag(c, &hook, cFlag.Hook, flag.DescKeyNotifyHook)
	flagbind.StringFlag(c, &variant, cFlag.Variant, flag.DescKeyNotifyVariant)

	c.AddCommand(setup.Cmd())
	c.AddCommand(test.Cmd())

	return c
}
