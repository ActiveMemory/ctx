//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

// Reserved for cluster mode: lastSequence will be called
// from replicateOnce to determine the sync cursor.
var _ = (*Store).lastSequence

// lastSequence returns the highest sequence in the store.
//
// Returns:
//   - bool: true if at least one entry exists
//   - uint64: highest sequence number, or 0 if empty
func (s *Store) lastSequence() (bool, uint64) {
	all := s.Query(nil, 0)
	if len(all) == 0 {
		return false, 0
	}
	return true, all[len(all)-1].Sequence
}
