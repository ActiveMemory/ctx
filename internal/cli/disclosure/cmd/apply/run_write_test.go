//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package apply_test

import (
	"os"
	"testing"

	errDisc "github.com/ActiveMemory/ctx/internal/err/disclosure"
)

// T12: a malformed root (two ## Themes) is refused and left
// byte-identical — the CLI never half-writes a file it could not validate.
func TestRun_MalformedRootUntouched(t *testing.T) {
	dir := t.TempDir()
	const bad = "# Learnings\n\n<!-- guide -->\n\n## Themes\n\n## Themes\n"
	rootPath := writeFile(t, dir, "LEARNINGS.md", bad)
	planPath := writePlan(t, dir, movePlan())

	before, _ := os.ReadFile(rootPath)
	_, err := runApply(t, rootPath, planPath, false)
	if err != errDisc.ErrMultipleThemes {
		t.Errorf("err = %v, want ErrMultipleThemes", err)
	}
	after, _ := os.ReadFile(rootPath)
	if string(before) != string(after) {
		t.Errorf("root mutated on refusal:\nbefore=%q\nafter=%q", before, after)
	}
}
