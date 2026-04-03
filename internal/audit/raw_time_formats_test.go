//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package audit

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
	"testing"
)

// timeFormatFuncs lists time package methods that take
// a layout string argument.
var timeFormatFuncs = map[string]bool{
	"Parse":        true,
	"Format":       true,
	"AppendFormat": true,
}

// goRefTime matches substrings of Go's reference time
// (Mon Jan 2 15:04:05 MST 2006) that indicate a raw
// time layout string.
var goRefTimeHints = []string{
	"2006", "01", "02", "15", "04", "05",
	"Jan", "Mon", "MST",
}

// TestNoRawTimeFormats ensures raw time layout strings
// in time.Parse, time.Format, and time.AppendFormat
// calls use config/time constants instead of inline
// format strings.
//
// Stdlib constants (time.RFC3339, etc.) are exempt
// because they are already named.
//
// The definition site (internal/config/time/) is exempt.
// Test files are exempt.
//
// See specs/ast-audit-tests.md for rationale.
func TestNoRawTimeFormats(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		if strings.Contains(pkg.PkgPath, "config/time") {
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

				if !timeFormatFuncs[sel.Sel.Name] {
					return true
				}

				// Find the layout argument: first arg
				// for Parse, first arg for Format/
				// AppendFormat (method on time.Time).
				var layoutArg ast.Expr
				switch sel.Sel.Name {
				case "Parse":
					// time.Parse(layout, value)
					if len(call.Args) >= 1 {
						layoutArg = call.Args[0]
					}
				case "Format", "AppendFormat":
					// t.Format(layout)
					if len(call.Args) >= 1 {
						layoutArg = call.Args[0]
					}
				}

				if layoutArg == nil {
					return true
				}

				lit, ok := layoutArg.(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					return true
				}

				s, err := strconv.Unquote(lit.Value)
				if err != nil {
					return true
				}

				if looksLikeTimeLayout(s) {
					violations = append(violations,
						posString(pkg.Fset, lit.Pos())+
							": raw time format "+
							lit.Value+
							" — use config/time",
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

// looksLikeTimeLayout reports whether s contains
// Go reference time components.
func looksLikeTimeLayout(s string) bool {
	for _, hint := range goRefTimeHints {
		if strings.Contains(s, hint) {
			return true
		}
	}
	return false
}
