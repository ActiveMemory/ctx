//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package notify

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	errcli "github.com/ActiveMemory/ctx/internal/err/cli"
	"github.com/spf13/cobra"

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

	short, long := desc.CommandDesc(cmd.DescKeyNotify)
	cmd := &cobra.Command{
		Use:   cmd.DescKeyNotify + " [message]",
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
		desc.FlagDesc(flag.DescKeyNotifyEvent),
	)
	cmd.Flags().StringVarP(&sessionID,
		"session-id", "s", "", desc.FlagDesc(flag.DescKeyNotifySessionId),
	)
	cmd.Flags().StringVar(&hook,
		"hook", "", desc.FlagDesc(flag.DescKeyNotifyHook),
	)
	cmd.Flags().StringVar(&variant,
		"variant", "", desc.FlagDesc(flag.DescKeyNotifyVariant),
	)

	cmd.AddCommand(setup.Cmd())
	cmd.AddCommand(test.Cmd())

	return cmd
}
