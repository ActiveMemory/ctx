//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	stdio "io"
)

// eof reports whether err is exactly io.EOF.
//
// The equality is deliberately strict. eof gates the clean
// end-of-stream path in client.Sync and client.Listen, where a true
// result means "the stream ended normally, treat what we have as
// complete." gRPC delivers a bare, unwrapped io.EOF only on an OK
// status, so == is correct today; keeping it strict means that if a
// future layer ever wraps a mid-stream (dirty) EOF, it is treated as
// a real error rather than silently accepted as a clean end. The
// replication warn-suppression path wants the opposite polarity and
// broadens to errors.Is inline at its own call site.
//
// Parameters:
//   - err: error to check
//
// Returns:
//   - bool: true if err is exactly io.EOF
func eof(err error) bool {
	return err == stdio.EOF
}
