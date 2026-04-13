//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package budget

import (
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/agent/core/score"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/index"
)

func makeBlock(date, title, body string) index.EntryBlock {
	header := "## [" + date + "-120000] " + title
	lines := []string{header}
	if body != "" {
		lines = append(lines, "", body)
	}
	return index.EntryBlock{
		Entry: entity.IndexEntry{
			Timestamp: date + "-120000",
			Date:      date,
			Title:     title,
		},
		Lines: lines,
	}
}

func TestFitItemsInBudget(t *testing.T) {
	tests := []struct {
		name   string
		items  []string
		budget int
		want   int
	}{
		{"all fit", []string{"short", "words"}, 1000, 2},
		{
			"none fit but one forced",
			[]string{"a very long item that exceeds budget"},
			1, 1,
		},
		{"empty items", nil, 1000, 0},
		{
			"partial fit",
			[]string{"a]longer item here", "another longer one", "third longer item"},
			6, 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FitItems(tt.items, tt.budget)
			if len(got) != tt.want {
				t.Errorf("FitItems() returned %d items, want %d", len(got), tt.want)
			}
		})
	}
}

func TestSplitBudget(t *testing.T) {
	tests := []struct {
		name       string
		total      int
		aEntries   []score.Entry
		bEntries   []score.Entry
		wantAMin   int
		wantBMin   int
		wantAMax   int
		wantBMax   int
		wantAExact int
		wantBExact int
		exact      bool
	}{
		{
			name:       "both empty",
			total:      1000,
			wantAExact: 0, wantBExact: 0, exact: true,
		},
		{
			name:  "a empty",
			total: 1000,
			bEntries: []score.Entry{
				{EntryBlock: makeBlock("2026-02-19", "B", "body"), Tokens: 100},
			},
			wantAExact: 0, wantBExact: 1000, exact: true,
		},
		{
			name:  "b empty",
			total: 1000,
			aEntries: []score.Entry{
				{EntryBlock: makeBlock("2026-02-19", "A", "body"), Tokens: 100},
			},
			wantAExact: 1000, wantBExact: 0, exact: true,
		},
		{
			name:  "both fit",
			total: 1000,
			aEntries: []score.Entry{
				{EntryBlock: makeBlock("2026-02-19", "A", "body"), Tokens: 200},
			},
			bEntries: []score.Entry{
				{EntryBlock: makeBlock("2026-02-19", "B", "body"), Tokens: 300},
			},
			wantAExact: 200, wantBExact: 300, exact: true,
		},
		{
			name:  "exceeds budget gets proportional split",
			total: 100,
			aEntries: []score.Entry{
				{EntryBlock: makeBlock("2026-02-19", "A", "body"), Tokens: 500},
			},
			bEntries: []score.Entry{
				{EntryBlock: makeBlock("2026-02-19", "B", "body"), Tokens: 500},
			},
			// Each gets at least 30%, split proportionally
			wantAMin: 30, wantAMax: 70,
			wantBMin: 30, wantBMax: 70,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotA, gotB := Split(tt.total, tt.aEntries, tt.bEntries)
			if tt.exact {
				if gotA != tt.wantAExact || gotB != tt.wantBExact {
					t.Errorf("Split() = (%d, %d), want (%d, %d)",
						gotA, gotB, tt.wantAExact, tt.wantBExact)
				}
			} else {
				if gotA < tt.wantAMin || gotA > tt.wantAMax {
					t.Errorf("Split() a = %d, want [%d, %d]",
						gotA, tt.wantAMin, tt.wantAMax)
				}
				if gotB < tt.wantBMin || gotB > tt.wantBMax {
					t.Errorf("Split() b = %d, want [%d, %d]",
						gotB, tt.wantBMin, tt.wantBMax)
				}
				if gotA+gotB != tt.total {
					t.Errorf("Split() a+b = %d, want %d", gotA+gotB, tt.total)
				}
			}
		})
	}
}

func TestFillSection(t *testing.T) {
	entries := []score.Entry{
		{
			EntryBlock: makeBlock(
				"2026-02-19", "High score",
				"important body content here",
			),
			Score:  1.8,
			Tokens: 10,
		},
		{
			EntryBlock: makeBlock("2026-02-10", "Medium score", "less important body"),
			Score:      1.0,
			Tokens:     8,
		},
		{
			EntryBlock: makeBlock("2025-10-01", "Low score", "old entry body text"),
			Score:      0.3,
			Tokens:     7,
		},
	}

	t.Run("all fit", func(t *testing.T) {
		full, summaries := FillSection(entries, 1000)
		if len(full) != 3 {
			t.Errorf("expected 3 full entries, got %d", len(full))
		}
		if len(summaries) != 0 {
			t.Errorf("expected 0 summaries, got %d", len(summaries))
		}
	})

	t.Run("partial fit with summaries", func(t *testing.T) {
		// Budget allows ~80% = 16 tokens for full entries
		// First entry (10 tokens) fits, second (8 tokens) = 18 total > 16
		full, summaries := FillSection(entries, 20)
		if len(full) != 1 {
			t.Errorf("expected 1 full entry, got %d", len(full))
		}
		if len(summaries) != 2 {
			t.Errorf("expected 2 summaries, got %d", len(summaries))
		}
	})

	t.Run("empty entries", func(t *testing.T) {
		full, summaries := FillSection(nil, 1000)
		if full != nil || summaries != nil {
			t.Error("expected nil for empty entries")
		}
	})

	t.Run("zero budget", func(t *testing.T) {
		full, summaries := FillSection(entries, 0)
		if full != nil || summaries != nil {
			t.Error("expected nil for zero budget")
		}
	})

	t.Run("superseded entries skipped", func(t *testing.T) {
		superseded := []score.Entry{
			{
				EntryBlock: index.EntryBlock{
					Entry: entity.IndexEntry{
						Timestamp: "2026-02-19-120000",
						Date:      "2026-02-19",
						Title:     "Old decision",
					},
					Lines: []string{
						"## [2026-02-19-120000] Old decision",
						"~~Superseded by newer~~",
					},
				},
				Score:  0.0,
				Tokens: 5,
			},
			entries[0],
		}
		full, _ := FillSection(superseded, 1000)
		if len(full) != 1 {
			t.Errorf("expected 1 full entry (superseded skipped), got %d", len(full))
		}
		if len(full) > 0 && !strings.Contains(full[0], "High score") {
			t.Errorf("expected high score entry, got %q", full[0])
		}
	})
}

func TestEstimateSliceTokens(t *testing.T) {
	items := []string{"hello", "world"}
	got := EstimateSliceTokens(items)
	if got <= 0 {
		t.Errorf("expected positive token estimate, got %d", got)
	}
}

func TestRenderMarkdownPacket(t *testing.T) {
	pkt := &AssembledPacket{
		ReadOrder:    []string{".context/CONSTITUTION.md"},
		Constitution: []string{"Never violate"},
		Tasks:        []string{"- [ ] Do something"},
		Conventions:  []string{"Use gofmt"},
		Decisions:    []string{"## [2026-02-19-120000] Use JWT\n\nFor auth."},
		Learnings: []string{
			"## [2026-02-19-130000] Hooks fail silently\n\nCheck stderr.",
		},
		Summaries:   []string{"Old learning about paths"},
		Instruction: "Confirm context reading.",
		Budget:      8000,
		TokensUsed:  2000,
	}

	output := RenderMarkdownPacket(pkt)

	checks := []string{
		"# Context Packet",
		"Budget: 8000",
		"## Read These Files",
		"CONSTITUTION.md",
		"## Constitution (NEVER VIOLATE)",
		"Never violate",
		"## Current Tasks",
		"Do something",
		"## Key Conventions",
		"Use gofmt",
		"## Recent Decisions",
		"Use JWT",
		"## Key Learnings",
		"Hooks fail silently",
		"## Also Noted",
		"Old learning about paths",
		"Confirm context reading.",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("output missing %q", check)
		}
	}
}

func TestRenderMarkdownPacket_Empty(t *testing.T) {
	pkt := &AssembledPacket{
		Instruction: "Do stuff.",
		Budget:      100,
	}

	output := RenderMarkdownPacket(pkt)

	if !strings.Contains(output, "# Context Packet") {
		t.Error("missing header")
	}
	if !strings.Contains(output, "Do stuff.") {
		t.Error("missing instruction")
	}
	// Should not contain section headers for empty sections
	if strings.Contains(output, "## Current Tasks") {
		t.Error("should not render empty tasks section")
	}
}

func TestRenderMarkdownPacket_WithSteering(t *testing.T) {
	pkt := &AssembledPacket{
		ReadOrder:    []string{".context/CONSTITUTION.md"},
		Constitution: []string{"Never violate"},
		Steering:     []string{"Use RESTful conventions", "Always return JSON"},
		Instruction:  "Confirm context reading.",
		Budget:       8000,
		TokensUsed:   500,
	}

	output := RenderMarkdownPacket(pkt)

	checks := []string{
		"## Steering",
		"Use RESTful conventions",
		"Always return JSON",
	}
	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("output missing %q", check)
		}
	}
}

func TestRenderMarkdownPacket_WithSkill(t *testing.T) {
	pkt := &AssembledPacket{
		ReadOrder:   []string{".context/CONSTITUTION.md"},
		Skill:       "# React Patterns\n\nUse functional components.",
		Instruction: "Confirm context reading.",
		Budget:      8000,
		TokensUsed:  500,
	}

	output := RenderMarkdownPacket(pkt)

	checks := []string{
		"## Skill",
		"React Patterns",
		"Use functional components.",
	}
	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("output missing %q", check)
		}
	}
}

func TestRenderMarkdownPacket_NoSteeringOrSkill(t *testing.T) {
	pkt := &AssembledPacket{
		ReadOrder:   []string{".context/CONSTITUTION.md"},
		Instruction: "Confirm context reading.",
		Budget:      8000,
		TokensUsed:  500,
	}

	output := RenderMarkdownPacket(pkt)

	if strings.Contains(output, "## Steering") {
		t.Error("should not render empty steering section")
	}
	if strings.Contains(output, "## Skill") {
		t.Error("should not render empty skill section")
	}
}

func TestAssemblePacket_WithSteering(t *testing.T) {
	ctx := &entity.Context{}
	bodies := []string{"Rule one", "Rule two"}

	pkt := AssemblePacket(ctx, 8000, bodies, "", nil)

	if len(pkt.Steering) == 0 {
		t.Error("expected steering files in packet")
	}
	if pkt.Steering[0] != "Rule one" {
		t.Errorf("expected first steering body %q, got %q", "Rule one", pkt.Steering[0])
	}
}

func TestAssemblePacket_WithSkill(t *testing.T) {
	ctx := &entity.Context{}
	skillBody := "# My Skill\n\nDo things."

	pkt := AssemblePacket(ctx, 8000, nil, skillBody, nil)

	if pkt.Skill != skillBody {
		t.Errorf("expected skill body %q, got %q", skillBody, pkt.Skill)
	}
}

func TestAssemblePacket_NoSteeringNoSkill(t *testing.T) {
	ctx := &entity.Context{}

	pkt := AssemblePacket(ctx, 8000, nil, "", nil)

	if len(pkt.Steering) != 0 {
		t.Errorf("expected no steering, got %d", len(pkt.Steering))
	}
	if pkt.Skill != "" {
		t.Errorf("expected empty skill, got %q", pkt.Skill)
	}
}

func TestAssemblePacket_SteeringRespectsBudget(t *testing.T) {
	ctx := &entity.Context{}
	// Use a very small budget so steering gets truncated
	bigBody := strings.Repeat("x", 5000)
	bodies := []string{bigBody, bigBody}

	pkt := AssemblePacket(ctx, 100, bodies, "", nil)

	// With a tiny budget, at most one steering body should fit
	// (FitItems always includes at least one)
	if len(pkt.Steering) > 1 {
		t.Errorf("expected at most 1 steering body with tiny budget, got %d", len(pkt.Steering))
	}
}

func TestAssemblePacket_SkillOmittedWhenBudgetExhausted(t *testing.T) {
	ctx := &entity.Context{}
	// Use a very small budget
	pkt := AssemblePacket(ctx, 1, nil, strings.Repeat("x", 5000), nil)

	// Skill should be omitted when budget is exhausted
	if pkt.Skill != "" {
		t.Error("expected skill to be omitted when budget exhausted")
	}
}
