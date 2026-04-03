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

	"golang.org/x/tools/go/packages"
)

// exemptStrings lists string values always acceptable.
var exemptStrings = map[string]bool{
	"":     true, // empty string
	" ":    true, // single space
	"  ":   true, // two-space indent
	"    ": true, // four-space indent
	". ":   true, // sentence separator
	": ":   true, // key-value separator
}

// exemptStringPackages lists package paths fully exempt
// from magic string checks.
var exemptStringPackages = []string{
	"internal/config/",
	"internal/config",
	"internal/assets/tpl",
	"internal/err/",
}

// TestNoMagicStrings flags magic string literals in non-test
// Go files under internal/.
//
// Exempt: empty string, single space, indentation strings,
// single characters, format verbs, regex replacements, HTML
// entities, URL scheme prefixes, config/tpl/err packages,
// file-level const/var definitions, import paths, struct tags.
//
// Test files are exempt.
//
// See specs/ast-audit-tests.md for rationale.
func TestNoMagicStrings(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		if isExemptStringPackage(pkg.PkgPath) {
			continue
		}

		for _, file := range pkg.Syntax {
			fpath := pkg.Fset.Position(file.Pos()).Filename
			if isTestFile(fpath) {
				continue
			}

			imports := importLitPositions(file)

			ast.Inspect(file, func(n ast.Node) bool {
				lit, ok := n.(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					return true
				}

				if imports[lit.Pos()] {
					return true
				}

				if isConstDef(file, lit) ||
					isVarDef(file, lit) {
					return true
				}

				if isStructTag(file, lit) {
					return true
				}

				checkMagicString(
					pkg, file, lit, &violations,
				)

				return true
			})
		}
	}

	if len(violations) > 0 {
		t.Errorf(
			"%d magic strings found:",
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

// checkMagicString flags non-exempt string literals.
func checkMagicString(
	pkg *packages.Package, _ *ast.File,
	lit *ast.BasicLit, violations *[]string,
) {
	s, err := strconv.Unquote(lit.Value)
	if err != nil {
		return
	}

	if exemptStrings[s] {
		return
	}

	// Format verbs ("%s", "%d %s", etc.).
	if isFormatString(s) {
		return
	}

	// Regex capture group references.
	if isRegexRef(s) {
		return
	}

	// HTML entities (&lt;, &gt;, etc.).
	if strings.HasPrefix(s, "&") &&
		strings.HasSuffix(s, ";") {
		return
	}

	// URL scheme prefixes.
	if strings.HasSuffix(s, "://") {
		return
	}

	*violations = append(*violations,
		posString(pkg.Fset, lit.Pos())+
			": magic string "+lit.Value,
	)
}

// isExemptStringPackage reports whether pkgPath matches
// an exempt package for string checks.
func isExemptStringPackage(pkgPath string) bool {
	for _, exempt := range exemptStringPackages {
		if strings.Contains(pkgPath, exempt) {
			return true
		}
	}
	return false
}

// isFormatString reports whether s looks like a printf
// format string (contains % followed by a verb char).
func isFormatString(s string) bool {
	if !strings.Contains(s, "%") {
		return false
	}
	// Any string containing a % verb is likely a format
	// string. Accept any string with at least one
	// standard format directive.
	for i := 0; i < len(s)-1; i++ {
		if s[i] != '%' {
			continue
		}
		next := s[i+1]
		if next == '%' {
			i++ // skip %%
			continue
		}
		// Skip flags/width/precision chars.
		j := i + 1
		for j < len(s) &&
			(s[j] == '+' || s[j] == '-' ||
				s[j] == '#' || s[j] == '0' ||
				s[j] == ' ' ||
				(s[j] >= '1' && s[j] <= '9') ||
				s[j] == '.') {
			j++
		}
		if j < len(s) {
			verb := s[j]
			if strings.ContainsRune(
				"sdvqwfegxXobctpT", rune(verb),
			) {
				return true
			}
		}
	}
	return false
}

// isRegexRef reports whether s contains regex capture
// group references ($1, $2, etc.).
func isRegexRef(s string) bool {
	return strings.Contains(s, "$1") ||
		strings.Contains(s, "$2") ||
		strings.Contains(s, "$3") ||
		strings.Contains(s, "${")
}
