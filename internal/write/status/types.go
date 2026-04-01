//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

// FileInfo holds prepared data for a single file in status output.
//
// Fields:
//   - Indicator: Status symbol (checkmark, warning, etc.)
//   - Name: Context file name
//   - Status: Human-readable status text
//   - Tokens: Estimated token count
//   - Size: File size in bytes
//   - Preview: First few content lines for verbose mode
type FileInfo struct {
	Indicator string
	Name      string
	Status    string
	Tokens    int
	Size      int64
	Preview   []string
}

// ActivityInfo holds prepared data for a recent activity entry.
//
// Fields:
//   - Name: Context file name
//   - Ago: Human-readable time since last modification
type ActivityInfo struct {
	Name string
	Ago  string
}
