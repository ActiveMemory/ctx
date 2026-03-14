//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"bytes"

	"github.com/ActiveMemory/ctx/internal/config/token"
)

// CountFileLines counts the number of newline characters in data.
//
// Parameters:
//   - data: raw file bytes
//
// Returns:
//   - int: number of newline characters
func CountFileLines(data []byte) int {
	return bytes.Count(data, []byte(token.NewlineLF))
}
