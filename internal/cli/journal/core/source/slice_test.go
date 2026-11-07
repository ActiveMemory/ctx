//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package source

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/entity"
)

func TestParseTurnRange(t *testing.T) {
	tests := []struct {
		name      string
		spec      string
		total     int
		wantStart int
		wantEnd   int
		wantErr   bool
	}{
		{
			name:      "valid range",
			spec:      "2:5",
			total:     10,
			wantStart: 2,
			wantEnd:   5,
		},
		{
			name:      "single turn shorthand",
			spec:      "7",
			total:     10,
			wantStart: 7,
			wantEnd:   7,
		},
		{
			name:      "clamps end to totalTurns",
			spec:      "8:25",
			total:     10,
			wantStart: 8,
			wantEnd:   10,
		},
		{
			name:    "rejects zero start",
			spec:    "0:5",
			total:   10,
			wantErr: true,
		},
		{
			name:    "rejects inverted range",
			spec:    "8:3",
			total:   10,
			wantErr: true,
		},
		{
			name:    "rejects non-numeric",
			spec:    "abc:def",
			total:   10,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd, err := ParseTurnRange(tt.spec, tt.total)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseTurnRange(%q) error = %v, wantErr %v",
					tt.spec, err, tt.wantErr)
			}
			if !tt.wantErr {
				if gotStart != tt.wantStart || gotEnd != tt.wantEnd {
					t.Errorf("ParseTurnRange(%q) = (%d, %d), want (%d, %d)",
						tt.spec, gotStart, gotEnd, tt.wantStart, tt.wantEnd)
				}
			}
		})
	}
}

func TestMatchTurns_AndSelectiveRender(t *testing.T) {
	ts := time.Date(2026, 11, 4, 14, 22, 10, 0, time.UTC)
	msgs := []entity.Message{
		{
			Role:      "user",
			Timestamp: ts,
			Text:      "Investigate pad test failure",
		},
		{
			Role:      "user",
			Timestamp: ts.Add(5 * time.Second),
			ToolResults: []entity.ToolResult{
				{
					Content: "1→FAIL internal/cli/pad: cipher: message authentication failed",
					IsError: true,
				},
			},
		},
		{
			Role:      "assistant",
			Timestamp: ts.Add(10 * time.Second),
			Text:      "The pad test overwrote the user key.",
		},
	}

	re, err := regex.CompileDynamic(regex.CaseInsensitivePrefix + "cipher: message")
	if err != nil {
		t.Fatal(err)
	}
	matches := MatchTurns(msgs, 1, 3, re)
	if len(matches) != 1 {
		t.Fatalf("MatchTurns() returned %d matches, want 1", len(matches))
	}
	if matches[0].Turn != 2 {
		t.Errorf("match turn = %d, want 2", matches[0].Turn)
	}
	if !strings.Contains(matches[0].Snippet, "cipher: message authentication") {
		t.Errorf("unexpected snippet %q", matches[0].Snippet)
	}

	// Verify bounded rendering and next: continuation banner.
	sess := &entity.Session{
		ID:       "96765858-aaaa-bbbb",
		Slug:     "pad-debug",
		Messages: msgs,
	}
	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&out)

	if err := runSelectiveShow(cmd, sess, Opts{
		Turns:    "1:2",
		MaxChars: 6000,
	}); err != nil {
		t.Fatalf("runSelectiveShow() error = %v", err)
	}
	rendered := out.String()
	if !strings.Contains(rendered, "next: ctx journal source --show 96765858 --turns 3:3") {
		t.Errorf("expected next continuation command in:\n%s", rendered)
	}
}
