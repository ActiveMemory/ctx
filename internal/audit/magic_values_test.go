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

// exemptIntLiterals lists integer values that are always acceptable.
var exemptIntLiterals = map[string]bool{
	"0":  true,
	"1":  true,
	"-1": true,
}

// strconvRadixBitsize lists numeric literals acceptable as strconv
// radix or bitsize arguments.
var strconvRadixBitsize = map[string]bool{
	"8":  true,
	"10": true,
	"16": true,
	"32": true,
	"64": true,
}

// strconvFuncs lists strconv function names whose radix/bitsize
// arguments are exempt.
var strconvFuncs = map[string]bool{
	"ParseInt":    true,
	"ParseUint":   true,
	"ParseFloat":  true,
	"FormatInt":   true,
	"FormatUint":  true,
	"FormatFloat": true,
	"AppendInt":   true,
	"AppendUint":  true,
	"AppendFloat": true,
}

// exemptPackagePaths lists package path substrings that are fully
// exempt from magic value checks — config definitions, template
// definitions, and error constructors.
var exemptPackagePaths = []string{
	"internal/config/",
	"internal/config",
	"internal/assets/tpl",
	"internal/err/",
}

// TestNoMagicValues flags magic numeric literals in non-test Go files
// under internal/. Walks ast.BasicLit nodes and checks parent context.
//
// Numeric exceptions: 0, 1, -1, 2-10 (small ints), strconv
// radix/bitsize arguments, octal permissions (handled by
// TestNoRawPermissions), const/var definition sites.
//
// Config packages and test files are exempt.
//
// See specs/ast-audit-tests.md for rationale.
func TestNoMagicValues(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		if isExemptPackage(pkg.PkgPath) {
			continue
		}

		for _, file := range pkg.Syntax {
			fpath := pkg.Fset.Position(file.Pos()).Filename
			if isTestFile(fpath) {
				continue
			}

			ast.Inspect(file, func(n ast.Node) bool {
				lit, ok := n.(*ast.BasicLit)
				if !ok || lit.Kind != token.INT {
					return true
				}

				// Skip const/var definition sites (file-level and
				// function-level).
				if isConstDef(file, lit) || isVarDef(file, lit) ||
					isLocalConstDef(file, lit) {
					return true
				}

				if exemptIntLiterals[lit.Value] {
					return true
				}

				// Small integers 2-10 are generally acceptable (arg
				// counts, field indices, slice bounds, arithmetic).
				if val, err := strconv.Atoi(lit.Value); err == nil &&
					val >= 2 && val <= 10 {
					return true
				}

				// Octal permissions are handled by TestNoRawPermissions.
				if strings.HasPrefix(lit.Value, "0o") ||
					strings.HasPrefix(lit.Value, "0O") {
					return true
				}

				// Strconv radix/bitsize arguments.
				if strconvRadixBitsize[lit.Value] &&
					isStrconvArg(file, lit) {
					return true
				}

				violations = append(violations,
					posString(pkg.Fset, lit.Pos())+
						": magic number "+lit.Value,
				)

				return true
			})
		}
	}

	if len(violations) > 0 {
		t.Errorf("%d magic values found:", len(violations))
	}
	limit := 30
	if len(violations) < limit {
		limit = len(violations)
	}
	for _, v := range violations[:limit] {
		t.Error(v)
	}
	if len(violations) > 30 {
		t.Errorf("... and %d more", len(violations)-30)
	}
}

// isExemptPackage reports whether pkgPath matches an exempt package.
func isExemptPackage(pkgPath string) bool {
	for _, exempt := range exemptPackagePaths {
		if strings.Contains(pkgPath, exempt) {
			return true
		}
	}
	return false
}

// isLocalConstDef reports whether lit appears inside a const or var
// declaration within a function body (at any nesting depth).
func isLocalConstDef(file *ast.File, lit *ast.BasicLit) bool {
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if found {
			return false
		}
		ds, ok := n.(*ast.DeclStmt)
		if !ok {
			return true
		}
		gd, ok := ds.Decl.(*ast.GenDecl)
		if !ok || (gd.Tok != token.CONST && gd.Tok != token.VAR) {
			return true
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, val := range vs.Values {
				if containsNode(val, lit) {
					found = true
					return false
				}
			}
		}
		return true
	})
	return found
}

// isVarDef reports whether lit appears inside a var declaration.
func isVarDef(file *ast.File, lit *ast.BasicLit) bool {
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
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

// isStrconvArg reports whether lit is an argument to a strconv
// function (radix or bitsize parameter).
func isStrconvArg(file *ast.File, lit *ast.BasicLit) bool {
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
		if ident.Name != "strconv" || !strconvFuncs[sel.Sel.Name] {
			return true
		}
		for _, arg := range call.Args {
			if arg == lit {
				found = true
				return false
			}
		}
		return true
	})
	return found
}
