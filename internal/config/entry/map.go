//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entry

import (
	"sync"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
)

// toCtxFile maps short names to actual file names.
// Guarded by mu; use [CtxFile] or [MustCtxFile] for access.
var (
	ctxFileMu sync.RWMutex
	toCtxFile = map[string]string{
		Decision:   ctx.Decision,
		Task:       ctx.Task,
		Learning:   ctx.Learning,
		Convention: ctx.Convention,
	}
)

// CtxFile returns the context filename for the given entry type.
// Returns ("", false) if the type is unknown. Thread-safe.
//
// Parameters:
//   - entryType: entry type key (decision, task, learning, convention)
//
// Returns:
//   - string: context filename (e.g. "DECISIONS.md")
//   - bool: true if the type was found
func CtxFile(entryType string) (string, bool) {
	ctxFileMu.RLock()
	f, ok := toCtxFile[entryType]
	ctxFileMu.RUnlock()
	return f, ok
}

// MustCtxFile returns the context filename for a validated entry type.
// Panics if the type is unknown — use only after validation. Thread-safe.
//
// Parameters:
//   - entryType: validated entry type key
//
// Returns:
//   - string: context filename (e.g. "DECISIONS.md")
func MustCtxFile(entryType string) string {
	ctxFileMu.RLock()
	f, ok := toCtxFile[entryType]
	ctxFileMu.RUnlock()
	if !ok {
		panic("unknown entry type: " + entryType)
	}
	return f
}
