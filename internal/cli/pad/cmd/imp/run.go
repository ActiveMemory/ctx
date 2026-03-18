//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package imp

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/pad"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/fs"
	io2 "github.com/ActiveMemory/ctx/internal/io"
	pad2 "github.com/ActiveMemory/ctx/internal/write/pad"
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
		f, err := io2.SafeOpenUserFile(file)
		if err != nil {
			return ctxerr.OpenFile(file, err)
		}
		defer func() {
			if cerr := f.Close(); cerr != nil {
				pad2.ErrPadImportCloseWarning(cmd, file, cerr)
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
		return ctxerr.ReadInput(scanErr)
	}

	if count == 0 {
		pad2.PadImportNone(cmd)
		return nil
	}

	if writeErr := core.WriteEntries(entries); writeErr != nil {
		return writeErr
	}

	pad2.PadImportDone(cmd, count)
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
		return ctxerr.StatPath(path, statErr)
	}
	if !info.IsDir() {
		return ctxerr.NotDirectory(path)
	}

	dirEntries, readErr := os.ReadDir(path)
	if readErr != nil {
		return ctxerr.ReadDirectory(path, readErr)
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

		data, fileErr := io2.SafeReadFile(path, name)
		if fileErr != nil {
			pad2.ErrPadImportBlobSkipped(cmd, name, fileErr)
			skipped++
			continue
		}

		if len(data) > pad.MaxBlobSize {
			pad2.ErrPadImportBlobTooLarge(cmd, name, pad.MaxBlobSize)
			skipped++
			continue
		}

		entries = append(entries, core.MakeBlob(name, data))
		pad2.PadImportBlobAdded(cmd, name)
		added++
	}

	if added > 0 {
		if writeErr := core.WriteEntries(entries); writeErr != nil {
			return writeErr
		}
	}

	pad2.PadImportBlobSummary(cmd, added, skipped)
	return nil
}
