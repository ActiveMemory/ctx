//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package lookup

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	Init()
	os.Exit(m.Run())
}

func TestPermAllowListDefault(t *testing.T) {
	allow := PermAllowListDefault()
	if len(allow) == 0 {
		t.Fatal("PermAllowListDefault() returned empty list")
	}

	allowSet := make(map[string]bool)
	for _, p := range allow {
		allowSet[p] = true
	}
	if !allowSet["Bash(ctx:*)"] {
		t.Error("allow list missing: Bash(ctx:*)")
	}
}

func TestPermDenyListDefault(t *testing.T) {
	deny := PermDenyListDefault()
	if len(deny) == 0 {
		t.Fatal("PermDenyListDefault() returned empty list")
	}

	denySet := make(map[string]bool)
	for _, p := range deny {
		denySet[p] = true
	}

	expected := []string{
		"Bash(sudo *)",
		"Bash(git push *)",
		"Bash(rm -rf /*)",
		"Bash(curl *)",
		"Read(**/.env)",
		"Edit(**/.env)",
	}
	for _, e := range expected {
		if !denySet[e] {
			t.Errorf("deny list missing: %s", e)
		}
	}
}

func TestStopWords(t *testing.T) {
	sw := StopWords()
	if len(sw) == 0 {
		t.Fatal("StopWords() returned empty map")
	}
	if !sw["the"] {
		t.Error("StopWords() missing common word 'the'")
	}
}
