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

	"github.com/ActiveMemory/ctx/internal/cli/recall/core/extract"
	"github.com/ActiveMemory/ctx/internal/cli/recall/core/format"
	"github.com/ActiveMemory/ctx/internal/cli/recall/core/validate"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/entity"
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

	expectedSubs := []string{"list", "show", "import", "lock", "unlock", "sync"}
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

func TestRecallImportCmd_Flags(t *testing.T) {
	cmd := Cmd()

	var importCmd *cobra.Command
	for _, sub := range cmd.Commands() {
		if sub.Name() == "import" {
			importCmd = sub
			break
		}
	}

	if importCmd == nil {
		t.Fatal("import subcommand not found")
	}

	flags := []string{
		"all", "all-projects", "regenerate", "keep-frontmatter",
		"yes", "dry-run",
	}
	for _, f := range flags {
		if importCmd.Flags().Lookup(f) == nil {
			t.Errorf("import subcommand missing --%s flag", f)
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
			got := extract.ExtractFrontmatter(tt.content)
			if got != tt.want {
				t.Errorf("ExtractFrontmatter() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatJournalFilename(t *testing.T) {
	session := &entity.Session{
		ID:        "abc12345-6789-0123-4567-890123456789",
		Slug:      "gleaming-wobbling-sutherland",
		StartTime: time.Date(2026, 1, 21, 14, 30, 0, 0, time.UTC),
	}

	filename := format.JournalFilename(session, "")

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
		msg  entity.Message
		want bool
	}{
		{
			name: "empty message",
			msg:  entity.Message{},
			want: true,
		},
		{
			name: "message with text",
			msg:  entity.Message{Text: "Hello"},
			want: false,
		},
		{
			name: "message with tool uses",
			msg: entity.Message{
				ToolUses: []entity.ToolUse{{Name: "Bash"}},
			},
			want: false,
		},
		{
			name: "message with tool results",
			msg: entity.Message{
				ToolResults: []entity.ToolResult{{Content: "output"}},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validate.EmptyMessage(tt.msg)
			if got != tt.want {
				t.Errorf("EmptyMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
