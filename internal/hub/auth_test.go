//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"strings"
	"testing"

	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
)

func TestGenerateAdminToken(t *testing.T) {
	tok, err := GenerateAdminToken()
	if err != nil {
		t.Fatalf("GenerateAdminToken: %v", err)
	}
	if !strings.HasPrefix(tok, cfgHub.AdminTokenPrefix) {
		t.Errorf("admin token should start with %q, got %q", cfgHub.AdminTokenPrefix, tok)
	}
	// prefix (8) + 64 hex chars = 72
	if len(tok) != 72 {
		t.Errorf("expected length 72, got %d", len(tok))
	}
}

func TestGenerateClientToken(t *testing.T) {
	tok, err := GenerateClientToken()
	if err != nil {
		t.Fatalf("GenerateClientToken: %v", err)
	}
	if !strings.HasPrefix(tok, cfgHub.ClientTokenPrefix) {
		t.Errorf("client token should start with %q, got %q", cfgHub.ClientTokenPrefix, tok)
	}
	// prefix (8) + 64 hex chars = 72
	if len(tok) != 72 {
		t.Errorf("expected length 72, got %d", len(tok))
	}
}

func TestTokenUniqueness(t *testing.T) {
	tok1, _ := GenerateAdminToken()
	tok2, _ := GenerateAdminToken()
	if tok1 == tok2 {
		t.Error("two generated tokens should not be identical")
	}
}
