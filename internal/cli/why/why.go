//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package why

import (
	"github.com/spf13/cobra"

	whyroot "github.com/ActiveMemory/ctx/internal/cli/why/cmd/root"
	"github.com/ActiveMemory/ctx/internal/config"
)

// Cmd returns the "ctx why" cobra command.
//
// Returns:
//   - *cobra.Command: Configured why command with document aliases
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:         "why [DOCUMENT]",
		Short:       "Read the philosophy behind ctx",
		Annotations: map[string]string{config.AnnotationSkipInit: ""},
		ValidArgs:   []string{"manifesto", "about", "invariants"},
		Long: `Surface ctx's philosophy documents in the terminal.

Documents:
  manifesto    The ctx Manifesto — creation, not code
  about        About ctx — what it is and why it exists
  invariants   Design invariants — properties that must hold

Usage:
  ctx why              Interactive numbered menu
  ctx why manifesto    Show the manifesto directly
  ctx why about        Show the about page
  ctx why invariants   Show the design invariants`,
		Args: cobra.MaximumNArgs(1),
		RunE: whyroot.Run,
	}

	return cmd
}
