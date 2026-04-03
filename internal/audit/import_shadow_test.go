//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package audit

import (
	"go/ast"
	"go/token"
	"strings"
	"testing"

	"golang.org/x/tools/go/packages"
)

// TestNoImportNameShadowing ensures local variable and
// parameter names do not collide with imported package
// names in the same file. For example, a variable named
// "session" when "session" is also an import alias
// makes the import inaccessible in that scope.
//
// Test files are exempt.
func TestNoImportNameShadowing(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			fpath := pkg.Fset.Position(file.Pos()).Filename
			if isTestFile(fpath) {
				continue
			}

			// Collect import names for this file.
			importNames := fileImportNames(file)
			if len(importNames) == 0 {
				continue
			}

			// Walk all declarations for local
			// variable/param names that collide.
			ast.Inspect(file, func(n ast.Node) bool {
				switch v := n.(type) {
				case *ast.AssignStmt:
					if v.Tok != token.DEFINE {
						return true
					}
					for _, lhs := range v.Lhs {
						checkIdent(
							pkg, importNames,
							lhs, &violations,
						)
					}

				case *ast.RangeStmt:
					if v.Key != nil {
						checkIdent(
							pkg, importNames,
							v.Key, &violations,
						)
					}
					if v.Value != nil {
						checkIdent(
							pkg, importNames,
							v.Value, &violations,
						)
					}

				case *ast.FuncDecl:
					checkFuncParams(
						pkg, importNames,
						v, &violations,
					)
				}

				return true
			})
		}
	}

	if len(violations) > 0 {
		t.Errorf(
			"%d import name shadows found:",
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

// fileImportNames returns the set of effective import
// names (aliases or last path element) for a file.
func fileImportNames(file *ast.File) map[string]bool {
	names := make(map[string]bool)
	for _, imp := range file.Imports {
		if imp.Name != nil {
			// Explicit alias.
			if imp.Name.Name != "_" &&
				imp.Name.Name != "." {
				names[imp.Name.Name] = true
			}
			continue
		}
		// Default: last path element.
		path := strings.Trim(imp.Path.Value, `"`)
		parts := strings.Split(path, "/")
		names[parts[len(parts)-1]] = true
	}
	return names
}

// checkIdent flags an identifier if its name matches
// an imported package name.
func checkIdent(
	pkg *packages.Package,
	importNames map[string]bool,
	expr ast.Expr,
	violations *[]string,
) {
	ident, ok := expr.(*ast.Ident)
	if !ok || ident.Name == "_" {
		return
	}

	if importNames[ident.Name] {
		*violations = append(*violations,
			posString(pkg.Fset, ident.Pos())+
				": var "+ident.Name+
				" shadows import",
		)
	}
}

// checkFuncParams flags function parameters and return
// names that shadow import names.
func checkFuncParams(
	pkg *packages.Package,
	importNames map[string]bool,
	fn *ast.FuncDecl,
	violations *[]string,
) {
	if fn.Type.Params != nil {
		for _, field := range fn.Type.Params.List {
			for _, name := range field.Names {
				if importNames[name.Name] {
					*violations = append(
						*violations,
						posString(
							pkg.Fset, name.Pos(),
						)+
							": param "+name.Name+
							" shadows import",
					)
				}
			}
		}
	}

	if fn.Type.Results != nil {
		for _, field := range fn.Type.Results.List {
			for _, name := range field.Names {
				if importNames[name.Name] {
					*violations = append(
						*violations,
						posString(
							pkg.Fset, name.Pos(),
						)+
							": return "+name.Name+
							" shadows import",
					)
				}
			}
		}
	}
}
