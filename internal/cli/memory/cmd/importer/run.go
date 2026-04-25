//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package importer

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/memory/core/resolve"
	"github.com/ActiveMemory/ctx/internal/config/entry"
	cfgFmt "github.com/ActiveMemory/ctx/internal/config/format"
	cfgMemory "github.com/ActiveMemory/ctx/internal/config/memory"
	"github.com/ActiveMemory/ctx/internal/entity"
	errState "github.com/ActiveMemory/ctx/internal/err/state"
	"github.com/ActiveMemory/ctx/internal/format"
	"github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/write/ctximport"
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
	contextDir, projectRoot, err := resolve.ContextAndRoot(cmd)
	if err != nil {
		return err
	}

	sourcePath, discoverErr := resolve.DiscoverSource(cmd, projectRoot)
	if discoverErr != nil {
		return discoverErr
	}
	sourceData, readErr := resolve.ReadSource(sourcePath)
	if readErr != nil {
		return readErr
	}

	entries := memory.Entries(string(sourceData))
	if len(entries) == 0 {
		ctximport.NoEntries(cmd, cfgMemory.Source)
		return nil
	}

	state, loadErr := memory.LoadState(contextDir)
	if loadErr != nil {
		return errState.Load(loadErr)
	}

	ctximport.ScanHeader(cmd, cfgMemory.Source, len(entries))

	var result entity.ImportResult

	for _, e := range entries {
		hash := memory.EntryHash(e.Text)

		if state.Imported(hash) {
			result.Dupes++
			continue
		}

		classification := memory.Classify(e)
		title := format.TruncateFirstLine(e.Text, cfgFmt.TruncateTitle)

		if classification.Target == cfgMemory.TargetSkip {
			result.Skipped++
			if dryRun {
				ctximport.EntrySkipped(cmd, title)
			}
			continue
		}

		targetFile := entry.MustCtxFile(classification.Target)

		if dryRun {
			ctximport.EntryClassified(cmd, title, targetFile, classification.Keywords)
		} else {
			if promoteErr := memory.Promote(e, classification); promoteErr != nil {
				ctximport.ErrPromote(cmd, targetFile, promoteErr)
				continue
			}
			state.MarkImported(hash, classification.Target)
			ctximport.EntryAdded(cmd, title, targetFile)
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

	ctximport.Summary(cmd, result, dryRun)

	if !dryRun && result.Total() > 0 {
		state.MarkImportedDone()
		if saveErr := memory.SaveState(contextDir, state); saveErr != nil {
			return errState.Save(saveErr)
		}
	}

	return nil
}
