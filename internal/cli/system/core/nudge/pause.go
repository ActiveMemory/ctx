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
	"github.com/ActiveMemory/ctx/internal/io"
)

// PauseMarkerPath returns the path to the session pause marker file.
//
// Parameters:
//   - sessionID: Session identifier
//
// Returns:
//   - string: Absolute path to the pause marker file
func PauseMarkerPath(sessionID string) string {
	return filepath.Join(state.StateDir(), "ctx-paused-"+sessionID)
}

// Paused checks if the session is paused. If paused, increments the
// turn counter and returns the current count. Returns 0 if not paused.
//
// Parameters:
//   - sessionID: Session identifier
//
// Returns:
//   - int: Turn count if paused, 0 if not paused
func Paused(sessionID string) int {
	path := PauseMarkerPath(sessionID)
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
	if turns <= 5 {
		return "ctx:paused"
	}
	return fmt.Sprintf(desc.Text(text.DescKeyWritePausedMessage), turns)
}

// Pause creates the session pause marker. Exported for use by the
// top-level ctx pause command.
//
// Parameters:
//   - sessionID: Session identifier
func Pause(sessionID string) {
	counter.Write(PauseMarkerPath(sessionID), 0)
}

// Resume removes the session pause marker. Exported for use by the
// top-level ctx resume command. No-op if not paused.
//
// Parameters:
//   - sessionID: Session identifier
func Resume(sessionID string) {
	_ = os.Remove(PauseMarkerPath(sessionID))
}
