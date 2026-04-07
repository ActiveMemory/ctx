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

// grandfatheredDocStructure is the number of pre-existing
// doc structure violations. New code must not add to this
// count. Reduce it as violations are fixed.
const grandfatheredDocStructure = 0

// TestDocCommentStructure verifies that all documented
// functions with parameters include a "Parameters:" section
// and functions with return values include a "Returns:"
// section in their doc comments, per CONVENTIONS.md.
//
// Applies to both exported and unexported functions.
// Test files are excluded. Functions without doc comments
// are skipped (caught by TestDocComments).
func TestDocCommentStructure(t *testing.T) {
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

			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}
				if fn.Doc == nil {
					continue
				}

				doc := fn.Doc.Text()
				hasParams := fnHasParams(fn)
				hasReturns := fnHasReturns(fn)

				if hasParams &&
					!strings.Contains(
						doc, "Parameters:",
					) {
					violations = append(
						violations,
						posString(
							pkg.Fset, fn.Pos(),
						)+": "+fn.Name.Name+
							" has parameters but"+
							" missing Parameters:"+
							" section",
					)
				}

				if hasReturns &&
					!strings.Contains(
						doc, "Returns:",
					) {
					violations = append(
						violations,
						posString(
							pkg.Fset, fn.Pos(),
						)+": "+fn.Name.Name+
							" has return values but"+
							" missing Returns:"+
							" section",
					)
				}
			}
		}
	}

	if len(violations) > grandfatheredDocStructure {
		t.Errorf(
			"%d doc structure violations "+
				"(grandfathered: %d, new: %d):",
			len(violations),
			grandfatheredDocStructure,
			len(violations)-grandfatheredDocStructure,
		)
		// Show all violations.
		for _, v := range violations {
			t.Error(v)
		}
	} else if len(violations) < grandfatheredDocStructure {
		t.Errorf(
			"violations dropped to %d — "+
				"update grandfatheredDocStructure "+
				"from %d to %d",
			len(violations),
			grandfatheredDocStructure,
			len(violations),
		)
	}
}

// fnHasParams reports whether fn has at least one
// named parameter (excluding the receiver).
func fnHasParams(fn *ast.FuncDecl) bool {
	if fn.Type.Params == nil {
		return false
	}
	for _, field := range fn.Type.Params.List {
		for _, name := range field.Names {
			if name.Name != "_" {
				return true
			}
		}
	}
	return false
}

// fnHasReturns reports whether fn declares at least
// one return value.
func fnHasReturns(fn *ast.FuncDecl) bool {
	return fn.Type.Results != nil &&
		len(fn.Type.Results.List) > 0
}
