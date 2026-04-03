//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package audit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// allowedCmdFiles lists the files permitted in a
// cmd/$sub/ directory (excluding test files).
var allowedCmdFiles = map[string]bool{
	"cmd.go": true,
	"run.go": true,
	"doc.go": true,
}

// cmdSubdirAllowlist lists known cmd/ subdirectories
// that have extra files pending migration to core/.
// Each entry should have a tracking task. Remove
// entries as files are migrated.
var cmdSubdirAllowlist = map[string]bool{
	"guide/cmd/root":         true,
	"journal/cmd/source":     true,
	"loop/cmd/root":          true,
	"pad/cmd/edit":           true,
	"pad/cmd/resolve":        true,
	"system/cmd/message":     true,
	"system/cmd/post_commit": true,
	"why/cmd/root":           true,
}

// TestCLICmdStructure enforces the cmd/ directory
// convention: each cmd/$sub/ directory should contain
// only cmd.go, run.go, doc.go, and test files. Extra
// files belong in the corresponding core/ package.
//
// Known violations are allowlisted with a tracking
// comment. The allowlist should shrink over time.
func TestCLICmdStructure(t *testing.T) {
	cliRoot := filepath.Join("..", "cli")
	var violations []string

	walkErr := filepath.WalkDir(
		cliRoot,
		func(
			path string, d os.DirEntry, err error,
		) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if !strings.HasSuffix(d.Name(), ".go") {
				return nil
			}
			if isTestFile(d.Name()) {
				return nil
			}

			// Only check files inside cmd/ directories.
			rel, relErr := filepath.Rel(cliRoot, path)
			if relErr != nil {
				return relErr
			}

			if !strings.Contains(rel, "cmd/") {
				return nil
			}

			if allowedCmdFiles[d.Name()] {
				return nil
			}

			// Check allowlist by cmd subdirectory.
			dir := filepath.Dir(rel)
			// Normalize: "journal/cmd/source" from
			// "journal/cmd/source/types.go"
			if cmdSubdirAllowlist[dir] {
				return nil
			}

			abs, absErr := filepath.Abs(path)
			if absErr != nil {
				return absErr
			}

			violations = append(violations,
				abs+": stray file "+d.Name()+
					" in cmd/ (move to core/)",
			)

			return nil
		},
	)
	if walkErr != nil {
		t.Fatalf("walk: %v", walkErr)
	}

	for _, v := range violations {
		t.Error(v)
	}
}
