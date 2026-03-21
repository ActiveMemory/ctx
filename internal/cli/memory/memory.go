//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/memory/cmd/diff"
	"github.com/ActiveMemory/ctx/internal/cli/memory/cmd/importer"
	"github.com/ActiveMemory/ctx/internal/cli/memory/cmd/publish"
	"github.com/ActiveMemory/ctx/internal/cli/memory/cmd/status"
	"github.com/ActiveMemory/ctx/internal/cli/memory/cmd/sync"
	"github.com/ActiveMemory/ctx/internal/cli/memory/cmd/unpublish"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx memory" parent command.
func Cmd() *cobra.Command {
	short, long := desc.CommandDesc(cmd.DescKeyMemory)
	c := &cobra.Command{
		Use:   cmd.UseMemory,
		Short: short,
		Long:  long,
	}

	c.AddCommand(
		sync.Cmd(),
		status.Cmd(),
		diff.Cmd(),
		importer.Cmd(),
		publish.Cmd(),
		unpublish.Cmd(),
	)

	return c
}
