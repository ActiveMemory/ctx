//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
)

// readJSONL is a generic helper that opens the file at path and decodes each
// line as a JSON value of type T. Malformed lines are silently skipped.
// Returns an empty (non-nil) slice when the file does not exist.
func readJSONL[T any](path string) ([]T, error) {
	//nolint:gosec // path built from trusted directory + constant filename by callers
	f, openErr := os.Open(filepath.Clean(path))
	if openErr != nil {
		if errors.Is(openErr, os.ErrNotExist) {
			return []T{}, nil
		}
		return nil, openErr
	}
	defer func() { _ = f.Close() }()

	var entries []T
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var e T
		if unmarshalErr := json.Unmarshal(scanner.Bytes(), &e); unmarshalErr != nil {
			continue // skip malformed lines
		}
		entries = append(entries, e)
	}

	if scanErr := scanner.Err(); scanErr != nil {
		return nil, scanErr
	}

	if entries == nil {
		entries = []T{}
	}

	return entries, nil
}

// appendJSONL marshals entry as JSON and appends it as a line to dir/filename.
// Creates the directory if needed.
func appendJSONL[T any](dir, filename string, entry T) error {
	if mkErr := os.MkdirAll(dir, cfgFs.PermRestrictedDir); mkErr != nil {
		return mkErr
	}

	line, marshalErr := json.Marshal(entry)
	if marshalErr != nil {
		return marshalErr
	}
	line = append(line, '\n')

	path := filepath.Join(dir, filename)
	//nolint:gosec // path built from trusted dir + constant filename by callers
	f, openErr := os.OpenFile(
		filepath.Clean(path),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		cfgFs.PermFile,
	)
	if openErr != nil {
		return openErr
	}
	defer func() { _ = f.Close() }()

	_, writeErr := f.Write(line)
	return writeErr
}
