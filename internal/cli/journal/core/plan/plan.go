//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package plan

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/tpl"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/index"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/lock"
	srcFmt "github.com/ActiveMemory/ctx/internal/cli/journal/core/source/format"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/validate"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/slug"
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
	jstate *state.State,
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
		numParts := (totalMsgs +
			journal.MaxMessagesPerPart - 1) /
			journal.MaxMessagesPerPart
		if numParts < 1 {
			numParts = 1
		}

		// Determine title-based slug.
		var existingTitle string
		if oldFile := index.LookupSessionFile(sessionIndex, s.ID); oldFile != "" {
			oldPath := filepath.Join(journalDir, oldFile)
			data, readErr := io.SafeReadUserFile(
				filepath.Clean(oldPath),
			)
			if readErr == nil {
				existingTitle = index.ExtractFrontmatterField(
					string(data), session.FrontmatterTitle,
				)
			}
		}
		slg, title := slug.ForTitle(s, existingTitle)

		baseFilename := srcFmt.JournalFilename(s, slg)
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

		// Growth detection (session-level). The unit of memory is the
		// source transcript, not the output file: a session whose
		// recorded mtime/size still match is Unchanged (skip); one that
		// grew is re-rendered part by part, guarded so a hand-edited
		// body is never clobbered.
		mtime, size, statOK := sourceStat(s.SourceFile)
		rec, seen := jstate.SessionSource(s.ID)
		forceRegen := singleSession || opts.Regenerate ||
			opts.DiscardFrontmatter()
		// Growth is any divergence from the recorded source: a different
		// transcript file (a richer resume copy was picked), a newer
		// mtime, or a larger size. Comparing the file path too means a
		// switch to a larger resume transcript counts as growth even if
		// its mtime/size happen not to differ from the old record.
		grown := seen && statOK &&
			(rec.SourceFile != s.SourceFile ||
				rec.SourceMtime != mtime || rec.SourceSize != size)

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

			_, statErr := io.SafeStat(path)
			fileExists := statErr == nil

			var action entity.ImportAction
			switch {
			case !fileExists:
				action = entity.ActionNew
				plan.NewCount++
			case jstate.Locked(filename):
				action = entity.ActionLocked
				plan.LockedCount++
			case lock.HasLocked(path):
				// Frontmatter says locked - promote to state so future
				// operations skip the file without reparsing.
				jstate.Mark(filename, session.FrontmatterLocked)
				action = entity.ActionLocked
				plan.LockedCount++
			case forceRegen:
				action = entity.ActionRegenerate
				plan.RegenCount++
			case grown:
				// The source grew. Re-render only if the existing body is
				// provably still ctx's last write; otherwise a human
				// edited it, so leave it untouched and warn. A grown
				// re-render is counted as GrownCount, not RegenCount, so it
				// does not trip the regenerate confirmation prompt: it is a
				// non-destructive splice that preserves frontmatter.
				if bodyEdited(path, jstate.RenderHash(filename)) {
					action = entity.ActionForeignEdit
					plan.SkipCount++
				} else {
					action = entity.ActionRegenerate
					plan.GrownCount++
				}
			default:
				// Unchanged (recorded stats match), adopted (file exists
				// but the session predates source tracking), or an
				// unreadable source: leave the file as-is.
				//
				// Adoption: a file that exists for a session never tracked
				// under source tracking (a pre-v2 import) is taken as
				// ctx-owned. Record its body hash now so a later growth
				// sweep can re-render it, instead of reading the absent
				// hash as a hand edit and stranding the growth (a
				// still-live session imported across the v1→v2 upgrade).
				if !seen && jstate.RenderHash(filename) == "" {
					if h := adoptRenderHash(path); h != "" {
						jstate.SetRenderHash(filename, h)
					}
				}
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

		// Stash the observed source stats in the plan, keyed by session.
		// execute commits them to journal state only AFTER the session's
		// writes succeed (see execute.Import): recording the stat at plan
		// time would advance it even when the write later fails, silently
		// forgetting the un-rendered growth. Adopted/unchanged sessions
		// have no write and are committed unconditionally by execute.
		if statOK {
			if plan.Sources == nil {
				plan.Sources = make(map[string]entity.SourceObservation)
			}
			plan.Sources[s.ID] = entity.SourceObservation{
				SourceFile: s.SourceFile,
				Mtime:      mtime,
				Size:       size,
			}
		}
	}

	return plan
}
