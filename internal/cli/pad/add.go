//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// addCmd returns the pad add subcommand.
//
// Returns:
//   - *cobra.Command: Configured add subcommand
func addCmd() *cobra.Command {
	var filePath string

	cmd := &cobra.Command{
		Use:   "add TEXT",
		Short: "Append a new entry to the scratchpad",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if filePath != "" {
				return runAddBlob(cmd, args[0], filePath)
			}
			return runAdd(cmd, args[0])
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "ingest a file as a blob entry")

	return cmd
}

// runAdd appends a new entry and prints confirmation.
func runAdd(cmd *cobra.Command, text string) error {
	entries, err := readEntries()
	if err != nil {
		return err
	}

	entries = append(entries, text)

	if err := writeEntries(entries); err != nil {
		return err
	}

	cmd.Printf("Added entry %d.\n", len(entries))
	return nil
}

// runAddBlob reads a file, encodes it as a blob entry, and appends it.
func runAddBlob(cmd *cobra.Command, label, filePath string) error {
	data, err := os.ReadFile(filePath) //nolint:gosec // user-provided path is intentional
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	if len(data) > MaxBlobSize {
		return fmt.Errorf("file too large: %d bytes (max %d)", len(data), MaxBlobSize)
	}

	entries, err := readEntries()
	if err != nil {
		return err
	}

	entries = append(entries, makeBlob(label, data))

	if err := writeEntries(entries); err != nil {
		return err
	}

	cmd.Printf("Added entry %d.\n", len(entries))
	return nil
}
