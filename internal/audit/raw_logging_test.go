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

// stderrWriteMethods lists direct write methods on os.Stderr that
// bypass internal/log/.
var stderrWriteMethods = map[string]bool{
	"Write":       true,
	"WriteString": true,
}

// stdlibLogFuncs lists stdlib log package functions that bypass
// internal/log/.
var stdlibLogFuncs = map[string]bool{
	"Print":   true,
	"Printf":  true,
	"Println": true,
	"Fatal":   true,
	"Fatalf":  true,
	"Fatalln": true,
	"Panic":   true,
	"Panicf":  true,
	"Panicln": true,
}

// TestNoRawLogging ensures direct logging to stderr or via the stdlib
// log package only appears in internal/log/**. All other packages must
// use log/warn.Warn (stderr warnings) or log/event.Append (structured
// event log).
//
// Detected patterns:
//   - fmt.Fprintf(os.Stderr, ...) / fmt.Fprintln(os.Stderr, ...) / fmt.Fprint(os.Stderr, ...)
//   - os.Stderr.Write / os.Stderr.WriteString
//   - log.Print* / log.Fatal* / log.Panic*
//
// Test files are exempt.
//
// See specs/ast-audit-tests.md for rationale.
func TestNoRawLogging(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		if strings.Contains(pkg.PkgPath, "internal/log/") ||
			strings.HasSuffix(pkg.PkgPath, "internal/log") {
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

				switch x := sel.X.(type) {
				case *ast.Ident:
					// log.Print*, log.Fatal*, log.Panic*
					if x.Name == "log" && stdlibLogFuncs[sel.Sel.Name] {
						violations = append(violations,
							posString(pkg.Fset, call.Pos())+
								": log."+sel.Sel.Name+
								"() must be in internal/log/",
						)
					}

				case *ast.SelectorExpr:
					// os.Stderr.Write / os.Stderr.WriteString
					ident, ok := x.X.(*ast.Ident)
					if !ok {
						break
					}
					if ident.Name == "os" && x.Sel.Name == "Stderr" &&
						stderrWriteMethods[sel.Sel.Name] {
						violations = append(violations,
							posString(pkg.Fset, call.Pos())+
								": os.Stderr."+sel.Sel.Name+
								"() must be in internal/log/",
						)
					}
				}

				// fmt.Fprint*(os.Stderr, ...)
				if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "fmt" {
					method := sel.Sel.Name
					if method == "Fprintf" || method == "Fprintln" || method == "Fprint" {
						if len(call.Args) > 0 {
							if isOsStderr(call.Args[0]) {
								violations = append(violations,
									posString(pkg.Fset, call.Pos())+
										": fmt."+method+
										"(os.Stderr, ...) must be in internal/log/",
								)
							}
						}
					}
				}

				return true
			})
		}
	}

	for _, v := range violations {
		t.Error(v)
	}
}

// isOsStderr reports whether expr is os.Stderr.
func isOsStderr(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "os" && sel.Sel.Name == "Stderr"
}
