//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entry

// checkRequired returns the names of any fields whose values are empty.
//
// Parameters:
//   - fields: Pairs of [name, value] to validate
//
// Returns:
//   - []string: Names of fields with empty values; nil when all are populated
func checkRequired(fields [][2]string) []string {
	var missing []string
	for _, f := range fields {
		if f[1] == "" {
			missing = append(missing, f[0])
		}
	}
	return missing
}
