//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

// FileInfo holds prepared data for a single file in status output.
type FileInfo struct {
	Indicator string
	Name      string
	Status    string
	Tokens    int
	Size      int64
	Preview   []string
}

// ActivityInfo holds prepared data for a recent activity entry.
type ActivityInfo struct {
	Name string
	Ago  string
}
