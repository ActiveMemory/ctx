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

	"github.com/ActiveMemory/ctx/internal/cli/pad/core/blob"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/store"
	"github.com/ActiveMemory/ctx/internal/config/pad"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
)

// FromReader reads non-empty lines from r, appends them to the current
// scratchpad entries, and returns the updated list with the count of
// new entries. The caller owns writing the result.
//
// Parameters:
//   - r: Input reader (file or stdin)
//
// Returns:
//   - []string: Updated entries list (existing + imported)
//   - int: Number of new entries added
//   - error: Non-nil on entry load or read failure
func FromReader(r io.Reader) ([]string, int, error) {
	entries, loadErr := store.ReadEntries()
	if loadErr != nil {
		return nil, 0, loadErr
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
		return nil, 0, errFs.ReadInput(scanErr)
	}

	return entries, count, nil
}

// FromDirectory reads first-level regular files from dir and imports
// each as a blob entry. Returns the updated entries list, count of
// added blobs, and per-file results for reporting. The caller owns
// writing the result.
//
// Parameters:
//   - dir: Directory path containing files to import
//
// Returns:
//   - []string: Updated entries list (existing + imported blobs)
//   - int: Number of blobs successfully imported
//   - []BlobResult: Per-file outcomes for caller to report
//   - error: Non-nil on directory read or entry load failure
func FromDirectory(dir string) ([]string, int, []BlobResult, error) {
	info, statErr := os.Stat(dir)
	if statErr != nil {
		return nil, 0, nil, errFs.StatPath(dir, statErr)
	}
	if !info.IsDir() {
		return nil, 0, nil, errFs.NotDirectory(dir)
	}

	dirEntries, readErr := os.ReadDir(dir)
	if readErr != nil {
		return nil, 0, nil, errFs.ReadDirectory(dir, readErr)
	}

	entries, loadErr := store.ReadEntries()
	if loadErr != nil {
		return nil, 0, nil, loadErr
	}

	var added int
	var results []BlobResult

	for _, de := range dirEntries {
		if !de.Type().IsRegular() {
			continue
		}

		name := de.Name()

		data, fileErr := internalIo.SafeReadFile(dir, name)
		if fileErr != nil {
			results = append(results, BlobResult{Name: name, Err: fileErr})
			continue
		}

		if len(data) > pad.MaxBlobSize {
			results = append(results, BlobResult{Name: name, TooLarge: true})
			continue
		}

		entries = append(entries, blob.Make(name, data))
		results = append(results, BlobResult{Name: name, Added: true})
		added++
	}

	return entries, added, results, nil
}
