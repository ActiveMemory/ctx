//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package importer

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/entry"
	cfgMemory "github.com/ActiveMemory/ctx/internal/config/memory"
	"github.com/ActiveMemory/ctx/internal/entity"
	errMemory "github.com/ActiveMemory/ctx/internal/err/memory"
	errState "github.com/ActiveMemory/ctx/internal/err/state"
	"github.com/ActiveMemory/ctx/internal/format"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write/ctximport"
	"github.com/ActiveMemory/ctx/internal/write/sync"
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

	sourcePath, discoverErr := memory.DiscoverPath(projectRoot)
	if discoverErr != nil {
		sync.ErrAutoMemoryNotActive(cmd, discoverErr)
		return errMemory.NotFound()
	}

	sourceData, readErr := io.SafeReadFile(
		filepath.Dir(sourcePath), filepath.Base(sourcePath),
	)
	if readErr != nil {
		return errMemory.Read(readErr)
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
		title := format.TruncateFirstLine(e.Text, 60)

		if classification.Target == memory.TargetSkip {
			result.Skipped++
			if dryRun {
				ctximport.EntrySkipped(cmd, title)
			}
			continue
		}

		targetFile := entry.ToCtxFile[classification.Target]

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
