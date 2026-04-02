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

// TestNoErrorsAs flags calls to errors.As() which should use the
// generic errors.AsType() (available since Go 1.23). errors.AsType
// avoids the need for a target pointer variable and is the preferred
// convention.
//
// See specs/ast-audit-tests.md for rationale.
func TestNoErrorsAs(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
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

				if ident.Name == "errors" && sel.Sel.Name == "As" {
					violations = append(violations,
						posString(pkg.Fset, call.Pos())+
							": use errors.AsType[T]() instead of errors.As()",
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
