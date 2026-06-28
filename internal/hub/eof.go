//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"errors"
	stdio "io"
)

// eof reports whether err is io.EOF.
//
// gRPC delivers an unwrapped io.EOF at clean end-of-stream, but
// errors.Is keeps the suppression correct if any layer ever wraps
// it — a bare == would let a wrapped EOF leak through as a noisy
// transport warning on every clean replication cycle.
//
// Parameters:
//   - err: error to check
//
// Returns:
//   - bool: true if err is io.EOF
func eof(err error) bool {
	return errors.Is(err, stdio.EOF)
}
