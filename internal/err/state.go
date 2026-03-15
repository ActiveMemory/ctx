//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// ReadingStateDir wraps a state directory read failure.
//
// Parameters:
//   - cause: the underlying error from reading the directory.
//
// Returns:
//   - error: "reading state directory: <cause>"
func ReadingStateDir(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrStateReadingStateDir), cause)
}

// LoadState wraps a state-loading failure.
//
// Parameters:
//   - cause: the underlying error from loading state.
//
// Returns:
//   - error: "loading state: <cause>"
func LoadState(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrStateLoadState), cause)
}

// SaveState wraps a state-saving failure.
//
// Parameters:
//   - cause: the underlying error from saving state.
//
// Returns:
//   - error: "saving state: <cause>"
func SaveState(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrStateSaveState), cause)
}
