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

// TargetSkip indicates an entry that doesn't match any classification rule.
const TargetSkip = "skip"

// ClassifyRule maps keyword patterns to a target entry type.
//
// Fields:
//   - Target: entry type constant (convention, decision, learning, task)
//   - Keywords: case-insensitive keyword patterns to match
type ClassifyRule struct {
	Target   string   `yaml:"target"`
	Keywords []string `yaml:"keywords"`
}

// DefaultClassifyRules are the built-in heuristic rules for classifying
// memory entries. Users can override this list via the classify_rules
// key in .ctxrc. Rules are evaluated in priority order.
var DefaultClassifyRules = []ClassifyRule{
	{
		Target: "convention",
		Keywords: []string{
			"always use", "prefer", "convention",
			"never use", "standard", "always ",
		},
	},
	{
		Target: "decision",
		Keywords: []string{
			"decided", "chose", "trade-off",
			"approach", "over", "instead of",
		},
	},
	{
		Target: "learning",
		Keywords: []string{
			"gotcha", "learned", "watch out",
			"bug", "caveat", "careful", "turns out",
		},
	},
	{
		Target:   "task",
		Keywords: []string{"todo", "need to", "follow up", "should", "task"},
	},
}
