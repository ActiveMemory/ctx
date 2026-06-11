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
	"github.com/ActiveMemory/ctx/internal/config/token"
	cfgWarn "github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/io"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
)

// readJSONL is a generic helper that opens the file at path and
// decodes each line as a JSON value of type T. Malformed lines
// are silently skipped. Returns an empty (non-nil) slice when
// the file does not exist.
//
// Parameters:
//   - path: filesystem path to the JSONL file
//
// Returns:
//   - []T: decoded entries, empty slice if file missing
//   - error: file-open or scanner failure
func readJSONL[T any](path string) ([]T, error) {
	f, openErr := io.SafeOpenUserFile(path)
	if openErr != nil {
		if errors.Is(openErr, os.ErrNotExist) {
			return []T{}, nil
		}
		return nil, openErr
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			logWarn.Warn(cfgWarn.Close, path, cerr)
		}
	}()

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

// appendJSONL marshals entry as JSON and appends it as a line to
// dir/filename. Creates the directory if needed.
//
// Parameters:
//   - dir: directory containing the JSONL file
//   - filename: name of the JSONL file within dir
//   - entry: value to marshal and append
//
// Returns:
//   - error: marshal, directory creation, or write failure
func appendJSONL[T any](dir, filename string, entry T) (err error) {
	if mkErr := io.SafeMkdirAll(dir, cfgFs.PermRestrictedDir); mkErr != nil {
		return mkErr
	}

	line, marshalErr := json.Marshal(entry)
	if marshalErr != nil {
		return marshalErr
	}
	line = append(line, token.NewlineLF...)

	path := filepath.Join(dir, filename)
	f, openErr := io.SafeAppendFile(path, cfgFs.PermFile)
	if openErr != nil {
		return openErr
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	_, writeErr := f.Write(line)
	return writeErr
}
