//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	internalConfig "github.com/ActiveMemory/ctx/internal/config"

	"github.com/ActiveMemory/ctx/internal/cli/config/core"
)

// Cmd returns the "ctx config status" subcommand.
//
// Returns:
//   - *cobra.Command: Configured status subcommand
func Cmd() *cobra.Command {
	short, _ := assets.CommandDesc("config.status")

	return &cobra.Command{
		Use:         "status",
		Short:       short,
		Annotations: map[string]string{internalConfig.AnnotationSkipInit: ""},
		Args:        cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			root, rootErr := core.GitRoot()
			if rootErr != nil {
				return rootErr
			}
			return Run(cmd, root)
		},
	}
}
