//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// ExamplesForType returns example usage strings for a given entry type.
//
// The examples are loaded from the embedded commands.yaml asset.
// Entry type keys (decision, task, learning, convention) are defined in
// config/entry and match the YAML keys in examples.yaml.
//
// Parameters:
//   - fileType: Entry type (e.g., "decision", "task", "learning", "convention")
//
// Returns:
//   - string: Formatted example commands; returns a generic example for
//     unrecognized types
func ExamplesForType(fileType string) string {
	if d := desc.Example(fileType); d != "" {
		return d
	}

	return desc.Example(cmd.ExampleKeyDefault)
}
