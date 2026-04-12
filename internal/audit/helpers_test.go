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
	"sync"
	"testing"

	"golang.org/x/tools/go/packages"
)

// cachedPkgs holds the result of loadPackages, populated once per test run.
var (
	pkgOnce    sync.Once
	cachedPkgs []*packages.Package
	cachedErr  error
)

// loadPackages loads and caches all packages matching the internal/...
// pattern with full syntax trees. The result is shared across all tests
// in a single run via sync.Once.
//
// Parameters:
//   - t: test context for fatal errors
//
// Returns:
//   - []*packages.Package: parsed packages with syntax and type info
func loadPackages(t *testing.T) []*packages.Package {
	t.Helper()

	pkgOnce.Do(func() {
		cfg := &packages.Config{
			Mode: packages.NeedName |
				packages.NeedFiles |
				packages.NeedSyntax |
				packages.NeedTypes |
				packages.NeedTypesInfo,
			Tests: false,
		}
		cachedPkgs, cachedErr = packages.Load(cfg, "github.com/ActiveMemory/ctx/internal/...")
	})

	if cachedErr != nil {
		t.Fatalf("packages.Load: %v", cachedErr)
	}

	return cachedPkgs
}

// isTestFile reports whether filename is a _test.go file.
//
// Parameters:
//   - filename: file path to check
//
// Returns:
//   - bool: true if the file ends with _test.go
func isTestFile(filename string) bool {
	return strings.HasSuffix(filename, "_test.go")
}

// posString formats a token position as file:line for error messages.
//
// Parameters:
//   - fset: file set for position resolution
//   - pos: token position to format
//
// Returns:
//   - string: formatted as "file:line"
func posString(fset *token.FileSet, pos token.Pos) string {
	p := fset.Position(pos)
	return p.String()
}

// importLitPositions returns the token positions of all
// string literals in import declarations.
func importLitPositions(
	file *ast.File,
) map[token.Pos]bool {
	positions := make(map[token.Pos]bool)
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.IMPORT {
			continue
		}
		for _, spec := range gd.Specs {
			imp, ok := spec.(*ast.ImportSpec)
			if ok && imp.Path != nil {
				positions[imp.Path.Pos()] = true
			}
		}
	}
	return positions
}

// isStructTag reports whether lit is a struct field tag.
func isStructTag(
	file *ast.File, lit *ast.BasicLit,
) bool {
	if lit.Kind != token.STRING {
		return false
	}
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if found {
			return false
		}
		field, ok := n.(*ast.Field)
		if !ok {
			return true
		}
		if field.Tag == lit {
			found = true
			return false
		}
		return true
	})
	return found
}
