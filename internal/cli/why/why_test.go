//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package why

import (
	"bytes"
	"strings"
	"testing"

	whyRoot "github.com/ActiveMemory/ctx/internal/cli/why/cmd/root"
	"github.com/spf13/cobra"
)

// newTestCmd creates a Cmd() wired to capture output in the returned buffer.
func newTestCmd() (*cobra.Command, *bytes.Buffer) {
	cmd := Cmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	return cmd, &buf
}

func TestShowDoc_Manifesto(t *testing.T) {
	cmd, buf := newTestCmd()

	showErr := whyRoot.ShowDoc(cmd, "manifesto")
	if showErr != nil {
		t.Fatalf("ShowDoc(manifesto) error = %v", showErr)
	}

	output := buf.String()
	if len(output) == 0 {
		t.Fatal("ShowDoc(manifesto) produced empty output")
	}
	if !strings.Contains(output, "Manifesto") {
		t.Errorf("output missing 'Manifesto', got %d bytes", len(output))
	}
}

func TestShowDoc_UnknownAlias(t *testing.T) {
	cmd, _ := newTestCmd()

	showErr := whyRoot.ShowDoc(cmd, "nonexistent")
	if showErr == nil {
		t.Fatal("expected error for unknown document alias")
	}
	if !strings.Contains(showErr.Error(), "unknown document") {
		t.Errorf("error = %q, want mention of 'unknown document'", showErr.Error())
	}
}

func TestRunWhy_DirectArg(t *testing.T) {
	cmd, buf := newTestCmd()
	cmd.SetArgs([]string{"manifesto"})

	execErr := cmd.Execute()
	if execErr != nil {
		t.Fatalf("Execute() error = %v", execErr)
	}

	output := buf.String()
	if len(output) == 0 {
		t.Fatal("expected non-empty output for 'ctx why manifesto'")
	}
	if !strings.Contains(output, "Manifesto") {
		t.Errorf("output missing 'Manifesto', got %d bytes", len(output))
	}
}

func TestExtractAdmonitionTitle(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "quoted title",
			input: `!!! info "Title"`,
			want:  "Title",
		},
		{
			name:  "no title",
			input: `!!! warning`,
			want:  "",
		},
		{
			name:  "note with spaces",
			input: `!!! note "Important Note"`,
			want:  "Important Note",
		},
		{
			name:  "single quote mark only",
			input: `!!! tip "`,
			want:  "",
		},
		{
			name:  "empty quotes",
			input: `!!! tip ""`,
			want:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := whyRoot.ExtractAdmonitionTitle(tc.input)
			if got != tc.want {
				t.Errorf("ExtractAdmonitionTitle(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestExtractTabTitle(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple tab",
			input: `=== "Tab 1"`,
			want:  "Tab 1",
		},
		{
			name:  "tab with spaces",
			input: `=== "Without ctx"`,
			want:  "Without ctx",
		},
		{
			name:  "no quotes",
			input: `=== plain`,
			want:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := whyRoot.ExtractTabTitle(tc.input)
			if got != tc.want {
				t.Errorf("ExtractTabTitle(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
