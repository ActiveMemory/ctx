//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package render

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	"github.com/ActiveMemory/ctx/internal/hub"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// WriteEntries renders hub entries as markdown and appends
// them to type-specific files in .context/hub/.
//
// Parameters:
//   - entries: hub entries to render
//
// Returns:
//   - error: non-nil if directory creation or write fails
func WriteEntries(entries []hub.EntryMsg) error {
	dir := filepath.Join(rc.ContextDir(), cfgHub.DirHub)
	if mkErr := io.SafeMkdirAll(
		dir, fs.PermKeyDir,
	); mkErr != nil {
		return mkErr
	}

	grouped := groupByType(entries)
	for entryType, group := range grouped {
		fPath := filepath.Join(
			dir, typedFileName(entryType),
		)
		content := toMarkdown(group)
		if appendErr := appendShared(
			fPath, content,
		); appendErr != nil {
			return appendErr
		}
	}
	return nil
}
