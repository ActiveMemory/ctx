//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

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
