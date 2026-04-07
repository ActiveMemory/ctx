//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import "strings"

// matchesCommit checks whether stored and query commit hashes
// match. Supports short hashes by checking whether either
// string is a prefix of the other.
//
// Parameters:
//   - stored: commit hash from the persisted record
//   - query: commit hash provided by the caller
//
// Returns:
//   - bool: true when either hash is a prefix of the other
func matchesCommit(stored, query string) bool {
	return strings.HasPrefix(stored, query) ||
		strings.HasPrefix(query, stored)
}
