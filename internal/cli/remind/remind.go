//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package remind

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// Cmd returns the remind command with subcommands.
//
// When invoked with arguments and no subcommand, it adds a reminder.
// When invoked with no arguments, it lists all reminders.
func Cmd() *cobra.Command {
	var afterFlag string

	cmd := &cobra.Command{
		Use:   "remind [TEXT]",
		Short: "Session-scoped reminders",
		Long: `Manage session-scoped reminders stored in .context/reminders.json.

Reminders surface verbatim at session start and repeat every session until
dismissed. Use --after to gate a reminder until a specific date.

When invoked with a text argument, adds a reminder (equivalent to "remind add").
When invoked with no arguments, lists all reminders.

Subcommands:
  add      Add a reminder (default action)
  list     Show all pending reminders
  dismiss  Dismiss one or all reminders`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return runAdd(cmd, args[0], afterFlag)
			}
			return runList(cmd)
		},
	}

	cmd.Flags().StringVarP(&afterFlag, "after", "a", "", "Don't surface until this date (YYYY-MM-DD)")

	cmd.AddCommand(addCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(dismissCmd())

	return cmd
}

func addCmd() *cobra.Command {
	var afterFlag string

	cmd := &cobra.Command{
		Use:   "add TEXT",
		Short: "Add a reminder",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(cmd, args[0], afterFlag)
		},
	}

	cmd.Flags().StringVarP(&afterFlag, "after", "a", "", "Don't surface until this date (YYYY-MM-DD)")

	return cmd
}

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Show all pending reminders",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runList(cmd)
		},
	}
}

func dismissCmd() *cobra.Command {
	var allFlag bool

	cmd := &cobra.Command{
		Use:     "dismiss [ID]",
		Aliases: []string{"rm"},
		Short:   "Dismiss one or all reminders",
		RunE: func(cmd *cobra.Command, args []string) error {
			if allFlag {
				return runDismissAll(cmd)
			}
			if len(args) == 0 {
				return fmt.Errorf("provide a reminder ID or use --all")
			}
			return runDismiss(cmd, args[0])
		},
	}

	cmd.Flags().BoolVar(&allFlag, "all", false, "Dismiss all reminders")

	return cmd
}

func runAdd(cmd *cobra.Command, message, after string) error {
	reminders, err := ReadReminders()
	if err != nil {
		return err
	}

	r := Reminder{
		ID:      NextID(reminders),
		Message: message,
		Created: time.Now().UTC().Format(time.RFC3339),
	}
	if after != "" {
		if _, err := time.Parse("2006-01-02", after); err != nil {
			return fmt.Errorf("invalid date %q (expected YYYY-MM-DD)", after)
		}
		r.After = &after
	}

	reminders = append(reminders, r)
	if err := WriteReminders(reminders); err != nil {
		return err
	}

	suffix := ""
	if r.After != nil {
		suffix = fmt.Sprintf("  (after %s)", *r.After)
	}
	cmd.Printf("  + [%d] %s%s\n", r.ID, r.Message, suffix)
	return nil
}

func runList(cmd *cobra.Command) error {
	reminders, err := ReadReminders()
	if err != nil {
		return err
	}

	if len(reminders) == 0 {
		cmd.Println("No reminders.")
		return nil
	}

	today := time.Now().Format("2006-01-02")
	for _, r := range reminders {
		annotation := ""
		if r.After != nil {
			if *r.After > today {
				annotation = fmt.Sprintf("  (after %s, not yet due)", *r.After)
			}
		}
		cmd.Printf("  [%d] %s%s\n", r.ID, r.Message, annotation)
	}

	return nil
}

func runDismiss(cmd *cobra.Command, idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid ID %q", idStr)
	}

	reminders, err := ReadReminders()
	if err != nil {
		return err
	}

	found := -1
	for i, r := range reminders {
		if r.ID == id {
			found = i
			break
		}
	}

	if found < 0 {
		return fmt.Errorf("no reminder with ID %d", id)
	}

	cmd.Printf("  - [%d] %s\n", reminders[found].ID, reminders[found].Message)
	reminders = append(reminders[:found], reminders[found+1:]...)
	return WriteReminders(reminders)
}

func runDismissAll(cmd *cobra.Command) error {
	reminders, err := ReadReminders()
	if err != nil {
		return err
	}

	if len(reminders) == 0 {
		cmd.Println("No reminders.")
		return nil
	}

	for _, r := range reminders {
		cmd.Printf("  - [%d] %s\n", r.ID, r.Message)
	}
	cmd.Printf("Dismissed %d reminders.\n", len(reminders))

	return WriteReminders([]Reminder{})
}
