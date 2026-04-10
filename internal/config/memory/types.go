//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

// ClassifyRule maps keyword patterns to a target entry type.
//
// Fields:
//   - Target: entry type constant (convention, decision, learning, task)
//   - Keywords: case-insensitive keyword patterns to match
type ClassifyRule struct {
	Target   string   `yaml:"target"`
	Keywords []string `yaml:"keywords"`
}
