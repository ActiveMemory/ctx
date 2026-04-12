//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import stdio "io"

// eof reports whether err is io.EOF.
//
// Parameters:
//   - err: error to check
//
// Returns:
//   - bool: true if err is io.EOF
func eof(err error) bool {
	return err == stdio.EOF
}
