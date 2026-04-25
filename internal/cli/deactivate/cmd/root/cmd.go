//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	embedFlag "github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the `ctx deactivate` cobra command.
//
// Accepts one flag:
//
//	--shell <name>   override auto-detection (defaults to $SHELL).
//
// # Stdout discipline (critical)
//
// Same eval-recursion hazard as `ctx activate`: stdout is consumed
// by `eval "$(ctx deactivate)"`, so cobra must never print Usage /
// Flags / Examples on stdout (the Examples block contains the eval
// invocation literally). [SilenceUsage] is set unconditionally
// below; errors keep going to stderr via the root [SilenceErrors]
// settings.
//
// Returns:
//   - *cobra.Command: configured deactivate command.
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyDeactivate)
	c := &cobra.Command{
		Use:     cmd.UseDeactivate,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyDeactivate),
		Args:    cobra.NoArgs,
		// Exempt from the global init / require-context-dir checks:
		// `unset CTX_DIR` must work regardless of current state.
		Annotations: map[string]string{cli.AnnotationSkipInit: cli.AnnotationTrue},
		// See the Stdout discipline note above.
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			shell, _ := cmd.Flags().GetString(cFlag.Shell)
			return Run(cmd, shell)
		},
	}
	c.Flags().String(cFlag.Shell, "",
		desc.Flag(embedFlag.DescKeyActivateShell),
	)
	return c
}
