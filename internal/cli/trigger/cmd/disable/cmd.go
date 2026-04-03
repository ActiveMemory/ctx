//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disable

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	errTrigger "github.com/ActiveMemory/ctx/internal/err/trigger"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/trigger"
	writeTrigger "github.com/ActiveMemory/ctx/internal/write/trigger"
)

// Cmd returns the "ctx hook disable" subcommand.
//
// Returns:
//   - *cobra.Command: Configured disable subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyTriggerDisable)

	return &cobra.Command{
		Use:   cmd.UseTriggerDisable,
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return Run(c, args[0])
		},
	}
}

// Run disables a hook by removing the executable permission bit.
//
// Parameters:
//   - c: The cobra command for output
//   - name: The hook name to disable
func Run(c *cobra.Command, name string) error {
	hooksDir := rc.HooksDir()

	h, findErr := trigger.FindByName(hooksDir, name)
	if findErr != nil {
		return findErr
	}

	if h == nil {
		return errTrigger.NotFound(name)
	}

	fi, statErr := os.Stat(h.Path)
	if statErr != nil {
		return errTrigger.Stat(statErr)
	}

	// Remove executable permission bit for user, group, and other.
	newMode := fi.Mode() &^ 0o111
	if chmodErr := os.Chmod(h.Path, newMode); chmodErr != nil {
		return errTrigger.Chmod(chmodErr)
	}

	writeTrigger.Disabled(c, name, h.Path)
	return nil
}
