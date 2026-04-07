//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import stdio "io"

// isEOF reports whether err is io.EOF.
func isEOF(err error) bool {
	return err == stdio.EOF
}
