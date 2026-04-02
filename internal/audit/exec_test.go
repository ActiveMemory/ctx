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

// TestNoExecOutsideExecPkg ensures exec.Command and exec.CommandContext
// calls only appear in internal/exec/** packages. Centralizing command
// execution keeps nolint:gosec annotations and per-command sanitization
// in one place.
//
// Test files are exempt.
//
// See specs/ast-audit-tests.md for rationale.
func TestNoExecOutsideExecPkg(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		if strings.Contains(pkg.PkgPath, "internal/exec/") ||
			strings.HasSuffix(pkg.PkgPath, "internal/exec") {
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

				if ident.Name == "exec" &&
					(sel.Sel.Name == "Command" || sel.Sel.Name == "CommandContext") {
					violations = append(violations,
						posString(pkg.Fset, call.Pos())+
							": exec."+sel.Sel.Name+"() must be in internal/exec/",
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
