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
	"github.com/ActiveMemory/ctx/internal/io"
)

// dirPerm is the permission for the hub data directory.
const dirPerm = fs.PermKeyDir

// Data file names within the hub directory.
const (
	// fileEntries is the append-only JSONL entry log.
	fileEntries = "entries.jsonl"
	// fileClients is the registered client registry.
	fileClients = "clients.json"
	// fileMeta is the hub-level metadata file.
	fileMeta = "meta.json"
)

// entriesPath returns the full path to the entries file.
func entriesPath(dir string) string {
	return filepath.Join(dir, fileEntries)
}

// clientsPath returns the full path to the clients file.
func clientsPath(dir string) string {
	return filepath.Join(dir, fileClients)
}

// metaPath returns the full path to the metadata file.
func metaPath(dir string) string {
	return filepath.Join(dir, fileMeta)
}

// loadJSON reads a JSON file into dst. Returns nil if the
// file does not exist.
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
func saveJSON(path string, src any) error {
	data, marshalErr := json.MarshalIndent(src, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	return io.SafeWriteFile(path, data, fs.PermFile)
}

// appendFile appends data to a file, creating it if needed.
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
