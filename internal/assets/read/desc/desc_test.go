//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package desc

import (
	"os"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
)

func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}

func TestTextDescKeysResolve(t *testing.T) {
	// Verify every TextDescKey constant resolves to a non-empty string.
	// This catches typos in constants or missing YAML entries.
	keys := collectTextDescKeys(t)
	if len(keys) == 0 {
		t.Fatal("no TextDescKey constants found")
	}

	for _, key := range keys {
		val := TextDesc(key)
		if val == "" {
			t.Errorf("TextDesc(%q) returned empty string — missing YAML entry?", key)
		}
	}
	t.Logf("verified %d TextDescKey constants", len(keys))
}

// collectTextDescKeys extracts all TextDescKey constant values from the
// text package by parsing lines matching the pattern: TextDescKey... = "..."
func collectTextDescKeys(t *testing.T) []string {
	t.Helper()
	data, err := os.ReadFile("../../../config/embed/text/text.go")
	if err != nil {
		t.Fatalf("read text.go: %v", err)
	}

	var keys []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "TextDescKey") {
			continue
		}
		idx := strings.Index(line, "\"")
		if idx < 0 {
			continue
		}
		end := strings.LastIndex(line, "\"")
		if end <= idx {
			continue
		}
		keys = append(keys, line[idx+1:end])
	}
	return keys
}
