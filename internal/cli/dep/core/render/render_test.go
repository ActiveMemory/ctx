//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package render

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMermaidID(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"internal/cli/dep", "internal_cli_dep"},
		{"github.com/foo/bar", "github_com_foo_bar"},
		{"my-pkg", "my_pkg"},
	}
	for _, tt := range tests {
		if got := MermaidID(tt.input); got != tt.want {
			t.Errorf(
				"MermaidID(%q) = %q, want %q",
				tt.input, got, tt.want,
			)
		}
	}
}

func TestMermaid(t *testing.T) {
	graph := map[string][]string{
		"cmd":          {"internal/cli"},
		"internal/cli": {"internal/config"},
	}

	out := Mermaid(graph)
	if !strings.HasPrefix(out, "graph TD\n") {
		t.Errorf(
			"Mermaid should start with 'graph TD\\n', got: %s",
			out,
		)
	}
	if !strings.Contains(
		out,
		`cmd["cmd"] --> internal_cli["internal/cli"]`,
	) {
		t.Errorf(
			"Mermaid missing expected edge, got: %s", out,
		)
	}
}

func TestTable(t *testing.T) {
	graph := map[string][]string{
		"cmd": {"internal/cli"},
	}

	out := Table(graph)
	if !strings.Contains(out, "Package") {
		t.Errorf(
			"Table should contain header 'Package', got: %s",
			out,
		)
	}
	if !strings.Contains(out, "cmd") {
		t.Errorf(
			"Table should contain 'cmd', got: %s", out,
		)
	}
}

func TestJSON(t *testing.T) {
	graph := map[string][]string{
		"cmd": {"internal/cli"},
	}

	out := JSON(graph)
	var parsed map[string][]string
	if err := json.Unmarshal(
		[]byte(out), &parsed,
	); err != nil {
		t.Fatalf("JSON produced invalid JSON: %v", err)
	}
	if len(parsed["cmd"]) != 1 ||
		parsed["cmd"][0] != "internal/cli" {
		t.Errorf("unexpected parsed result: %v", parsed)
	}
}
