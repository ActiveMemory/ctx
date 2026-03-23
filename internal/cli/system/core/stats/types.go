//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stats

import "github.com/ActiveMemory/ctx/internal/entity"

// Entry is a Stats with the source session ID for display.
type Entry struct {
	entity.Stats
	Session string `json:"session"`
}
