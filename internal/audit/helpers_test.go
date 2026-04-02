//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package audit

import (
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
