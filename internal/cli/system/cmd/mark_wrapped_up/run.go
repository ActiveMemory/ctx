//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mark_wrapped_up

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/config/wrap"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
	"github.com/ActiveMemory/ctx/internal/write/session"
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
	initialized, initErr := state.Initialized()
	if initErr != nil {
		logWarn.Warn(warn.StateInitializedProbe, initErr)
		return nil
	}
	if !initialized {
		return nil
	}

	stateDir, dirErr := state.Dir()
	if dirErr != nil {
		logWarn.Warn(warn.StateDirProbe, dirErr)
		return nil
	}
	markerPath := filepath.Join(stateDir, wrap.Marker)

	if writeErr := ctxIo.SafeWriteFile(
		markerPath, []byte(wrap.Content), fs.PermSecret,
	); writeErr != nil {
		return writeErr
	}

	session.WrappedUp(cmd)
	return nil
}
