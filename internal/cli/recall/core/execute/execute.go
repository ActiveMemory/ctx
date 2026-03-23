//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package execute

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/recall/core/extract"
	"github.com/ActiveMemory/ctx/internal/cli/recall/core/format"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/write/err"
	"github.com/ActiveMemory/ctx/internal/write/recall"
)

// Export writes files according to the plan.
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
func Export(
		cmd *cobra.Command,
		plan entity.ExportPlan,
		jstate *state.JournalState,
		opts entity.ExportOpts,
) (exported, updated, skipped int) {
	for _, fa := range plan.Actions {
		if fa.Action == entity.ActionLocked {
			skipped++
			recall.SkipFile(cmd, fa.Filename, session.FrontmatterLocked)
			continue
		}
		if fa.Action == entity.ActionSkip {
			skipped++
			recall.SkipFile(cmd, fa.Filename, desc.Text(text.DescKeyLabelReasonExists))
			continue
		}

		// Generate content, sanitizing any invalid UTF-8.
		content := strings.ToValidUTF8(
			format.JournalEntryPart(
				fa.Session, fa.Messages[fa.StartIdx:fa.EndIdx],
				fa.StartIdx, fa.Part, fa.TotalParts, fa.BaseName, fa.Title,
			),
			token.Ellipsis,
		)

		fileExists := fa.Action == entity.ActionRegenerate

		// Preserve enriched YAML frontmatter from the existing file.
		discard := opts.DiscardFrontmatter()
		if fileExists && !discard {
			existing, readErr := os.ReadFile(filepath.Clean(fa.Path))
			if readErr == nil {
				if fm := extract.ExtractFrontmatter(string(existing)); fm != "" {
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
			exported++
		}

		// Write file.
		if writeErr := os.WriteFile(
			fa.Path, []byte(content), fs.PermFile,
		); writeErr != nil {
			err.WarnFile(cmd, fa.Filename, writeErr)
			continue
		}

		jstate.MarkExported(fa.Filename)

		if fileExists && !discard {
			recall.ExportedFile(
				cmd, fa.Filename, desc.Text(text.DescKeyLabelReasonUpdated),
			)
		} else {
			recall.ExportedFile(cmd, fa.Filename, "")
		}
	}

	return exported, updated, skipped
}
