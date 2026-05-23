//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package placeholders

// file is the on-disk YAML shape for a placeholder locale
// file. The sole top-level key is `placeholders` (a list
// of strings). Future fields (e.g. attribution, version)
// should be added here, not as new top-level files.
type file struct {
	Placeholders []string `yaml:"placeholders"`
}
