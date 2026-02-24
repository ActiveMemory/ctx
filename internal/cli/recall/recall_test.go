//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/recall/parser"
)

func TestCmd(t *testing.T) {
	cmd := Cmd()

	if cmd == nil {
		t.Fatal("Cmd() returned nil")
	}

	if cmd.Use != "recall" {
		t.Errorf("Cmd().Use = %q, want %q", cmd.Use, "recall")
	}

	if cmd.Short == "" {
		t.Error("Cmd().Short is empty")
	}

	if cmd.Long == "" {
		t.Error("Cmd().Long is empty")
	}
}

func TestCmd_HasSubcommands(t *testing.T) {
	cmd := Cmd()

	expectedSubs := []string{"list", "show", "export", "lock", "unlock"}
	subs := make(map[string]bool)

	for _, sub := range cmd.Commands() {
		subs[sub.Name()] = true
	}

	for _, exp := range expectedSubs {
		if !subs[exp] {
			t.Errorf("missing subcommand: %s", exp)
		}
	}
}

func TestRecallListCmd_Flags(t *testing.T) {
	cmd := Cmd()

	var listCmd *cobra.Command
	for _, sub := range cmd.Commands() {
		if sub.Name() == "list" {
			listCmd = sub
			break
		}
	}

	if listCmd == nil {
		t.Fatal("list subcommand not found")
	}

	// Check flags
	flags := []string{"limit", "project", "tool", "all-projects"}
	for _, f := range flags {
		if listCmd.Flags().Lookup(f) == nil {
			t.Errorf("list subcommand missing --%s flag", f)
		}
	}
}

func TestRecallShowCmd_Flags(t *testing.T) {
	cmd := Cmd()

	var showCmd *cobra.Command
	for _, sub := range cmd.Commands() {
		if sub.Name() == "show" {
			showCmd = sub
			break
		}
	}

	if showCmd == nil {
		t.Fatal("show subcommand not found")
	}

	// Check flags
	flags := []string{"latest", "full", "all-projects"}
	for _, f := range flags {
		if showCmd.Flags().Lookup(f) == nil {
			t.Errorf("show subcommand missing --%s flag", f)
		}
	}
}

func TestRecallExportCmd_Flags(t *testing.T) {
	cmd := Cmd()

	var exportCmd *cobra.Command
	for _, sub := range cmd.Commands() {
		if sub.Name() == "export" {
			exportCmd = sub
			break
		}
	}

	if exportCmd == nil {
		t.Fatal("export subcommand not found")
	}

	// Check flags (includes deprecated flags for backward compatibility).
	flags := []string{
		"all", "all-projects", "regenerate", "keep-frontmatter",
		"yes", "dry-run", "force", "skip-existing",
	}
	for _, f := range flags {
		if exportCmd.Flags().Lookup(f) == nil {
			t.Errorf("export subcommand missing --%s flag", f)
		}
	}
}

func TestExtractFrontmatter(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "with frontmatter",
			content: "---\ntitle: Test\ntopics:\n  - go\n---\n\n# Heading\n",
			want:    "---\ntitle: Test\ntopics:\n  - go\n---\n",
		},
		{
			name:    "no frontmatter",
			content: "# Just a heading\n\nSome content.\n",
			want:    "",
		},
		{
			name:    "unclosed frontmatter",
			content: "---\ntitle: Test\nno closing delimiter\n",
			want:    "",
		},
		{
			name:    "empty frontmatter",
			content: "---\n---\n\n# Heading\n",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractFrontmatter(tt.content)
			if got != tt.want {
				t.Errorf("extractFrontmatter() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatJournalFilename(t *testing.T) {
	session := &parser.Session{
		ID:        "abc12345-6789-0123-4567-890123456789",
		Slug:      "gleaming-wobbling-sutherland",
		StartTime: time.Date(2026, 1, 21, 14, 30, 0, 0, time.UTC),
	}

	filename := formatJournalFilename(session, "")

	// Should contain slug
	if !strings.Contains(filename, "gleaming-wobbling-sutherland") {
		t.Errorf("filename missing slug: %q", filename)
	}

	// Should contain short ID (first 8 chars)
	if !strings.Contains(filename, "abc12345") {
		t.Errorf("filename missing short ID: %q", filename)
	}

	// Should end with .md
	if !strings.HasSuffix(filename, ".md") {
		t.Errorf("filename missing .md extension: %q", filename)
	}
}

func TestIsEmptyMessage(t *testing.T) {
	tests := []struct {
		name string
		msg  parser.Message
		want bool
	}{
		{
			name: "empty message",
			msg:  parser.Message{},
			want: true,
		},
		{
			name: "message with text",
			msg:  parser.Message{Text: "Hello"},
			want: false,
		},
		{
			name: "message with tool uses",
			msg: parser.Message{
				ToolUses: []parser.ToolUse{{Name: "Bash"}},
			},
			want: false,
		},
		{
			name: "message with tool results",
			msg: parser.Message{
				ToolResults: []parser.ToolResult{{Content: "output"}},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := emptyMessage(tt.msg)
			if got != tt.want {
				t.Errorf("emptyMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
