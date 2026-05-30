//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package jsonpayload

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"

	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// writePayload writes content to a temp file and returns its path.
func writePayload(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "payload.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write payload: %v", err)
	}
	return path
}

// TestLoadDecodesAllFields verifies a full payload, including the
// provenance envelope, decodes into the typed fields.
func TestLoadDecodesAllFields(t *testing.T) {
	path := writePayload(t, `{
		"title": "Install into PATH",
		"context": "the binary lives at /usr/local/bin/ctx",
		"rationale": "system PATH install is the documented route",
		"consequence": "users run ctx from anywhere",
		"provenance": {"session_id": "abc12345", "branch": "main", "commit": "deadbeef"}
	}`)

	p, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if p.Title != "Install into PATH" {
		t.Errorf("title = %q", p.Title)
	}
	if p.Context != "the binary lives at /usr/local/bin/ctx" {
		t.Errorf("context = %q", p.Context)
	}
	if p.Provenance.SessionID != "abc12345" {
		t.Errorf("session_id = %q", p.Provenance.SessionID)
	}
	if p.Provenance.Commit != "deadbeef" {
		t.Errorf("commit = %q", p.Provenance.Commit)
	}
}

// TestLoadRejectsUnknownField verifies strict decoding so a typo'd
// key is a hard error instead of a silently dropped field.
func TestLoadRejectsUnknownField(t *testing.T) {
	path := writePayload(t, `{"title": "x", "rationail": "typo"}`)
	if _, err := Load(path); err == nil {
		t.Fatal("expected error on unknown field")
	}
}

// TestLoadMissingFile surfaces a read error rather than a zero value.
func TestLoadMissingFile(t *testing.T) {
	if _, err := Load(filepath.Join(t.TempDir(), "absent.json")); err == nil {
		t.Fatal("expected error for missing file")
	}
}

// TestContentTitleOnly returns the trimmed title.
func TestContentTitleOnly(t *testing.T) {
	p := Payload{Title: "  Use Postgres  "}
	if got := p.Content(); got != "Use Postgres" {
		t.Errorf("Content() = %q, want %q", got, "Use Postgres")
	}
}

// TestContentTitleAndBody space-joins a task body onto the title.
func TestContentTitleAndBody(t *testing.T) {
	p := Payload{Title: "Add flag", Body: "for JSON ingest"}
	if got := p.Content(); got != "Add flag for JSON ingest" {
		t.Errorf("Content() = %q", got)
	}
}

// TestContentEmpty returns empty so callers fall through.
func TestContentEmpty(t *testing.T) {
	if got := (Payload{}).Content(); got != "" {
		t.Errorf("Content() = %q, want empty", got)
	}
}

// newAddCmd builds a cobra command carrying the add string flags
// that OverlayFlags writes to.
func newAddCmd() *cobra.Command {
	c := &cobra.Command{Use: "add"}
	for _, name := range []string{
		cFlag.JSONFile, cFlag.Context, cFlag.Rationale,
		cFlag.Consequence, cFlag.Lesson, cFlag.Application,
		cFlag.Priority, cFlag.Section,
		cFlag.SessionID, cFlag.Branch, cFlag.Commit,
	} {
		c.Flags().String(name, "", "")
	}
	return c
}

// TestOverlayFlagsSupersedes verifies a JSON payload overrides the
// individually-supplied flags and folds in provenance.
func TestOverlayFlagsSupersedes(t *testing.T) {
	path := writePayload(t, `{
		"context": "json context",
		"rationale": "json rationale",
		"consequence": "json consequence",
		"provenance": {"session_id": "s1", "branch": "b1", "commit": "c1"}
	}`)

	c := newAddCmd()
	if err := c.Flags().Set(cFlag.JSONFile, path); err != nil {
		t.Fatalf("set json-file: %v", err)
	}
	if err := c.Flags().Set(cFlag.Rationale, "cli rationale"); err != nil {
		t.Fatalf("set rationale: %v", err)
	}

	if err := OverlayFlags(c); err != nil {
		t.Fatalf("OverlayFlags: %v", err)
	}

	cases := map[string]string{
		cFlag.Context:     "json context",
		cFlag.Rationale:   "json rationale",
		cFlag.Consequence: "json consequence",
		cFlag.SessionID:   "s1",
		cFlag.Branch:      "b1",
		cFlag.Commit:      "c1",
	}
	for name, want := range cases {
		got, _ := c.Flags().GetString(name)
		if got != want {
			t.Errorf("flag %s = %q, want %q", name, got, want)
		}
	}
}

// TestOverlayFlagsNoFile is a no-op when --json-file is unset, leaving
// CLI-supplied values intact.
func TestOverlayFlagsNoFile(t *testing.T) {
	c := newAddCmd()
	if err := c.Flags().Set(cFlag.Rationale, "cli only"); err != nil {
		t.Fatalf("set rationale: %v", err)
	}
	if err := OverlayFlags(c); err != nil {
		t.Fatalf("OverlayFlags: %v", err)
	}
	if got, _ := c.Flags().GetString(cFlag.Rationale); got != "cli only" {
		t.Errorf("rationale = %q, want %q", got, "cli only")
	}
}

// TestOverlayFlagsEmptyPayloadKeepsFlag verifies that an absent payload
// field leaves a CLI-supplied flag untouched.
func TestOverlayFlagsEmptyPayloadKeepsFlag(t *testing.T) {
	path := writePayload(t, `{"context": "json context"}`)

	c := newAddCmd()
	if err := c.Flags().Set(cFlag.JSONFile, path); err != nil {
		t.Fatalf("set json-file: %v", err)
	}
	if err := c.Flags().Set(cFlag.Rationale, "cli rationale"); err != nil {
		t.Fatalf("set rationale: %v", err)
	}
	if err := OverlayFlags(c); err != nil {
		t.Fatalf("OverlayFlags: %v", err)
	}
	if got, _ := c.Flags().GetString(cFlag.Rationale); got != "cli rationale" {
		t.Errorf("rationale = %q, want %q (unchanged)", got, "cli rationale")
	}
}
