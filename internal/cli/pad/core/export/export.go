//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package export

import (
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core/blob"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/store"
)

// Plan extracts blob entries from the scratchpad and computes output
// paths, handling filename collisions with timestamp prefixes.
//
// Parameters:
//   - dir: Target directory for exported files
//   - force: When true, overwrite existing files (no collision rename)
//
// Returns:
//   - []Item: Blobs to export with resolved paths
//   - error: Non-nil on scratchpad read failure
func Plan(dir string, force bool) ([]Item, error) {
	entries, readErr := store.ReadEntries()
	if readErr != nil {
		return nil, readErr
	}

	var items []Item
	for _, entry := range entries {
		label, data, ok := blob.Split(entry)
		if !ok {
			continue
		}

		outPath := filepath.Join(dir, label)
		item := Item{Label: label, Data: data, OutPath: outPath}

		if !force {
			if _, statErr := os.Stat(outPath); statErr == nil {
				item.Exists = true
				item.AltName = tsWithLabel(label)
				item.OutPath = filepath.Join(dir, item.AltName)
			}
		}

		items = append(items, item)
	}

	return items, nil
}
