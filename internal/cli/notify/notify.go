//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package notify

import (
	"strings"

	errcli "github.com/ActiveMemory/ctx/internal/err/cli"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/notify/cmd/setup"
	"github.com/ActiveMemory/ctx/internal/cli/notify/cmd/test"
	notifylib "github.com/ActiveMemory/ctx/internal/notify"
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

	short, long := assets.CommandDesc(assets.CmdDescKeyNotify)
	cmd := &cobra.Command{
		Use:   "notify [message]",
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if event == "" {
				return errcli.FlagRequired("event")
			}
			if len(args) == 0 {
				return errcli.ArgRequired("message")
			}
			message := strings.Join(args, " ")
			var ref *notifylib.TemplateRef
			if hook != "" {
				ref = notifylib.NewTemplateRef(hook, variant, nil)
			}
			return notifylib.Send(event, message, sessionID, ref)
		},
	}

	cmd.Flags().StringVarP(&event,
		"event", "e", "",
		assets.FlagDesc(assets.FlagDescKeyNotifyEvent),
	)
	cmd.Flags().StringVarP(&sessionID,
		"session-id", "s", "", assets.FlagDesc(assets.FlagDescKeyNotifySessionId),
	)
	cmd.Flags().StringVar(&hook,
		"hook", "", assets.FlagDesc(assets.FlagDescKeyNotifyHook),
	)
	cmd.Flags().StringVar(&variant,
		"variant", "", assets.FlagDesc(assets.FlagDescKeyNotifyVariant),
	)

	cmd.AddCommand(setup.Cmd())
	cmd.AddCommand(test.Cmd())

	return cmd
}
