//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// defaultCooldown is the default cooldown duration between context packet
// emissions within the same session.
const defaultCooldown = 10 * time.Minute

// tombstonePrefix is the filename prefix for cooldown tombstone files.
const tombstonePrefix = "ctx-agent-"

// cooldownActive checks whether the cooldown tombstone for the given
// session is still fresh.
//
// Parameters:
//   - session: session identifier (typically the caller's PID)
//   - cooldown: duration to suppress repeated output
//
// Returns:
//   - bool: true if tombstone exists and is within the cooldown window
func cooldownActive(session string, cooldown time.Duration) bool {
	if session == "" || cooldown <= 0 {
		return false
	}
	info, err := os.Stat(tombstonePath(session))
	if err != nil {
		return false
	}
	return time.Since(info.ModTime()) < cooldown
}

// touchTombstone creates or updates the tombstone file for the given
// session, marking the current time as the last emission.
//
// Parameters:
//   - session: session identifier (typically the caller's PID)
func touchTombstone(session string) {
	if session == "" {
		return
	}
	_ = os.WriteFile(tombstonePath(session), nil, 0o600)
}

// tombstonePath returns the filesystem path for a session's tombstone.
//
// Parameters:
//   - session: session identifier
//
// Returns:
//   - string: absolute path in the system temp directory
func tombstonePath(session string) string {
	return filepath.Join(
		os.TempDir(), fmt.Sprintf("%s%s", tombstonePrefix, session),
	)
}
