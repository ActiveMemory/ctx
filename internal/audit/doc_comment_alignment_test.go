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

// TestDocCommentAlignment ensures that the first line of every
// function doc comment starts with the function name, per Go
// convention and CONVENTIONS.md:
//
//	// FunctionName does X.
//	func FunctionName(...)
//
// This catches copy-paste errors where the comment belongs to
// a different function.
func TestDocCommentAlignment(t *testing.T) {
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
				if !ok || fn.Doc == nil {
					continue
				}

				name := fn.Name.Name

				// The first line of the doc comment must
				// start with the function name.
				firstLine := fn.Doc.List[0].Text
				// Strip the leading "// " or "//" prefix.
				text := strings.TrimPrefix(firstLine, "// ")
				text = strings.TrimPrefix(text, "//")

				if !strings.HasPrefix(text, name) {
					violations = append(violations,
						posString(pkg.Fset, fn.Doc.Pos())+
							": func "+name+
							" doc comment starts with "+
							truncate(text, 50)+
							" (must start with "+name+")",
					)
				}
			}
		}
	}

	if len(violations) > 0 {
		t.Errorf(
			"%d doc comments do not start with the function name:",
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
		t.Errorf("... and %d more", len(violations)-20)
	}
}

// truncate shortens s to n characters, appending "..." if
// truncated.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
