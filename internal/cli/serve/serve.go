//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package serve

import (
	"github.com/spf13/cobra"
)

// Cmd returns the serve command.
//
// Serves a static site by invoking zensical serve on the specified directory.
//
// Returns:
//   - *cobra.Command: The serve command
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve [directory]",
		Short: "Serve a static site locally via zensical",
		Long: `Serve a static site using zensical.

If no directory is specified, serves the journal site (.context/journal-site).

Requires zensical to be installed:
  pipx install zensical

Examples:
  ctx serve                           # Serve journal site
  ctx serve .context/journal-site     # Serve specific directory
  ctx serve ./docs                    # Serve docs folder`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return runServe(args)
		},
	}

	return cmd
}
