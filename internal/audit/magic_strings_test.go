//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package audit

import (
	"go/ast"
	"go/token"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/tools/go/packages"
)

// DO NOT add entries here to make tests pass. New code must
// conform to the check. Widening requires a dedicated PR with
// justification for each entry.
//
// exemptStrings lists string values always acceptable.
var exemptStrings = map[string]bool{
	"": true, // empty string
}

// DO NOT add entries here to make tests pass. New code must
// conform to the check. Widening requires a dedicated PR with
// justification for each entry.
//
// exemptStringPackages lists package paths fully exempt
// from magic string checks.
var exemptStringPackages = []string{
	"internal/config/",
	"internal/config",
	"internal/assets/tpl",
	"internal/hub",
	"internal/err/hub",
	"internal/err/serve",
	"internal/cli/agent/core/budget",
	"internal/cli/agent/core/shared",
	"internal/cli/connect",
	"internal/cli/hub",
	"internal/cli/serve/core/shared",
	"internal/cli/system/core/hubsync",
	"internal/cli/system/cmd/check_hub_sync",
	"internal/exec/daemon",
	"internal/exec/sysinfo",
	"internal/sysinfo",
	"internal/write/hub",
	"internal/write/connect",
	"internal/write/serve",
}

// TestNoMagicStrings flags magic string literals in non-test
// Go files under internal/.
//
// Exempt: empty string, single space, indentation strings,
// regex capture references, config/tpl/err packages,
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

				// Const/var definitions in exempt packages
				// are already skipped (line 61). Outside
				// those packages, string constants are
				// magic strings that belong in config/.
				//
				// DO NOT re-add a blanket isConstDef
				// exemption. It masks constants defined
				// in the wrong package.

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

	// Regex capture group references.
	if isRegexRef(s) {
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

// regexRef matches regex capture group references.
var regexRef = regexp.MustCompile(`\$\d|\$\{`)

// isRegexRef reports whether s contains regex capture
// group references ($1, $2, etc.).
func isRegexRef(s string) bool {
	return regexRef.MatchString(s)
}
