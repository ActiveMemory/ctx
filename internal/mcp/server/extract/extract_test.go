//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package extract

import (
	"os"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
	"github.com/ActiveMemory/ctx/internal/config/mcp/cfg"
)

func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}

func TestEntryArgsValid(t *testing.T) {
	args := map[string]interface{}{
		"type":    "decision",
		"content": "Use Go",
	}
	typ, content, err := EntryArgs(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if typ != "decision" {
		t.Errorf("type = %q, want decision", typ)
	}
	if content != "Use Go" {
		t.Errorf("content = %q, want Use Go", content)
	}
}

func TestEntryArgsMissingType(t *testing.T) {
	args := map[string]interface{}{"content": "ok"}
	_, _, err := EntryArgs(args)
	if err == nil {
		t.Fatal("expected error for missing type")
	}
}

func TestEntryArgsMissingContent(t *testing.T) {
	args := map[string]interface{}{"type": "decision"}
	_, _, err := EntryArgs(args)
	if err == nil {
		t.Fatal("expected error for missing content")
	}
}

func TestEntryArgsTooLong(t *testing.T) {
	args := map[string]interface{}{
		"type":    "decision",
		"content": strings.Repeat("x", cfg.MaxContentLen+1),
	}
	_, _, err := EntryArgs(args)
	if err == nil {
		t.Fatal("expected error for content too long")
	}
}

func TestOptsAllFields(t *testing.T) {
	args := map[string]interface{}{
		"priority":    "high",
		"context":     "ctx",
		"rationale":   "because",
		"consequence": "result",
		"lesson":      "learned",
		"application": "apply",
	}
	opts := Opts(args)
	if opts.Priority != "high" {
		t.Errorf("priority = %q", opts.Priority)
	}
	if opts.Context != "ctx" {
		t.Errorf("context = %q", opts.Context)
	}
	if opts.Rationale != "because" {
		t.Errorf("rationale = %q", opts.Rationale)
	}
	if opts.Consequence != "result" {
		t.Errorf("consequence = %q", opts.Consequence)
	}
	if opts.Lesson != "learned" {
		t.Errorf("lesson = %q", opts.Lesson)
	}
	if opts.Application != "apply" {
		t.Errorf("application = %q", opts.Application)
	}
}

func TestOptsEmpty(t *testing.T) {
	opts := Opts(map[string]interface{}{})
	if opts.Priority != "" {
		t.Error("expected empty priority")
	}
}

func TestSanitizedOpts(t *testing.T) {
	args := map[string]interface{}{
		"context":   "safe text",
		"rationale": "good reason",
	}
	opts, err := SanitizedOpts(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Context != "safe text" {
		t.Errorf("context = %q", opts.Context)
	}
}

func TestSanitizedOptsEmpty(t *testing.T) {
	opts, err := SanitizedOpts(map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Context != "" || opts.Rationale != "" {
		t.Errorf("expected empty opts, got %+v", opts)
	}
}

func TestSanitizedOptsContextTooLong(t *testing.T) {
	args := map[string]interface{}{
		"context": strings.Repeat("x", cfg.MaxOptsFieldLen+1),
	}
	_, err := SanitizedOpts(args)
	if err == nil {
		t.Fatal("expected InputTooLong error for context")
	}
	if !strings.Contains(err.Error(), "context") {
		t.Errorf("error should name field 'context', got %q", err.Error())
	}
}

func TestSanitizedOptsRationaleTooLong(t *testing.T) {
	args := map[string]interface{}{
		"rationale": strings.Repeat("y", cfg.MaxOptsFieldLen+1),
	}
	_, err := SanitizedOpts(args)
	if err == nil {
		t.Fatal("expected InputTooLong error for rationale")
	}
	if !strings.Contains(err.Error(), "rationale") {
		t.Errorf("error should name field 'rationale', got %q", err.Error())
	}
}

func TestSanitizedOptsConsequenceTooLong(t *testing.T) {
	args := map[string]interface{}{
		"consequence": strings.Repeat("z", cfg.MaxOptsFieldLen+1),
	}
	_, err := SanitizedOpts(args)
	if err == nil {
		t.Fatal("expected InputTooLong error for consequence")
	}
}

func TestSanitizedOptsLessonTooLong(t *testing.T) {
	args := map[string]interface{}{
		"lesson": strings.Repeat("a", cfg.MaxOptsFieldLen+1),
	}
	_, err := SanitizedOpts(args)
	if err == nil {
		t.Fatal("expected InputTooLong error for lesson")
	}
}

func TestSanitizedOptsApplicationTooLong(t *testing.T) {
	args := map[string]interface{}{
		"application": strings.Repeat("b", cfg.MaxOptsFieldLen+1),
	}
	_, err := SanitizedOpts(args)
	if err == nil {
		t.Fatal("expected InputTooLong error for application")
	}
}

func TestSanitizedOptsAtBoundary(t *testing.T) {
	// Exactly MaxOptsFieldLen bytes — must be accepted.
	args := map[string]interface{}{
		"rationale": strings.Repeat("ok", cfg.MaxOptsFieldLen/2),
	}
	_, err := SanitizedOpts(args)
	if err != nil {
		t.Errorf("boundary value rejected: %v", err)
	}
}
