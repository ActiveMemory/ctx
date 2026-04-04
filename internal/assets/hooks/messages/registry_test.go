//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package messages

import (
	"testing"

	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
)

func TestRegistryCount(t *testing.T) {
	entries := Registry()
	if registryErr != nil {
		t.Fatalf("Registry() parse error: %v", registryErr)
	}
	if len(entries) != 32 {
		t.Errorf("Registry() returned %d entries, want 32", len(entries))
	}
}

func TestRegistryYAMLParses(t *testing.T) {
	if parseErr := RegistryError(); parseErr != nil {
		t.Fatalf("RegistryError() = %v, want nil", parseErr)
	}

	for i, entry := range Registry() {
		if entry.Hook == "" {
			t.Errorf("entry %d: empty hook", i)
		}
		if entry.Variant == "" {
			t.Errorf("entry %d: empty variant", i)
		}
		validCategory := entry.Category == cfgHook.CategoryCustomizable ||
			entry.Category == cfgHook.CategoryCtxSpecific
		if !validCategory {
			t.Errorf("entry %d (%s/%s): invalid category %q",
				i, entry.Hook, entry.Variant, entry.Category)
		}
		if entry.Description == "" {
			t.Errorf("entry %d (%s/%s): empty description",
				i, entry.Hook, entry.Variant)
		}
	}
}

func TestLookupKnownEntry(t *testing.T) {
	info := Lookup("check-persistence", "nudge")
	if info == nil {
		t.Fatal("Lookup(check-persistence, nudge) = nil, want non-nil")
	}
	if info.Category != cfgHook.CategoryCustomizable {
		t.Errorf("category = %q, want %q", info.Category, cfgHook.CategoryCustomizable)
	}
	if info.Description != "Context persistence nudge" {
		t.Errorf("description = %q, want %q",
			info.Description,
			"Context persistence nudge")
	}
	if len(info.TemplateVars) != 1 || info.TemplateVars[0] != "PromptsSinceNudge" {
		t.Errorf("vars = %v, want [PromptsSinceNudge]", info.TemplateVars)
	}
}

func TestLookupUnknown(t *testing.T) {
	info := Lookup("nonexistent-hook", "nonexistent-variant")
	if info != nil {
		t.Errorf("Lookup(nonexistent) = %+v, want nil", info)
	}
}

func TestHooksReturnsUniqueNames(t *testing.T) {
	hooks := hooks()
	if len(hooks) == 0 {
		t.Fatal("Hooks() returned empty list")
	}
	seen := make(map[string]bool)
	for _, h := range hooks {
		if seen[h] {
			t.Errorf("Hooks() returned duplicate: %q", h)
		}
		seen[h] = true
	}
}

func TestVariantsKnownHook(t *testing.T) {
	variants := Variants("check-context-size")
	if len(variants) != 4 {
		t.Errorf("Variants(check-context-size) = %d entries, want 4", len(variants))
	}
}

func TestVariantsUnknownHook(t *testing.T) {
	variants := Variants("nonexistent")
	if variants != nil {
		t.Errorf("Variants(nonexistent) = %v, want nil", variants)
	}
}
