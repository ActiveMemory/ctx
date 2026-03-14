//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/recall/parser"
)

// PlanExport builds an ExportPlan without writing any files.
//
// Parameters:
//   - sessions: sessions to plan for.
//   - journalDir: absolute path to the journal output directory.
//   - sessionIndex: map of session ID to existing filename.
//   - jstate: journal processing state for lock checks.
//   - opts: export flag values.
//   - singleSession: true when exporting a single session by ID.
//
// Returns:
//   - ExportPlan: the planned actions, counters, and pending renames.
func PlanExport(
	sessions []*parser.Session,
	journalDir string,
	sessionIndex map[string]string,
	jstate *state.JournalState,
	opts ExportOpts,
	singleSession bool,
) ExportPlan {
	var plan ExportPlan

	for _, s := range sessions {
		// Collect non-empty messages.
		var nonEmptyMsgs []parser.Message
		for _, msg := range s.Messages {
			if !EmptyMessage(msg) {
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
		if oldFile := LookupSessionFile(sessionIndex, s.ID); oldFile != "" {
			oldPath := filepath.Join(journalDir, oldFile)
			if data, readErr := os.ReadFile(filepath.Clean(oldPath)); readErr == nil {
				existingTitle = ExtractFrontmatterField(
					string(data), assets.FrontmatterTitle,
				)
			}
		}
		slug, title := TitleSlug(s, existingTitle)

		baseFilename := FormatJournalFilename(s, slug)
		baseName := strings.TrimSuffix(baseFilename, file.ExtMarkdown)

		// Detect renames (dedup: old slug → new slug).
		if oldFile := LookupSessionFile(sessionIndex, s.ID); oldFile != "" {
			oldBase := strings.TrimSuffix(oldFile, file.ExtMarkdown)
			if oldBase != baseName {
				plan.RenameOps = append(plan.RenameOps, RenameOp{
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
				filename = fmt.Sprintf(assets.TplRecallPartFilename, baseName, part)
			}
			path := filepath.Join(journalDir, filename)

			startIdx := (part - 1) * journal.MaxMessagesPerPart
			endIdx := startIdx + journal.MaxMessagesPerPart
			if endIdx > totalMsgs {
				endIdx = totalMsgs
			}

			_, statErr := os.Stat(path)
			fileExists := statErr == nil

			var action ExportAction
			switch {
			case !fileExists:
				action = ActionNew
				plan.NewCount++
			case jstate.Locked(filename):
				action = ActionLocked
				plan.LockedCount++
			case FrontmatterHasLocked(path):
				// Frontmatter says locked — promote to state so future
				// operations skip the file without reparsing.
				jstate.Mark(filename, assets.FrontmatterLocked)
				action = ActionLocked
				plan.LockedCount++
			case singleSession || opts.Regenerate || opts.DiscardFrontmatter():
				action = ActionRegenerate
				plan.RegenCount++
			default:
				action = ActionSkip
				plan.SkipCount++
			}

			plan.Actions = append(plan.Actions, FileAction{
				Session:    s,
				Filename:   filename,
				Path:       path,
				Part:       part,
				TotalParts: numParts,
				StartIdx:   startIdx,
				EndIdx:     endIdx,
				Action:     action,
				Messages:   nonEmptyMsgs,
				Slug:       slug,
				Title:      title,
				BaseName:   baseName,
			})
		}
	}

	return plan
}
