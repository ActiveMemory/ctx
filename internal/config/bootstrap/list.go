//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

// Numbered list parsing constants.
const (
	// NumberedListSep is the separator between the number and
	// text in numbered lists (e.g. "1. item").
	NumberedListSep = ". "
	// NumberedListMaxDigits is the maximum index position for
	// the separator to be recognized as a prefix.
	NumberedListMaxDigits = 2
)
