//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tags

import (
	"encoding/json"
	"sort"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core/store"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/tag"
	writePad "github.com/ActiveMemory/ctx/internal/write/pad"
)

// Run lists all tags across scratchpad entries with counts.
//
// Parameters:
//   - cmd: Cobra command for output
//   - jsonOut: When true, output as JSON array
//
// Returns:
//   - error: Non-nil on read failure or JSON marshal error
func Run(cmd *cobra.Command, jsonOut bool) error {
	entries, err := store.ReadEntries()
	if err != nil {
		return err
	}

	counts := make(map[string]int)
	for _, entry := range entries {
		for _, t := range tag.Extract(entry) {
			counts[t]++
		}
	}

	if len(counts) == 0 {
		writePad.TagsNone(cmd)
		return nil
	}

	names := make([]string, 0, len(counts))
	for name := range counts {
		names = append(names, name)
	}
	sort.Strings(names)

	if jsonOut {
		items := make([]tag.Count, len(names))
		for i, name := range names {
			items[i] = tag.Count{Tag: name, Count: counts[name]}
		}
		data, marshalErr := json.Marshal(items)
		if marshalErr != nil {
			return marshalErr
		}
		writePad.TagsJSON(cmd, data)
		return nil
	}

	for _, name := range names {
		writePad.TagsItem(cmd, name, counts[name])
	}
	return nil
}
