//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package audit

import (
	"go/ast"
	"strings"
	"testing"
)

// printfMethods lists cobra cmd.Printf-style methods
// that bypass the write package formatting pipeline.
var printfMethods = map[string]bool{
	"Printf":    true,
	"PrintErrf": true,
}

// TestNoPrintfCalls ensures cmd.Printf and
// cmd.PrintErrf are not used anywhere. All formatted
// output must go through internal/write/ which uses
// cmd.Print/cmd.Println with pre-formatted strings
// from desc.Text().
//
// Test files are exempt.
//
// See specs/ast-audit-tests.md for rationale.
func TestNoPrintfCalls(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			fpath := pkg.Fset.Position(file.Pos()).Filename
			if isTestFile(fpath) {
				continue
			}

			// Allow calls inside internal/write/.
			if strings.Contains(pkg.PkgPath, "internal/write/") ||
				strings.HasSuffix(pkg.PkgPath, "internal/write") {
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

				if !printfMethods[sel.Sel.Name] {
					return true
				}

				ident, ok := sel.X.(*ast.Ident)
				if !ok {
					return true
				}

				if ident.Name == "cmd" {
					violations = append(violations,
						posString(pkg.Fset, call.Pos())+
							": cmd."+sel.Sel.Name+
							"() — use write/ helpers",
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
