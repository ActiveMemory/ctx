//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write"
)

// runSync discovers MEMORY.md, mirrors it into .context/memory/, and
// updates the sync state. In dry-run mode it reports what would happen
// without writing any files.
//
// Parameters:
//   - cmd: Cobra command for output routing.
//   - dryRun: when true, report the plan without writing.
//
// Returns:
//   - error: on discovery failure, sync failure, or state persistence failure.
func runSync(cmd *cobra.Command, dryRun bool) error {
	contextDir := rc.ContextDir()
	projectRoot := filepath.Dir(contextDir)

	sourcePath, discoverErr := memory.DiscoverMemoryPath(projectRoot)
	if discoverErr != nil {
		write.ErrAutoMemoryNotActive(cmd, discoverErr)
		return ctxerr.MemoryNotFound()
	}

	if dryRun {
		write.DryRun(cmd)
		write.Source(cmd, sourcePath)
		write.Mirror(cmd, config.PathMemoryMirror)
		if memory.HasDrift(contextDir, sourcePath) {
			write.StatusDrift(cmd)
		} else {
			write.StatusNoDrift(cmd)
		}
		return nil
	}

	result, syncErr := memory.Sync(contextDir, sourcePath)
	if syncErr != nil {
		return ctxerr.SyncFailed(syncErr)
	}

	if result.ArchivedTo != "" {
		write.Archived(cmd, filepath.Base(result.ArchivedTo))
	}

	write.Synced(cmd, config.FileMemorySource, config.PathMemoryMirror)
	write.Source(cmd, result.SourcePath)
	write.Lines(cmd, result.SourceLines, result.MirrorLines)

	if result.SourceLines > result.MirrorLines {
		write.NewContent(cmd, result.SourceLines-result.MirrorLines)
	}

	// Update sync state
	state, loadErr := memory.LoadState(contextDir)
	if loadErr != nil {
		return ctxerr.LoadState(loadErr)
	}
	state.MarkSynced()
	if saveErr := memory.SaveState(contextDir, state); saveErr != nil {
		return ctxerr.SaveState(saveErr)
	}

	return nil
}

// runUnpublish removes the ctx-managed marker block from MEMORY.md,
// preserving all Claude-owned content outside the markers.
//
// Parameters:
//   - cmd: Cobra command for output routing.
//
// Returns:
//   - error: on discovery, read, or write failure.
func runUnpublish(cmd *cobra.Command) error {
	contextDir := rc.ContextDir()
	projectRoot := filepath.Dir(contextDir)

	memoryPath, discoverErr := memory.DiscoverMemoryPath(projectRoot)
	if discoverErr != nil {
		write.ErrAutoMemoryNotActive(cmd, discoverErr)
		return ctxerr.MemoryNotFound()
	}

	data, readErr := os.ReadFile(memoryPath) //nolint:gosec // discovered path
	if readErr != nil {
		return ctxerr.ReadMemory(readErr)
	}

	cleaned, found := memory.RemovePublished(string(data))
	if !found {
		cmd.Println("No published block found in " + config.FileMemorySource + ".")
		return nil
	}

	if writeErr := os.WriteFile(memoryPath, []byte(cleaned), config.PermFile); writeErr != nil {
		return ctxerr.WriteMemory(writeErr)
	}

	cmd.Println("Removed published block from " + config.FileMemorySource + ".")
	return nil
}

// runImport parses MEMORY.md entries, classifies them by heuristic keyword
// matching, deduplicates against prior imports, and promotes new entries
// into the appropriate .context/ files.
//
// Parameters:
//   - cmd: Cobra command for output routing.
//   - dryRun: when true, show the classification plan without writing.
//
// Returns:
//   - error: on discovery, read, state, or promotion failure.
func runImport(cmd *cobra.Command, dryRun bool) error {
	contextDir := rc.ContextDir()
	projectRoot := filepath.Dir(contextDir)

	sourcePath, discoverErr := memory.DiscoverMemoryPath(projectRoot)
	if discoverErr != nil {
		write.ErrAutoMemoryNotActive(cmd, discoverErr)
		return ctxerr.MemoryNotFound()
	}

	sourceData, readErr := os.ReadFile(sourcePath) //nolint:gosec // discovered path
	if readErr != nil {
		return ctxerr.ReadMemory(readErr)
	}

	entries := memory.ParseEntries(string(sourceData))
	if len(entries) == 0 {
		cmd.Println("No entries found in " + config.FileMemorySource + ".")
		return nil
	}

	state, loadErr := memory.LoadState(contextDir)
	if loadErr != nil {
		return ctxerr.LoadState(loadErr)
	}

	cmd.Println("Scanning " + config.FileMemorySource + " for new entries...")
	cmd.Println(fmt.Sprintf("  Found %d entries", len(entries)))
	cmd.Println()

	var result importResult

	for _, entry := range entries {
		hash := memory.EntryHash(entry.Text)

		if state.Imported(hash) {
			result.dupes++
			continue
		}

		classification := memory.Classify(entry)
		title := truncate(entry.Text, 60)

		if classification.Target == memory.TargetSkip {
			result.skipped++
			if dryRun {
				cmd.Println(fmt.Sprintf("  -> %q", title))
				cmd.Println("     Classified: skip")
				cmd.Println()
			}
			continue
		}

		targetFile := config.FileType[classification.Target]

		if dryRun {
			cmd.Println(fmt.Sprintf("  -> %q", title))
			cmd.Println(fmt.Sprintf("     Classified: %s (keywords: %s)",
				targetFile, strings.Join(classification.Keywords, ", ")))
			cmd.Println()
		} else {
			if promoteErr := memory.Promote(entry, classification); promoteErr != nil {
				cmd.PrintErrln(fmt.Sprintf("  Error promoting to %s: %v", targetFile, promoteErr))
				continue
			}
			state.MarkImported(hash, classification.Target)
			cmd.Println(fmt.Sprintf("  -> %q", title))
			cmd.Println(fmt.Sprintf("     Added to %s", targetFile))
			cmd.Println()
		}

		switch classification.Target {
		case config.EntryConvention:
			result.conventions++
		case config.EntryDecision:
			result.decisions++
		case config.EntryLearning:
			result.learnings++
		case config.EntryTask:
			result.tasks++
		}
	}

	// Summary
	var summary string
	if dryRun {
		summary = fmt.Sprintf("Dry run — would import: %d entries", result.total())
	} else {
		summary = fmt.Sprintf("Imported: %d entries", result.total())
	}

	var parts []string
	if result.conventions > 0 {
		parts = append(parts, fmt.Sprintf("%d convention", result.conventions))
	}
	if result.decisions > 0 {
		parts = append(parts, fmt.Sprintf("%d decision", result.decisions))
	}
	if result.learnings > 0 {
		parts = append(parts, fmt.Sprintf("%d learning", result.learnings))
	}
	if result.tasks > 0 {
		parts = append(parts, fmt.Sprintf("%d task", result.tasks))
	}
	if len(parts) > 0 {
		summary += fmt.Sprintf(" (%s)", strings.Join(parts, ", "))
	}
	cmd.Println(summary)

	if result.skipped > 0 {
		cmd.Println(fmt.Sprintf("Skipped: %d entries (session notes/unclassified)", result.skipped))
	}
	if result.dupes > 0 {
		cmd.Println(fmt.Sprintf("Duplicates: %d entries (already imported)", result.dupes))
	}

	if !dryRun && result.total() > 0 {
		state.MarkImportedDone()
		if saveErr := memory.SaveState(contextDir, state); saveErr != nil {
			return ctxerr.SaveState(saveErr)
		}
	}

	return nil
}

// truncate returns the first line of s, capped at max characters.
//
// Parameters:
//   - s: input string (may be multi-line).
//   - max: maximum length including ellipsis.
//
// Returns:
//   - string: truncated first line.
func truncate(s string, max int) string {
	line := strings.SplitN(s, config.NewlineLF, 2)[0]
	if len(line) <= max {
		return line
	}
	return line[:max-3] + "..."
}
