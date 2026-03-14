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

	"github.com/ActiveMemory/ctx/internal/config/wrap"
)

// WrappedUpExpiry is how long the marker suppresses nudges.
const WrappedUpExpiry = 2 * time.Hour

// WrappedUpRecently checks whether the wrap-up marker exists and is
// less than WrappedUpExpiry old.
//
// Returns true if nudges should be suppressed.
//
// Returns:
//   - bool: True if wrap-up marker is fresh
func WrappedUpRecently() bool {
	markerPath := filepath.Join(StateDir(), wrap.WrappedUpMarker)

	info, statErr := os.Stat(markerPath)
	if statErr != nil {
		return false
	}

	return time.Since(info.ModTime()) < WrappedUpExpiry
}
