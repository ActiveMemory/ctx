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
	"strings"

	"github.com/spf13/cobra"
)

// importCmd returns the pad import subcommand.
func importCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import FILE",
		Short: "Bulk-import lines from a file into the scratchpad",
		Long: `Import lines from a file into the scratchpad. Each non-empty line
becomes a separate entry. Use "-" to read from stdin.

Examples:
  ctx pad import notes.txt
  grep pattern file | ctx pad import -`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runImport(cmd, args[0])
		},
	}

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
