//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	cfgTrigger "github.com/ActiveMemory/ctx/internal/config/trigger"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/trigger"
	writeTrigger "github.com/ActiveMemory/ctx/internal/write/trigger"
)

// Cmd returns the "ctx hook list" subcommand.
//
// Returns:
//   - *cobra.Command: Configured list subcommand
func Cmd() *cobra.Command {
	short, _ := desc.Command(cmd.DescKeyTriggerList)

	return &cobra.Command{
		Use:     cmd.UseTriggerList,
		Short:   short,
		Example: desc.Example(cmd.DescKeyTriggerList),
		Args:    cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			return Run(c)
		},
	}
}

// Run lists all hooks grouped by hook type with name, enabled/disabled
// status, and file path.
//
// Parameters:
//   - c: The cobra command for output
func Run(c *cobra.Command) error {
	hooksDir := rc.HooksDir()

	all, err := trigger.Discover(hooksDir)
	if err != nil {
		return err
	}

	total := 0
	for _, ht := range trigger.ValidTypes() {
		hooks := all[ht]
		if len(hooks) == 0 {
			continue
		}

		writeTrigger.TypeHeader(c, ht)
		for _, h := range hooks {
			status := cfgTrigger.StatusEnabled
			if !h.Enabled {
				status = cfgTrigger.StatusDisabled
			}
			writeTrigger.Entry(c, h.Name, status, h.Path)
			total++
		}
		writeTrigger.BlankLine(c)
	}

	if total == 0 {
		writeTrigger.NoHooksFound(c)
		return nil
	}

	writeTrigger.Count(c, total)
	return nil
}
