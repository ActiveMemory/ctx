//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package state

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// ReadingDir wraps a state directory read failure.
//
// Parameters:
//   - cause: the underlying error from reading the directory.
//
// Returns:
//   - error: "reading state directory: <cause>"
func ReadingDir(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrStateReadingStateDir), cause,
	)
}

// Load wraps a state-loading failure.
//
// Parameters:
//   - cause: the underlying error from loading state.
//
// Returns:
//   - error: "loading state: <cause>"
func Load(cause error) error {
	return fmt.Errorf(desc.TextDesc(text.DescKeyErrStateLoadState), cause)
}

// Save wraps a state-saving failure.
//
// Parameters:
//   - cause: the underlying error from saving state.
//
// Returns:
//   - error: "saving state: <cause>"
func Save(cause error) error {
	return fmt.Errorf(desc.TextDesc(text.DescKeyErrStateSaveState), cause)
}
