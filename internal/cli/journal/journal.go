//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/journal/cmd/importer"
	"github.com/ActiveMemory/ctx/internal/cli/journal/cmd/lock"
	"github.com/ActiveMemory/ctx/internal/cli/journal/cmd/obsidian"
	journalSchema "github.com/ActiveMemory/ctx/internal/cli/journal/cmd/schema"
	"github.com/ActiveMemory/ctx/internal/cli/journal/cmd/site"
	"github.com/ActiveMemory/ctx/internal/cli/journal/cmd/source"
	journalSync "github.com/ActiveMemory/ctx/internal/cli/journal/cmd/sync"
	"github.com/ActiveMemory/ctx/internal/cli/journal/cmd/unlock"
	"github.com/ActiveMemory/ctx/internal/cli/parent"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the journal command with subcommands.
//
// The journal system provides tools for importing, analyzing,
// enriching, and publishing AI session history in
// .context/journal/.
//
// Returns:
//   - *cobra.Command: The journal command with subcommands
func Cmd() *cobra.Command {
	return parent.Cmd(cmd.DescKeyJournal, cmd.UseJournal,
		source.Cmd(),
		importer.Cmd(),
		journalSchema.Cmd(),
		lock.Cmd(),
		unlock.Cmd(),
		journalSync.Cmd(),
		site.Cmd(),
		obsidian.Cmd(),
	)
}
