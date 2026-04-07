//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// stateFile is the sync state persistence filename.
const stateFile = ".sync-state.json"

// loadState reads sync state from .context/shared/.
func loadState() (state, error) {
	var s state
	path := filepath.Join(
		rc.ContextDir(), "shared", stateFile,
	)
	data, readErr := io.SafeReadUserFile(path)
	if os.IsNotExist(readErr) {
		return s, nil
	}
	if readErr != nil {
		return s, readErr
	}
	if len(data) == 0 {
		return s, nil
	}
	return s, json.Unmarshal(data, &s)
}

// saveState writes sync state to .context/shared/.
func saveState(s state) error {
	dir := filepath.Join(rc.ContextDir(), "shared")
	if mkErr := io.SafeMkdirAll(
		dir, fs.PermKeyDir,
	); mkErr != nil {
		return mkErr
	}
	data, marshalErr := json.MarshalIndent(
		s, "", "  ",
	)
	if marshalErr != nil {
		return marshalErr
	}
	path := filepath.Join(dir, stateFile)
	return io.SafeWriteFile(path, data, fs.PermFile)
}
