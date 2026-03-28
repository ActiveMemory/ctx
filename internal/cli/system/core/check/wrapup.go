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
// Returns true if nudges should be suppressed.
//
// Returns:
//   - bool: True if wrap-up marker is fresh
func WrappedUpRecently() bool {
	markerPath := filepath.Join(state.Dir(), wrap.Marker)

	info, statErr := os.Stat(markerPath)
	if statErr != nil {
		return false
	}

	return time.Since(info.ModTime()) < wrap.ExpiryHours*time.Hour
}
