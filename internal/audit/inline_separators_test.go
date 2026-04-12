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
	"go/token"
	"strings"
	"testing"
)

// TestNoInlineSeparators ensures strings.Join calls use token constants
// for the separator argument, not string literals. All separator strings
// must come from internal/config/token/.
//
// The definition site (internal/config/token/) is exempt.
// Test files are exempt.
//
// See specs/ast-audit-tests.md for rationale.
func TestNoInlineSeparators(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		if strings.Contains(pkg.PkgPath, "internal/config/token") {
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

				if ident.Name != "strings" || sel.Sel.Name != "Join" {
					return true
				}

				// strings.Join takes 2 args; the second is the separator.
				if len(call.Args) < 2 {
					return true
				}

				lit, ok := call.Args[1].(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					return true
				}

				violations = append(violations,
					posString(pkg.Fset, lit.Pos())+
						": strings.Join with literal separator "+lit.Value+
						" — use config/token constants",
				)

				return true
			})
		}
	}

	for _, v := range violations {
		t.Error(v)
	}
}
