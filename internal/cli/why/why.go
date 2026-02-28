//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package why

import (
	"github.com/spf13/cobra"
)

// docAliases maps user-facing names to embedded asset names.
var docAliases = map[string]string{
	"manifesto":  "manifesto",
	"about":      "about",
	"invariants": "design-invariants",
}

// docOrder defines the display order for the interactive menu.
var docOrder = []struct {
	alias string
	label string
}{
	{"manifesto", "The ctx Manifesto"},
	{"about", "About ctx"},
	{"invariants", "Design Invariants"},
}

// Cmd returns the "ctx why" cobra command.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "why [DOCUMENT]",
		Short:     "Read the philosophy behind ctx",
		ValidArgs: []string{"manifesto", "about", "invariants"},
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
		RunE: runWhy,
	}

	return cmd
}
