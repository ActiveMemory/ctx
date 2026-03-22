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

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/pad"
	ctxErr "github.com/ActiveMemory/ctx/internal/err/fs"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
	writePad "github.com/ActiveMemory/ctx/internal/write/pad"
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
	if file == cli.StdinSentinel {
		r = os.Stdin
	} else {
		f, err := internalIo.SafeOpenUserFile(file)
		if err != nil {
			return ctxErr.OpenFile(file, err)
		}
		defer func() {
			if cErr := f.Close(); cErr != nil {
				writePad.ErrImportCloseWarning(cmd, file, cErr)
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
		return ctxErr.ReadInput(scanErr)
	}

	if count == 0 {
		writePad.ImportNone(cmd)
		return nil
	}

	if writeErr := core.WriteEntries(entries); writeErr != nil {
		return writeErr
	}

	writePad.ImportDone(cmd, count)
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
		return ctxErr.StatPath(path, statErr)
	}
	if !info.IsDir() {
		return ctxErr.NotDirectory(path)
	}

	dirEntries, readErr := os.ReadDir(path)
	if readErr != nil {
		return ctxErr.ReadDirectory(path, readErr)
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

		data, fileErr := internalIo.SafeReadFile(path, name)
		if fileErr != nil {
			writePad.ErrImportBlobSkipped(cmd, name, fileErr)
			skipped++
			continue
		}

		if len(data) > pad.MaxBlobSize {
			writePad.ErrImportBlobTooLarge(cmd, name, pad.MaxBlobSize)
			skipped++
			continue
		}

		entries = append(entries, core.MakeBlob(name, data))
		writePad.ImportBlobAdded(cmd, name)
		added++
	}

	if added > 0 {
		if writeErr := core.WriteEntries(entries); writeErr != nil {
			return writeErr
		}
	}

	writePad.ImportBlobSummary(cmd, added, skipped)
	return nil
}
