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

	cmd.Println("IMPORTANT: Relay these reminders to the user VERBATIM before answering their question.")
	cmd.Println()
	cmd.Println("┌─ Reminders ──────────────────────────────────────")
	for _, r := range due {
		cmd.Printf("│  [%d] %s\n", r.ID, r.Message)
	}
	cmd.Println("│")
	cmd.Println("│ Dismiss: ctx remind dismiss <id>")
	cmd.Println("│ Dismiss all: ctx remind dismiss --all")
	cmd.Println("└──────────────────────────────────────────────────")

	_ = notify.Send("nudge", fmt.Sprintf("You have %d pending reminders", len(due)), input.SessionID, "")

	return nil
}
