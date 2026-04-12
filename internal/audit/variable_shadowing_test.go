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
	"go/token"
	"testing"
)

// TestNoVariableShadowing detects two forms of variable
// shadowing in non-test Go files:
//
// (a) Bare "err" reuse: multiple := assignments to
// the unadorned name "err" in the same function body.
// The convention requires descriptive names (readErr,
// writeErr, parseErr) so each error site is
// independently identifiable.
//
// (b) General inner-scope shadowing (any variable):
// already caught by golangci-lint's shadow checker,
// which is enabled in .golangci.yml.
//
// Test files are exempt.
func TestNoVariableShadowing(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			fpath := pkg.Fset.Position(file.Pos()).Filename
			if isTestFile(fpath) {
				continue
			}

			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || fn.Body == nil {
					continue
				}

				sites := collectBareErrDefines(
					pkg.Fset, fn.Body,
				)

				if len(sites) > 1 {
					for _, pos := range sites {
						violations = append(
							violations,
							pos+
								": bare err := in "+
								fn.Name.Name+
								" (use descriptive "+
								"name)",
						)
					}
				}
			}
		}
	}

	if len(violations) > 0 {
		t.Errorf(
			"%d bare err := reuses found:",
			len(violations),
		)
	}
	limit := 30
	if len(violations) < limit {
		limit = len(violations)
	}
	for _, v := range violations[:limit] {
		t.Error(v)
	}
	if len(violations) > 30 {
		t.Errorf(
			"... and %d more",
			len(violations)-30,
		)
	}
}

// collectBareErrDefines walks a function body and
// returns positions of := assignments where "err" is
// on the LHS.
func collectBareErrDefines(
	fset *token.FileSet, body *ast.BlockStmt,
) []string {
	var sites []string

	ast.Inspect(body, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		if assign.Tok != token.DEFINE {
			return true
		}

		for _, lhs := range assign.Lhs {
			ident, ok := lhs.(*ast.Ident)
			if !ok {
				continue
			}

			if ident.Name == "err" {
				sites = append(sites,
					posString(fset, assign.Pos()),
				)
				break
			}
		}

		return true
	})

	return sites
}
