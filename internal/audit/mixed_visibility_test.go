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
	"go/ast"
	"testing"
	"unicode"
)

// TestNoMixedVisibility ensures that files containing
// exported functions do not also contain unexported
// functions. Private helpers belong in their own file
// to keep public API files focused and short.
//
// Exempt: doc.go files, test files, and files with
// only one function total (too small to split).
func TestNoMixedVisibility(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			fpath := pkg.Fset.Position(
				file.Pos(),
			).Filename
			if isTestFile(fpath) {
				continue
			}

			var exported []string
			var unexported []string

			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}

				name := fn.Name.Name
				if unicode.IsUpper(rune(name[0])) {
					exported = append(
						exported, name,
					)
				} else {
					unexported = append(
						unexported, name,
					)
				}
			}

			// Skip if only one function total.
			total := len(exported) + len(unexported)
			if total <= 1 {
				continue
			}

			// Flag files with both exported and
			// unexported functions.
			if len(exported) > 0 &&
				len(unexported) > 0 {
				for _, name := range unexported {
					violations = append(
						violations,
						fpath+": unexported "+
							name+"() in file "+
							"with exported funcs",
					)
				}
			}
		}
	}

	if len(violations) > 0 {
		t.Errorf(
			"%d mixed visibility issues:",
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
