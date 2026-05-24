//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parse_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/ActiveMemory/ctx/internal/cli/audit/core/parse"
)

func TestFrontmatter_RoundTrip(t *testing.T) {
	in := []byte(`---
kind: surface
status: findings
commit-range: main..HEAD
generated-at: 2026-05-24T14:30:12Z
generator: /ctx-surface-audit
digest: abc123
---

Body line 1
Body line 2
`)
	hdr, body, err := parse.Frontmatter(in)
	if err != nil {
		t.Fatalf("Frontmatter error: %v", err)
	}
	if hdr.Kind != "surface" {
		t.Errorf("Kind = %q, want surface", hdr.Kind)
	}
	if hdr.Status != "findings" {
		t.Errorf("Status = %q, want findings", hdr.Status)
	}
	if hdr.CommitRange != "main..HEAD" {
		t.Errorf("CommitRange = %q, want main..HEAD", hdr.CommitRange)
	}
	if hdr.Digest != "abc123" {
		t.Errorf("Digest = %q, want abc123", hdr.Digest)
	}
	wantTime := time.Date(2026, 5, 24, 14, 30, 12, 0, time.UTC)
	if !hdr.GeneratedAt.Equal(wantTime) {
		t.Errorf(
			"GeneratedAt = %v, want %v", hdr.GeneratedAt, wantTime,
		)
	}
	if !strings.Contains(body, "Body line 1") {
		t.Errorf("body missing line 1: %q", body)
	}
	if !strings.Contains(body, "Body line 2") {
		t.Errorf("body missing line 2: %q", body)
	}
}

func TestFrontmatter_NoFrontmatter(t *testing.T) {
	_, _, err := parse.Frontmatter([]byte("no header here\n"))
	if !errors.Is(err, parse.ErrNoFrontmatter) {
		t.Errorf("err = %v, want ErrNoFrontmatter", err)
	}
}

func TestFrontmatter_Unterminated(t *testing.T) {
	in := []byte(`---
kind: surface
status: findings
no closing delimiter here
`)
	_, _, err := parse.Frontmatter(in)
	if !errors.Is(err, parse.ErrUnterminatedFrontmatter) {
		t.Errorf("err = %v, want ErrUnterminatedFrontmatter", err)
	}
}

func TestFrontmatter_PreservesBodyBytes(t *testing.T) {
	body := "  • indent matters\n```\ncode block stays\n```\nend."
	in := []byte("---\nkind: x\nstatus: clean\n---\n" + body)
	_, got, err := parse.Frontmatter(in)
	if err != nil {
		t.Fatalf("Frontmatter error: %v", err)
	}
	if got != body {
		t.Errorf("body mangled\ngot:  %q\nwant: %q", got, body)
	}
}
