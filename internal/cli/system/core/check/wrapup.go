//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check

import (
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/wrap"
)

// WrappedUpRecently checks whether the wrap-up marker exists and is
// less than the configured expiry old.
//
// Returns false when the state directory cannot be resolved: hooks
// that gate on this are already downstream of [state.Initialized],
// where the resolver failure surfaced once; here we fail-closed
// (assume not wrapped up, let nudges fire) rather than silently
// suppress everything.
//
// Returns:
//   - bool: True if wrap-up marker is fresh
func WrappedUpRecently() bool {
	stateDir, dirErr := state.Dir()
	if dirErr != nil {
		return false
	}
	markerPath := filepath.Join(stateDir, wrap.Marker)

	info, statErr := os.Stat(markerPath)
	if statErr != nil {
		return false
	}

	return time.Since(info.ModTime()) < wrap.ExpiryHours*time.Hour
}
