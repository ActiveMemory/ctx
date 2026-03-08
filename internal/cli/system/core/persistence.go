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

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/validation"
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
	data, readErr := validation.SafeReadFile(filepath.Dir(path), filepath.Base(path))
	if readErr != nil {
		return PersistenceState{}, false
	}

	var ps PersistenceState
	for _, line := range strings.Split(strings.TrimSpace(string(data)), config.NewlineLF) {
		parts := strings.SplitN(line, config.KeyValueSep, 2)
		if len(parts) != 2 {
			continue
		}
		switch parts[0] {
		case config.PersistenceKeyCount:
			n, parseErr := strconv.Atoi(parts[1])
			if parseErr == nil {
				ps.Count = n
			}
		case config.PersistenceKeyLastNudge:
			n, parseErr := strconv.Atoi(parts[1])
			if parseErr == nil {
				ps.LastNudge = n
			}
		case config.PersistenceKeyLastMtime:
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
	content := fmt.Sprintf(assets.TextDesc(assets.TextDescKeyCheckPersistenceStateFormat),
		s.Count, s.LastNudge, s.LastMtime)
	_ = os.WriteFile(path, []byte(content), config.PermSecret)
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
	if count >= config.PersistenceEarlyMin && count <= config.PersistenceEarlyMax && sinceNudge >= config.PersistenceEarlyInterval {
		return true
	}
	if count > config.PersistenceEarlyMax && sinceNudge >= config.PersistenceLateInterval {
		return true
	}
	return false
}
