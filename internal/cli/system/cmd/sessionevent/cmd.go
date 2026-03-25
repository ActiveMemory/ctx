//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sessionevent

import (
	"fmt"

	"github.com/spf13/cobra"

	coreState "github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/log"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// Cmd returns the "ctx system session-event" subcommand.
//
// Returns:
//   - *cobra.Command: Configured session-event subcommand
func Cmd() *cobra.Command {
	var eventType string
	var caller string

	cmd := &cobra.Command{
		Use:   "session-event",
		Short: "Record session start or end",
		Long: `Records a session lifecycle event (start or end) to the event log.
Called by editor integrations when a workspace is opened or closed.

Examples:
  ctx system session-event --type start --caller vscode
  ctx system session-event --type end --caller vscode`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runSessionEvent(cmd, eventType, caller)
		},
	}

	cmd.Flags().StringVar(&eventType, "type", "", "Event type: start or end")
	cmd.Flags().StringVar(&caller, "caller", "", "Calling editor (e.g., vscode)")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("caller")

	return cmd
}

func runSessionEvent(cmd *cobra.Command, eventType, caller string) error {
	if !coreState.Initialized() {
		return nil
	}

	if eventType != "start" && eventType != "end" {
		return fmt.Errorf("--type must be 'start' or 'end', got %q", eventType)
	}

	msg := fmt.Sprintf("session-%s: %s", eventType, caller)
	ref := notify.NewTemplateRef("session-event", eventType,
		map[string]any{"Caller": caller})

	log.AppendEvent("session", msg, "", ref)
	_ = notify.Send("session", msg, "", ref)

	cmd.Println(msg)
	return nil
}
