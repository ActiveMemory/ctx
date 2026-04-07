//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

import (
	"errors"

	cfgSchema "github.com/ActiveMemory/ctx/internal/config/schema"
)

// ErrDrift indicates schema drift was detected.
var ErrDrift = errors.New(cfgSchema.ErrMsgDrift)

// Drift returns a schema drift error.
//
// Returns:
//   - error: the drift sentinel error
func Drift() error {
	return ErrDrift
}
