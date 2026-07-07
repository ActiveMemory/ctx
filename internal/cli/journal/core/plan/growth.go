//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package plan

import (
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/cli/journal/core/extract"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/journal/state"
)

// sourceStat returns the mtime (Unix seconds) and byte size of a
// session's source transcript, and whether the stat succeeded. A
// missing or unreadable source disables growth detection for that
// session, which then degrades to skip-existing.
//
// Parameters:
//   - path: absolute path to the source transcript
//
// Returns:
//   - mtime: Unix mtime in seconds (0 when ok is false)
//   - size: byte size (0 when ok is false)
//   - ok: true if the file was stat-able
func sourceStat(path string) (mtime, size int64, ok bool) {
	if path == "" {
		return 0, 0, false
	}
	info, err := os.Stat(path)
	if err != nil {
		return 0, 0, false
	}
	return info.ModTime().Unix(), info.Size(), true
}

// adoptRenderHash computes the render hash of an existing entry's body,
// used to adopt a session imported before source tracking existed (v1).
// Recording this hash at adoption time lets a later growth sweep prove
// the entry is still ctx-owned and re-render it, instead of treating
// the empty recorded hash as a hand edit and stranding the growth.
//
// Parameters:
//   - path: absolute path to the existing entry file
//
// Returns:
//   - string: render hash of the body, or "" if the file is unreadable
func adoptRenderHash(path string) string {
	data, err := io.SafeReadUserFile(filepath.Clean(path))
	if err != nil {
		return ""
	}
	return state.HashRender(extract.StripFrontmatter(string(data)))
}

// bodyEdited reports whether an existing entry's body differs from the
// hash ctx recorded for its last authored write. An entry with no
// recorded hash (pre-v2, or never written under v2) is treated as
// edited: ctx cannot prove it owns the body, so it must not clobber it.
//
// Parameters:
//   - path: absolute path to the existing entry file
//   - recordedHash: the render hash ctx recorded for this entry
//
// Returns:
//   - bool: true if the body was edited outside ctx (or cannot be read)
func bodyEdited(path, recordedHash string) bool {
	if recordedHash == "" {
		return true
	}
	data, err := io.SafeReadUserFile(filepath.Clean(path))
	if err != nil {
		return true
	}
	body := extract.StripFrontmatter(string(data))
	return state.HashRender(body) != recordedHash
}
