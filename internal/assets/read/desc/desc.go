//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package desc

import "github.com/ActiveMemory/ctx/internal/assets/read/lookup"

// Command returns the Short and Long descriptions for a command.
//
// Keys use dot notation: "pad", "pad.show", "system.bootstrap".
// Returns empty strings if the key is not found.
//
// Parameters:
//   - key: Command key in dot notation
//
// Returns:
//   - short: One-line description
//   - long: Multi-line help text (may be empty)
func Command(key string) (short, long string) {
	entry, ok := lookup.CommandsMap[key]
	if !ok {
		return "", ""
	}
	return entry.Short, entry.Long
}

// Flag returns the description for a flag.
//
// Keys use dot notation: "add.file", "context-dir".
// Returns an empty string if the key is not found.
//
// Parameters:
//   - name: Flag key in dot notation
//
// Returns:
//   - string: Flag description
func Flag(name string) string {
	entry, ok := lookup.FlagsMap[name]
	if !ok {
		return ""
	}
	return entry.Short
}

// Example returns example usage text for a given key.
//
// Keys match entry types: "decision", "learning", "task", "convention".
// Returns an empty string if the key is not found.
//
// Parameters:
//   - name: Entry type key
//
// Returns:
//   - string: Example text
func Example(name string) string {
	entry, ok := lookup.ExamplesMap[name]
	if !ok {
		return ""
	}
	return entry.Short
}

// Text returns a user-facing text string by key.
//
// Keys use dot notation: "agent.instruction", "backup.run-hint".
// Returns an empty string if the key is not found.
//
// Parameters:
//   - name: Text key in dot notation
//
// Returns:
//   - string: Text content
func Text(name string) string {
	entry, ok := lookup.TextMap[name]
	if !ok {
		return ""
	}
	return entry.Short
}
