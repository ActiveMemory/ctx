//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"errors"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// Violations returns an error when drift detection found violations.
//
// Returns:
//   - error: "drift detection found violations"
func Violations() error {
	return errors.New(
		desc.Text(text.DescKeyErrValidateDriftViolations),
	)
}
