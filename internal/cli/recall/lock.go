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

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// lockedFrontmatterLine is the YAML line inserted into frontmatter when
// a journal entry is locked.
const lockedFrontmatterLine = "locked: true  # managed by ctx"

// recallLockCmd returns the "ctx recall lock" subcommand.
//
// Protects journal entries from being overwritten by export --regenerate.
// Locked entries are skipped during export regardless of flags.
//
// Returns:
//   - *cobra.Command: Command for locking journal entries
func recallLockCmd() *cobra.Command {
	var all bool

	cmd := &cobra.Command{
		Use:   "lock <pattern>",
		Short: "Protect journal entries from export regeneration",
		Long: `Lock journal entries to prevent export --regenerate from overwriting them.

Locked entries are skipped during export regardless of --regenerate or --force.
Use "ctx recall unlock" to remove the protection.

The pattern matches against filenames by slug, date, or short ID (same
matching as export). Locking a multi-part entry locks all parts.

The lock is recorded in .context/journal/.state.json (source of truth) and
a "locked: true" line is added to the file's YAML frontmatter for visibility.

Examples:
  ctx recall lock 2026-01-21-session-abc12345.md
  ctx recall lock abc12345
  ctx recall lock --all`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLockUnlock(cmd, args, all, true)
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "Lock all journal entries")

	return cmd
}

// recallUnlockCmd returns the "ctx recall unlock" subcommand.
//
// Removes lock protection from journal entries, allowing export
// --regenerate to overwrite them again.
//
// Returns:
//   - *cobra.Command: Command for unlocking journal entries
func recallUnlockCmd() *cobra.Command {
	var all bool

	cmd := &cobra.Command{
		Use:   "unlock <pattern>",
		Short: "Remove lock protection from journal entries",
		Long: `Unlock journal entries to allow export --regenerate to overwrite them.

The pattern matches against filenames by slug, date, or short ID (same
matching as export). Unlocking a multi-part entry unlocks all parts.

Examples:
  ctx recall unlock 2026-01-21-session-abc12345.md
  ctx recall unlock abc12345
  ctx recall unlock --all`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLockUnlock(cmd, args, all, false)
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "Unlock all journal entries")

	return cmd
}

// runLockUnlock handles both lock and unlock commands.
//
// Parameters:
//   - cmd: Cobra command for output
//   - args: Patterns to match against journal filenames
//   - all: If true, apply to all journal entries
//   - lock: True for lock, false for unlock
//
// Returns:
//   - error: Non-nil on validation or I/O failure
func runLockUnlock(
	cmd *cobra.Command,
	args []string,
	all, lock bool,
) error {
	if len(args) == 0 && !all {
		return cmd.Help()
	}
	if len(args) > 0 && all {
		return fmt.Errorf(
			"cannot use --all with a pattern; use one or the other",
		)
	}

	journalDir := filepath.Join(rc.ContextDir(), config.DirJournal)

	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		return fmt.Errorf("load journal state: %w", loadErr)
	}

	// Collect matching .md files.
	files, matchErr := matchJournalFiles(journalDir, args, all)
	if matchErr != nil {
		return matchErr
	}
	if len(files) == 0 {
		if all {
			cmd.Println("No journal entries found.")
		} else {
			return fmt.Errorf(
				"no journal entries match: %s", strings.Join(args, ", "),
			)
		}
		return nil
	}

	green := color.New(color.FgGreen).SprintFunc()
	verb := "locked"
	if !lock {
		verb = "unlocked"
	}

	count := 0
	for _, filename := range files {
		alreadyLocked := jstate.Locked(filename)
		if lock && alreadyLocked {
			continue
		}
		if !lock && !alreadyLocked {
			continue
		}

		// Update state.
		if lock {
			jstate.Mark(filename, "locked")
		} else {
			jstate.Clear(filename, "locked")
		}

		// Update frontmatter for human visibility.
		path := filepath.Join(journalDir, filename)
		updateLockFrontmatter(path, lock)

		cmd.Printf("  %s %s (%s)\n", green("✓"), filename, verb)
		count++
	}

	if saveErr := jstate.Save(journalDir); saveErr != nil {
		return fmt.Errorf("save journal state: %w", saveErr)
	}

	if count == 0 {
		cmd.Printf("No changes — all matched entries already %s.\n", verb)
	} else {
		cmd.Printf("\n%s %d entry(s).\n", strings.Title(verb), count) //nolint:staticcheck // strings.Title is fine for single words
	}

	return nil
}

// matchJournalFiles returns journal .md filenames matching the given
// patterns. If all is true, returns every .md file in the directory.
// Multi-part files (base + -pN parts) are included when the base matches.
//
// Parameters:
//   - journalDir: Path to the journal directory
//   - patterns: Slug, date, or short-ID substrings to match
//   - all: If true, return all .md files
//
// Returns:
//   - []string: Matching filenames (basename only)
//   - error: Non-nil on I/O failure
func matchJournalFiles(
	journalDir string,
	patterns []string,
	all bool,
) ([]string, error) {
	entries, readErr := os.ReadDir(journalDir)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return nil, nil
		}
		return nil, fmt.Errorf("read journal directory: %w", readErr)
	}

	// Collect all .md filenames.
	var mdFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), config.ExtMarkdown) {
			mdFiles = append(mdFiles, e.Name())
		}
	}

	if all {
		return mdFiles, nil
	}

	// Build a set of matching base names, then expand to include parts.
	matchedBases := make(map[string]bool)
	for _, f := range mdFiles {
		lower := strings.ToLower(f)
		for _, pat := range patterns {
			if strings.Contains(lower, strings.ToLower(pat)) {
				base := multipartBase(f)
				matchedBases[base] = true
			}
		}
	}

	// Expand: include all files sharing a matched base.
	var result []string
	for _, f := range mdFiles {
		base := multipartBase(f)
		if matchedBases[base] {
			result = append(result, f)
		}
	}

	return result, nil
}

// multipartBase returns the base name for a potentially multi-part file.
// For "2026-01-21-slug-abc12345-p2.md" it returns
// "2026-01-21-slug-abc12345.md". For non-multipart files, returns the
// filename as-is.
//
// Parameters:
//   - filename: Journal entry filename
//
// Returns:
//   - string: Base filename (without -pN suffix)
func multipartBase(filename string) string {
	base := strings.TrimSuffix(filename, config.ExtMarkdown)
	if idx := strings.LastIndex(base, "-p"); idx > 0 {
		suffix := base[idx+2:]
		allDigits := true
		for _, r := range suffix {
			if r < '0' || r > '9' {
				allDigits = false
				break
			}
		}
		if allDigits && len(suffix) > 0 {
			return base[:idx] + config.ExtMarkdown
		}
	}
	return filename
}

// updateLockFrontmatter inserts or removes the "locked: true" line in
// a journal file's YAML frontmatter. The state file is the source of
// truth; this is for human visibility only.
//
// Parameters:
//   - path: Absolute path to the journal .md file
//   - lock: True to insert, false to remove
func updateLockFrontmatter(path string, lock bool) {
	data, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		return
	}
	content := string(data)

	nl := config.NewlineLF
	fmOpen := config.Separator + nl

	if !strings.HasPrefix(content, fmOpen) {
		// No frontmatter — nothing to modify.
		return
	}

	closeIdx := strings.Index(content[len(fmOpen):], nl+config.Separator+nl)
	if closeIdx < 0 {
		return
	}

	fmEnd := len(fmOpen) + closeIdx // index of the newline before closing ---
	fmBlock := content[len(fmOpen):fmEnd]

	if lock {
		// Already has locked line?
		if strings.Contains(fmBlock, "locked:") {
			return
		}
		// Insert before closing ---.
		updated := content[:fmEnd] + nl + lockedFrontmatterLine +
			content[fmEnd:]
		_ = os.WriteFile(path, []byte(updated), config.PermFile)
	} else {
		// Remove the locked line.
		lines := strings.Split(fmBlock, nl)
		var filtered []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "locked:") {
				continue
			}
			filtered = append(filtered, line)
		}
		newFM := strings.Join(filtered, nl)
		updated := content[:len(fmOpen)] + newFM + content[fmEnd:]
		_ = os.WriteFile(path, []byte(updated), config.PermFile)
	}
}
