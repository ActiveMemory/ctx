//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

import (
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/index"
)

// ScoredEntry is an entry block with a computed relevance score.
//
// Fields:
//   - EntryBlock: The parsed entry block from a knowledge file
//   - Score: Combined recency + relevance score (0.0–2.0)
//   - Tokens: Pre-computed token estimate of the full body
type ScoredEntry struct {
	index.EntryBlock
	Score  float64
	Tokens int
}

// recencyScore returns a score based on the entry's age.
//
// Scoring brackets:
//   - 0–7 days: 1.0
//   - 8–30 days: 0.7
//   - 31–90 days: 0.4
//   - 90+ days: 0.2
//
// Parameters:
//   - eb: Entry block to score
//   - now: Current time for age calculation
//
// Returns:
//   - float64: Recency score between 0.2 and 1.0
func recencyScore(eb *index.EntryBlock, now time.Time) float64 {
	entryDate, err := time.ParseInLocation("2006-01-02", eb.Entry.Date, time.Local)
	if err != nil {
		return 0.2
	}
	days := int(now.Sub(entryDate).Hours() / 24)
	switch {
	case days <= 7:
		return 1.0
	case days <= 30:
		return 0.7
	case days <= 90:
		return 0.4
	default:
		return 0.2
	}
}

// relevanceScore computes keyword overlap between an entry and active tasks.
//
// Counts how many task keywords appear in the entry's title and body.
// Normalized to 1.0 at 3+ matches.
//
// Parameters:
//   - eb: Entry block to score
//   - keywords: Lowercase keywords extracted from active tasks
//
// Returns:
//   - float64: Relevance score between 0.0 and 1.0
func relevanceScore(eb *index.EntryBlock, keywords []string) float64 {
	if len(keywords) == 0 {
		return 0.0
	}
	text := strings.ToLower(eb.BlockContent())
	matches := 0
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			matches++
		}
	}
	if matches >= 3 {
		return 1.0
	}
	return float64(matches) / 3.0
}

// scoreEntry computes the combined relevance score for an entry block.
//
// Superseded entries always get score 0.0.
// All other entries get recency + task relevance (range 0.0–2.0).
//
// Parameters:
//   - eb: Entry block to score
//   - keywords: Task keywords for relevance matching
//   - now: Current time for recency calculation
//
// Returns:
//   - float64: Combined score (0.0–2.0), or 0.0 if superseded
func scoreEntry(eb *index.EntryBlock, keywords []string, now time.Time) float64 {
	if eb.IsSuperseded() {
		return 0.0
	}
	return recencyScore(eb, now) + relevanceScore(eb, keywords)
}

// stopWords is a set of common English words to exclude from keyword extraction.
var stopWords = map[string]bool{
	"the": true, "and": true, "for": true, "that": true, "this": true,
	"with": true, "from": true, "are": true, "was": true, "were": true,
	"been": true, "have": true, "has": true, "had": true, "but": true,
	"not": true, "you": true, "all": true, "can": true, "her": true,
	"his": true, "she": true, "its": true, "our": true, "they": true,
	"will": true, "each": true, "make": true, "like": true, "use": true,
	"way": true, "may": true, "any": true, "into": true, "when": true,
	"which": true, "their": true, "about": true, "would": true,
	"there": true, "what": true, "also": true, "should": true,
	"after": true, "before": true, "than": true, "then": true,
	"them": true, "could": true, "more": true, "some": true,
	"other": true, "only": true, "just": true, "see": true,
	"add": true, "new": true, "update": true, "how": true,
}

// extractTaskKeywords extracts meaningful keywords from task text.
//
// Splits task text on whitespace and punctuation, lowercases, and filters
// out stop words and words shorter than 3 characters. Deduplicates results.
//
// Parameters:
//   - tasks: Active task strings (e.g., "- [ ] Implement feature X")
//
// Returns:
//   - []string: Unique lowercase keywords
func extractTaskKeywords(tasks []string) []string {
	seen := make(map[string]bool)
	var keywords []string
	for _, t := range tasks {
		// Split on whitespace and common punctuation
		words := strings.FieldsFunc(strings.ToLower(t), func(r rune) bool {
			isAlnum := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
			return !isAlnum && r != '-' && r != '_'
		})
		for _, w := range words {
			if len(w) < 3 || stopWords[w] || seen[w] {
				continue
			}
			seen[w] = true
			keywords = append(keywords, w)
		}
	}
	return keywords
}

// scoreEntries scores and sorts entry blocks by relevance.
//
// Parameters:
//   - blocks: Parsed entry blocks from a knowledge file
//   - keywords: Task keywords for relevance matching
//   - now: Current time for recency scoring
//
// Returns:
//   - []ScoredEntry: Entries sorted by score descending, with token estimates
func scoreEntries(blocks []index.EntryBlock, keywords []string, now time.Time) []ScoredEntry {
	scored := make([]ScoredEntry, 0, len(blocks))
	for i := range blocks {
		s := scoreEntry(&blocks[i], keywords, now)
		tokens := context.EstimateTokensString(blocks[i].BlockContent())
		scored = append(scored, ScoredEntry{
			EntryBlock: blocks[i],
			Score:      s,
			Tokens:     tokens,
		})
	}
	// Sort by score descending (stable for equal scores — preserves file order)
	for i := 1; i < len(scored); i++ {
		for j := i; j > 0 && scored[j].Score > scored[j-1].Score; j-- {
			scored[j], scored[j-1] = scored[j-1], scored[j]
		}
	}
	return scored
}
