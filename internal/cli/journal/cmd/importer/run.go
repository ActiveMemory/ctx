//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package importer

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/journal/core/confirm"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/execute"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/index"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/plan"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/query"
	coreSchema "github.com/ActiveMemory/ctx/internal/cli/journal/core/schema"
	srcFmt "github.com/ActiveMemory/ctx/internal/cli/journal/core/source/format"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/validate"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/entity"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errJournal "github.com/ActiveMemory/ctx/internal/err/journal"
	errSession "github.com/ActiveMemory/ctx/internal/err/session"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/journal/schema"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write/err"
	writeRecall "github.com/ActiveMemory/ctx/internal/write/journal"
	writeSchema "github.com/ActiveMemory/ctx/internal/write/schema"
)

// Run handles the journal import command.
//
// Parameters:
//   - cmd: Cobra command for output.
//   - args: positional arguments (optional session ID).
//   - opts: import flag values.
//
// Returns:
//   - error: non-nil on validation, scan, or write failures.
func Run(cmd *cobra.Command, args []string, opts entity.ImportOpts) error {
	// --keep-frontmatter=false implies --regenerate
	// (can't discard without regenerating).
	if !opts.KeepFrontmatter {
		opts.Regenerate = true
	}

	// 1. Validate flags.
	if validateErr := validate.ImportFlags(args, opts); validateErr != nil {
		return validateErr
	}

	// 2. Bare import (no args, no --all) → show help (T2.8).
	if len(args) == 0 && !opts.All {
		return cmd.Help()
	}

	// 3. Resolve sessions.
	sessions, scanErr := query.FindSessions(opts.AllProjects)
	if scanErr != nil {
		return errSession.Find(scanErr)
	}

	if len(sessions) == 0 {
		writeRecall.NoSessionsForProject(cmd, opts.AllProjects)
		return nil
	}

	var toImport []*entity.Session
	singleSession := false
	if opts.All {
		toImport = sessions
	} else {
		qry := strings.ToLower(args[0])
		for _, s := range sessions {
			if strings.HasPrefix(strings.ToLower(s.ID), qry) ||
				strings.Contains(strings.ToLower(s.Slug), qry) {
				toImport = append(toImport, s)
			}
		}
		if len(toImport) == 0 {
			return errSession.NotFound(args[0])
		}
		if len(toImport) > 1 {
			lines := srcFmt.SessionMatchLines(toImport)
			writeRecall.AmbiguousSessionMatch(cmd, args[0], lines)
			return errSession.AmbiguousQuery()
		}
		singleSession = true
	}

	// 4. Ensure journal directory exists.
	journalDir := filepath.Join(rc.ContextDir(), dir.Journal)
	if mkErr := ctxIo.SafeMkdirAll(journalDir, fs.PermExec); mkErr != nil {
		return errFs.Mkdir(dir.Journal, mkErr)
	}

	// 5. Load state + build index.
	jState, loadErr := state.Load(journalDir)
	if loadErr != nil {
		return errJournal.LoadState(loadErr)
	}
	sessionIndex := index.Session(journalDir)

	// 6. Build the plan.
	importPlan := plan.Import(
		toImport, journalDir, sessionIndex, jState, opts, singleSession,
	)

	// 7. Execute renames.
	renamed := 0
	for _, rop := range importPlan.RenameOps {
		index.RenameJournalFiles(journalDir, rop.OldBase, rop.NewBase, rop.NumParts)
		jState.Rename(
			rop.OldBase+file.ExtMarkdown, rop.NewBase+file.ExtMarkdown,
		)
		renamed++
	}

	// 8. Dry-run → print summary and return.
	if opts.DryRun {
		writeRecall.ImportSummary(
			cmd, importPlan.NewCount, importPlan.RegenCount,
			importPlan.SkipCount, importPlan.LockedCount, true,
		)
		return nil
	}

	// 9. Confirmation prompt for regeneration.
	if importPlan.RegenCount > 0 && !opts.Yes && !singleSession {
		ok, promptErr := confirm.Import(cmd, importPlan)
		if promptErr != nil {
			return promptErr
		}
		if !ok {
			writeRecall.Aborted(cmd)
			return nil
		}
	}

	// 10. Execute the import.
	imported, updated, skipped := execute.Import(cmd, importPlan, jState, opts)

	// 11. Persist journal state.
	if saveErr := jState.Save(journalDir); saveErr != nil {
		err.WarnFile(cmd, journal.File, saveErr)
	}

	// 12. Schema drift check on imported source files.
	c := coreSchema.CheckSessions(toImport)
	if c.Drift() {
		writeSchema.DriftSummary(
			cmd, schema.Summary(c),
		)
	}

	// 13. Print final summary.
	writeRecall.ImportFinalSummary(
		cmd, imported, updated, renamed, skipped,
	)

	return nil
}
