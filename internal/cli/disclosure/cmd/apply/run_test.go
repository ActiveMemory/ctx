//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package apply_test

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
	"github.com/ActiveMemory/ctx/internal/cli/disclosure/cmd/apply"
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

// writeFile writes content to <dir>/<name> and returns the path.
func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return p
}

// writePlan marshals a plan to <dir>/plan.json and returns the path.
func writePlan(t *testing.T, dir string, plan disclosure.Plan) string {
	t.Helper()
	b, err := json.Marshal(plan)
	if err != nil {
		t.Fatal(err)
	}
	return writeFile(t, dir, "plan.json", string(b))
}

// runApply invokes apply.Run with buffered I/O and returns stdout + err.
func runApply(t *testing.T, path, planPath string, jsonOut bool) (string, error) {
	t.Helper()
	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	err := apply.Run(cmd, path, planPath, jsonOut)
	return out.String(), err
}

func movePlan() disclosure.Plan {
	return disclosure.Plan{
		Kind: "learning",
		Assignments: []disclosure.Assignment{{
			Theme: "context", Slug: "context", Gist: "context stuff",
			Entries: []disclosure.StagedEntry{
				{Timestamp: "2026-07-15-120000", Title: "a staged entry"},
			},
		}},
	}
}

// T10: apply moves the staged entry into its theme file, rewrites the
// root, and reports the result.
func TestRun_Apply(t *testing.T) {
	dir := t.TempDir()
	rootPath := writeFile(t, dir, "LEARNINGS.md", fixtureRoot)
	planPath := writePlan(t, dir, movePlan())

	out, err := runApply(t, rootPath, planPath, false)
	if err != nil {
		t.Fatalf("Run error = %v", err)
	}
	if !strings.Contains(out, "Moved 1 entries into 1 themes: context") {
		t.Errorf("summary missing; got: %q", out)
	}

	root, _ := os.ReadFile(rootPath)
	if strings.Contains(string(root), "## [2026-07-15-120000]") {
		t.Errorf("staged entry still in root:\n%s", root)
	}
	body, err := os.ReadFile(filepath.Join(dir, "learnings", "context.md"))
	if err != nil || !strings.Contains(string(body), "a staged entry") {
		t.Errorf("context.md = %q, err = %v", body, err)
	}
}

// T10: --json emits an ApplyResult that decodes to the expected values.
func TestRun_ApplyJSON(t *testing.T) {
	dir := t.TempDir()
	rootPath := writeFile(t, dir, "LEARNINGS.md", fixtureRoot)
	planPath := writePlan(t, dir, movePlan())

	out, err := runApply(t, rootPath, planPath, true)
	if err != nil {
		t.Fatalf("Run error = %v", err)
	}
	var res disclosure.ApplyResult
	if uErr := json.Unmarshal([]byte(out), &res); uErr != nil {
		t.Fatalf("output not JSON: %v\n%s", uErr, out)
	}
	if res.Moved != 1 || len(res.Themes) != 1 || res.Themes[0] != "context" {
		t.Errorf("ApplyResult = %+v", res)
	}
}

// T11: a non-knowledge file is rejected with the typed sentinel, before
// the plan is even read.
func TestRun_RejectsNonKnowledgeFile(t *testing.T) {
	dir := t.TempDir()
	path := writeFile(t, dir, "README.md", "# readme\n")
	planPath := writePlan(t, dir, movePlan())

	_, err := runApply(t, path, planPath, false)
	if !errors.Is(err, errDisc.ErrNotAKnowledgeFile) {
		t.Errorf("Run(README.md) err = %v, want ErrNotAKnowledgeFile", err)
	}
}
