//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/crypto"
)

// mergeCmd returns the pad merge subcommand.
//
// Returns:
//   - *cobra.Command: Configured merge subcommand
func mergeCmd() *cobra.Command {
	var keyFile string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "merge FILE...",
		Short: "Merge entries from scratchpad files into the current pad",
		Long: `Merge entries from one or more scratchpad files into the current pad.

Each input file is auto-detected as encrypted or plaintext: decryption is
attempted first, and on failure the file is parsed as plain text. Entries
are deduplicated by exact content — position does not matter.

Use --key to provide a key file for encrypted pads from other projects.

Examples:
  ctx pad merge worktree/.context/scratchpad.enc
  ctx pad merge notes.md backup.enc
  ctx pad merge --key /other/.context.key foreign.enc
  ctx pad merge --dry-run pad-a.enc pad-b.md`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMerge(cmd, args, keyFile, dryRun)
		},
	}

	cmd.Flags().StringVarP(&keyFile, "key", "k", "",
		"path to key file for decrypting input files")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false,
		"print what would be merged without writing")

	return cmd
}

// runMerge reads entries from input files, deduplicates against the current
// pad, and writes the merged result.
//
// Parameters:
//   - cmd: Cobra command for output
//   - files: Input file paths to merge
//   - keyFile: Optional path to key file (empty = use project key)
//   - dryRun: If true, print summary without writing
//
// Returns:
//   - error: Non-nil on read/write failures
func runMerge(
	cmd *cobra.Command,
	files []string,
	keyFile string,
	dryRun bool,
) error {
	current, readErr := readEntries()
	if readErr != nil {
		return readErr
	}

	key := loadMergeKey(keyFile)

	seen := make(map[string]bool, len(current))
	for _, e := range current {
		seen[e] = true
	}

	blobLabels := buildBlobLabelMap(current)

	var added, dupes int
	var newEntries []string

	for _, file := range files {
		entries, fileErr := readFileEntries(file, key)
		if fileErr != nil {
			return fmt.Errorf("open %s: %w", file, fileErr)
		}

		warnIfBinary(cmd, file, entries)

		for _, entry := range entries {
			if seen[entry] {
				dupes++
				cmd.Printf(
					"  = %-40s (duplicate, skipped)\n",
					displayEntry(entry),
				)
				continue
			}
			seen[entry] = true
			checkBlobConflict(cmd, entry, blobLabels)
			newEntries = append(newEntries, entry)
			added++
			cmd.Printf(
				"  + %-40s (from %s)\n",
				displayEntry(entry),
				file,
			)
		}
	}

	if added == 0 && dupes == 0 {
		cmd.Println("No entries to merge.")
		return nil
	}

	if added == 0 {
		cmd.Printf(
			"No new entries to merge (%d %s skipped).\n",
			dupes,
			pluralize("duplicate", dupes),
		)
		return nil
	}

	if dryRun {
		cmd.Printf(
			"Would merge %d new %s (%d %s skipped).\n",
			added,
			pluralize("entry", added),
			dupes,
			pluralize("duplicate", dupes),
		)
		return nil
	}

	merged := make([]string, 0, len(current)+len(newEntries))
	merged = append(merged, current...)
	merged = append(merged, newEntries...)
	if writeErr := writeEntries(merged); writeErr != nil {
		return writeErr
	}

	cmd.Printf(
		"Merged %d new %s (%d %s skipped).\n",
		added,
		pluralize("entry", added),
		dupes,
		pluralize("duplicate", dupes),
	)
	return nil
}

// readFileEntries reads a scratchpad file, attempting decryption first.
//
// If a key is provided and decryption succeeds, the decrypted entries are
// returned. Otherwise the file is parsed as plaintext.
//
// Parameters:
//   - path: Path to the scratchpad file
//   - key: Encryption key (nil to skip decryption attempt)
//
// Returns:
//   - []string: Parsed entries
//   - error: Non-nil if the file cannot be read
func readFileEntries(path string, key []byte) ([]string, error) {
	data, readErr := os.ReadFile(path) //nolint:gosec // user-provided path is intentional
	if readErr != nil {
		return nil, readErr
	}

	if len(data) == 0 {
		return nil, nil
	}

	if key != nil {
		plaintext, decErr := crypto.Decrypt(key, data)
		if decErr == nil {
			return parseEntries(plaintext), nil
		}
	}

	return parseEntries(data), nil
}

// loadMergeKey loads the encryption key for merge input decryption.
//
// Priority: explicit --key flag > project scratchpad key > nil (no key).
//
// Parameters:
//   - keyFile: Explicit key file path (empty string = use project key)
//
// Returns:
//   - []byte: The loaded key, or nil if no key is available
func loadMergeKey(keyFile string) []byte {
	if keyFile != "" {
		key, loadErr := crypto.LoadKey(keyFile)
		if loadErr != nil {
			return nil
		}
		return key
	}

	key, loadErr := crypto.LoadKey(keyPath())
	if loadErr != nil {
		return nil
	}
	return key
}

// buildBlobLabelMap creates a map of blob labels to their full entry strings
// from the given entries.
//
// Parameters:
//   - entries: Scratchpad entries to scan
//
// Returns:
//   - map[string]string: Blob label → full entry string
func buildBlobLabelMap(entries []string) map[string]string {
	labels := make(map[string]string)
	for _, entry := range entries {
		if label, _, ok := splitBlob(entry); ok {
			labels[label] = entry
		}
	}
	return labels
}

// checkBlobConflict warns if a blob entry has the same label as an existing
// blob but different content.
//
// Parameters:
//   - cmd: Cobra command for output
//   - entry: The new entry to check
//   - blobLabels: Map of existing blob labels to their full entry strings
func checkBlobConflict(
	cmd *cobra.Command,
	entry string,
	blobLabels map[string]string,
) {
	label, _, ok := splitBlob(entry)
	if !ok {
		return
	}

	existing, found := blobLabels[label]
	if found && existing != entry {
		cmd.Printf(
			"  ! blob %q has different content across sources; both kept\n",
			label,
		)
	}

	blobLabels[label] = entry
}

// warnIfBinary prints a warning if any entries contain non-UTF-8 bytes,
// which suggests the file is encrypted but was parsed as plaintext.
//
// Parameters:
//   - cmd: Cobra command for output
//   - file: The source file path (for the warning message)
//   - entries: The parsed entries to check
func warnIfBinary(cmd *cobra.Command, file string, entries []string) {
	for _, entry := range entries {
		if !utf8.ValidString(entry) {
			cmd.Printf(
				"  ! %s appears to contain binary data;"+
					" it may be encrypted (use --key)\n",
				file,
			)
			return
		}
	}
}

// pluralize returns the singular or plural form of a word.
//
// Parameters:
//   - word: The singular form
//   - count: The count to check
//
// Returns:
//   - string: Singular form if count == 1, otherwise plural
func pluralize(word string, count int) string {
	if count == 1 {
		return word
	}
	if strings.HasSuffix(word, "y") {
		return word[:len(word)-1] + "ies"
	}
	return word + "s"
}
