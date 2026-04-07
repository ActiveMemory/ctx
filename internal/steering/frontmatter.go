//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

import (
	cfgSteering "github.com/ActiveMemory/ctx/internal/config/steering"
)

// applyDefaults sets default values for fields not present in
// the parsed frontmatter.
//
// Parameters:
//   - sf: steering file to populate with defaults
func applyDefaults(sf *SteeringFile) {
	if sf.Inclusion == "" {
		sf.Inclusion = defaultInclusion
	}
	if sf.Priority == 0 {
		sf.Priority = cfgSteering.DefaultPriority
	}
	// Tools: nil means all tools — no default needed.
}
