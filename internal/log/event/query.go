//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package event

import (
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// matchesFilter reports whether an event passes all non-empty query
// filters. Empty filter fields are treated as wildcards.
//
// Parameters:
//   - e: the event payload to test
//   - opts: query filters to match against
//
// Returns:
//   - bool: true if the event matches all non-empty filters
func matchesFilter(e notify.Payload, opts entity.EventQueryOpts) bool {
	if opts.Event != "" && e.Event != opts.Event {
		return false
	}
	if opts.Session != "" && e.SessionID != opts.Session {
		return false
	}
	if opts.Hook != "" {
		if e.Detail == nil || e.Detail.Hook != opts.Hook {
			return false
		}
	}
	return true
}
