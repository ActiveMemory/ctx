//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sourcemap_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	"github.com/ActiveMemory/ctx/internal/write/kb/sourcemap"
)

func TestAppend_CreatesHeaderOnFirstWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, cfgKB.SourceMap)

	row := sourcemap.Row{
		ShortName:          "WIDGET-SPEC",
		Kind:               "url",
		Locator:            "https://example.org/widget/v2/spec",
		AdmissionStatus:    "admitted",
		AdmissionRationale: "Canonical upstream spec; in scope.",
		Dated:              "2026-05-16",
	}
	if err := sourcemap.Append(path, row); err != nil {
		t.Fatalf("Append: %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	text := string(got)
	if !strings.Contains(text, "| Short name | Kind |") {
		t.Errorf("header missing: %q", text)
	}
	if !strings.Contains(text, "WIDGET-SPEC") {
		t.Errorf("short-name missing: %q", text)
	}
	if !strings.Contains(text, "admitted") {
		t.Errorf("admission missing: %q", text)
	}
}

func TestAppend_DatedOmittedRendersEmptyCell(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, cfgKB.SourceMap)
	row := sourcemap.Row{
		ShortName: "X", Kind: "doc", Locator: "/path",
		AdmissionStatus: "admitted", AdmissionRationale: "ok",
	}
	if err := sourcemap.Append(path, row); err != nil {
		t.Fatalf("Append: %v", err)
	}
	got, _ := os.ReadFile(path)
	// The row should still be a well-formed table line.
	rows := strings.Split(strings.TrimRight(string(got), "\n"), "\n")
	if len(rows) != 3 {
		t.Errorf("expected 3 lines; got %d: %q", len(rows), string(got))
	}
	// Last cell empty: row ends with "|  |".
	if !strings.HasSuffix(rows[2], "|  |") {
		t.Errorf("expected trailing empty Dated cell; got %q", rows[2])
	}
}
