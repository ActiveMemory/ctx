//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package checkreminders

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	remindcore "github.com/ActiveMemory/ctx/internal/cli/remind/core"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// Cmd returns the "ctx system check-reminders" subcommand.
//
// Returns:
//   - *cobra.Command: Configured check-reminders subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:    "check-reminders",
		Short:  "Surface pending reminders at session start",
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCheckReminders(cmd, os.Stdin)
		},
	}
}

func runCheckReminders(cmd *cobra.Command, stdin *os.File) error {
	if !core.IsInitialized() {
		return nil
	}

	input := core.ReadInput(stdin)

	sessionID := input.SessionID
	if sessionID == "" {
		sessionID = core.SessionUnknown
	}
	if core.Paused(sessionID) > 0 {
		return nil
	}

	reminders, readErr := remindcore.ReadReminders()
	if readErr != nil {
		return nil // non-fatal: don't break session start
	}

	today := time.Now().Format("2006-01-02")
	var due []remindcore.Reminder
	for _, r := range reminders {
		if r.After == nil || *r.After <= today {
			due = append(due, r)
		}
	}

	if len(due) == 0 {
		return nil
	}

	// Build pre-formatted reminder list for the template variable
	var reminderList string
	for _, r := range due {
		reminderList += fmt.Sprintf(" [%d] %s\n", r.ID, r.Message)
	}

	fallback := reminderList +
		"\nDismiss: ctx remind dismiss <id>\n" +
		"Dismiss all: ctx remind dismiss --all"
	content := core.LoadMessage("check-reminders", "reminders",
		map[string]any{"ReminderList": reminderList}, fallback)
	if content == "" {
		return nil
	}

	msg := "IMPORTANT: Relay these reminders to the user VERBATIM before answering their question.\n\n" +
		"┌─ Reminders ──────────────────────────────────────\n"
	msg += core.BoxLines(content)
	msg += config.NudgeBoxBottom
	cmd.Println(msg)

	ref := notify.NewTemplateRef("check-reminders", "reminders",
		map[string]any{"ReminderList": reminderList})
	nudgeMsg := fmt.Sprintf("You have %d pending reminders", len(due))
	_ = notify.Send("nudge", nudgeMsg, input.SessionID, ref)
	eventlog.Append("nudge", nudgeMsg, input.SessionID, ref)

	return nil
}
