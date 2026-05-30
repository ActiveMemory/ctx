//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package assets

import (
	"path"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/config/asset"
)

func TestHookMessageRegistry(t *testing.T) {
	data, readErr := FS.ReadFile(asset.PathMessageRegistry)
	if readErr != nil {
		t.Fatalf("unexpected error: %v", readErr)
	}
	if len(data) == 0 {
		t.Fatal("returned empty data")
	}

	var entries []map[string]any
	if parseErr := yaml.Unmarshal(data, &entries); parseErr != nil {
		t.Fatalf("invalid YAML: %v", parseErr)
	}
	for i, entry := range entries {
		if _, ok := entry["hook"]; !ok {
			t.Errorf("entry %d missing 'hook' key", i)
		}
		if _, ok := entry["variant"]; !ok {
			t.Errorf("entry %d missing 'variant' key", i)
		}
	}
}

func TestListHookMessages(t *testing.T) {
	entries, listErr := FS.ReadDir(asset.DirHooksMessages)
	if listErr != nil {
		t.Fatalf("unexpected error: %v", listErr)
	}
	if len(entries) == 0 {
		t.Fatal("returned empty list")
	}

	hookSet := make(map[string]bool)
	for _, h := range entries {
		if h.IsDir() {
			hookSet[h.Name()] = true
		}
	}
	wantHooks := []string{
		"qa-reminder",
		"check-context-size",
		"block-non-path-ctx",
	}
	for _, exp := range wantHooks {
		if !hookSet[exp] {
			t.Errorf("missing expected hook: %s", exp)
		}
	}
}

func TestHookMessage_ReadVariant(t *testing.T) {
	gatePath := path.Join(
		asset.DirHooksMessages,
		"qa-reminder", "gate.txt",
	)
	content, readErr := FS.ReadFile(gatePath)
	if readErr != nil {
		t.Fatalf("unexpected error: %v", readErr)
	}
	if len(content) == 0 {
		t.Fatal("returned empty content")
	}
}
