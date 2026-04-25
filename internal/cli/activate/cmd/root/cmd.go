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

// Cmd returns the `ctx activate` cobra command.
//
// Args-free under the single-source-anchor model
// (specs/single-source-context-anchor.md). Activation is always
// project-local discovery via [rc.ScanCandidates] from CWD; the
// explicit-path mode that previously accepted an argument was
// removed because hub-client / hub-server scenarios store at
// `~/.ctx/hub-data/` and never read `.context/` directly, so they
// activate from the project root like everyone else.
//
// One flag remains:
//
//	--shell <name>   override auto-detection (defaults to $SHELL).
//
// # Stdout discipline (critical)
//
// Activate's stdout is consumed by `eval "$(ctx activate)"`. Every
// byte must be either valid shell or empty. Usage / Flags /
// Examples blocks must NEVER reach stdout, because cobra's
// Examples for this command literally contain
// `eval "$(ctx activate)"`, which would re-execute activate inside
// the eval and trigger an infinite loop on any error path.
//
// SilenceUsage is therefore set unconditionally below (rather than
// only after a return) so cobra renders only the error to stderr
// when something fails. SilenceErrors stays at the root level so
// errors keep going to stderr (visible to the user) without being
// captured by the eval.
//
// Returns:
//   - *cobra.Command: configured activate command.
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyActivate)
	c := &cobra.Command{
		Use:     cmd.UseActivate,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyActivate),
		Args:    cobra.NoArgs,
		// Exempt from the global init / require-context-dir checks:
		// activate's whole purpose is to help the user declare the
		// context directory in the first place.
		Annotations: map[string]string{cli.AnnotationSkipInit: cli.AnnotationTrue},
		// See the Stdout discipline note above. Without this, an
		// error path (multi-candidate, no-candidates, etc.) prints
		// Usage+Examples to stdout, gets captured by `$(...)`, and
		// the embedded `eval "$(ctx activate)"` example re-runs the
		// command. Loop.
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			shell, _ := cmd.Flags().GetString(cFlag.Shell)
			return Run(cmd, shell)
		},
	}
	c.Flags().String(cFlag.Shell, "",
		desc.Flag(embedFlag.DescKeyActivateShell),
	)
	return c
}
