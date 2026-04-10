//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package lookup

// stopWordsMap holds words excluded from search indexing.
var stopWordsMap map[string]bool

var (
	// CommandsMap maps command description keys to their YAML entries.
	CommandsMap map[string]commandEntry

	// FlagsMap maps flag description keys to their YAML entries.
	FlagsMap map[string]commandEntry

	// TextMap maps general text description keys to their YAML entries.
	TextMap map[string]commandEntry

	// ExamplesMap maps example description keys to their YAML entries.
	ExamplesMap map[string]commandEntry
)
