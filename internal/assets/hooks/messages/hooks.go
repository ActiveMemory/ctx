//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package messages

// registryError returns any error encountered while
// parsing the embedded registry.yaml. Nil on success.
//
// Returns:
//   - error: Parse error, or nil on success
func registryError() error {
	Registry() // ensure sync.Once has run
	return registryErr
}

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
