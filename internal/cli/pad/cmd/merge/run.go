//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package merge

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core"
	"github.com/ActiveMemory/ctx/internal/crypto"
)

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
	current, readErr := core.ReadEntries()
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
				cmd.Println(fmt.Sprintf(
					"  = %-40s (duplicate, skipped)\n",
					core.DisplayEntry(entry),
				))
				continue
			}
			seen[entry] = true
			checkBlobConflict(cmd, entry, blobLabels)
			newEntries = append(newEntries, entry)
			added++
			cmd.Println(fmt.Sprintf(
				"  + %-40s (from %s)\n",
				core.DisplayEntry(entry),
				file,
			))
		}
	}

	if added == 0 && dupes == 0 {
		cmd.Println("No entries to merge.")
		return nil
	}

	if added == 0 {
		cmd.Println(fmt.Sprintf(
			"No new entries to merge (%d %s skipped).\n",
			dupes,
			pluralize("duplicate", dupes),
		))
		return nil
	}

	if dryRun {
		cmd.Println(fmt.Sprintf(
			"Would merge %d new %s (%d %s skipped).\n",
			added,
			pluralize("entry", added),
			dupes,
			pluralize("duplicate", dupes),
		))
		return nil
	}

	merged := make([]string, 0, len(current)+len(newEntries))
	merged = append(merged, current...)
	merged = append(merged, newEntries...)
	if writeErr := core.WriteEntries(merged); writeErr != nil {
		return writeErr
	}

	cmd.Println(fmt.Sprintf(
		"Merged %d new %s (%d %s skipped).\n",
		added,
		pluralize("entry", added),
		dupes,
		pluralize("duplicate", dupes),
	))
	return nil
}

// readFileEntries reads a scratchpad file, attempting decryption first.
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
			return core.ParseEntries(plaintext), nil
		}
	}

	return core.ParseEntries(data), nil
}

// loadMergeKey loads the encryption key for merge input decryption.
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

	key, loadErr := crypto.LoadKey(core.KeyPath())
	if loadErr != nil {
		return nil
	}
	return key
}

// buildBlobLabelMap creates a map of blob labels to their full entry strings.
//
// Parameters:
//   - entries: Scratchpad entries to scan
//
// Returns:
//   - map[string]string: Blob label to full entry string
func buildBlobLabelMap(entries []string) map[string]string {
	labels := make(map[string]string)
	for _, entry := range entries {
		if label, _, ok := core.SplitBlob(entry); ok {
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
	label, _, ok := core.SplitBlob(entry)
	if !ok {
		return
	}

	existing, found := blobLabels[label]
	if found && existing != entry {
		cmd.Println(fmt.Sprintf(
			"  ! blob %q has different content across sources; both kept\n",
			label,
		))
	}

	blobLabels[label] = entry
}

// warnIfBinary prints a warning if any entries contain non-UTF-8 bytes.
//
// Parameters:
//   - cmd: Cobra command for output
//   - file: The source file path (for the warning message)
//   - entries: The parsed entries to check
func warnIfBinary(cmd *cobra.Command, file string, entries []string) {
	for _, entry := range entries {
		if !utf8.ValidString(entry) {
			cmd.Println(fmt.Sprintf(
				"  ! %s appears to contain binary data;"+
					" it may be encrypted (use --key)\n",
				file,
			))
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
