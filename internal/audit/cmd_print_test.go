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

// printMethods lists cobra cmd.Print* methods that must be routed
// through internal/write/ packages.
var printMethods = map[string]bool{
	"Print":      true,
	"Printf":     true,
	"Println":    true,
	"PrintErr":   true,
	"PrintErrf":  true,
	"PrintErrln": true,
}

// TestNoCmdPrintOutsideWrite ensures cmd.Print*, cmd.Printf, etc.
// calls only appear in internal/write/** packages. All other packages
// must delegate output through the corresponding write/ subpackage.
//
// Test files are exempt.
//
// See specs/ast-audit-tests.md for rationale.
func TestNoCmdPrintOutsideWrite(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		// Allow calls inside internal/write/ packages.
		if strings.Contains(pkg.PkgPath, "internal/write/") ||
			strings.HasSuffix(pkg.PkgPath, "internal/write") {
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

				if !printMethods[sel.Sel.Name] {
					return true
				}

				// Check if the receiver is named "cmd".
				ident, ok := sel.X.(*ast.Ident)
				if !ok {
					return true
				}

				if ident.Name == "cmd" {
					violations = append(violations,
						posString(pkg.Fset, call.Pos())+
							": cmd."+sel.Sel.Name+"() must be in internal/write/",
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
