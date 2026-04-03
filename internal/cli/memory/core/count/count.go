//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package count

import (
	"bytes"

	"github.com/ActiveMemory/ctx/internal/config/token"
)

// FileLines counts the number of newline characters in data.
//
// Parameters:
//   - data: raw file bytes
//
// Returns:
//   - int: number of newline characters
func FileLines(data []byte) int {
	return bytes.Count(data, []byte(token.NewlineLF))
}
