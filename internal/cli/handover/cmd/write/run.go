//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	handoverPath "github.com/ActiveMemory/ctx/internal/cli/handover/core/path"
	kbPath "github.com/ActiveMemory/ctx/internal/cli/kb/core/path"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/validate"
	"github.com/ActiveMemory/ctx/internal/write/handover"
)

// Run executes the handover write command.
//
// Parameters:
//   - cobraCmd: cobra command for output.
//   - title: handover title (passed as positional arg).
//   - summary: past-tense summary; validated non-placeholder.
//   - next: next-session first action; validated
//     non-placeholder.
//   - highlights: notable artifacts.
//   - openQuestions: undecided items.
//   - commit: optional commit-override.
//   - noFold: when true, skip closeout fold + archive.
//
// Returns:
//   - error: rejection or wrapped writer errors.
func Run(
	cobraCmd *cobra.Command,
	title, summary, next, highlights, openQuestions, commit string,
	noFold bool,
) error {
	if rejectErr := validate.RejectPlaceholder(
		cFlag.Summary, summary,
	); rejectErr != nil {
		cobraCmd.SilenceUsage = true
		return rejectErr
	}
	if rejectErr := validate.RejectPlaceholder(
		cFlag.Next, next,
	); rejectErr != nil {
		cobraCmd.SilenceUsage = true
		return rejectErr
	}

	handoversDir, hoDirErr := handoverPath.Dir()
	if hoDirErr != nil {
		return hoDirErr
	}
	closeoutsDir, coDirErr := kbPath.CloseoutsDir()
	if coDirErr != nil {
		return coDirErr
	}
	archiveDir, archiveErr := kbPath.ArchiveCloseoutsDir()
	if archiveErr != nil {
		return archiveErr
	}
	ctxDir, ctxErr := rc.ContextDir()
	if ctxErr != nil {
		return ctxErr
	}
	projectRoot := filepath.Dir(ctxDir)

	res, writeErr := handover.Write(
		handoversDir, closeoutsDir, archiveDir, projectRoot,
		handover.Entry{
			Title:          title,
			Summary:        summary,
			Next:           next,
			Highlights:     highlights,
			OpenQuestions:  openQuestions,
			CommitOverride: commit,
			NoFold:         noFold,
		},
	)
	if writeErr != nil {
		cobraCmd.SilenceUsage = true
		return writeErr
	}

	io.SafeFprintf(
		cobraCmd.OutOrStdout(),
		desc.Text(text.DescKeyWriteHandoverWrote),
		res.File.Path,
	)
	if len(res.FoldedCloseouts) > 0 {
		io.SafeFprintf(
			cobraCmd.OutOrStdout(),
			desc.Text(text.DescKeyWriteHandoverFolded),
			len(res.FoldedCloseouts), archiveDir,
		)
	}
	if len(res.MalformedCloseouts) > 0 {
		io.SafeFprintf(
			cobraCmd.ErrOrStderr(),
			desc.Text(text.DescKeyWriteHandoverMalformedWarning),
			len(res.MalformedCloseouts),
		)
		for _, p := range res.MalformedCloseouts {
			io.SafeFprintf(
				cobraCmd.ErrOrStderr(),
				desc.Text(text.DescKeyWriteHandoverMalformedLine),
				p,
			)
		}
	}
	return nil
}
