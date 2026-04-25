//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package activate

import (
	"errors"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// NoCandidates returns the error used when `ctx activate`
// finds zero .context/ directories on the upward path from CWD.
//
// Returns:
//   - error: multi-line message prompting `ctx init`.
func NoCandidates() error {
	return errors.New(desc.Text(text.DescKeyErrActivateNoCandidates))
}
