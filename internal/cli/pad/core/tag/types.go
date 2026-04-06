//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tag

// Count holds a tag name and its occurrence count.
type Count struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}
