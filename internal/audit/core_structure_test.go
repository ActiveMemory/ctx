//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0
//
// ================================================================
// STOP — Read internal/audit/README.md before editing this file.
//
// These tests enforce project conventions. The codebase is clean:
// all checks pass with zero violations, zero exceptions.
//
// If a test fails after your change, fix the code under test.
// Do NOT add allowlist entries, bump grandfathered counters, or
// weaken checks. Exceptions require a dedicated PR with
// justification for every entry. See README.md for the full policy.
// ================================================================

package audit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// allowedCoreFiles lists the files permitted directly
// in a core/ directory (not in subdirectories).
var allowedCoreFiles = map[string]bool{
	"doc.go": true,
}

// TestCoreStructure ensures core/ directories contain
// only doc.go and test files at the top level. All
// domain logic must live in subpackages (e.g.
// core/budget/, core/score/). This prevents core/
// from becoming a god package.
//
// Test files are exempt.
func TestCoreStructure(t *testing.T) {
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

			rel, relErr := filepath.Rel(cliRoot, path)
			if relErr != nil {
				return relErr
			}

			// Only check files directly in core/
			// directories, not in core/subpkg/.
			dir := filepath.Dir(rel)
			if !strings.HasSuffix(dir, "/core") &&
				dir != "core" {
				return nil
			}

			if allowedCoreFiles[d.Name()] {
				return nil
			}

			abs, absErr := filepath.Abs(path)
			if absErr != nil {
				return absErr
			}

			violations = append(violations,
				abs+": "+d.Name()+
					" in core/ (move to subpackage)",
			)

			return nil
		},
	)
	if walkErr != nil {
		t.Fatalf("walk: %v", walkErr)
	}

	if len(violations) > 0 {
		t.Errorf(
			"%d core/ structure violations:",
			len(violations),
		)
	}
	limit := 30
	if len(violations) < limit {
		limit = len(violations)
	}
	for _, v := range violations[:limit] {
		t.Error(v)
	}
	if len(violations) > 30 {
		t.Errorf(
			"... and %d more",
			len(violations)-30,
		)
	}
}
