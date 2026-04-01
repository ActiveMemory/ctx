//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package reduce

import (
	"strings"
	"testing"
)

func TestStripFences(t *testing.T) {
	content := "Some text\n```go\nfunc main() {}\n```\nMore text"
	got := StripFences(content, false)

	if strings.Contains(got, "```") {
		t.Error("fence markers should be removed")
	}
	if !strings.Contains(got, "func main() {}") {
		t.Error("inner content should be preserved")
	}
}

func TestStripFences_PreservesFrontmatter(t *testing.T) {
	content := "---\ntitle: test\n---\nSome ```fenced``` text"
	got := StripFences(content, false)

	if !strings.HasPrefix(got, "---\ntitle: test\n---\n") {
		t.Error("frontmatter should be preserved")
	}
}

func TestStripFences_SkipsFencesVerified(t *testing.T) {
	content := "```go\ncode\n```"
	got := StripFences(content, true)

	if got != content {
		t.Error("fencesVerified=true should return content unchanged")
	}
}

func TestStripSystemReminders(t *testing.T) {
	content := "Before\n<system-reminder>\nSecret stuff\n</system-reminder>\nAfter"
	got := StripSystemReminders(content)

	if strings.Contains(got, "Secret stuff") {
		t.Error("system reminder content should be removed")
	}
	if !strings.Contains(got, "Before") || !strings.Contains(got, "After") {
		t.Error("surrounding content should be preserved")
	}
}

func TestStripSystemReminders_BoldStyle(t *testing.T) {
	content := "Before\n" +
		"**System Reminder**: Some reminder text\n" +
		"Continued on next line\n\nAfter"
	got := StripSystemReminders(content)

	if strings.Contains(got, "System Reminder") {
		t.Error("bold system reminder should be removed")
	}
	if !strings.Contains(got, "Before") || !strings.Contains(got, "After") {
		t.Error("surrounding content should be preserved")
	}
}

func TestStripSystemReminders_CompactionSummary(t *testing.T) {
	content := "Before\n<summary>\n" +
		"Compaction summary content\n" +
		"More summary\n</summary>\nAfter"
	got := StripSystemReminders(content)

	if strings.Contains(got, "Compaction summary content") {
		t.Error("compaction summary should be removed")
	}
	if !strings.Contains(got, "Before") || !strings.Contains(got, "After") {
		t.Error("surrounding content should be preserved")
	}
}

func TestStripSystemReminders_SingleLineSummaryPreserved(t *testing.T) {
	content := "Before\n<summary>5 lines</summary>\nAfter"
	got := StripSystemReminders(content)

	if !strings.Contains(got, "<summary>5 lines</summary>") {
		t.Error("single-line summary should be preserved")
	}
}

func TestCleanToolOutputJSON(t *testing.T) {
	content := "### 1. Tool Output (10:00:00)\n\n" +
		"[{\"type\":\"text\",\"text\":\"hello world\"}]\n\n" +
		"### 2. Assistant (10:00:01)\n\nhi"
	got := CleanToolOutputJSON(content)

	if strings.Contains(got, `"type"`) {
		t.Error("JSON should be replaced with plain text")
	}
	if !strings.Contains(got, "hello world") {
		t.Error("extracted text should be present")
	}
	if !strings.Contains(got, "### 2. Assistant") {
		t.Error("next turn should survive")
	}
}

func TestCleanToolOutputJSON_NonJSON(t *testing.T) {
	content := "### 1. Tool Output (10:00:00)\n\n" +
		"just plain text\n\n" +
		"### 2. Assistant (10:00:01)\n\nhi"
	got := CleanToolOutputJSON(content)

	if !strings.Contains(got, "just plain text") {
		t.Error("non-JSON content should be preserved")
	}
}
