//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

import (
	"testing"
	"time"

	"github.com/ActiveMemory/ctx/internal/index"
)

func makeBlock(date, title, body string) index.EntryBlock {
	header := "## [" + date + "-120000] " + title
	lines := []string{header}
	if body != "" {
		lines = append(lines, "", body)
	}
	return index.EntryBlock{
		Entry: index.Entry{
			Timestamp: date + "-120000",
			Date:      date,
			Title:     title,
		},
		Lines: lines,
	}
}

func TestRecencyScore(t *testing.T) {
	now := time.Date(2026, 2, 19, 12, 0, 0, 0, time.Local)
	tests := []struct {
		name string
		date string
		want float64
	}{
		{"today", "2026-02-19", 1.0},
		{"3 days ago", "2026-02-16", 1.0},
		{"7 days ago", "2026-02-12", 1.0},
		{"10 days ago", "2026-02-09", 0.7},
		{"30 days ago", "2026-01-20", 0.7},
		{"45 days ago", "2026-01-05", 0.4},
		{"90 days ago", "2025-11-21", 0.4},
		{"120 days ago", "2025-10-22", 0.2},
		{"invalid date", "not-a-date", 0.2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eb := makeBlock(tt.date, "Test", "")
			got := recencyScore(&eb, now)
			if got != tt.want {
				t.Errorf("recencyScore(%s) = %v, want %v", tt.date, got, tt.want)
			}
		})
	}
}

func TestRelevanceScore(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		keywords []string
		want     float64
	}{
		{"no keywords", "some content about hooks", nil, 0.0},
		{"no matches", "some content about hooks", []string{"database", "cache"}, 0.0},
		{"one match", "fix the hook edge case", []string{"hook", "cache"}, 1.0 / 3.0},
		{"two matches", "fix the hook edge case in agent", []string{"hook", "agent", "cache"}, 2.0 / 3.0},
		{"three matches", "fix hook in agent scoring", []string{"hook", "agent", "scoring"}, 1.0},
		{"more than three", "hook agent scoring budget", []string{"hook", "agent", "scoring", "budget"}, 1.0},
		{"case insensitive", "Hook AGENT Scoring", []string{"hook", "agent", "scoring"}, 1.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eb := makeBlock("2026-02-19", "Test", tt.body)
			got := relevanceScore(&eb, tt.keywords)
			if got != tt.want {
				t.Errorf("relevanceScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScoreEntry_Superseded(t *testing.T) {
	now := time.Date(2026, 2, 19, 12, 0, 0, 0, time.Local)
	eb := index.EntryBlock{
		Entry: index.Entry{
			Timestamp: "2026-02-19-120000",
			Date:      "2026-02-19",
			Title:     "Old decision",
		},
		Lines: []string{
			"## [2026-02-19-120000] Old decision",
			"",
			"~~Superseded by [2026-02-19-130000] New decision~~",
		},
	}
	got := scoreEntry(&eb, []string{"decision"}, now)
	if got != 0.0 {
		t.Errorf("superseded entry score = %v, want 0.0", got)
	}
}

func TestScoreEntry_Combined(t *testing.T) {
	now := time.Date(2026, 2, 19, 12, 0, 0, 0, time.Local)
	// Recent + relevant = high score
	eb := makeBlock("2026-02-19", "Hook edge cases", "hooks fail silently in agent mode")
	got := scoreEntry(&eb, []string{"hook", "agent", "scoring"}, now)
	// recency = 1.0, relevance = 2/3 ≈ 1.667
	if got < 1.66 || got > 1.67 {
		t.Errorf("scoreEntry() = %v, want ~1.667", got)
	}
}

func TestExtractTaskKeywords(t *testing.T) {
	tasks := []string{
		"- [ ] Implement hook scoring for the agent",
		"- [ ] Fix budget allocation in ctx agent",
	}
	keywords := extractTaskKeywords(tasks)

	// Should include meaningful words
	kwSet := make(map[string]bool)
	for _, kw := range keywords {
		kwSet[kw] = true
	}

	if !kwSet["hook"] {
		t.Error("expected 'hook' in keywords")
	}
	if !kwSet["scoring"] {
		t.Error("expected 'scoring' in keywords")
	}
	if !kwSet["budget"] {
		t.Error("expected 'budget' in keywords")
	}
	if !kwSet["agent"] {
		t.Error("expected 'agent' in keywords")
	}
	if !kwSet["allocation"] {
		t.Error("expected 'allocation' in keywords")
	}
	if !kwSet["implement"] {
		t.Error("expected 'implement' in keywords")
	}

	// Should exclude stop words and short words
	if kwSet["the"] {
		t.Error("'the' should be excluded (stop word)")
	}
	if kwSet["for"] {
		t.Error("'for' should be excluded (stop word)")
	}
	if kwSet["in"] {
		t.Error("'in' should be excluded (short word)")
	}

	// Should deduplicate
	count := 0
	for _, kw := range keywords {
		if kw == "agent" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("'agent' appears %d times, want 1", count)
	}
}

func TestExtractTaskKeywords_Empty(t *testing.T) {
	keywords := extractTaskKeywords(nil)
	if len(keywords) != 0 {
		t.Errorf("expected empty keywords, got %v", keywords)
	}
}

func TestScoreEntries_Ordering(t *testing.T) {
	now := time.Date(2026, 2, 19, 12, 0, 0, 0, time.Local)
	blocks := []index.EntryBlock{
		makeBlock("2025-10-01", "Old irrelevant", "something unrelated"),
		makeBlock("2026-02-19", "Recent relevant", "hook scoring for agent"),
		makeBlock("2026-02-10", "Medium age", "hook configuration"),
	}
	keywords := []string{"hook", "scoring", "agent"}
	scored := scoreEntries(blocks, keywords, now)

	if len(scored) != 3 {
		t.Fatalf("expected 3 scored entries, got %d", len(scored))
	}
	// First should be the recent+relevant one
	if scored[0].Entry.Title != "Recent relevant" {
		t.Errorf("expected 'Recent relevant' first, got %q", scored[0].Entry.Title)
	}
	// Scores should be descending
	for i := 1; i < len(scored); i++ {
		if scored[i].Score > scored[i-1].Score {
			t.Errorf("scored[%d].Score (%v) > scored[%d].Score (%v)",
				i, scored[i].Score, i-1, scored[i-1].Score)
		}
	}
}

func TestScoreEntries_Empty(t *testing.T) {
	now := time.Now()
	scored := scoreEntries(nil, nil, now)
	if len(scored) != 0 {
		t.Errorf("expected empty scored entries, got %d", len(scored))
	}
}

func TestScoreEntries_TokenEstimate(t *testing.T) {
	now := time.Now()
	blocks := []index.EntryBlock{
		makeBlock("2026-02-19", "Test entry", "This is some body content for testing tokens."),
	}
	scored := scoreEntries(blocks, nil, now)
	if scored[0].Tokens <= 0 {
		t.Error("expected positive token estimate")
	}
}
