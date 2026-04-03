//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package audit

import (
	"go/ast"
	"strings"
	"testing"
)

// descLookupFuncs maps desc lookup function names to
// the embed package that should supply their DescKey.
var descLookupFuncs = map[string]string{
	"Text":    "config/embed/text",
	"Flag":    "config/embed/flag",
	"Command": "config/embed/cmd",
}

// TestUseConstantsOnlyInCobraUse ensures Use* constants
// from config/embed/cmd/ only appear in cobra Use:
// struct field assignments.
//
// Test files are exempt.
func TestUseConstantsOnlyInCobraUse(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		// Skip the definition site.
		if strings.Contains(
			pkg.PkgPath, "config/embed/cmd",
		) {
			continue
		}

		for _, file := range pkg.Syntax {
			fpath := pkg.Fset.Position(file.Pos()).Filename
			if isTestFile(fpath) {
				continue
			}

			ast.Inspect(file, func(n ast.Node) bool {
				sel, ok := n.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				// Match cmd.UseXxx selectors.
				if !strings.HasPrefix(
					sel.Sel.Name, "Use",
				) {
					return true
				}

				ident, ok := sel.X.(*ast.Ident)
				if !ok {
					return true
				}

				// Check the import resolves to
				// config/embed/cmd.
				if !isEmbedCmdImport(file, ident.Name) {
					return true
				}

				// Verify parent is a cobra struct field
				// (Use:, Aliases:, or Short:/Long: for
				// parent commands that use Use* in group
				// headers).
				if !isCobraField(file, sel) {
					violations = append(violations,
						posString(pkg.Fset, sel.Pos())+
							": "+ident.Name+"."+
							sel.Sel.Name+
							" used outside cobra field",
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

// TestDescKeyOnlyInLookupCalls ensures DescKey*
// constants are only passed to desc.Text(),
// desc.Flag(), or desc.Command(). No other function
// should receive a DescKey as an argument.
//
// Test files are exempt.
func TestDescKeyOnlyInLookupCalls(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		// Skip definition sites and flagbind (which
		// receives descKey as a parameter, not a
		// constant).
		if strings.Contains(
			pkg.PkgPath, "config/embed/",
		) || strings.Contains(
			pkg.PkgPath, "flagbind",
		) {
			continue
		}

		for _, file := range pkg.Syntax {
			fpath := pkg.Fset.Position(file.Pos()).Filename
			if isTestFile(fpath) {
				continue
			}

			ast.Inspect(file, func(n ast.Node) bool {
				sel, ok := n.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				if !strings.HasPrefix(
					sel.Sel.Name, "DescKey",
				) {
					return true
				}

				// Check if the selector resolves to
				// an embed package.
				ident, ok := sel.X.(*ast.Ident)
				if !ok {
					return true
				}

				if !isEmbedImport(file, ident.Name) {
					return true
				}

				// Check that the parent is a desc
				// lookup call, a flagbind call, or a
				// data structure (struct/map/slice
				// literal).
				if !isDescLookupArg(file, sel) &&
					!isFlagbindArg(file, sel) &&
					!isDataStructure(file, sel) {
					violations = append(violations,
						posString(pkg.Fset, sel.Pos())+
							": "+ident.Name+"."+
							sel.Sel.Name+
							" not in desc/flagbind call",
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

// TestNoWrongNamespaceLookup ensures DescKeys from
// config/embed/text are passed to desc.Text(), from
// config/embed/cmd to desc.Command(), and from
// config/embed/flag to desc.Flag(). No cross-namespace
// usage.
//
// Test files are exempt.
func TestNoWrongNamespaceLookup(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		if strings.Contains(
			pkg.PkgPath, "config/embed/",
		) {
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

				// Match desc.Text/Flag/Command calls.
				callSel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				callIdent, ok := callSel.X.(*ast.Ident)
				if !ok {
					return true
				}

				if callIdent.Name != "desc" {
					return true
				}

				expectedPkg, isLookup :=
					descLookupFuncs[callSel.Sel.Name]
				if !isLookup {
					return true
				}

				// Check each argument for DescKey
				// selectors from the wrong package.
				for _, arg := range call.Args {
					argSel, ok :=
						arg.(*ast.SelectorExpr)
					if !ok {
						continue
					}

					if !strings.HasPrefix(
						argSel.Sel.Name, "DescKey",
					) {
						continue
					}

					argIdent, ok :=
						argSel.X.(*ast.Ident)
					if !ok {
						continue
					}

					argPkg := resolveImportPath(
						file, argIdent.Name,
					)
					if argPkg == "" {
						continue
					}

					if !strings.Contains(
						argPkg, expectedPkg,
					) {
						violations = append(
							violations,
							posString(
								pkg.Fset,
								arg.Pos(),
							)+
								": "+argIdent.Name+
								"."+argSel.Sel.Name+
								" passed to desc."+
								callSel.Sel.Name+
								"() (wrong namespace)",
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

// isFlagbindArg reports whether sel appears as an
// argument to a flagbind.* function call.
func isFlagbindArg(
	file *ast.File, target *ast.SelectorExpr,
) bool {
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if found {
			return false
		}

		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}

		if ident.Name != "flagbind" {
			return true
		}

		for _, arg := range call.Args {
			if containsNode(arg, target) {
				found = true
				return false
			}
		}

		return true
	})

	return found
}

// isDataStructure reports whether sel appears inside
// a composite literal (struct, map, or slice), or as
// a function call argument (covers err constructors,
// ceremony configs, etc.).
func isDataStructure(
	file *ast.File, target *ast.SelectorExpr,
) bool {
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if found {
			return false
		}

		switch v := n.(type) {
		case *ast.CompositeLit:
			if containsNode(v, target) {
				found = true
				return false
			}
		case *ast.CallExpr:
			for _, arg := range v.Args {
				if containsNode(arg, target) {
					found = true
					return false
				}
			}
		case *ast.AssignStmt:
			for _, rhs := range v.Rhs {
				if containsNode(rhs, target) {
					found = true
					return false
				}
			}
		}

		return true
	})

	return found
}

// isEmbedCmdImport reports whether alias resolves to
// config/embed/cmd in the file's imports.
func isEmbedCmdImport(
	file *ast.File, alias string,
) bool {
	path := resolveImportPath(file, alias)
	return strings.Contains(path, "config/embed/cmd")
}

// isEmbedImport reports whether alias resolves to any
// config/embed/ subpackage.
func isEmbedImport(
	file *ast.File, alias string,
) bool {
	path := resolveImportPath(file, alias)
	return strings.Contains(path, "config/embed/")
}

// resolveImportPath returns the import path for the
// given alias in the file, or "" if not found.
func resolveImportPath(
	file *ast.File, alias string,
) string {
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)

		if imp.Name != nil {
			if imp.Name.Name == alias {
				return path
			}
			continue
		}

		// Default alias: last path element.
		parts := strings.Split(path, "/")
		if parts[len(parts)-1] == alias {
			return path
		}
	}

	return ""
}

// cobraFields lists struct field names where Use*
// constants are legitimate.
var cobraFields = map[string]bool{
	"Use":         true,
	"Aliases":     true,
	"GroupID":     true,
	"Annotations": true,
}

// isCobraField reports whether sel appears as the
// value of a cobra command struct field (Use:,
// Aliases:, GroupID:, etc.) or as a function argument.
func isCobraField(
	file *ast.File, target *ast.SelectorExpr,
) bool {
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if found {
			return false
		}

		// KeyValue in struct literal.
		kv, isKV := n.(*ast.KeyValueExpr)
		if !isKV {
			// Also accept as function call argument
			// (e.g. AddGroup, AddCommand).
			call, isCall := n.(*ast.CallExpr)
			if isCall {
				for _, arg := range call.Args {
					if containsNode(arg, target) {
						found = true
						return false
					}
				}
			}
			return true
		}

		ident, ok := kv.Key.(*ast.Ident)
		if !ok {
			return true
		}

		if cobraFields[ident.Name] &&
			containsNode(kv.Value, target) {
			found = true
			return false
		}

		return true
	})

	return found
}

// isDescLookupArg reports whether sel appears as an
// argument to desc.Text(), desc.Flag(), or
// desc.Command().
func isDescLookupArg(
	file *ast.File, target *ast.SelectorExpr,
) bool {
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if found {
			return false
		}

		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}

		if ident.Name != "desc" {
			return true
		}

		_, isLookup := descLookupFuncs[sel.Sel.Name]
		if !isLookup {
			return true
		}

		for _, arg := range call.Args {
			if containsNode(arg, target) {
				found = true
				return false
			}
		}

		return true
	})

	return found
}
