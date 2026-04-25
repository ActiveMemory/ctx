//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package nudge

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/counter"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	cfgNudge "github.com/ActiveMemory/ctx/internal/config/nudge"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/io"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
)

// PauseMarkerPath returns the path to the session pause marker file.
//
// Parameters:
//   - sessionID: Session identifier
//
// Returns:
//   - string: Absolute path to the pause marker file
//   - error: non-nil when the state directory cannot be resolved
func PauseMarkerPath(sessionID string) (string, error) {
	stateDir, dirErr := state.Dir()
	if dirErr != nil {
		return "", dirErr
	}
	return filepath.Join(stateDir, hook.PrefixPauseMarker+sessionID), nil
}

// Paused checks if the session is paused. If paused, increments the
// turn counter and returns the current count. Returns 0 if not paused
// or if the state directory cannot be resolved (silent bail keeps the
// calling hook lean; the resolver-failure warning fires once at the
// [state.Initialized] gate those hooks check first).
//
// Parameters:
//   - sessionID: Session identifier
//
// Returns:
//   - int: Turn count if paused, 0 if not paused
func Paused(sessionID string) int {
	path, pathErr := PauseMarkerPath(sessionID)
	if pathErr != nil {
		return 0
	}
	data, readErr := io.SafeReadUserFile(path)
	if readErr != nil {
		return 0
	}
	count, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	count++
	counter.Write(path, count)
	return count
}

// PausedMessage returns the appropriate pause indicator for the given
// turn count, or empty string if not paused (turns == 0).
//
// Parameters:
//   - turns: Number of paused turns
//
// Returns:
//   - string: Pause message, or empty string
func PausedMessage(turns int) string {
	if turns == 0 {
		return ""
	}
	if turns <= cfgNudge.PauseTurnThreshold {
		return hook.LabelPaused
	}
	return fmt.Sprintf(desc.Text(text.DescKeyWritePausedMessage), turns)
}

// Pause creates the session pause marker. Exported for use by the
// ctx hook pause command.
//
// Parameters:
//   - sessionID: Session identifier
//
// Returns:
//   - error: non-nil when the state directory cannot be resolved
func Pause(sessionID string) error {
	path, pathErr := PauseMarkerPath(sessionID)
	if pathErr != nil {
		return pathErr
	}
	counter.Write(path, 0)
	return nil
}

// Resume removes the session pause marker. Exported for use by the
// ctx hook resume command. No-op if not paused.
//
// Parameters:
//   - sessionID: Session identifier
//
// Returns:
//   - error: non-nil when the state directory cannot be resolved
func Resume(sessionID string) error {
	p, pathErr := PauseMarkerPath(sessionID)
	if pathErr != nil {
		return pathErr
	}
	if removeErr := os.Remove(p); removeErr != nil {
		ctxLog.Warn(warn.Remove, p, removeErr)
	}
	return nil
}
