//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

import (
	"strings"
	"testing"
)

func TestParse_FullFrontmatter(t *testing.T) {
	input := `---
name: api-standards
description: REST API design conventions
inclusion: auto
tools:
  - claude
  - cursor
priority: 10
---
# API Standards
Use RESTful conventions.
`
	sf, err := Parse([]byte(input), "steering/api-standards.md")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if sf.Name != "api-standards" {
		t.Errorf("Name = %q, want %q", sf.Name, "api-standards")
	}
	if sf.Description != "REST API design conventions" {
		t.Errorf("Description = %q, want %q", sf.Description, "REST API design conventions")
	}
	if sf.Inclusion != InclusionAuto {
		t.Errorf("Inclusion = %q, want %q", sf.Inclusion, InclusionAuto)
	}
	if len(sf.Tools) != 2 || sf.Tools[0] != "claude" || sf.Tools[1] != "cursor" {
		t.Errorf("Tools = %v, want [claude cursor]", sf.Tools)
	}
	if sf.Priority != 10 {
		t.Errorf("Priority = %d, want %d", sf.Priority, 10)
	}
	if sf.Path != "steering/api-standards.md" {
		t.Errorf("Path = %q, want %q", sf.Path, "steering/api-standards.md")
	}
	if !strings.Contains(sf.Body, "# API Standards") {
		t.Errorf("Body missing expected content, got %q", sf.Body)
	}
}

func TestParse_DefaultValues(t *testing.T) {
	input := `---
name: minimal
---
Some body content.
`
	sf, err := Parse([]byte(input), "test.md")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if sf.Inclusion != InclusionManual {
		t.Errorf("Inclusion = %q, want default %q", sf.Inclusion, InclusionManual)
	}
	if sf.Tools != nil {
		t.Errorf("Tools = %v, want nil (all tools)", sf.Tools)
	}
	if sf.Priority != 50 {
		t.Errorf("Priority = %d, want default %d", sf.Priority, 50)
	}
}

func TestParse_EmptyBody(t *testing.T) {
	input := `---
name: empty-body
---
`
	sf, err := Parse([]byte(input), "test.md")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if sf.Body != "" {
		t.Errorf("Body = %q, want empty", sf.Body)
	}
}

func TestParse_InvalidYAML(t *testing.T) {
	input := `---
name: [invalid
  yaml: {broken
---
body
`
	_, err := Parse([]byte(input), "bad-file.md")
	if err == nil {
		t.Fatal("Parse() expected error for invalid YAML, got nil")
	}
	if !strings.Contains(err.Error(), "bad-file.md") {
		t.Errorf("error should identify file path, got: %v", err)
	}
	if !strings.Contains(err.Error(), "invalid YAML frontmatter") {
		t.Errorf("error should describe YAML failure, got: %v", err)
	}
}

func TestParse_MissingOpeningDelimiter(t *testing.T) {
	input := `name: no-delimiters
---
body
`
	_, err := Parse([]byte(input), "test.md")
	if err == nil {
		t.Fatal("Parse() expected error for missing opening delimiter")
	}
	if !strings.Contains(err.Error(), "missing opening frontmatter delimiter") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParse_MissingClosingDelimiter(t *testing.T) {
	input := `---
name: no-close
`
	_, err := Parse([]byte(input), "test.md")
	if err == nil {
		t.Fatal("Parse() expected error for missing closing delimiter")
	}
	if !strings.Contains(err.Error(), "missing closing frontmatter delimiter") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParse_InclusionAlways(t *testing.T) {
	input := `---
name: always-on
inclusion: always
---
Always included.
`
	sf, err := Parse([]byte(input), "test.md")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if sf.Inclusion != InclusionAlways {
		t.Errorf("Inclusion = %q, want %q", sf.Inclusion, InclusionAlways)
	}
}

func TestPrint_RoundTrip(t *testing.T) {
	input := `---
name: round-trip
description: Test round-trip
inclusion: auto
tools:
  - kiro
priority: 25
---
# Round Trip
Content here.
`
	sf1, err := Parse([]byte(input), "test.md")
	if err != nil {
		t.Fatalf("first Parse() error = %v", err)
	}

	printed := Print(sf1)

	sf2, err := Parse(printed, "test.md")
	if err != nil {
		t.Fatalf("second Parse() error = %v", err)
	}

	if sf1.Name != sf2.Name {
		t.Errorf("Name mismatch: %q vs %q", sf1.Name, sf2.Name)
	}
	if sf1.Description != sf2.Description {
		t.Errorf("Description mismatch: %q vs %q", sf1.Description, sf2.Description)
	}
	if sf1.Inclusion != sf2.Inclusion {
		t.Errorf("Inclusion mismatch: %q vs %q", sf1.Inclusion, sf2.Inclusion)
	}
	if sf1.Priority != sf2.Priority {
		t.Errorf("Priority mismatch: %d vs %d", sf1.Priority, sf2.Priority)
	}
	if len(sf1.Tools) != len(sf2.Tools) {
		t.Errorf("Tools length mismatch: %d vs %d", len(sf1.Tools), len(sf2.Tools))
	}
	if sf1.Body != sf2.Body {
		t.Errorf("Body mismatch:\n  got:  %q\n  want: %q", sf2.Body, sf1.Body)
	}
}

func TestPrint_MinimalFile(t *testing.T) {
	sf := &SteeringFile{
		Name:      "minimal",
		Inclusion: InclusionManual,
		Priority:  50,
		Body:      "Hello.\n",
	}

	out := Print(sf)
	result := string(out)

	if !strings.HasPrefix(result, "---\n") {
		t.Error("Print output should start with ---")
	}
	if !strings.Contains(result, "name: minimal") {
		t.Error("Print output should contain name field")
	}
	if !strings.HasSuffix(result, "Hello.\n") {
		t.Errorf("Print output should end with body, got %q", result)
	}
}

func TestPrint_NilToolsOmitted(t *testing.T) {
	sf := &SteeringFile{
		Name:      "no-tools",
		Inclusion: InclusionManual,
		Priority:  50,
	}

	out := string(Print(sf))
	if strings.Contains(out, "tools:") {
		t.Error("Print should omit tools field when nil")
	}
}
