//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core"
)

// runAdd appends a new entry and prints confirmation.
//
// Parameters:
//   - cmd: Cobra command for output
//   - text: Entry text to add
//
// Returns:
//   - error: Non-nil on read/write failure
func runAdd(cmd *cobra.Command, text string) error {
	entries, err := core.ReadEntries()
	if err != nil {
		return err
	}

	entries = append(entries, text)

	if writeErr := core.WriteEntries(entries); writeErr != nil {
		return writeErr
	}

	cmd.Println(fmt.Sprintf("Added entry %d.", len(entries)))
	return nil
}

// runAddBlob reads a file, encodes it as a blob entry, and appends it.
//
// Parameters:
//   - cmd: Cobra command for output
//   - label: Blob label (filename)
//   - filePath: Path to the file to ingest
//
// Returns:
//   - error: Non-nil on read/write failure or file too large
func runAddBlob(cmd *cobra.Command, label, filePath string) error {
	data, err := os.ReadFile(filePath) //nolint:gosec // user-provided path is intentional
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	if len(data) > core.MaxBlobSize {
		return fmt.Errorf("file too large: %d bytes (max %d)", len(data), core.MaxBlobSize)
	}

	entries, readErr := core.ReadEntries()
	if readErr != nil {
		return readErr
	}

	entries = append(entries, core.MakeBlob(label, data))

	if writeErr := core.WriteEntries(entries); writeErr != nil {
		return writeErr
	}

	cmd.Println(fmt.Sprintf("Added entry %d.", len(entries)))
	return nil
}
