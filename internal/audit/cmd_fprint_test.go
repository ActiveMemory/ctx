//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

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

// fmtFprintMethods lists fmt.Fprint-family helpers that, when
// pointed at a user-facing stream (cmd.OutOrStdout / cmd.OutOrStderr
// / os.Stdout / os.Stderr), bypass the internal/write/ formatting
// pipeline. The cmd_print and printf_calls tests catch the direct
// `cmd.Print*(...)` form; this test closes the indirect form
// `fmt.Fprint*(stream, ...)`.
var fmtFprintMethods = map[string]bool{
	"Fprint":   true,
	"Fprintf":  true,
	"Fprintln": true,
}

// TestNoFmtFprintToUserStream catches `fmt.Fprint*(stream, ...)`
// calls where stream is a user-facing destination
// (cmd.OutOrStdout / cmd.OutOrStderr / os.Stdout / os.Stderr) made
// outside internal/write/. Same intent as TestNoCmdPrintOutsideWrite:
// every user-visible write must route through write/ so output
// formatting stays consistent and template-driven.
//
// Calls writing to in-memory destinations (strings.Builder,
// bytes.Buffer, json.Encoder targets, etc.) are unaffected because
// those arguments are neither cmd.OutOr* calls nor os.Std* idents.
//
// Test files are exempt.
//
// See specs/ast-audit-tests.md for rationale.
func TestNoFmtFprintToUserStream(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		// Allow calls inside internal/write/ — that is precisely
		// where these patterns belong.
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

				// Must be the fmt package.
				pkgIdent, ok := sel.X.(*ast.Ident)
				if !ok || pkgIdent.Name != "fmt" {
					return true
				}

				if !fmtFprintMethods[sel.Sel.Name] {
					return true
				}

				if len(call.Args) == 0 {
					return true
				}

				if !isUserFacingStream(call.Args[0]) {
					return true
				}

				violations = append(violations,
					posString(pkg.Fset, call.Pos())+
						": fmt."+sel.Sel.Name+
						"(<user stream>, ...) — must go through internal/write/",
				)
				return true
			})
		}
	}

	for _, v := range violations {
		t.Error(v)
	}
}

// isUserFacingStream reports whether expr is one of the
// user-visible writers we forbid bypassing.
//
// Recognized shapes:
//   - cmd.OutOrStdout() — cobra's stdout writer (SetOut or
//     os.Stdout fallback).
//   - cmd.OutOrStderr() — confusingly-named cobra accessor that
//     returns the SetOut writer with **stderr** as fallback. Still
//     a user-visible stream; route through internal/write/.
//   - cmd.ErrOrStderr() — cobra's stderr writer (SetErr or
//     os.Stderr fallback). The actual "write to stderr"
//     accessor; covered here to keep the rule total.
//   - os.Stdout / os.Stderr — direct *os.File globals.
//
// The receiver name "cmd" is the project convention; a non-"cmd"
// receiver is allowed as a calculated escape hatch (rare and
// would show up in review).
//
// Anything else (strings.Builder, bytes.Buffer, json.Encoder
// targets, custom io.Writer, etc.) is in-memory string assembly
// and is not a concern of this test.
//
// Parameters:
//   - expr: AST expression in the first-argument slot of a
//     fmt.Fprint*-family call.
//
// Returns:
//   - bool: true when expr is one of the recognized user streams.
func isUserFacingStream(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.CallExpr:
		sel, ok := e.Fun.(*ast.SelectorExpr)
		if !ok {
			return false
		}
		if sel.Sel.Name != "OutOrStdout" &&
			sel.Sel.Name != "OutOrStderr" &&
			sel.Sel.Name != "ErrOrStderr" {
			return false
		}
		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			return false
		}
		return ident.Name == "cmd"
	case *ast.SelectorExpr:
		ident, ok := e.X.(*ast.Ident)
		if !ok {
			return false
		}
		if ident.Name != "os" {
			return false
		}
		return e.Sel.Name == "Stdout" || e.Sel.Name == "Stderr"
	}
	return false
}
