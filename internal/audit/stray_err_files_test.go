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

// TestNoStrayErrFiles ensures err.go files only exist under
// internal/err/. An err.go anywhere else indicates error construction
// that should be consolidated into the centralized error package.
//
// See specs/ast-audit-tests.md for rationale.
func TestNoStrayErrFiles(t *testing.T) {
	var violations []string

	walkErr := filepath.WalkDir("../", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() != "err.go" {
			return nil
		}

		abs, absErr := filepath.Abs(path)
		if absErr != nil {
			return absErr
		}

		// Allow files inside internal/err/.
		if strings.Contains(abs, filepath.Join("internal", "err")+string(filepath.Separator)) {
			return nil
		}

		violations = append(violations, abs+": err.go must be in internal/err/")

		return nil
	})
	if walkErr != nil {
		t.Fatalf("filepath.WalkDir: %v", walkErr)
	}

	for _, v := range violations {
		t.Error(v)
	}
}
