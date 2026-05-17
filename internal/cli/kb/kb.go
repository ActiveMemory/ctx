//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package kb

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/kb/cmd/ask"
	"github.com/ActiveMemory/ctx/internal/cli/kb/cmd/ground"
	"github.com/ActiveMemory/ctx/internal/cli/kb/cmd/ingest"
	"github.com/ActiveMemory/ctx/internal/cli/kb/cmd/note"
	kbReindex "github.com/ActiveMemory/ctx/internal/cli/kb/cmd/reindex"
	"github.com/ActiveMemory/ctx/internal/cli/kb/cmd/sitereview"
	"github.com/ActiveMemory/ctx/internal/cli/kb/cmd/topic"
	"github.com/ActiveMemory/ctx/internal/cli/parent"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the `ctx kb` parent command with the editorial
// pipeline subcommands registered.
//
// Returns:
//   - *cobra.Command: kb parent with topic, ingest, ask,
//     site-review, ground, note, and reindex.
func Cmd() *cobra.Command {
	return parent.Cmd(cmd.DescKeyKB, cmd.UseKB,
		topic.Cmd(),
		ingest.Cmd(),
		ask.Cmd(),
		sitereview.Cmd(),
		ground.Cmd(),
		note.Cmd(),
		kbReindex.Cmd(),
	)
}
