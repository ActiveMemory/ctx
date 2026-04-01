//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/trace/cmd/collect"
	"github.com/ActiveMemory/ctx/internal/cli/trace/cmd/file"
	"github.com/ActiveMemory/ctx/internal/cli/trace/cmd/hook"
	"github.com/ActiveMemory/ctx/internal/cli/trace/cmd/show"
	"github.com/ActiveMemory/ctx/internal/cli/trace/cmd/tag"
)

// Cmd returns the trace command with all registered subcommands.
//
// Returns:
//   - *cobra.Command: The trace command
func Cmd() *cobra.Command {
	c := show.Cmd()
	c.AddCommand(collect.Cmd())
	c.AddCommand(file.Cmd())
	c.AddCommand(hook.Cmd())
	c.AddCommand(tag.Cmd())
	return c
}
