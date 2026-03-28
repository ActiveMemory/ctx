//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import "github.com/ActiveMemory/ctx/internal/entity"

// CountUnique counts distinct session filenames across all topics.
//
// Parameters:
//   - topics: Topic data with associated journal entries
//
// Returns:
//   - int: Number of unique sessions (by filename)
func CountUnique(topics []entity.TopicData) int {
	seen := make(map[string]bool)
	for _, t := range topics {
		for _, e := range t.Entries {
			seen[e.Filename] = true
		}
	}
	return len(seen)
}
