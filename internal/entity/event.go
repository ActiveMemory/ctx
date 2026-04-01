//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

// EventQueryOpts controls event filtering and pagination.
type EventQueryOpts struct {
	Hook           string // filter by hook name (from detail)
	Session        string // filter by session ID
	Event          string // filter by event type
	Last           int    // return last N events (0 = all)
	IncludeRotated bool   // also read events.1.jsonl
}
