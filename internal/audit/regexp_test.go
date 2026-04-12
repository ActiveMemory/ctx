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
)

// TestNoInlineRegexpCompile ensures regexp.MustCompile and
// regexp.Compile calls appear only at package level (var
// declarations), never inside function bodies. The project
// centralizes compiled patterns in internal/config/regex/ as
// package-level vars. Inline compilation causes per-call overhead
// and scatters pattern definitions.
//
// Test files are exempt.
//
// See specs/ast-audit-tests.md for rationale.
func TestNoInlineRegexpCompile(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			fpath := pkg.Fset.Position(file.Pos()).Filename
			if isTestFile(fpath) {
				continue
			}

			// Walk only inside function bodies — skip package-level
			// var declarations which are the correct location.
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || fn.Body == nil {
					continue
				}

				ast.Inspect(fn.Body, func(n ast.Node) bool {
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

					if ident.Name == "regexp" &&
						(sel.Sel.Name == "MustCompile" || sel.Sel.Name == "Compile") {
						violations = append(violations,
							posString(pkg.Fset, call.Pos())+
								": regexp."+sel.Sel.Name+"() must be a package-level var, not inside a function",
						)
					}

					return true
				})
			}
		}
	}

	for _, v := range violations {
		t.Error(v)
	}
}
