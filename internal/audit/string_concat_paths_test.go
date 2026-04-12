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
	"strings"
	"testing"
	"unicode"
)

// pathVarHints lists substrings that, when present in a variable name
// (case-insensitive), suggest the variable holds a filesystem path.
var pathVarHints = []string{"path", "dir", "folder", "file"}

// DO NOT add entries here to make tests pass. New code must
// conform to the check. Widening requires a dedicated PR with
// justification for each entry.
//
// stringConcatPathAllowlist lists known false positives where string
// concatenation is used for non-path purposes (e.g. substring search
// patterns, extension appending).
var stringConcatPathAllowlist = map[string]bool{
	// Builds a substring search pattern for strings.Contains, not a
	// filesystem path.
	"subagentPath": true,
	// Appends file extension, not building a path with separators.
	"filename": true,
}

// TestNoStringConcatPaths ensures variables whose names suggest they
// hold filesystem paths (containing "path", "dir", "folder", "file")
// are not assigned via string concatenation. Use filepath.Join instead.
//
// Exempt: internal/config/ (constant definitions), test files, and
// names in the allowlist.
//
// See specs/ast-audit-tests.md for rationale.
func TestNoStringConcatPaths(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		if strings.Contains(pkg.PkgPath, "internal/config/") ||
			strings.HasSuffix(pkg.PkgPath, "internal/config") {
			continue
		}

		for _, file := range pkg.Syntax {
			fpath := pkg.Fset.Position(file.Pos()).Filename
			if isTestFile(fpath) {
				continue
			}

			ast.Inspect(file, func(n ast.Node) bool {
				assign, ok := n.(*ast.AssignStmt)
				if !ok {
					return true
				}

				for i, lhs := range assign.Lhs {
					ident, ok := lhs.(*ast.Ident)
					if !ok {
						continue
					}

					if stringConcatPathAllowlist[ident.Name] {
						continue
					}

					if !looksLikePathVar(ident.Name) {
						continue
					}

					if i >= len(assign.Rhs) {
						continue
					}

					if containsStringConcat(assign.Rhs[i]) {
						violations = append(violations,
							posString(pkg.Fset, assign.Pos())+
								": "+ident.Name+
								" built with string concat — use filepath.Join",
						)
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

// looksLikePathVar reports whether name contains a path-related
// substring at a word boundary (case-insensitive).
func looksLikePathVar(name string) bool {
	lower := strings.ToLower(name)
	for _, hint := range pathVarHints {
		idx := strings.Index(lower, hint)
		if idx < 0 {
			continue
		}

		// Check word boundary: start of string or preceded by a
		// lowercase→uppercase transition (PascalCase) or underscore.
		if idx == 0 || !unicode.IsLetter(rune(lower[idx-1])) {
			return true
		}
		// PascalCase boundary: the original char at idx is uppercase.
		if idx < len(name) && unicode.IsUpper(rune(name[idx])) {
			return true
		}
	}

	return false
}

// containsStringConcat reports whether expr contains a binary + with
// at least one string operand (BasicLit STRING or CHAR).
func containsStringConcat(expr ast.Expr) bool {
	found := false

	ast.Inspect(expr, func(n ast.Node) bool {
		if found {
			return false
		}

		bin, ok := n.(*ast.BinaryExpr)
		if !ok || bin.Op != token.ADD {
			return true
		}

		if isStringExpr(bin.X) || isStringExpr(bin.Y) {
			found = true
			return false
		}

		return true
	})

	return found
}

// isStringExpr reports whether expr is a string literal, a call to
// string(), or a selector on filepath/os that returns a string.
func isStringExpr(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return e.Kind == token.STRING
	case *ast.CallExpr:
		// string(x) conversion
		if ident, ok := e.Fun.(*ast.Ident); ok && ident.Name == "string" {
			return true
		}
	}

	return false
}
