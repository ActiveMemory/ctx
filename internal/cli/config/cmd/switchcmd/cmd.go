//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package switchcmd

import (
	"github.com/spf13/cobra"

	internalConfig "github.com/ActiveMemory/ctx/internal/config"

	"github.com/ActiveMemory/ctx/internal/cli/config/core"
)

// Cmd returns the "ctx config switch" subcommand.
//
// Returns:
//   - *cobra.Command: Configured switch subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:         "switch [dev|base]",
		Short:       "Switch .ctxrc profile",
		Annotations: map[string]string{internalConfig.AnnotationSkipInit: ""},
		Long: `Switch between .ctxrc configuration profiles.

With no argument, toggles between dev and base.
Accepts "prod" as an alias for "base".

Source files (.ctxrc.base, .ctxrc.dev) are committed to git.
The working copy (.ctxrc) is gitignored.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root, rootErr := core.GitRoot()
			if rootErr != nil {
				return rootErr
			}
			return RunSwitch(cmd, root, args)
		},
	}
}
