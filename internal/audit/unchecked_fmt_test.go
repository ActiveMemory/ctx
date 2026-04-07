//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package audit

import (
	"go/ast"
	"testing"
)

// fmtWriteFuncs lists fmt package functions that write to
// an io.Writer and return (int, error). Discarding the error
// violates the "zero silent error discard" convention even
// when the writer never fails (e.g. strings.Builder).
var fmtWriteFuncs = map[string]bool{
	"Fprintf":  true,
	"Fprintln": true,
	"Fprint":   true,
}

// TestNoUncheckedFmtWrite flags fmt.Fprintf, fmt.Fprintln,
// and fmt.Fprint calls whose return values are discarded
// (used as bare statements instead of assigned).
//
// Per CONVENTIONS.md: "Zero silent error discard — handle
// every error, never suppress with _ = or //nolint:errcheck."
//
// Test files are exempt.
func TestNoUncheckedFmtWrite(t *testing.T) {
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

			ast.Inspect(file, func(n ast.Node) bool {
				stmt, ok := n.(*ast.ExprStmt)
				if !ok {
					return true
				}

				call, ok := stmt.X.(*ast.CallExpr)
				if !ok {
					return true
				}

				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				if !fmtWriteFuncs[sel.Sel.Name] {
					return true
				}

				ident, ok := sel.X.(*ast.Ident)
				if !ok {
					return true
				}

				if ident.Name == "fmt" {
					violations = append(violations,
						posString(
							pkg.Fset, call.Pos(),
						)+": fmt."+sel.Sel.Name+
							"() return value "+
							"discarded",
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
