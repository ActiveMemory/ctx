//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package lookup

var (
	stopWordsMap map[string]bool
)

var (
	CommandsMap map[string]commandEntry
	FlagsMap    map[string]commandEntry
	TextMap     map[string]commandEntry
	ExamplesMap map[string]commandEntry
)

type commandEntry struct {
	Short string `yaml:"short"`
	Long  string `yaml:"long"`
}
