//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config"
)

// Cmd returns the "ctx why" cobra command.
//
// Returns:
//   - *cobra.Command: Configured why command with document aliases
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc("why")

	cmd := &cobra.Command{
		Use:         "why [DOCUMENT]",
		Short:       short,
		Annotations: map[string]string{config.AnnotationSkipInit: ""},
		ValidArgs:   []string{"manifesto", "about", "invariants"},
		Long:        long,
		Args:        cobra.MaximumNArgs(1),
		RunE:        Run,
	}

	return cmd
}
