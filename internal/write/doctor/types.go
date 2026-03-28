//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package doctor

// ResultItem holds the display data for a single doctor check result.
type ResultItem struct {
	Category string
	Status   string
	Message  string
}
