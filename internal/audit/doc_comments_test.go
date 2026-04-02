//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package audit

import (
	"go/ast"
	"go/token"
	"testing"
)

// TestDocComments ensures all functions (exported and unexported),
// structs, and package-level variables have doc comments.
//
// Test files are exempt.
//
// See specs/ast-audit-tests.md for rationale.
func TestDocComments(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			fpath := pkg.Fset.Position(file.Pos()).Filename
			if isTestFile(fpath) {
				continue
			}

			for _, decl := range file.Decls {
				switch d := decl.(type) {
				case *ast.FuncDecl:
					if d.Doc == nil {
						violations = append(violations,
							posString(pkg.Fset, d.Pos())+
								": func "+d.Name.Name+" missing doc comment",
						)
					}

				case *ast.GenDecl:
					// Only check type and var declarations, not
					// import or const.
					if d.Tok != token.TYPE && d.Tok != token.VAR {
						continue
					}

					for _, spec := range d.Specs {
						switch s := spec.(type) {
						case *ast.TypeSpec:
							// Use the GenDecl doc if the TypeSpec
							// has none (common with grouped decls).
							doc := s.Doc
							if doc == nil {
								doc = d.Doc
							}
							if doc == nil {
								violations = append(violations,
									posString(pkg.Fset, s.Pos())+
										": type "+s.Name.Name+" missing doc comment",
								)
							}

						case *ast.ValueSpec:
							// Package-level var. Use GenDecl doc if
							// the ValueSpec has none.
							doc := s.Doc
							if doc == nil {
								doc = d.Doc
							}
							if doc == nil {
								name := "var"
								if len(s.Names) > 0 {
									name = "var " + s.Names[0].Name
								}
								// Skip blank identifiers.
								if len(s.Names) > 0 && s.Names[0].Name == "_" {
									continue
								}
								violations = append(violations,
									posString(pkg.Fset, s.Pos())+
										": "+name+" missing doc comment",
								)
							}
						}
					}
				}
			}
		}
	}

	// Report count first for orientation.
	if len(violations) > 0 {
		t.Errorf("%d declarations missing doc comments:", len(violations))
	}
	// Cap output to avoid noise — show first 20.
	limit := 20
	if len(violations) < limit {
		limit = len(violations)
	}
	for _, v := range violations[:limit] {
		t.Error(v)
	}
	if len(violations) > 20 {
		remaining := len(violations) - 20
		_ = remaining
		t.Errorf("... and %d more", len(violations)-20)
	}
}
