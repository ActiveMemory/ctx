//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trigger

import (
	"path/filepath"
	"strings"
)

// stripExt removes the file extension from a filename.
//
// Parameters:
//   - filename: file name possibly containing an extension
//
// Returns:
//   - string: filename with the extension removed
func stripExt(filename string) string {
	ext := filepath.Ext(filename)
	return strings.TrimSuffix(filename, ext)
}
