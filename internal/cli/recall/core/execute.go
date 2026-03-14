//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/write"
)

// ExecuteExport writes files according to the plan.
//
// Parameters:
//   - cmd: Cobra command for output.
//   - plan: the export plan with file actions.
//   - jstate: journal state to update as files are exported.
//   - opts: export flag values.
//
// Returns:
//   - exported: number of new files written.
//   - updated: number of existing files updated (frontmatter preserved).
//   - skipped: number of files skipped (existing or locked).
func ExecuteExport(
	cmd *cobra.Command,
	plan ExportPlan,
	jstate *state.JournalState,
	opts ExportOpts,
) (exported, updated, skipped int) {
	for _, fa := range plan.Actions {
		if fa.Action == ActionLocked {
			skipped++
			write.SkipFile(cmd, fa.Filename, assets.FrontmatterLocked)
			continue
		}
		if fa.Action == ActionSkip {
			skipped++
			write.SkipFile(cmd, fa.Filename, assets.ReasonExists)
			continue
		}

		// Generate content, sanitizing any invalid UTF-8.
		content := strings.ToValidUTF8(
			FormatJournalEntryPart(
				fa.Session, fa.Messages[fa.StartIdx:fa.EndIdx],
				fa.StartIdx, fa.Part, fa.TotalParts, fa.BaseName, fa.Title,
			),
			token.Ellipsis,
		)

		fileExists := fa.Action == ActionRegenerate

		// Preserve enriched YAML frontmatter from the existing file.
		discard := opts.DiscardFrontmatter()
		if fileExists && !discard {
			existing, readErr := os.ReadFile(filepath.Clean(fa.Path))
			if readErr == nil {
				if fm := ExtractFrontmatter(string(existing)); fm != "" {
					content = fm + token.NewlineLF + StripFrontmatter(content)
				}
			}
		}
		if fileExists && discard {
			jstate.ClearEnriched(fa.Filename)
		}
		if fileExists && !discard {
			updated++
		} else {
			exported++
		}

		// Write file.
		if writeErr := os.WriteFile(
			fa.Path, []byte(content), fs.PermFile,
		); writeErr != nil {
			write.WarnFileErr(cmd, fa.Filename, writeErr)
			continue
		}

		jstate.MarkExported(fa.Filename)

		if fileExists && !discard {
			write.ExportedFile(cmd, fa.Filename, assets.ReasonUpdated)
		} else {
			write.ExportedFile(cmd, fa.Filename, "")
		}
	}

	return exported, updated, skipped
}
