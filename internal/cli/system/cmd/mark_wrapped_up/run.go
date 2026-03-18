//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mark_wrapped_up

import (
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/wrap"
	"github.com/ActiveMemory/ctx/internal/write/session"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
)

// Run creates or updates the wrap-up marker file.
//
// Writes the marker so that nudge hooks (ceremonies, persistence, etc.)
// are suppressed for WrappedUpExpiry after a wrap-up ceremony completes.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil if the marker file cannot be written
func Run(cmd *cobra.Command) error {
	markerPath := filepath.Join(core.StateDir(), wrap.WrappedUpMarker)

	if writeErr := os.WriteFile(
		markerPath, []byte(wrap.WrappedUpContent), fs.PermSecret,
	); writeErr != nil {
		return writeErr
	}

	session.SessionWrappedUp(cmd)
	return nil
}
