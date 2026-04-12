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
	"go/types"
	"strings"
	"testing"
	"unicode"

	"golang.org/x/tools/go/packages"
)

// TestNoDeadExports flags exported constants, variables,
// functions, and types in internal/ that have zero
// references outside their definition file.
//
// Test files are exempt (both as definition and usage
// sites).
//
// Unexported symbols are skipped: they are package-
// internal and may be used via reflection or are
// genuinely file-scoped helpers.

// DO NOT add entries here to make tests pass. New code must
// conform to the check. Widening requires a dedicated PR with
// justification for each entry.
//
// linuxOnlyExports lists exported symbols used only from
// _linux.go source files. These appear dead on non-Linux
// builds because go/packages loads only the current
// platform's file set.
var linuxOnlyExports = map[string]bool{
	"github.com/ActiveMemory/ctx/internal/config/sysinfo.ProcLoadavg":       true,
	"github.com/ActiveMemory/ctx/internal/config/sysinfo.ProcMeminfo":       true,
	"github.com/ActiveMemory/ctx/internal/config/sysinfo.LoadavgFmt":        true,
	"github.com/ActiveMemory/ctx/internal/config/sysinfo.MemInfoSuffix":     true,
	"github.com/ActiveMemory/ctx/internal/config/sysinfo.BytesPerKB":        true,
	"github.com/ActiveMemory/ctx/internal/config/sysinfo.FieldMemTotal":     true,
	"github.com/ActiveMemory/ctx/internal/config/sysinfo.FieldMemAvailable": true,
	"github.com/ActiveMemory/ctx/internal/config/sysinfo.FieldMemFree":      true,
	"github.com/ActiveMemory/ctx/internal/config/sysinfo.FieldBuffers":      true,
	"github.com/ActiveMemory/ctx/internal/config/sysinfo.FieldCached":       true,
	"github.com/ActiveMemory/ctx/internal/config/sysinfo.FieldSwapTotal":    true,
	"github.com/ActiveMemory/ctx/internal/config/sysinfo.FieldSwapFree":     true,
}

func TestNoDeadExports(t *testing.T) {
	pkgs := loadPackages(t)

	// Also load cmd/ packages to catch cross-boundary
	// usage (cmd/ctx/main.go calls internal/ exports).
	cmdPkgs := loadCmdPackages(t)
	allPkgs := make([]*packages.Package, 0, len(pkgs)+len(cmdPkgs))
	allPkgs = append(allPkgs, pkgs...)
	allPkgs = append(allPkgs, cmdPkgs...)

	// Phase 1: collect all exported definitions.
	// Key: "pkgPath.Name" (stable across type-checker
	// instances). Value: definition metadata.
	type defInfo struct {
		label string // e.g. "const config/dep.BuilderGo"
		pos   string // file:line
		file  string // definition filename
	}
	defs := make(map[string]defInfo)

	for _, pkg := range pkgs {
		for ident, obj := range pkg.TypesInfo.Defs {
			if obj == nil {
				continue
			}
			if !isExported(ident.Name) {
				continue
			}

			pos := pkg.Fset.Position(ident.Pos())
			if isTestFile(pos.Filename) {
				continue
			}

			kind := objectKind(obj)
			if kind == "" {
				continue
			}

			key := obj.Pkg().Path() + "." + obj.Name()
			defs[key] = defInfo{
				label: kind + " " +
					shortPkg(pkg.PkgPath) +
					"." + ident.Name,
				pos:  pos.String(),
				file: pos.Filename,
			}
		}
	}

	// Phase 2: collect all usage sites. Remove any
	// def that has at least one use outside its own
	// definition file. Scan both internal/ and cmd/.
	for _, pkg := range allPkgs {
		for ident, obj := range pkg.TypesInfo.Uses {
			if obj == nil || obj.Pkg() == nil {
				continue
			}

			pos := pkg.Fset.Position(ident.Pos())
			if isTestFile(pos.Filename) {
				continue
			}

			key := obj.Pkg().Path() + "." + obj.Name()
			_, defined := defs[key]
			if !defined {
				continue
			}

			// Any use (same or different package)
			// means the symbol is alive.
			delete(defs, key)
		}
	}

	// Phase 2.5: remove symbols used cross-package in
	// test files. If a test in package B imports a
	// symbol from package A, the symbol is test
	// infrastructure — not dead. Same-package test
	// usage does not count (those should be unexported).
	testPkgs := loadTestPackages(t)
	for _, pkg := range testPkgs {
		for ident, obj := range pkg.TypesInfo.Uses {
			if obj == nil || obj.Pkg() == nil {
				continue
			}
			pos := pkg.Fset.Position(ident.Pos())
			if !isTestFile(pos.Filename) {
				continue
			}
			// Cross-package: the test's package path
			// differs from the symbol's defining package.
			if pkg.PkgPath == obj.Pkg().Path() {
				continue
			}
			key := obj.Pkg().Path() + "." +
				obj.Name()
			delete(defs, key)
		}
	}

	// Phase 3b: remove Linux-only exports (used from
	// _linux.go files not loaded on this platform).
	for key := range linuxOnlyExports {
		delete(defs, key)
	}

	// Phase 4: report survivors as dead exports.
	var violations []string
	for _, info := range defs {
		violations = append(violations,
			info.pos+
				": dead export: "+info.label,
		)
	}

	if len(violations) == 0 {
		return
	}

	t.Errorf(
		"%d dead exports found:", len(violations),
	)
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

// loadCmdPackages loads cmd/ packages for cross-
// boundary usage detection.
func loadCmdPackages(t *testing.T) []*packages.Package {
	t.Helper()
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo,
		Tests: false,
	}
	pkgs, err := packages.Load(
		cfg,
		"github.com/ActiveMemory/ctx/cmd/...",
	)
	if err != nil {
		t.Fatalf("packages.Load cmd: %v", err)
	}
	return pkgs
}

// loadTestPackages loads internal/ packages WITH test
// files for cross-package test usage detection.
func loadTestPackages(
	t *testing.T,
) []*packages.Package {
	t.Helper()
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo,
		Tests: true,
	}
	pkgs, loadErr := packages.Load(
		cfg,
		"github.com/ActiveMemory/ctx/internal/...",
	)
	if loadErr != nil {
		t.Fatalf("packages.Load tests: %v", loadErr)
	}
	return pkgs
}

// isExported reports whether name starts with an
// uppercase letter.
func isExported(name string) bool {
	if name == "" {
		return false
	}
	return unicode.IsUpper(rune(name[0]))
}

// objectKind returns a human-readable kind string for
// a types.Object, or "" to skip.
func objectKind(obj types.Object) string {
	switch o := obj.(type) {
	case *types.Const:
		return "const"
	case *types.Var:
		// Skip struct fields and function parameters.
		// Only flag package-level vars.
		if obj.Parent() == nil {
			return ""
		}
		return "var"
	case *types.Func:
		// Skip methods (have receivers) — they may
		// implement interfaces via dynamic dispatch.
		if o.Type().(*types.Signature).Recv() != nil {
			return ""
		}
		return "func"
	case *types.TypeName:
		return "type"
	default:
		return ""
	}
}

// shortPkg returns the last two path elements of a
// package path for readable labels.
func shortPkg(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) <= 2 {
		return path
	}
	return strings.Join(parts[len(parts)-2:], "/")
}
