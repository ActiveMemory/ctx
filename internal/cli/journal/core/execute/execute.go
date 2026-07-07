//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package execute

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/extract"
	srcFmt "github.com/ActiveMemory/ctx/internal/cli/journal/core/source/format"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/write/err"
	writeRecall "github.com/ActiveMemory/ctx/internal/write/journal"
)

// Import writes files according to the plan.
//
// Parameters:
//   - cmd: Cobra command for output.
//   - plan: the import plan with file actions.
//   - jstate: journal state to update as files are imported.
//   - opts: import flag values.
//
// Returns:
//   - imported: number of new files written.
//   - updated: number of existing files updated (frontmatter preserved).
//   - skipped: number of files skipped (existing or locked).
func Import(
	cmd *cobra.Command,
	plan entity.ImportPlan,
	jstate *state.State,
	opts entity.ImportOpts,
) (imported, updated, skipped int) {
	// Track per-session outcomes so the recorded source stat is advanced
	// only for sessions that fully caught up to their transcript. A write
	// failure or a foreign-edited part leaves the entry behind the
	// source, so its stat must not move — the next sweep has to retry.
	failed := make(map[string]bool)
	foreign := make(map[string]bool)

	for _, fa := range plan.Actions {
		if fa.Action == entity.ActionLocked {
			skipped++
			writeRecall.SkipFile(cmd, fa.Filename, session.FrontmatterLocked)
			continue
		}
		if fa.Action == entity.ActionSkip {
			skipped++
			writeRecall.SkipFile(
				cmd, fa.Filename,
				desc.Text(text.DescKeyLabelReasonExists),
			)
			continue
		}
		if fa.Action == entity.ActionForeignEdit {
			skipped++
			foreign[fa.Session.ID] = true
			writeRecall.SkipFile(
				cmd, fa.Filename,
				desc.Text(text.DescKeyLabelReasonEdited),
			)
			continue
		}

		// Generate content, sanitizing any invalid UTF-8.
		content := strings.ToValidUTF8(
			srcFmt.JournalEntryPart(
				fa.Session, fa.Messages[fa.StartIdx:fa.EndIdx],
				fa.StartIdx, fa.Part, fa.TotalParts, fa.BaseName, fa.Title,
			),
			token.Ellipsis,
		)

		fileExists := fa.Action == entity.ActionRegenerate

		// Preserve enriched YAML frontmatter from the existing file.
		discard := opts.DiscardFrontmatter()
		if fileExists && !discard {
			existing, readErr := io.SafeReadUserFile(filepath.Clean(fa.Path))
			if readErr == nil {
				if fm := extract.Frontmatter(string(existing)); fm != "" {
					content = fm + token.NewlineLF + extract.StripFrontmatter(content)
				}
			}
		}
		if fileExists && discard {
			jstate.ClearEnriched(fa.Filename)
		}
		if fileExists && !discard {
			updated++
		} else {
			imported++
		}

		// Write the entry atomically (temp + rename). Import is wired into
		// a SessionEnd hook that can be killed during teardown; a torn
		// write would leave a truncated entry that the next sweep records
		// as unchanged and never repairs.
		if writeErr := io.SafeWriteFileAtomic(
			fa.Path, []byte(content), fs.PermFile,
		); writeErr != nil {
			err.WarnFile(cmd, fa.Filename, writeErr)
			failed[fa.Session.ID] = true
			continue
		}

		jstate.MarkImported(fa.Filename)

		// Record the hash of the body ctx just wrote, so a later growth
		// sweep can prove the file is still ctx-owned before re-rendering.
		jstate.SetRenderHash(
			fa.Filename, state.HashRender(extract.StripFrontmatter(content)),
		)

		if fileExists && !discard {
			writeRecall.ImportedFile(
				cmd, fa.Filename, desc.Text(text.DescKeyLabelReasonUpdated),
			)
		} else {
			writeRecall.ImportedFile(cmd, fa.Filename, "")
		}
	}

	// Commit the observed source stats for sessions that caught up. A
	// session with any failed write or foreign-edited part is skipped, so
	// its recorded stat stays behind and the next sweep re-detects the
	// growth instead of silently forgetting it.
	for sid, obs := range plan.Sources {
		if failed[sid] || foreign[sid] {
			continue
		}
		jstate.MarkSource(sid, obs.SourceFile, obs.Mtime, obs.Size)
	}

	return imported, updated, skipped
}
