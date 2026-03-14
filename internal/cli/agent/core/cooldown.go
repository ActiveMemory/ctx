//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/agent"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// DefaultCooldown is the default cooldown duration between context packet
// emissions within the same session.
const DefaultCooldown = agent.DefaultCooldown

// CooldownActive checks whether the cooldown tombstone for the given
// session is still fresh.
//
// Parameters:
//   - session: session identifier (typically the caller's PID)
//   - cooldown: duration to suppress repeated output
//
// Returns:
//   - bool: true if tombstone exists and is within the cooldown window
func CooldownActive(session string, cooldown time.Duration) bool {
	if session == "" || cooldown <= 0 {
		return false
	}
	info, err := os.Stat(TombstonePath(session))
	if err != nil {
		return false
	}
	return time.Since(info.ModTime()) < cooldown
}

// TouchTombstone creates or updates the tombstone file for the given
// session, marking the current time as the last emission.
//
// Parameters:
//   - session: session identifier (typically the caller's PID)
func TouchTombstone(session string) {
	if session == "" {
		return
	}
	_ = os.WriteFile(TombstonePath(session), nil, 0o600)
}

// TombstonePath returns the filesystem path for a session's tombstone.
//
// Parameters:
//   - session: session identifier
//
// Returns:
//   - string: absolute path in the system temp directory
func TombstonePath(session string) string {
	stateDir := filepath.Join(rc.ContextDir(), dir.State)
	_ = os.MkdirAll(stateDir, 0o750)
	return filepath.Join(stateDir, agent.TombstonePrefix+session)
}
