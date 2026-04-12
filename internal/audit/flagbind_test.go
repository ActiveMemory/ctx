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

// flagMethods lists cobra flag registration methods
// that must go through internal/flagbind/.
var flagMethods = map[string]bool{
	"StringVar":      true,
	"StringVarP":     true,
	"BoolVar":        true,
	"BoolVarP":       true,
	"IntVar":         true,
	"IntVarP":        true,
	"IntP":           true,
	"BoolP":          true,
	"StringP":        true,
	"DurationVar":    true,
	"DurationVarP":   true,
	"StringSliceVar": true,
}

// TestNoFlagBindOutsideFlagbind ensures direct cobra
// flag registration (.Flags().StringVar, etc.) only
// appears in internal/flagbind/. All other packages
// must use the flagbind helpers.
//
// Test files are exempt.
//
// See specs/ast-audit-tests.md for rationale.
func TestNoFlagBindOutsideFlagbind(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		if strings.Contains(pkg.PkgPath, "flagbind") {
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

				// Match: x.Flags().Method() or
				// x.PersistentFlags().Method()
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				if !flagMethods[sel.Sel.Name] {
					return true
				}

				// The receiver should be a call to
				// Flags() or PersistentFlags().
				innerCall, ok := sel.X.(*ast.CallExpr)
				if !ok {
					return true
				}

				innerSel, ok := innerCall.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				method := innerSel.Sel.Name
				if method != "Flags" &&
					method != "PersistentFlags" {
					return true
				}

				violations = append(violations,
					posString(pkg.Fset, call.Pos())+
						": ."+method+"()."+
						sel.Sel.Name+
						"() must use flagbind",
				)

				return true
			})
		}
	}

	if len(violations) > 0 {
		t.Errorf(
			"%d direct flag registrations:",
			len(violations),
		)
	}
	limit := 20
	if len(violations) < limit {
		limit = len(violations)
	}
	for _, v := range violations[:limit] {
		t.Error(v)
	}
	if len(violations) > 20 {
		t.Errorf(
			"... and %d more",
			len(violations)-20,
		)
	}
}
