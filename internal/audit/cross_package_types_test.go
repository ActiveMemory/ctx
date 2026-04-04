//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package audit

import (
	"go/types"
	"strings"
	"testing"

	"golang.org/x/tools/go/packages"
)

// typeExemptPackages lists packages where exported
// types are expected to be used cross-package by
// design (entity, config, proto, etc.).
var typeExemptPackages = map[string]bool{
	"entity": true,
	"proto":  true,
	"core":   true,
}

// TestCrossPackageTypes flags exported type
// definitions that are used from other packages but
// are not in internal/entity/ or other exempt
// packages. Cross-cutting types should live in
// internal/entity/ for discoverability.
//
// Test files are exempt.
func TestCrossPackageTypes(t *testing.T) {
	pkgs := loadPackages(t)
	cmdPkgs := loadCmdPackages(t)
	allPkgs := make(
		[]*packages.Package,
		0, len(pkgs)+len(cmdPkgs),
	)
	allPkgs = append(allPkgs, pkgs...)
	allPkgs = append(allPkgs, cmdPkgs...)

	// Phase 1: collect all exported type definitions
	// outside exempt packages.
	type typeDef struct {
		pkg  string
		name string
		pos  string
	}
	defs := make(map[string]typeDef) // key: pkgPath.Name

	for _, pkg := range pkgs {
		// Skip exempt packages.
		parts := strings.Split(pkg.PkgPath, "/")
		lastPart := parts[len(parts)-1]
		if typeExemptPackages[lastPart] {
			continue
		}
		// Skip core/ subpackages (e.g. core/check,
		// core/python) — types there serve their
		// parent CLI module by design.
		if isCoreSubpackage(pkg.PkgPath) {
			continue
		}
		// Skip config/ (types there are config
		// structs, not domain types).
		if strings.Contains(
			pkg.PkgPath, "internal/config/",
		) {
			continue
		}

		for ident, obj := range pkg.TypesInfo.Defs {
			if obj == nil {
				continue
			}
			_, isTypeName := obj.(*types.TypeName)
			if !isTypeName {
				continue
			}
			if !isExported(ident.Name) {
				continue
			}
			pos := pkg.Fset.Position(ident.Pos())
			if isTestFile(pos.Filename) {
				continue
			}

			key := obj.Pkg().Path() + "." + obj.Name()
			defs[key] = typeDef{
				pkg:  shortPkg(pkg.PkgPath),
				name: obj.Name(),
				pos:  pos.String(),
			}
		}
	}

	// Phase 2: find types used cross-package.
	crossPkgUse := make(map[string]string) // key → using pkg

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
			if _, defined := defs[key]; !defined {
				continue
			}

			// Cross-package if user's package differs
			// from definition's package.
			if pkg.PkgPath == obj.Pkg().Path() {
				continue
			}
			// Skip when the consumer is a core
			// subpackage — these are internal to
			// their CLI module by design.
			if isCoreSubpackage(pkg.PkgPath) {
				continue
			}
			// Skip same-module usage. Types shared
			// within a module (e.g. mcp/handler →
			// mcp/server) are module-internal.
			if sameModule(
				pkg.PkgPath, obj.Pkg().Path(),
			) {
				continue
			}
			crossPkgUse[key] = shortPkg(
				pkg.PkgPath,
			)
		}
	}

	// Phase 3: report types used cross-package that
	// are not in entity/.
	var violations []string
	for key, usingPkg := range crossPkgUse {
		def := defs[key]
		violations = append(violations,
			def.pos+": type "+def.pkg+"."+
				def.name+" used from "+usingPkg+
				" (consider entity/)",
		)
	}

	if len(violations) == 0 {
		return
	}

	t.Errorf(
		"%d cross-package types outside entity/:",
		len(violations),
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

// isCoreSubpackage returns true if pkgPath is a
// subpackage of a core/ directory (e.g.
// ".../cli/doctor/core/check").
func isCoreSubpackage(pkgPath string) bool {
	return strings.Contains(pkgPath, "/core/")
}

// domainAliases maps write/ package names to their
// corresponding domain module when the names differ.
var domainAliases = map[string]string{
	"resource": "sysinfo",
}

// sameModule returns true if two package paths share
// the same domain. Handles: same module root,
// cli/<X> consuming internal/<X>, err/<X> consumed
// from cli/<X>, and domain aliases.
func sameModule(a, b string) bool {
	ma := canonicalModule(moduleRoot(a))
	mb := canonicalModule(moduleRoot(b))
	if ma == "" || mb == "" {
		return false
	}
	if ma == mb {
		return true
	}
	// cli/* consuming any domain module is the
	// standard consumer layer pattern.
	if isConsumerLayer(ma) && !isConsumerLayer(mb) {
		return true
	}
	if isConsumerLayer(mb) && !isConsumerLayer(ma) {
		return true
	}
	// err/<X> consumed from cli/<X> or <X>.
	if strings.HasPrefix(ma, "err/") {
		base := ma[len("err/"):]
		if mb == base ||
			mb == "cli/"+base {
			return true
		}
	}
	if strings.HasPrefix(mb, "err/") {
		base := mb[len("err/"):]
		if ma == base ||
			ma == "cli/"+base {
			return true
		}
	}
	return false
}

// canonicalModule resolves domain aliases.
func canonicalModule(mod string) string {
	if alias, ok := domainAliases[mod]; ok {
		return alias
	}
	return mod
}

// moduleRoot extracts the first path segment after
// "internal/" as the module root. For cli/ packages,
// uses the CLI subcommand (e.g. "cli/doctor").
// For write/<X>, uses X to match with internal/<X>.
func moduleRoot(pkgPath string) string {
	const prefix = "ctx/internal/"
	idx := strings.Index(pkgPath, prefix)
	if idx < 0 {
		return ""
	}
	rest := pkgPath[idx+len(prefix):]

	// write/<X> → X
	if strings.HasPrefix(rest, "write/") {
		parts := strings.SplitN(
			rest[len("write/"):], "/", 2,
		)
		return parts[0]
	}

	// cli/<X> → cli/<X>
	if strings.HasPrefix(rest, "cli/") {
		parts := strings.SplitN(
			rest[len("cli/"):], "/", 2,
		)
		return "cli/" + parts[0]
	}

	// err/<X> → err/<X>
	if strings.HasPrefix(rest, "err/") {
		parts := strings.SplitN(
			rest[len("err/"):], "/", 2,
		)
		return "err/" + parts[0]
	}

	// Top-level: mcp, trace, notify, sysinfo, etc.
	parts := strings.SplitN(rest, "/", 2)
	return parts[0]
}

// isConsumerLayer returns true if the module root is a
// consumer layer that naturally imports domain types.
func isConsumerLayer(mod string) bool {
	return strings.HasPrefix(mod, "cli/")
}
