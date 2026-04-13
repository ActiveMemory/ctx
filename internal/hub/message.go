//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

// entryToMsg converts a store Entry to a wire EntryMsg.
//
// Parameters:
//   - e: store entry to convert
//
// Returns:
//   - *EntryMsg: wire-format entry for streaming RPCs
func entryToMsg(e *Entry) *EntryMsg {
	return &EntryMsg{
		ID:        e.ID,
		Type:      e.Type,
		Content:   e.Content,
		Origin:    e.Origin,
		Meta:      e.Meta,
		Timestamp: e.Timestamp.Unix(),
		Sequence:  e.Sequence,
	}
}
