//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package importer

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/entry"
	memory2 "github.com/ActiveMemory/ctx/internal/config/memory"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/memory/core"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/validation"
	"github.com/ActiveMemory/ctx/internal/write"
)

// Run parses MEMORY.md entries, classifies them by heuristic keyword
// matching, deduplicates against prior imports, and promotes new entries
// into the appropriate .context/ files.
//
// Parameters:
//   - cmd: Cobra command for output routing.
//   - dryRun: when true, show the classification plan without writing.
//
// Returns:
//   - error: on discovery, read, state, or promotion failure.
func Run(cmd *cobra.Command, dryRun bool) error {
	contextDir := rc.ContextDir()
	projectRoot := filepath.Dir(contextDir)

	sourcePath, discoverErr := memory.DiscoverMemoryPath(projectRoot)
	if discoverErr != nil {
		write.ErrAutoMemoryNotActive(cmd, discoverErr)
		return ctxerr.MemoryNotFound()
	}

	sourceData, readErr := validation.SafeReadFile(
		filepath.Dir(sourcePath), filepath.Base(sourcePath),
	)
	if readErr != nil {
		return ctxerr.ReadMemory(readErr)
	}

	entries := memory.ParseEntries(string(sourceData))
	if len(entries) == 0 {
		write.ImportNoEntries(cmd, memory2.MemorySource)
		return nil
	}

	state, loadErr := memory.LoadState(contextDir)
	if loadErr != nil {
		return ctxerr.LoadState(loadErr)
	}

	write.ImportScanHeader(cmd, memory2.MemorySource, len(entries))

	var result core.ImportResult

	for _, e := range entries {
		hash := memory.EntryHash(e.Text)

		if state.Imported(hash) {
			result.Dupes++
			continue
		}

		classification := memory.Classify(e)
		title := core.Truncate(e.Text, 60)

		if classification.Target == memory.TargetSkip {
			result.Skipped++
			if dryRun {
				write.ImportEntrySkipped(cmd, title)
			}
			continue
		}

		targetFile := entry.ToCtxFile[classification.Target]

		if dryRun {
			write.ImportEntryClassified(cmd, title, targetFile, classification.Keywords)
		} else {
			if promoteErr := memory.Promote(e, classification); promoteErr != nil {
				write.ErrImportPromote(cmd, targetFile, promoteErr)
				continue
			}
			state.MarkImported(hash, classification.Target)
			write.ImportEntryAdded(cmd, title, targetFile)
		}

		switch classification.Target {
		case entry.Convention:
			result.Conventions++
		case entry.Decision:
			result.Decisions++
		case entry.Learning:
			result.Learnings++
		case entry.Task:
			result.Tasks++
		}
	}

	write.ImportSummary(cmd, write.ImportCounts{
		Conventions: result.Conventions,
		Decisions:   result.Decisions,
		Learnings:   result.Learnings,
		Tasks:       result.Tasks,
		Skipped:     result.Skipped,
		Dupes:       result.Dupes,
	}, dryRun)

	if !dryRun && result.Total() > 0 {
		state.MarkImportedDone()
		if saveErr := memory.SaveState(contextDir, state); saveErr != nil {
			return ctxerr.SaveState(saveErr)
		}
	}

	return nil
}
