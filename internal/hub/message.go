//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

// entryToMsg converts a store Entry to a wire EntryMsg.
func entryToMsg(e *Entry) *EntryMsg {
	return &EntryMsg{
		ID:        e.ID,
		Type:      e.Type,
		Content:   e.Content,
		Origin:    e.Origin,
		Author:    e.Author,
		Timestamp: e.Timestamp.Unix(),
		Sequence:  e.Sequence,
	}
}
