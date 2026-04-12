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
	"strings"
	"testing"
)

// nakedErrorFuncs maps package-qualified function names that construct
// errors inline. All error construction must go through internal/err/.
var nakedErrorFuncs = map[string]map[string]bool{
	"fmt":    {"Errorf": true},
	"errors": {"New": true},
}

// TestNoNakedErrors ensures fmt.Errorf and errors.New calls only appear
// in internal/err/** packages. All other packages must use the
// corresponding error constructors from internal/err/.
//
// Test files are exempt.
//
// See specs/ast-audit-tests.md for rationale.
func TestNoNakedErrors(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		if strings.Contains(pkg.PkgPath, "internal/err/") ||
			strings.HasSuffix(pkg.PkgPath, "internal/err") {
			continue
		}

		for _, file := range pkg.Syntax {
			fpath := pkg.Fset.Position(file.Pos()).Filename
			if isTestFile(fpath) {
				continue
			}

			ast.Inspect(file, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				ident, ok := sel.X.(*ast.Ident)
				if !ok {
					return true
				}

				methods, found := nakedErrorFuncs[ident.Name]
				if !found {
					return true
				}

				if methods[sel.Sel.Name] {
					violations = append(violations,
						posString(pkg.Fset, call.Pos())+
							": "+ident.Name+"."+sel.Sel.Name+
							"() must be in internal/err/",
					)
				}

				return true
			})
		}
	}

	for _, v := range violations {
		t.Error(v)
	}
}
