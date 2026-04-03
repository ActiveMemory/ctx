//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package audit

import (
	"go/ast"
	"testing"
	"unicode"
)

// TestNoMixedVisibility ensures that files containing
// exported functions do not also contain unexported
// functions. Private helpers belong in their own file
// to keep public API files focused and short.
//
// Exempt: doc.go files, test files, and files with
// only one function total (too small to split).
func TestNoMixedVisibility(t *testing.T) {
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

			var exported []string
			var unexported []string

			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}

				name := fn.Name.Name
				if unicode.IsUpper(rune(name[0])) {
					exported = append(
						exported, name,
					)
				} else {
					unexported = append(
						unexported, name,
					)
				}
			}

			// Skip if only one function total.
			total := len(exported) + len(unexported)
			if total <= 1 {
				continue
			}

			// Flag files with both exported and
			// unexported functions.
			if len(exported) > 0 &&
				len(unexported) > 0 {
				for _, name := range unexported {
					violations = append(
						violations,
						fpath+": unexported "+
							name+"() in file "+
							"with exported funcs",
					)
				}
			}
		}
	}

	if len(violations) > 0 {
		t.Errorf(
			"%d mixed visibility issues:",
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
