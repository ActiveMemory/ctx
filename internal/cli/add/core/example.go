//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config"
)

// ExamplesForType returns example usage strings for a given entry type.
//
// The examples are loaded from the embedded commands.yaml asset.
//
// Parameters:
//   - fileType: Entry type (e.g., "decision", "task", "learning", "convention")
//
// Returns:
//   - string: Formatted example commands; returns a generic example for
//     unrecognized types
func ExamplesForType(fileType string) string {
	const defaultKeyName = "default"

	key := config.UserInputToEntry(fileType)

	if key == "" {
		key = defaultKeyName
	}

	if desc := assets.ExampleDesc(key); desc != "" {
		return desc
	}

	return assets.ExampleDesc(defaultKeyName)
}
