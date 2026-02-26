//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package notify implements the "ctx notify" CLI command for webhook
// notifications.
package notify

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	notifylib "github.com/ActiveMemory/ctx/internal/notify"
)

// Cmd returns the "ctx notify" parent command.
func Cmd() *cobra.Command {
	var event string
	var sessionID string
	var detail string

	cmd := &cobra.Command{
		Use:   "notify [message]",
		Short: "Send a webhook notification",
		Long: `Send a fire-and-forget webhook notification.

Requires a configured webhook URL (see "ctx notify setup").
Silent noop when no webhook is configured or the event is filtered.

Examples:
  ctx notify --event loop "Loop completed after 5 iterations"
  ctx notify -e nudge -s session-abc "Context checkpoint at prompt #20"
  ctx notify -e relay -d "full hook output here" "Short summary"`,
		Args: cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if event == "" {
				return fmt.Errorf("required flag \"event\" not set")
			}
			if len(args) == 0 {
				return fmt.Errorf("message argument is required")
			}
			message := strings.Join(args, " ")
			return notifylib.Send(event, message, sessionID, detail)
		},
	}

	cmd.Flags().StringVarP(&event, "event", "e", "", "Event name (required)")
	cmd.Flags().StringVarP(&sessionID, "session-id", "s", "", "Session ID (optional)")
	cmd.Flags().StringVarP(&detail, "detail", "d", "", "Full payload detail (optional, truncated to 1000 chars)")

	cmd.AddCommand(setupCmd())
	cmd.AddCommand(testCmd())

	return cmd
}
