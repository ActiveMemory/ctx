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
)

// TestDocComments ensures all functions (exported and unexported),
// structs, and package-level variables have doc comments.
//
// Test files are exempt.
//
// See specs/ast-audit-tests.md for rationale.
// configPackage reports whether pkgPath is a config-style
// package where group doc comments are sufficient for
// const blocks.
func configPackage(pkgPath string) bool {
	return strings.Contains(pkgPath, "internal/config/") ||
		strings.HasSuffix(pkgPath, "internal/config")
}

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
					// Skip import declarations.
					if d.Tok == token.IMPORT {
						continue
					}

					// Singleton declarations (no parens) use
					// the GenDecl doc for the single spec.
					singleton := !d.Lparen.IsValid()
					isCfg := configPackage(pkg.PkgPath)

					// Config const/var blocks (excluding
					// embed/): group doc covers all specs
					// because names are self-documenting
					// (Dir*, File*, Perm*, etc.).
					//
					// DO NOT widen this exemption. New code
					// must have per-constant doc comments.
					// Widening requires a dedicated PR with
					// justification — not a drive-by allowlist
					// change to make tests pass.
					if isCfg && d.Lparen.IsValid() &&
						!strings.Contains(
							pkg.PkgPath, "config/embed/",
						) &&
						(d.Tok == token.CONST ||
							d.Tok == token.VAR) {
						if d.Doc == nil {
							violations = append(violations,
								posString(pkg.Fset, d.Pos())+
									": const/var block"+
									" missing group doc",
							)
						}
						continue
					}

					for _, spec := range d.Specs {
						switch s := spec.(type) {
						case *ast.TypeSpec:
							doc := s.Doc
							if doc == nil && singleton {
								doc = d.Doc
							}
							if doc == nil {
								violations = append(violations,
									posString(pkg.Fset, s.Pos())+
										": type "+s.Name.Name+
										" missing doc comment",
								)
							}

						case *ast.ValueSpec:
							// Skip blank identifiers.
							if len(s.Names) > 0 &&
								s.Names[0].Name == "_" {
								continue
							}

							doc := s.Doc
							if doc == nil && singleton {
								doc = d.Doc
							}
							if doc == nil {
								tok := "var"
								if d.Tok == token.CONST {
									tok = "const"
								}
								name := tok
								if len(s.Names) > 0 {
									name = tok + " " +
										s.Names[0].Name
								}
								violations = append(violations,
									posString(pkg.Fset, s.Pos())+
										": "+name+
										" missing doc comment",
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
		t.Errorf(
			"... and %d more",
			len(violations)-20,
		)
	}
}
