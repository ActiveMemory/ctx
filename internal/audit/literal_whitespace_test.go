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
	"strconv"
	"strings"
	"testing"
)

// literalWhitespace maps unquoted string/rune values to the
// config/token constant that should be used instead.
var literalWhitespace = map[string]string{
	"\n":   "token.NewlineLF",
	"\r\n": "token.NewlineCRLF",
	"\r":   `token.NewlineCRLF or token.NewlineLF`,
	"\t":   "token.Tab",
}

// TestNoLiteralWhitespace ensures bare whitespace string and byte
// literals ("\n", "\r\n", "\r", "\t", '\n', '\r', '\t') only appear
// in internal/config/token/ constant definitions. All other packages
// must use token.NewlineLF, token.NewlineCRLF, token.Tab, etc.
//
// Test files are exempt.
//
// See specs/ast-audit-tests.md for rationale.
func TestNoLiteralWhitespace(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		// Allow the constant definition site.
		if strings.Contains(pkg.PkgPath, "internal/config/token") {
			continue
		}

		for _, file := range pkg.Syntax {
			fpath := pkg.Fset.Position(file.Pos()).Filename
			if isTestFile(fpath) {
				continue
			}

			ast.Inspect(file, func(n ast.Node) bool {
				lit, ok := n.(*ast.BasicLit)
				if !ok {
					return true
				}

				var unquoted string

				switch lit.Kind {
				case token.STRING:
					s, err := strconv.Unquote(lit.Value)
					if err != nil {
						return true
					}
					unquoted = s

				case token.CHAR:
					r, _, _, err := strconv.UnquoteChar(lit.Value[1:], '\'')
					if err != nil {
						return true
					}
					unquoted = string(r)

				default:
					return true
				}

				suggestion, found := literalWhitespace[unquoted]
				if !found {
					return true
				}

				// Skip constant definition sites (const blocks).
				if isConstDef(file, lit) {
					return true
				}

				violations = append(violations,
					posString(pkg.Fset, lit.Pos())+
						": literal whitespace "+lit.Value+
						" — use "+suggestion,
				)

				return true
			})
		}
	}

	for _, v := range violations {
		t.Error(v)
	}
}

// isConstDef reports whether lit appears inside a const declaration.
func isConstDef(file *ast.File, lit *ast.BasicLit) bool {
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.CONST {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, val := range vs.Values {
				if containsNode(val, lit) {
					return true
				}
			}
		}
	}

	return false
}

// containsNode reports whether root contains target.
func containsNode(root ast.Node, target ast.Node) bool {
	found := false

	ast.Inspect(root, func(n ast.Node) bool {
		if n == target {
			found = true
			return false
		}

		return !found
	})

	return found
}
