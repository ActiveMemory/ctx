//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package messages

// hooks returns a deduplicated list of hook names in
// the registry.
//
// Returns:
//   - []string: Hook names in alphabetical order
func hooks() []string {
	seen := make(map[string]bool)
	var result []string
	for _, info := range Registry() {
		if !seen[info.Hook] {
			seen[info.Hook] = true
			result = append(result, info.Hook)
		}
	}
	return result
}
