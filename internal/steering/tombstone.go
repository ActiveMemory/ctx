//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

import (
	"strings"

	cfgSteering "github.com/ActiveMemory/ctx/internal/config/steering"
)

// Tombstone re-exports the canonical tombstone marker from
// [cfgSteering.Tombstone] so callers within this package and
// nearby tests can reference it without an additional import.
// The single source of truth is the config package.
const Tombstone = cfgSteering.Tombstone

// HasTombstone reports whether the given steering file body
// still contains the [Tombstone] marker, indicating it is an
// unmodified placeholder from `ctx steering init`.
//
// Files containing the tombstone are excluded from:
//   - the agent context packet assembled by `ctx agent`
//   - MCP `ctx_steering_get` results
//   - native-tool exports via `ctx steering sync` (Cursor,
//     Cline, Kiro)
//
// Parameters:
//   - body: steering file body content (the markdown after
//     the YAML frontmatter)
//
// Returns:
//   - bool: true if the tombstone marker is present in the body
func HasTombstone(body string) bool {
	return strings.Contains(body, cfgSteering.Tombstone)
}
