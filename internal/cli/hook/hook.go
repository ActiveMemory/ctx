//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/event"
	"github.com/ActiveMemory/ctx/internal/cli/message"
	"github.com/ActiveMemory/ctx/internal/cli/notify"
	"github.com/ActiveMemory/ctx/internal/cli/pause"
	"github.com/ActiveMemory/ctx/internal/cli/resume"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx hook" parent command.
//
// Consolidates hook-related user-facing commands: message,
// notify, pause, resume, and event.
//
// Returns:
//   - *cobra.Command: Parent command with hook subcommands
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyHook)
	c := &cobra.Command{
		Use:     cmd.UseHook,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyHook),
	}
	c.AddCommand(
		event.Cmd(),
		message.Cmd(),
		notify.Cmd(),
		pause.Cmd(),
		resume.Cmd(),
	)
	return c
}
