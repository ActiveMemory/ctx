//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stats

import "github.com/ActiveMemory/ctx/internal/cli/system/core/session"

// Entry is a SessionStats with the source session ID for display.
type Entry struct {
	session.SessionStats
	Session string `json:"session"`
}
