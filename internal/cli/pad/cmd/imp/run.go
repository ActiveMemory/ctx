//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package imp

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core"
)

// runImport reads lines from a file (or stdin) and appends them as entries.
//
// Parameters:
//   - cmd: Cobra command for output
//   - file: File path or "-" for stdin
//
// Returns:
//   - error: Non-nil on read/write failure
func runImport(cmd *cobra.Command, file string) error {
	var r io.Reader
	if file == "-" {
		r = os.Stdin
	} else {
		f, err := os.Open(file) //nolint:gosec // user-provided path is intentional
		if err != nil {
			return fmt.Errorf("open %s: %w", file, err)
		}
		defer func() {
			if cerr := f.Close(); cerr != nil {
				fmt.Fprintf(os.Stderr, "warning: close %s: %v\n", file, cerr)
			}
		}()
		r = f
	}

	entries, err := core.ReadEntries()
	if err != nil {
		return err
	}

	var count int
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		entries = append(entries, line)
		count++
	}
	if scanErr := scanner.Err(); scanErr != nil {
		return fmt.Errorf("read input: %w", scanErr)
	}

	if count == 0 {
		cmd.Println("No entries to import.")
		return nil
	}

	if writeErr := core.WriteEntries(entries); writeErr != nil {
		return writeErr
	}

	cmd.Println(fmt.Sprintf("Imported %d entries.", count))
	return nil
}

// runImportBlobs reads first-level files from a directory and imports
// each as a blob entry.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: Directory path containing files to import
//
// Returns:
//   - error: Non-nil on read/write failure
func runImportBlobs(cmd *cobra.Command, path string) error {
	info, statErr := os.Stat(path)
	if statErr != nil {
		return fmt.Errorf("stat %s: %w", path, statErr)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	dirEntries, readErr := os.ReadDir(path)
	if readErr != nil {
		return fmt.Errorf("read directory %s: %w", path, readErr)
	}

	entries, loadErr := core.ReadEntries()
	if loadErr != nil {
		return loadErr
	}

	var added, skipped int
	for _, de := range dirEntries {
		if !de.Type().IsRegular() {
			continue
		}

		name := de.Name()
		filePath := filepath.Join(path, name)

		data, fileErr := os.ReadFile(filePath) //nolint:gosec // user-provided path is intentional
		if fileErr != nil {
			cmd.PrintErrln(fmt.Sprintf("  ! skipped: %s (%v)", name, fileErr))
			skipped++
			continue
		}

		if len(data) > core.MaxBlobSize {
			cmd.PrintErrln(fmt.Sprintf("  ! skipped: %s (exceeds %d byte limit)",
				name, core.MaxBlobSize))
			skipped++
			continue
		}

		entries = append(entries, core.MakeBlob(name, data))
		cmd.Println(fmt.Sprintf("  + %s", name))
		added++
	}

	if added == 0 && skipped == 0 {
		cmd.Println("No files to import.")
		return nil
	}

	if added > 0 {
		if writeErr := core.WriteEntries(entries); writeErr != nil {
			return writeErr
		}
	}

	cmd.Println(fmt.Sprintf("Done. Added %d, skipped %d.", added, skipped))
	return nil
}
