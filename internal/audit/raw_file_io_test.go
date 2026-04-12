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

// rawIOFuncs lists os package functions that must be routed through
// internal/io/ Safe* wrappers.
var rawIOFuncs = map[string]bool{
	"ReadFile":  true,
	"WriteFile": true,
	"Open":      true,
	"OpenFile":  true,
	"Create":    true,
	"MkdirAll":  true,
}

// TestNoRawFileIO ensures direct os.ReadFile, os.WriteFile, os.Open,
// os.OpenFile, os.Create, and os.MkdirAll calls only appear in
// internal/io/. All other packages must use the Safe* wrappers which
// centralize path validation, sanitization, and nolint:gosec
// suppression.
//
// Test files are exempt.
//
// See specs/ast-audit-tests.md for rationale.
func TestNoRawFileIO(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		if strings.Contains(pkg.PkgPath, "internal/io") {
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

				if ident.Name == "os" && rawIOFuncs[sel.Sel.Name] {
					violations = append(violations,
						posString(pkg.Fset, call.Pos())+
							": os."+sel.Sel.Name+"() must be in internal/io/",
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
