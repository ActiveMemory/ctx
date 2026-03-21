//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx why" cobra command.
//
// Returns:
//   - *cobra.Command: Configured why command with document aliases
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyWhy)

	cmd := &cobra.Command{
		Use:         cmd.UseWhy,
		Short:       short,
		Annotations: map[string]string{cli.AnnotationSkipInit: ""},
		ValidArgs:   []string{"manifesto", "about", "invariants"},
		Long:        long,
		Args:        cobra.MaximumNArgs(1),
		RunE:        Run,
	}

	return cmd
}
