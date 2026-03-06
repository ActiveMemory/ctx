//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/recall/parser"
)

// planExport builds an exportPlan without writing any files.
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
//   - exportPlan: the planned actions, counters, and pending renames.
func planExport(
	sessions []*parser.Session,
	journalDir string,
	sessionIndex map[string]string,
	jstate *state.JournalState,
	opts exportOpts,
	singleSession bool,
) exportPlan {
	var plan exportPlan

	for _, s := range sessions {
		// Collect non-empty messages.
		var nonEmptyMsgs []parser.Message
		for _, msg := range s.Messages {
			if !emptyMessage(msg) {
				nonEmptyMsgs = append(nonEmptyMsgs, msg)
			}
		}

		totalMsgs := len(nonEmptyMsgs)
		numParts := (totalMsgs + config.MaxMessagesPerPart - 1) / config.MaxMessagesPerPart
		if numParts < 1 {
			numParts = 1
		}

		// Determine title-based slug.
		var existingTitle string
		if oldFile := lookupSessionFile(sessionIndex, s.ID); oldFile != "" {
			oldPath := filepath.Join(journalDir, oldFile)
			if data, readErr := os.ReadFile(filepath.Clean(oldPath)); readErr == nil {
				existingTitle = extractFrontmatterField(
					string(data), config.FrontmatterTitle,
				)
			}
		}
		slug, title := titleSlug(s, existingTitle)

		baseFilename := formatJournalFilename(s, slug)
		baseName := strings.TrimSuffix(baseFilename, config.ExtMarkdown)

		// Detect renames (dedup: old slug → new slug).
		if oldFile := lookupSessionFile(sessionIndex, s.ID); oldFile != "" {
			oldBase := strings.TrimSuffix(oldFile, config.ExtMarkdown)
			if oldBase != baseName {
				plan.renameOps = append(plan.renameOps, renameOp{
					oldBase:  oldBase,
					newBase:  baseName,
					numParts: numParts,
				})
			}
		}

		// Plan each part.
		for part := 1; part <= numParts; part++ {
			filename := baseFilename
			if numParts > 1 && part > 1 {
				filename = fmt.Sprintf(config.TplRecallPartFilename, baseName, part)
			}
			path := filepath.Join(journalDir, filename)

			startIdx := (part - 1) * config.MaxMessagesPerPart
			endIdx := startIdx + config.MaxMessagesPerPart
			if endIdx > totalMsgs {
				endIdx = totalMsgs
			}

			_, statErr := os.Stat(path)
			fileExists := statErr == nil

			var action exportAction
			switch {
			case !fileExists:
				action = actionNew
				plan.newCount++
			case jstate.Locked(filename):
				action = actionLocked
				plan.lockedCount++
			case frontmatterHasLocked(path):
				// Frontmatter says locked — promote to state so future
				// operations skip the file without reparsing.
				jstate.Mark(filename, config.FrontmatterLocked)
				action = actionLocked
				plan.lockedCount++
			case singleSession || opts.regenerate || opts.discardFrontmatter():
				action = actionRegenerate
				plan.regenCount++
			default:
				action = actionSkip
				plan.skipCount++
			}

			plan.actions = append(plan.actions, fileAction{
				session:    s,
				filename:   filename,
				path:       path,
				part:       part,
				totalParts: numParts,
				startIdx:   startIdx,
				endIdx:     endIdx,
				action:     action,
				messages:   nonEmptyMsgs,
				slug:       slug,
				title:      title,
				baseName:   baseName,
			})
		}
	}

	return plan
}
