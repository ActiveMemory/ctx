//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package timeline_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	"github.com/ActiveMemory/ctx/internal/write/kb/timeline"
)

func TestAppend_CreatesHeaderOnFirstWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, cfgKB.Timeline)

	row := timeline.Row{
		Date:          "2026-05-16",
		Event:         "Widget contract v2 published.",
		SourceEV:      []string{"EV-042", "EV-051"},
		RelatedTopics: []string{"widgets", "widget-composition"},
	}
	if err := timeline.Append(path, row); err != nil {
		t.Fatalf("Append: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	text := string(got)
	if !strings.Contains(text, "| Date | Event |") {
		t.Errorf("header missing: %q", text)
	}
	if !strings.Contains(text, "2026-05-16") {
		t.Errorf("date missing: %q", text)
	}
	if !strings.Contains(text, "EV-042, EV-051") {
		t.Errorf("EV list not joined: %q", text)
	}
	if !strings.Contains(text, "widgets, widget-composition") {
		t.Errorf("topics not joined: %q", text)
	}
}

func TestAppend_SecondCallDoesNotRewriteHeader(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, cfgKB.Timeline)

	for _, date := range []string{"2026-05-15", "2026-05-16"} {
		row := timeline.Row{
			Date: date, Event: "x", SourceEV: []string{"EV-001"},
		}
		if err := timeline.Append(path, row); err != nil {
			t.Fatalf("Append: %v", err)
		}
	}
	got, _ := os.ReadFile(path)
	if c := strings.Count(string(got), "| Date | Event |"); c != 1 {
		t.Errorf("header written %d times; want 1", c)
	}
}
