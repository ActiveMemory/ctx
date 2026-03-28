//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/dir"
)

// Memory bridge file constants for .context/memory/ directory.
const (
	// Source is the Claude Code auto memory filename.
	Source = "MEMORY.md"
	// Mirror is the raw copy of Claude Code's MEMORY.md.
	Mirror = "mirror.md"
	// State is the sync/import tracking state file.
	State = "memory-import.json"
)

// PathMemoryMirror is the relative path from the project root to the
// memory mirror file. Constructed from directory and file constants.
var PathMemoryMirror = filepath.Join(dir.Context, dir.Memory, Mirror)
