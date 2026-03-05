//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
)

// Classification is the result of heuristic entry classification.
type Classification struct {
	Target   string   // config.Entry* constant or "skip"
	Keywords []string // Keywords that triggered the match
}

// TargetSkip indicates an entry that doesn't match any classification rule.
const TargetSkip = "skip"

// classRule maps keyword patterns to a target file type.
type classRule struct {
	target   string
	keywords []string
}

// rules are evaluated in priority order: conventions > decisions > learnings > tasks.
var rules = []classRule{
	{
		target:   config.EntryConvention,
		keywords: []string{"always use", "prefer", "convention", "never use", "standard", "always "},
	},
	{
		target:   config.EntryDecision,
		keywords: []string{"decided", "chose", "trade-off", "approach", "over", "instead of"},
	},
	{
		target:   config.EntryLearning,
		keywords: []string{"gotcha", "learned", "watch out", "bug", "caveat", "careful", "turns out"},
	},
	{
		target:   config.EntryTask,
		keywords: []string{"todo", "need to", "follow up", "should", "task"},
	},
}

// Classify assigns a target file type to an entry based on keyword heuristics.
//
// Matching is case-insensitive. The first rule with a keyword match wins
// (priority: conventions > decisions > learnings > tasks > skip).
func Classify(entry Entry) Classification {
	lower := strings.ToLower(entry.Text)

	for _, rule := range rules {
		var matched []string
		for _, kw := range rule.keywords {
			if strings.Contains(lower, kw) {
				matched = append(matched, kw)
			}
		}
		if len(matched) > 0 {
			return Classification{
				Target:   rule.target,
				Keywords: matched,
			}
		}
	}

	return Classification{Target: TargetSkip}
}
