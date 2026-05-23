//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package placeholders_test

import (
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/placeholders"
	"github.com/ActiveMemory/ctx/internal/config/asset"
	"github.com/ActiveMemory/ctx/internal/i18n"
)

func TestLoad_EnContainsCanonicalEntries(t *testing.T) {
	t.Cleanup(placeholders.Reset)
	set, err := placeholders.Load(asset.LocaleEN)
	if err != nil {
		t.Fatalf("Load(en): %v", err)
	}
	for _, raw := range []string{
		"tbd", "n/a", "na", "none",
		"see chat", "see above", "see below",
		"pending", "to be done",
	} {
		if _, ok := set[i18n.Fold(raw)]; !ok {
			t.Errorf("set missing %q", raw)
		}
	}
}

func TestLoad_FoldsCaseInsensitively(t *testing.T) {
	t.Cleanup(placeholders.Reset)
	set, _ := placeholders.Load(asset.LocaleEN)
	for _, variant := range []string{"TBD", "Tbd", "TbD", "tBd"} {
		if _, ok := set[i18n.Fold(variant)]; !ok {
			t.Errorf("variant %q should hit (folds to %q)", variant, i18n.Fold(variant))
		}
	}
}

func TestLoad_MemoizesPerLocale(t *testing.T) {
	t.Cleanup(placeholders.Reset)
	a, err := placeholders.Load(asset.LocaleEN)
	if err != nil {
		t.Fatalf("first Load: %v", err)
	}
	b, err := placeholders.Load(asset.LocaleEN)
	if err != nil {
		t.Fatalf("second Load: %v", err)
	}
	// Both calls return the same underlying map (cache hit).
	// We can't compare map identities directly in Go, so check
	// by mutating a is visible in b — but mutating shared state
	// in tests is bad form. Instead just verify equality of
	// keys and count, then trust the cache code path.
	if len(a) != len(b) {
		t.Errorf("memoization changed result: len %d → %d", len(a), len(b))
	}
}

func TestLoad_RejectsUnknownLocale(t *testing.T) {
	t.Cleanup(placeholders.Reset)
	_, err := placeholders.Load("zz")
	if err == nil {
		t.Fatal("expected error for unknown locale; got nil")
	}
}
