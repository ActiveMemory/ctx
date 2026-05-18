//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

import (
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// ErrDrift indicates schema drift was detected.
const ErrDrift = entity.Sentinel(text.DescKeyErrSchemaDrift)

// Drift returns a schema drift error.
//
// Returns:
//   - error: the drift sentinel error
func Drift() error {
	return ErrDrift
}
