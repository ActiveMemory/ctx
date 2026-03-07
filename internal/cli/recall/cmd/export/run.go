//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package export

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/recall/core"
	"github.com/ActiveMemory/ctx/internal/config"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/recall/parser"
	"github.com/ActiveMemory/ctx/internal/write"
)

// executeExport writes files according to the plan.
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
func executeExport(
	cmd *cobra.Command,
	plan core.ExportPlan,
	jstate *state.JournalState,
	opts core.ExportOpts,
) (exported, updated, skipped int) {
	for _, fa := range plan.Actions {
		if fa.Action == core.ActionLocked {
			skipped++
			write.SkipFile(cmd, fa.Filename, config.FrontmatterLocked)
			continue
		}
		if fa.Action == core.ActionSkip {
			skipped++
			write.SkipFile(cmd, fa.Filename, config.ReasonExists)
			continue
		}

		// Generate content, sanitizing any invalid UTF-8.
		content := strings.ToValidUTF8(
			core.FormatJournalEntryPart(
				fa.Session, fa.Messages[fa.StartIdx:fa.EndIdx],
				fa.StartIdx, fa.Part, fa.TotalParts, fa.BaseName, fa.Title,
			),
			config.Ellipsis,
		)

		fileExists := fa.Action == core.ActionRegenerate

		// Preserve enriched YAML frontmatter from the existing file.
		discard := opts.DiscardFrontmatter()
		if fileExists && !discard {
			existing, readErr := os.ReadFile(filepath.Clean(fa.Path))
			if readErr == nil {
				if fm := core.ExtractFrontmatter(string(existing)); fm != "" {
					content = fm + config.NewlineLF + core.StripFrontmatter(content)
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
			fa.Path, []byte(content), config.PermFile,
		); writeErr != nil {
			write.WarnFileErr(cmd, fa.Filename, writeErr)
			continue
		}

		jstate.MarkExported(fa.Filename)

		if fileExists && !discard {
			write.ExportedFile(cmd, fa.Filename, config.ReasonUpdated)
		} else {
			write.ExportedFile(cmd, fa.Filename, "")
		}
	}

	return exported, updated, skipped
}

// runExport handles the recall export command.
//
// Parameters:
//   - cmd: Cobra command for output.
//   - args: positional arguments (optional session ID).
//   - opts: export flag values.
//
// Returns:
//   - error: non-nil on validation, scan, or write failures.
func runExport(cmd *cobra.Command, args []string, opts core.ExportOpts) error {
	// --keep-frontmatter=false implies --regenerate
	// (can't discard without regenerating).
	if !opts.KeepFrontmatter {
		opts.Regenerate = true
	}

	// 1. Validate flags.
	if validateErr := core.ValidateExportFlags(args, opts); validateErr != nil {
		return validateErr
	}

	// 2. Bare export (no args, no --all) → show help (T2.8).
	if len(args) == 0 && !opts.All {
		return cmd.Help()
	}

	// 3. Resolve sessions.
	sessions, scanErr := core.FindSessions(opts.AllProjects)
	if scanErr != nil {
		return ctxerr.FindSessions(scanErr)
	}

	if len(sessions) == 0 {
		write.NoSessionsForProject(cmd, opts.AllProjects)
		return nil
	}

	var toExport []*parser.Session
	singleSession := false
	if opts.All {
		toExport = sessions
	} else {
		query := strings.ToLower(args[0])
		for _, s := range sessions {
			if strings.HasPrefix(strings.ToLower(s.ID), query) ||
				strings.Contains(strings.ToLower(s.Slug), query) {
				toExport = append(toExport, s)
			}
		}
		if len(toExport) == 0 {
			return ctxerr.SessionNotFound(args[0])
		}
		if len(toExport) > 1 {
			lines := core.FormatSessionMatchLines(toExport)
			write.AmbiguousSessionMatch(cmd, args[0], lines)
			return ctxerr.AmbiguousQuery()
		}
		singleSession = true
	}

	// 4. Ensure journal directory exists.
	journalDir := filepath.Join(rc.ContextDir(), config.DirJournal)
	if mkErr := os.MkdirAll(journalDir, config.PermExec); mkErr != nil {
		return ctxerr.Mkdir(config.DirJournal, mkErr)
	}

	// 5. Load state + build index.
	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		return ctxerr.LoadJournalState(loadErr)
	}
	sessionIndex := core.BuildSessionIndex(journalDir)

	// 6. Build the plan.
	plan := core.PlanExport(toExport, journalDir, sessionIndex, jstate, opts, singleSession)

	// 7. Execute renames.
	renamed := 0
	for _, rop := range plan.RenameOps {
		core.RenameJournalFiles(journalDir, rop.OldBase, rop.NewBase, rop.NumParts)
		jstate.Rename(
			rop.OldBase+config.ExtMarkdown, rop.NewBase+config.ExtMarkdown,
		)
		renamed++
	}

	// 8. Dry-run → print summary and return.
	if opts.DryRun {
		write.ExportSummary(cmd, core.PlanCounts(plan), true)
		return nil
	}

	// 9. Confirmation prompt for regeneration.
	if plan.RegenCount > 0 && !opts.Yes && !singleSession {
		ok, promptErr := core.ConfirmExport(cmd, plan)
		if promptErr != nil {
			return promptErr
		}
		if !ok {
			write.Aborted(cmd)
			return nil
		}
	}

	// 10. Execute the export.
	exported, updated, skipped := executeExport(cmd, plan, jstate, opts)

	// 11. Persist journal state.
	if saveErr := jstate.Save(journalDir); saveErr != nil {
		write.WarnFileErr(cmd, config.FileJournalState, saveErr)
	}

	// 12. Print final summary.
	write.ExportFinalSummary(cmd, exported, updated, renamed, skipped)

	return nil
}
