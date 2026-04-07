//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

import (
	"slices"
	"strings"

	cfgSteering "github.com/ActiveMemory/ctx/internal/config/steering"
)

// matchInclusion checks whether a steering file should be
// included based on its inclusion mode.
//
// Parameters:
//   - sf: steering file to evaluate
//   - promptLower: lowercased user prompt for auto-match
//   - manualNames: explicitly requested steering file names
//
// Returns:
//   - bool: true when the file matches the inclusion criteria
func matchInclusion(
	sf *SteeringFile, promptLower string,
	manualNames []string,
) bool {
	switch sf.Inclusion {
	case cfgSteering.InclusionAlways:
		return true
	case cfgSteering.InclusionAuto:
		if sf.Description == "" {
			return false
		}
		return strings.Contains(promptLower, strings.ToLower(sf.Description))
	case cfgSteering.InclusionManual:
		return slices.Contains(manualNames, sf.Name)
	default:
		return false
	}
}

// matchTool checks whether a steering file applies to the given
// tool. When the file's Tools list is nil or empty, it applies
// to all tools. When tool is empty, no filtering is applied.
//
// Parameters:
//   - sf: steering file whose tool scope to check
//   - tool: tool identifier to match against
//
// Returns:
//   - bool: true when the file applies to the tool
func matchTool(sf *SteeringFile, tool string) bool {
	if tool == "" {
		return true
	}
	if len(sf.Tools) == 0 {
		return true
	}
	return slices.Contains(sf.Tools, tool)
}
