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
	cfgWhy "github.com/ActiveMemory/ctx/internal/config/why"
)

// Cmd returns the "ctx why" cobra command.
//
// Returns:
//   - *cobra.Command: Configured why command with document aliases
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyWhy)

	c := &cobra.Command{
		Use:         cmd.UseWhy,
		Short:       short,
		Annotations: map[string]string{cli.AnnotationSkipInit: ""},
		ValidArgs: []string{
			cfgWhy.DocManifesto, cfgWhy.DocAbout, cfgWhy.DocInvariants,
		},
		Long: long,
		Args: cobra.MaximumNArgs(1),
		RunE: Run,
	}

	return c
}
