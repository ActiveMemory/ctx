//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/remind"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// checkRemindersCmd returns the "ctx system check-reminders" command.
//
// Surfaces pending reminders at session start via VERBATIM relay.
// No throttle — reminders fire every session until dismissed.
func checkRemindersCmd() *cobra.Command {
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
	if !isInitialized() {
		return nil
	}

	input := readInput(stdin)

	reminders, err := remind.ReadReminders()
	if err != nil {
		return nil // non-fatal: don't break session start
	}

	today := time.Now().Format("2006-01-02")
	var due []remind.Reminder
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
	content := loadMessage("check-reminders", "reminders",
		map[string]any{"ReminderList": reminderList}, fallback)
	if content == "" {
		return nil
	}

	msg := "IMPORTANT: Relay these reminders to the user VERBATIM before answering their question.\n\n" +
		"┌─ Reminders ──────────────────────────────────────\n"
	msg += boxLines(content)
	msg += "└──────────────────────────────────────────────────"
	cmd.Println(msg)

	_ = notify.Send("nudge", fmt.Sprintf("You have %d pending reminders", len(due)), input.SessionID, msg)

	return nil
}
