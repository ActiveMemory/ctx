//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// importCmd returns the pad import subcommand.
func importCmd() *cobra.Command {
	var blobs bool

	cmd := &cobra.Command{
		Use:   "import FILE",
		Short: "Bulk-import lines from a file into the scratchpad",
		Long: `Import lines from a file into the scratchpad. Each non-empty line
becomes a separate entry. Use "-" to read from stdin.

With --blobs, import all first-level files from a directory as blob entries.
Each file becomes a blob with the filename as its label. Subdirectories and
non-regular files are skipped.

Examples:
  ctx pad import notes.txt
  grep pattern file | ctx pad import -
  ctx pad import --blobs ./ideas/`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if blobs {
				return runImportBlobs(cmd, args[0])
			}
			return runImport(cmd, args[0])
		},
	}

	cmd.Flags().BoolVar(&blobs, "blobs", false,
		"import first-level files from a directory as blob entries")

	return cmd
}

// runImport reads lines from a file (or stdin) and appends them as entries.
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

	entries, err := readEntries()
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
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read input: %w", err)
	}

	if count == 0 {
		cmd.Println("No entries to import.")
		return nil
	}

	if err := writeEntries(entries); err != nil {
		return err
	}

	cmd.Printf("Imported %d entries.\n", count)
	return nil
}

// runImportBlobs reads first-level files from a directory and imports
// each as a blob entry. All entries are written in a single read/write
// cycle.
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

	entries, loadErr := readEntries()
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
			cmd.PrintErrf("  ! skipped: %s (%v)\n", name, fileErr)
			skipped++
			continue
		}

		if len(data) > MaxBlobSize {
			cmd.PrintErrf("  ! skipped: %s (exceeds %d byte limit)\n",
				name, MaxBlobSize)
			skipped++
			continue
		}

		entries = append(entries, makeBlob(name, data))
		cmd.Printf("  + %s\n", name)
		added++
	}

	if added == 0 && skipped == 0 {
		cmd.Println("No files to import.")
		return nil
	}

	if added > 0 {
		if writeErr := writeEntries(entries); writeErr != nil {
			return writeErr
		}
	}

	cmd.Printf("Done. Added %d, skipped %d.\n", added, skipped)
	return nil
}
