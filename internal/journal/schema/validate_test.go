//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func testdataPath(name string) string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "testdata", name)
}

func TestValidateFileClean(t *testing.T) {
	s := Default()
	c := NewCollector(s.Version)

	validateErr := ValidateFile(s, testdataPath("valid.jsonl"), c)
	if validateErr != nil {
		t.Fatalf("unexpected error: %v", validateErr)
	}

	if c.Drift() {
		findings := c.SortedFindings()
		for _, f := range findings {
			t.Errorf("unexpected finding: type=%d name=%s count=%d",
				f.Type, f.Name, f.Count)
		}
	}

	if c.Meta.FilesScanned != 1 {
		t.Fatalf("expected 1 file scanned, got %d", c.Meta.FilesScanned)
	}
	// 3 lines total: user, assistant, progress.
	if c.Meta.LinesScanned != 3 {
		t.Fatalf("expected 3 lines scanned, got %d", c.Meta.LinesScanned)
	}
}

func TestValidateFileDrift(t *testing.T) {
	s := Default()
	c := NewCollector(s.Version)

	validateErr := ValidateFile(s, testdataPath("drift.jsonl"), c)
	if validateErr != nil {
		t.Fatalf("unexpected error: %v", validateErr)
	}

	if !c.Drift() {
		t.Fatal("expected drift findings")
	}

	findings := c.SortedFindings()

	// Expect: unknown field "brandNewField", unknown block "hologram",
	// unknown record type "quantum-snapshot".
	findingNames := make(map[string]bool)
	for _, f := range findings {
		key := findingKey(f.Type, f.Name)
		findingNames[key] = true
	}

	expected := []string{
		"field:brandNewField",
		"block:hologram",
		"record:quantum-snapshot",
	}
	for _, exp := range expected {
		if !findingNames[exp] {
			t.Errorf("missing expected finding: %s", exp)
		}
	}
}

func TestValidateFileMalformed(t *testing.T) {
	s := Default()
	c := NewCollector(s.Version)

	validateErr := ValidateFile(s, testdataPath("malformed.jsonl"), c)
	if validateErr != nil {
		t.Fatalf("unexpected error: %v", validateErr)
	}

	if c.Meta.MalformedLines != 1 {
		t.Fatalf("expected 1 malformed line, got %d", c.Meta.MalformedLines)
	}

	if !c.Drift() {
		t.Fatal("expected drift from malformed line")
	}
}

func TestCollectorMergesFindings(t *testing.T) {
	s := Default()
	c := NewCollector(s.Version)

	// Validate two files with drift.
	validateErr := ValidateFile(s, testdataPath("drift.jsonl"), c)
	if validateErr != nil {
		t.Fatalf("unexpected error: %v", validateErr)
	}
	validateErr = ValidateFile(s, testdataPath("malformed.jsonl"), c)
	if validateErr != nil {
		t.Fatalf("unexpected error: %v", validateErr)
	}

	if c.Meta.FilesScanned != 2 {
		t.Fatalf("expected 2 files scanned, got %d", c.Meta.FilesScanned)
	}
}

func TestValidateFileNotFound(t *testing.T) {
	s := Default()
	c := NewCollector(s.Version)

	validateErr := ValidateFile(s, "/nonexistent/file.jsonl", c)
	if validateErr == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestReport(t *testing.T) {
	s := Default()
	c := NewCollector(s.Version)

	validateErr := ValidateFile(s, testdataPath("drift.jsonl"), c)
	if validateErr != nil {
		t.Fatalf("unexpected error: %v", validateErr)
	}

	report := Report(c)
	if report == "" {
		t.Fatal("expected non-empty report")
	}

	// Report should contain key sections.
	for _, expected := range []string{
		"# Schema Drift Report",
		"Unknown Fields",
		"`brandNewField`",
		"Unknown Block Types",
		"`hologram`",
		"Unknown Record Types",
		"`quantum-snapshot`",
	} {
		if !strings.Contains(report, expected) {
			t.Errorf("report missing expected content: %s", expected)
		}
	}
}

func TestReportEmpty(t *testing.T) {
	c := NewCollector("1.0.0")
	report := Report(c)
	if report != "" {
		t.Fatal("expected empty report for no findings")
	}
}

func TestSummary(t *testing.T) {
	s := Default()
	c := NewCollector(s.Version)

	validateErr := ValidateFile(s, testdataPath("drift.jsonl"), c)
	if validateErr != nil {
		t.Fatalf("unexpected error: %v", validateErr)
	}

	summary := Summary(c)
	if summary == "" {
		t.Fatal("expected non-empty summary")
	}
	if !strings.Contains(summary, "Schema drift detected") {
		t.Error("summary should mention schema drift")
	}
	if !strings.Contains(summary, "ctx journal schema check") {
		t.Error("summary should reference the check command")
	}
}

func TestSummaryEmpty(t *testing.T) {
	c := NewCollector("1.0.0")
	summary := Summary(c)
	if summary != "" {
		t.Fatal("expected empty summary for no findings")
	}
}
