//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"strings"

	cfgMemory "github.com/ActiveMemory/ctx/internal/config/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Classify assigns a target file type to an entry based on keyword heuristics.
//
// Rules are loaded from .ctxrc (classify_rules) with fallback to built-in
// defaults. Matching is case-insensitive. The first rule with a keyword
// match wins (default priority: conventions > decisions > learnings > tasks > skip).
//
// Parameters:
//   - entry: Parsed memory entry to classify
func Classify(entry Entry) Classification {
	lower := strings.ToLower(entry.Text)

	for _, rule := range rc.ClassifyRules() {
		var matched []string
		for _, kw := range rule.Keywords {
			if strings.Contains(lower, kw) {
				matched = append(matched, kw)
			}
		}
		if len(matched) > 0 {
			return Classification{
				Target:   rule.Target,
				Keywords: matched,
			}
		}
	}

	return Classification{Target: cfgMemory.TargetSkip}
}
