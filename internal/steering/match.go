//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

import (
	"slices"
	"strings"
)

// matchInclusion checks whether a steering file should be included
// based on its inclusion mode.
func matchInclusion(
	sf *SteeringFile, promptLower string,
	manualNames []string,
) bool {
	switch sf.Inclusion {
	case InclusionAlways:
		return true
	case InclusionAuto:
		if sf.Description == "" {
			return false
		}
		return strings.Contains(promptLower, strings.ToLower(sf.Description))
	case InclusionManual:
		return slices.Contains(manualNames, sf.Name)
	default:
		return false
	}
}

// matchTool checks whether a steering file applies to the given tool.
// When the file's Tools list is nil or empty, it applies to all tools.
// When tool is empty, no filtering is applied.
func matchTool(sf *SteeringFile, tool string) bool {
	if tool == "" {
		return true
	}
	if len(sf.Tools) == 0 {
		return true
	}
	return slices.Contains(sf.Tools, tool)
}
