//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package doctor provides the "ctx doctor" command for structural
// health checks across context, hooks, and configuration.
package doctor

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/doctor/cmd/root"
)

// Cmd returns the "ctx doctor" command.
//
// Returns:
//   - *cobra.Command: The doctor command with subcommands registered
func Cmd() *cobra.Command {
	return root.Cmd()
}
