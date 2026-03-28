//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

// SessionInfo holds pre-formatted session metadata for display.
type SessionInfo struct {
	Slug      string
	ID        string
	Tool      string
	Project   string
	Branch    string // empty to omit
	Model     string // empty to omit
	Started   string
	Duration  string
	Turns     int
	Messages  int
	TokensIn  string
	TokensOut string
	TokensAll string
}
