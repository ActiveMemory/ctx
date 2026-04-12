//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	"github.com/ActiveMemory/ctx/internal/io"
)

// entriesPath returns the full path to the entries file.
//
// Parameters:
//   - dir: hub data directory
//
// Returns:
//   - string: absolute path to entries JSONL file
func entriesPath(dir string) string {
	return filepath.Join(dir, cfgHub.FileEntries)
}

// clientsPath returns the full path to the clients file.
//
// Parameters:
//   - dir: hub data directory
//
// Returns:
//   - string: absolute path to clients JSON file
func clientsPath(dir string) string {
	return filepath.Join(dir, cfgHub.FileClients)
}

// metaPath returns the full path to the metadata file.
//
// Parameters:
//   - dir: hub data directory
//
// Returns:
//   - string: absolute path to metadata JSON file
func metaPath(dir string) string {
	return filepath.Join(dir, cfgHub.FileMeta)
}

// loadJSON reads a JSON file into dst. Returns nil if the
// file does not exist.
//
// Parameters:
//   - path: file path to read
//   - dst: target to unmarshal into
//
// Returns:
//   - error: non-nil if read or unmarshal fails
func loadJSON(path string, dst any) error {
	data, readErr := io.SafeReadUserFile(path)
	if os.IsNotExist(readErr) {
		return nil
	}
	if readErr != nil {
		return readErr
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, dst)
}

// saveJSON marshals src and writes it to path.
//
// Parameters:
//   - path: file path to write
//   - src: value to marshal as JSON
//
// Returns:
//   - error: non-nil if marshal or write fails
func saveJSON(path string, src any) error {
	data, marshalErr := json.MarshalIndent(
		src, "", cfgHub.JSONIndent,
	)
	if marshalErr != nil {
		return marshalErr
	}
	return io.SafeWriteFile(path, data, fs.PermFile)
}

// appendFile appends data to a file, creating it if needed.
//
// Parameters:
//   - path: file path to append to
//   - data: bytes to append
//
// Returns:
//   - error: non-nil if read or write fails
func appendFile(path string, data []byte) error {
	existing, readErr := io.SafeReadUserFile(path)
	if readErr != nil && !os.IsNotExist(readErr) {
		return readErr
	}
	return io.SafeWriteFile(
		path, append(existing, data...), fs.PermFile,
	)
}

// loadEntries reads the JSONL entry log into the slice.
//
// Parameters:
//   - dir: hub data directory
//   - dst: slice to append loaded entries into
//
// Returns:
//   - error: non-nil if read or unmarshal fails
func loadEntries(dir string, dst *[]Entry) error {
	data, readErr := io.SafeReadUserFile(entriesPath(dir))
	if os.IsNotExist(readErr) {
		return nil
	}
	if readErr != nil {
		return readErr
	}
	if len(data) == 0 {
		return nil
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		var e Entry
		if decErr := json.Unmarshal(
			scanner.Bytes(), &e,
		); decErr != nil {
			return decErr
		}
		*dst = append(*dst, e)
	}
	return scanner.Err()
}
