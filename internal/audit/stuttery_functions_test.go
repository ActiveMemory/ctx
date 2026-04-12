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
	"unicode"
)

// TestNoStutteryFunctions ensures function names do not redundantly
// include their package name as a PascalCase word boundary.
//
// Examples of stutter:
//   - write.WriteJournal → Write matches package write
//   - parse.ParseLine → Parse matches package parse
//
// Identity functions (write.Write, write.write) are exempt because
// the entire name IS the package name.
//
// Test files are exempt.
//
// See specs/ast-audit-tests.md for rationale.
func TestNoStutteryFunctions(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		// Use the last element of the package path as the package name.
		parts := strings.Split(pkg.PkgPath, "/")
		pkgName := parts[len(parts)-1]

		for _, file := range pkg.Syntax {
			fpath := pkg.Fset.Position(file.Pos()).Filename
			if isTestFile(fpath) {
				continue
			}

			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}

				// Skip methods — only check package-level functions.
				if fn.Recv != nil {
					continue
				}

				name := fn.Name.Name

				// Identity exemption: the whole name IS the package name.
				if strings.EqualFold(name, pkgName) {
					continue
				}

				if stutters(name, pkgName) {
					violations = append(violations,
						posString(pkg.Fset, fn.Pos())+
							": "+pkgName+"."+name+
							" stutters — remove "+pkgName+" from function name",
					)
				}
			}
		}
	}

	for _, v := range violations {
		t.Error(v)
	}
}

// stutters reports whether funcName contains pkgName as a PascalCase
// word (case-insensitive). It splits funcName at uppercase boundaries
// and checks each word.
func stutters(funcName, pkgName string) bool {
	words := splitPascalCase(funcName)
	lower := strings.ToLower(pkgName)

	for _, w := range words {
		if strings.ToLower(w) == lower {
			return true
		}
	}

	return false
}

// splitPascalCase splits a PascalCase or camelCase identifier into
// words. For example:
//
//	"WriteJournal"          → ["Write", "Journal"]
//	"journalWriteSilent"    → ["journal", "Write", "Silent"]
//	"HTMLParser"            → ["HTML", "Parser"]
//	"oversizeNudgeContent"  → ["oversize", "Nudge", "Content"]
func splitPascalCase(s string) []string {
	if s == "" {
		return nil
	}

	var words []string
	start := 0

	for i := 1; i < len(s); i++ {
		cur := rune(s[i])
		prev := rune(s[i-1])

		// Split at lowercase→uppercase boundary: "writeJournal"
		if unicode.IsLower(prev) && unicode.IsUpper(cur) {
			words = append(words, s[start:i])
			start = i
			continue
		}

		// Split at uppercase run end: "HTMLParser" → "HTML", "Parser"
		if i+1 < len(s) && unicode.IsUpper(prev) && unicode.IsUpper(cur) && unicode.IsLower(rune(s[i+1])) {
			words = append(words, s[start:i])
			start = i
			continue
		}
	}

	words = append(words, s[start:])

	return words
}
