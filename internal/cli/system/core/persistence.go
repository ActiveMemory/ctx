//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/nudge"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/io"
)

// ReadPersistenceState reads a persistence state file and returns the
// parsed state. Returns ok=false if the file does not exist or cannot
// be read.
//
// Parameters:
//   - path: absolute path to the state file
//
// Returns:
//   - PersistenceState: parsed counter state
//   - bool: true if the file was read successfully
func ReadPersistenceState(path string) (PersistenceState, bool) {
	data, readErr := io.SafeReadFile(filepath.Dir(path), filepath.Base(path))
	if readErr != nil {
		return PersistenceState{}, false
	}

	var ps PersistenceState
	for _, line := range strings.Split(strings.TrimSpace(string(data)), token.NewlineLF) {
		parts := strings.SplitN(line, token.KeyValueSep, 2)
		if len(parts) != 2 {
			continue
		}
		switch parts[0] {
		case nudge.PersistenceKeyCount:
			n, parseErr := strconv.Atoi(parts[1])
			if parseErr == nil {
				ps.Count = n
			}
		case nudge.PersistenceKeyLastNudge:
			n, parseErr := strconv.Atoi(parts[1])
			if parseErr == nil {
				ps.LastNudge = n
			}
		case nudge.PersistenceKeyLastMtime:
			n, parseErr := strconv.ParseInt(parts[1], 10, 64)
			if parseErr == nil {
				ps.LastMtime = n
			}
		}
	}
	return ps, true
}

// WritePersistenceState writes the persistence state to the given file.
//
// Parameters:
//   - path: absolute path to the state file
//   - s: state to persist
func WritePersistenceState(path string, s PersistenceState) {
	content := fmt.Sprintf(desc.TextDesc(text.DescKeyCheckPersistenceStateFormat),
		s.Count, s.LastNudge, s.LastMtime)
	_ = os.WriteFile(path, []byte(content), fs.PermSecret)
}

// PersistenceNudgeNeeded determines whether a persistence nudge should
// fire based on prompt count and the number of prompts since the last nudge.
//
// Parameters:
//   - count: total prompt count for the session
//   - sinceNudge: number of prompts since the last nudge or context update
//
// Returns:
//   - bool: true if a nudge should be emitted
func PersistenceNudgeNeeded(count, sinceNudge int) bool {
	if count >= nudge.PersistenceEarlyMin && count <= nudge.PersistenceEarlyMax && sinceNudge >= nudge.PersistenceEarlyInterval {
		return true
	}
	if count > nudge.PersistenceEarlyMax && sinceNudge >= nudge.PersistenceLateInterval {
		return true
	}
	return false
}
