//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package notify

import (
	"fmt"
	"strings"

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

	short, long := assets.CommandDesc("notify")
	cmd := &cobra.Command{
		Use:   "notify [message]",
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if event == "" {
				return fmt.Errorf("required flag \"event\" not set")
			}
			if len(args) == 0 {
				return fmt.Errorf("message argument is required")
			}
			message := strings.Join(args, " ")
			var ref *notifylib.TemplateRef
			if hook != "" {
				ref = notifylib.NewTemplateRef(hook, variant, nil)
			}
			return notifylib.Send(event, message, sessionID, ref)
		},
	}

	cmd.Flags().StringVarP(&event, "event", "e", "", "Event name (required)")
	cmd.Flags().StringVarP(&sessionID, "session-id", "s", "", "Session ID (optional)")
	cmd.Flags().StringVar(&hook, "hook", "", "Hook name for structured detail (optional)")
	cmd.Flags().StringVar(&variant, "variant", "", "Template variant for structured detail (optional)")

	cmd.AddCommand(setup.Cmd())
	cmd.AddCommand(test.Cmd())

	return cmd
}
