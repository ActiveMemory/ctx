//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package preview

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

// Cmd returns the "ctx steering preview" subcommand.
//
// Returns:
//   - *cobra.Command: Configured preview subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySteeringPreview)

	return &cobra.Command{
		Use:   cmd.UseSteeringPreview,
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return Run(c, args[0])
		},
	}
}

// Run shows which steering files would be included for the given
// prompt text, respecting inclusion mode rules.
//
// Parameters:
//   - c: The cobra command for output
//   - prompt: The prompt text to match against
func Run(c *cobra.Command, prompt string) error {
	steeringDir := rc.SteeringDir()

	files, err := steering.LoadAll(steeringDir)
	if err != nil {
		return err
	}

	// Filter with no manual names — preview only shows always + auto matches.
	matched := steering.Filter(files, prompt, nil, "")

	if len(matched) == 0 {
		writeSteering.NoFilesMatch(c)
		return nil
	}

	writeSteering.PreviewHeader(c, prompt)
	for _, sf := range matched {
		tools := cfgSteering.LabelAllTools
		if len(sf.Tools) > 0 {
			tools = strings.Join(sf.Tools, token.CommaSpace)
		}
		writeSteering.PreviewEntry(
			c, sf.Name,
			string(sf.Inclusion), sf.Priority, tools,
		)
	}

	writeSteering.PreviewCount(c, len(matched))
	return nil
}
