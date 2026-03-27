//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package plan

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/tpl"
	"github.com/ActiveMemory/ctx/internal/cli/recall/core/format"
	"github.com/ActiveMemory/ctx/internal/cli/recall/core/index"
	"github.com/ActiveMemory/ctx/internal/cli/recall/core/lock"
	"github.com/ActiveMemory/ctx/internal/cli/recall/core/slug"
	"github.com/ActiveMemory/ctx/internal/cli/recall/core/validate"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/journal/state"
)

// Import builds an ImportPlan without writing any files.
//
// Parameters:
//   - sessions: sessions to plan for.
//   - journalDir: absolute path to the journal output directory.
//   - sessionIndex: map of session ID to existing filename.
//   - jstate: journal processing state for lock checks.
//   - opts: import flag values.
//   - singleSession: true when importing a single session by ID.
//
// Returns:
//   - ImportPlan: the planned actions, counters, and pending renames.
func Import(
	sessions []*entity.Session,
	journalDir string,
	sessionIndex map[string]string,
	jstate *state.JournalState,
	opts entity.ImportOpts,
	singleSession bool,
) entity.ImportPlan {
	var plan entity.ImportPlan

	for _, s := range sessions {
		// Collect non-empty messages.
		var nonEmptyMsgs []entity.Message
		for _, msg := range s.Messages {
			if !validate.EmptyMessage(msg) {
				nonEmptyMsgs = append(nonEmptyMsgs, msg)
			}
		}

		totalMsgs := len(nonEmptyMsgs)
		numParts := (totalMsgs + journal.MaxMessagesPerPart - 1) / journal.MaxMessagesPerPart
		if numParts < 1 {
			numParts = 1
		}

		// Determine title-based slug.
		var existingTitle string
		if oldFile := index.LookupSessionFile(sessionIndex, s.ID); oldFile != "" {
			oldPath := filepath.Join(journalDir, oldFile)
			if data, readErr := os.ReadFile(filepath.Clean(oldPath)); readErr == nil {
				existingTitle = index.ExtractFrontmatterField(
					string(data), session.FrontmatterTitle,
				)
			}
		}
		slg, title := slug.TitleSlug(s, existingTitle)

		baseFilename := format.JournalFilename(s, slg)
		baseName := strings.TrimSuffix(baseFilename, file.ExtMarkdown)

		// Detect renames (dedup: old slug → new slug).
		if oldFile := index.LookupSessionFile(sessionIndex, s.ID); oldFile != "" {
			oldBase := strings.TrimSuffix(oldFile, file.ExtMarkdown)
			if oldBase != baseName {
				plan.RenameOps = append(plan.RenameOps, entity.RenameOp{
					OldBase:  oldBase,
					NewBase:  baseName,
					NumParts: numParts,
				})
			}
		}

		// Plan each part.
		for part := 1; part <= numParts; part++ {
			filename := baseFilename
			if numParts > 1 && part > 1 {
				filename = fmt.Sprintf(tpl.RecallPartFilename, baseName, part)
			}
			path := filepath.Join(journalDir, filename)

			startIdx := (part - 1) * journal.MaxMessagesPerPart
			endIdx := startIdx + journal.MaxMessagesPerPart
			if endIdx > totalMsgs {
				endIdx = totalMsgs
			}

			_, statErr := os.Stat(path)
			fileExists := statErr == nil

			var action entity.ImportAction
			switch {
			case !fileExists:
				action = entity.ActionNew
				plan.NewCount++
			case jstate.Locked(filename):
				action = entity.ActionLocked
				plan.LockedCount++
			case lock.FrontmatterHasLocked(path):
				// Frontmatter says locked - promote to state so future
				// operations skip the file without reparsing.
				jstate.Mark(filename, session.FrontmatterLocked)
				action = entity.ActionLocked
				plan.LockedCount++
			case singleSession || opts.Regenerate || opts.DiscardFrontmatter():
				action = entity.ActionRegenerate
				plan.RegenCount++
			default:
				action = entity.ActionSkip
				plan.SkipCount++
			}

			plan.Actions = append(plan.Actions, entity.FileAction{
				Session:    s,
				Filename:   filename,
				Path:       path,
				Part:       part,
				TotalParts: numParts,
				StartIdx:   startIdx,
				EndIdx:     endIdx,
				Action:     action,
				Messages:   nonEmptyMsgs,
				Slug:       slg,
				Title:      title,
				BaseName:   baseName,
			})
		}
	}

	return plan
}
