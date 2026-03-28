//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package importer

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/journal/core/confirm"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/execute"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/index"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/plan"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/query"
	sourceFormat "github.com/ActiveMemory/ctx/internal/cli/journal/core/source/format"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/validate"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/entity"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errJournal "github.com/ActiveMemory/ctx/internal/err/journal"
	errSession "github.com/ActiveMemory/ctx/internal/err/session"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write/err"
	"github.com/ActiveMemory/ctx/internal/write/recall"
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
		recall.NoSessionsForProject(cmd, opts.AllProjects)
		return nil
	}

	var toImport []*entity.Session
	singleSession := false
	if opts.All {
		toImport = sessions
	} else {
		query := strings.ToLower(args[0])
		for _, s := range sessions {
			if strings.HasPrefix(strings.ToLower(s.ID), query) ||
				strings.Contains(strings.ToLower(s.Slug), query) {
				toImport = append(toImport, s)
			}
		}
		if len(toImport) == 0 {
			return errSession.NotFound(args[0])
		}
		if len(toImport) > 1 {
			lines := sourceFormat.SessionMatchLines(toImport)
			recall.AmbiguousSessionMatch(cmd, args[0], lines)
			return errSession.AmbiguousQuery()
		}
		singleSession = true
	}

	// 4. Ensure journal directory exists.
	journalDir := filepath.Join(rc.ContextDir(), dir.Journal)
	if mkErr := os.MkdirAll(journalDir, fs.PermExec); mkErr != nil {
		return errFs.Mkdir(dir.Journal, mkErr)
	}

	// 5. Load state + build index.
	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		return errJournal.LoadState(loadErr)
	}
	sessionIndex := index.SessionIndex(journalDir)

	// 6. Build the plan.
	plan := plan.Import(
		toImport, journalDir, sessionIndex, jstate, opts, singleSession,
	)

	// 7. Execute renames.
	renamed := 0
	for _, rop := range plan.RenameOps {
		index.RenameJournalFiles(journalDir, rop.OldBase, rop.NewBase, rop.NumParts)
		jstate.Rename(
			rop.OldBase+file.ExtMarkdown, rop.NewBase+file.ExtMarkdown,
		)
		renamed++
	}

	// 8. Dry-run → print summary and return.
	if opts.DryRun {
		recall.ImportSummary(
			cmd, plan.NewCount, plan.RegenCount,
			plan.SkipCount, plan.LockedCount, true,
		)
		return nil
	}

	// 9. Confirmation prompt for regeneration.
	if plan.RegenCount > 0 && !opts.Yes && !singleSession {
		ok, promptErr := confirm.Import(cmd, plan)
		if promptErr != nil {
			return promptErr
		}
		if !ok {
			recall.Aborted(cmd)
			return nil
		}
	}

	// 10. Execute the import.
	imported, updated, skipped := execute.Import(cmd, plan, jstate, opts)

	// 11. Persist journal state.
	if saveErr := jstate.Save(journalDir); saveErr != nil {
		err.WarnFile(cmd, journal.File, saveErr)
	}

	// 12. Print final summary.
	recall.ImportFinalSummary(cmd, imported, updated, renamed, skipped)

	return nil
}
