//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package turn

import (
	"strings"
	"testing"
)

func TestMergeConsecutiveTurns(t *testing.T) {
	content := strings.Join([]string{
		"### 1. Assistant (10:00:00)",
		"",
		"First assistant text",
		"",
		"### 2. Assistant (10:00:01)",
		"",
		"Second assistant text",
		"",
		"### 3. User (10:00:02)",
		"",
		"User text",
	}, "\n")

	got := MergeConsecutive(content)

	if strings.Contains(got, "### 2. Assistant") {
		t.Error("consecutive same-role header should be merged")
	}
	if !strings.Contains(got, "### 1. Assistant") {
		t.Error("first header should be kept")
	}
	if !strings.Contains(got, "First assistant text") {
		t.Error("first body should be preserved")
	}
	if !strings.Contains(got, "Second assistant text") {
		t.Error("second body should be preserved")
	}
	if !strings.Contains(got, "### 3. User") {
		t.Error("different role header should be preserved")
	}
}

func TestMergeConsecutiveTurns_DifferentRoles(t *testing.T) {
	content := strings.Join([]string{
		"### 1. Assistant (10:00:00)",
		"",
		"Assistant text",
		"",
		"### 2. User (10:00:01)",
		"",
		"User text",
		"",
		"### 3. Assistant (10:00:02)",
		"",
		"Another assistant text",
	}, "\n")

	got := MergeConsecutive(content)

	if !strings.Contains(got, "### 1. Assistant") {
		t.Error("first assistant header missing")
	}
	if !strings.Contains(got, "### 2. User") {
		t.Error("user header missing")
	}
	if !strings.Contains(got, "### 3. Assistant") {
		t.Error("second assistant header missing")
	}
}
