//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package inspect_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
	"github.com/ActiveMemory/ctx/internal/cli/disclosure/cmd/inspect"
	"github.com/ActiveMemory/ctx/internal/disclosure"
	errDisc "github.com/ActiveMemory/ctx/internal/err/disclosure"
)

// TestMain initializes the asset lookup so desc.Text resolves the output
// labels and error text (see the "uninitialized desc.Text" learning).
func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}

const fixtureRoot = "# Learnings\n\n<!-- guide -->\n\n" +
	"## [2026-07-15-120000] a staged entry\n\n**Context**: x.\n\n---\n\n" +
	"## Themes\n\n- hooks — hook mechanics → [hooks](learnings/hooks.md)\n"

// writeFixture creates <dir>/<name> with content and returns its path.
func writeFixture(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}

// run invokes inspect.Run with an output buffer and returns stdout + err.
func run(t *testing.T, path string, jsonOut bool) (string, error) {
	t.Helper()
	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	err := inspect.Run(cmd, path, jsonOut)
	return out.String(), err
}

// T05: human output lists the kind, staged entries, and themes.
func TestRun_HumanOutput(t *testing.T) {
	path := writeFixture(t, "LEARNINGS.md", fixtureRoot)
	out, err := run(t, path, false)
	if err != nil {
		t.Fatalf("Run error = %v", err)
	}
	for _, want := range []string{
		"learning", "Staged (1)", "a staged entry", "Themes (1)",
		"hooks", "learnings/hooks.md",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("human output missing %q; got:\n%s", want, out)
		}
	}
}

// T06: --json output decodes into an Inspection with expected values.
func TestRun_JSONOutput(t *testing.T) {
	path := writeFixture(t, "LEARNINGS.md", fixtureRoot)
	out, err := run(t, path, true)
	if err != nil {
		t.Fatalf("Run error = %v", err)
	}
	var got disclosure.Inspection
	if uErr := json.Unmarshal([]byte(out), &got); uErr != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", uErr, out)
	}
	if got.Kind != "learning" {
		t.Errorf("Kind = %q, want learning", got.Kind)
	}
	if len(got.Staging) != 1 || got.Staging[0].Title != "a staged entry" {
		t.Errorf("Staging = %+v", got.Staging)
	}
	if len(got.Themes) != 1 || got.Themes[0].Link != "learnings/hooks.md" {
		t.Errorf("Themes = %+v", got.Themes)
	}
}

// T07: inspect writes nothing — the file is byte-identical after.
func TestRun_WritesNothing(t *testing.T) {
	path := writeFixture(t, "DECISIONS.md",
		strings.Replace(fixtureRoot, "# Learnings", "# Decisions", 1))
	before, _ := os.ReadFile(path)
	if _, err := run(t, path, false); err != nil {
		t.Fatalf("Run error = %v", err)
	}
	after, _ := os.ReadFile(path)
	if !bytes.Equal(before, after) {
		t.Errorf("inspect modified the file:\nbefore=%q\nafter=%q", before, after)
	}
}

// T08: a non-knowledge file is rejected with the typed sentinel.
func TestRun_RejectsNonKnowledgeFile(t *testing.T) {
	path := writeFixture(t, "README.md", "# readme\n")
	_, err := run(t, path, false)
	if !errors.Is(err, errDisc.ErrNotAKnowledgeFile) {
		t.Errorf("Run(README.md) err = %v, want ErrNotAKnowledgeFile", err)
	}
}
