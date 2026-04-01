//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package lookup

// ConfigPattern pairs a glob pattern with its localizable topic description.
//
// Fields:
//   - Pattern: Glob pattern to match file paths
//   - Topic: Localized topic description key
type ConfigPattern struct {
	Pattern string
	Topic   string
}
