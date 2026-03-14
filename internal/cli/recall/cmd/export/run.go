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

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/recall/core"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/recall/parser"
	"github.com/ActiveMemory/ctx/internal/write"
)

// Run handles the recall export command.
//
// Parameters:
//   - cmd: Cobra command for output.
//   - args: positional arguments (optional session ID).
//   - opts: export flag values.
//
// Returns:
//   - error: non-nil on validation, scan, or write failures.
func Run(cmd *cobra.Command, args []string, opts core.ExportOpts) error {
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
	journalDir := filepath.Join(rc.ContextDir(), dir.Journal)
	if mkErr := os.MkdirAll(journalDir, fs.PermExec); mkErr != nil {
		return ctxerr.Mkdir(dir.Journal, mkErr)
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
			rop.OldBase+file.ExtMarkdown, rop.NewBase+file.ExtMarkdown,
		)
		renamed++
	}

	// 8. Dry-run → print summary and return.
	if opts.DryRun {
		write.ExportSummary(cmd, plan.NewCount, plan.RegenCount, plan.SkipCount, plan.LockedCount, true)
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
	exported, updated, skipped := core.ExecuteExport(cmd, plan, jstate, opts)

	// 11. Persist journal state.
	if saveErr := jstate.Save(journalDir); saveErr != nil {
		write.WarnFileErr(cmd, journal.FileState, saveErr)
	}

	// 12. Print final summary.
	write.ExportFinalSummary(cmd, exported, updated, renamed, skipped)

	return nil
}
