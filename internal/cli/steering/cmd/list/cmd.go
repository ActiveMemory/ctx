//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	cfgSteering "github.com/ActiveMemory/ctx/internal/config/steering"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/steering"
	writeSteering "github.com/ActiveMemory/ctx/internal/write/steering"
)

// Cmd returns the "ctx steering list" subcommand.
//
// Returns:
//   - *cobra.Command: Configured list subcommand
func Cmd() *cobra.Command {
	short, _ := desc.Command(cmd.DescKeySteeringList)

	return &cobra.Command{
		Use:     cmd.UseSteeringList,
		Short:   short,
		Example: desc.Example(cmd.DescKeySteeringList),
		Args:    cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			return Run(c)
		},
	}
}

// Run lists all steering files with name, inclusion mode, priority,
// and target tools.
//
// Parameters:
//   - c: The cobra command for output
func Run(c *cobra.Command) error {
	steeringDir := rc.SteeringDir()

	files, err := steering.LoadAll(steeringDir)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		writeSteering.NoFilesFound(c)
		return nil
	}

	for _, sf := range files {
		tools := cfgSteering.LabelAllTools
		if len(sf.Tools) > 0 {
			tools = strings.Join(sf.Tools, token.CommaSpace)
		}
		writeSteering.FileEntry(
			c, sf.Name,
			string(sf.Inclusion), sf.Priority, tools,
		)
	}

	writeSteering.FileCount(c, len(files))
	return nil
}
