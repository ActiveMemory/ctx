//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package persistence

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
	"github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/io"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
)

// ReadState reads a persistence state file and returns the
// parsed state. Returns ok=false if the file does not exist or cannot
// be read.
//
// Parameters:
//   - path: absolute path to the state file
//
// Returns:
//   - PersistenceState: parsed counter state
//   - bool: true if the file was read successfully
func ReadState(path string) (State, bool) {
	data, readErr := io.SafeReadFile(filepath.Dir(path), filepath.Base(path))
	if readErr != nil {
		return State{}, false
	}

	var ps State
	lines := strings.Split(
		strings.TrimSpace(string(data)), token.NewlineLF,
	)
	for _, line := range lines {
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
		case nudge.KeyLastNudge:
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

// WriteState writes the persistence state to the given file.
//
// Parameters:
//   - path: absolute path to the state file
//   - s: state to persist
func WriteState(path string, s State) {
	content := fmt.Sprintf(desc.Text(text.DescKeyCheckPersistenceStateFormat),
		s.Count, s.LastNudge, s.LastMtime)
	if writeErr := os.WriteFile(
		path, []byte(content), fs.PermSecret,
	); writeErr != nil {
		ctxLog.Warn(warn.Write, path, writeErr)
	}
}

// NudgeNeeded determines whether a persistence nudge should
// fire based on prompt count and the number of prompts since the last nudge.
//
// Parameters:
//   - count: total prompt count for the session
//   - sinceNudge: number of prompts since the last nudge or context update
//
// Returns:
//   - bool: true if a nudge should be emitted
func NudgeNeeded(count, sinceNudge int) bool {
	earlyRange := count >= nudge.PersistenceEarlyMin &&
		count <= nudge.PersistenceEarlyMax
	if earlyRange && sinceNudge >= nudge.PersistenceEarlyInterval {
		return true
	}
	lateRange := count > nudge.PersistenceEarlyMax
	if lateRange && sinceNudge >= nudge.PersistenceLateInterval {
		return true
	}
	return false
}
