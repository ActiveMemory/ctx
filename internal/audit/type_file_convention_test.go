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
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// DO NOT add entries here to make tests pass. New code must
// conform to the check. Widening requires a dedicated PR with
// justification for each entry.
//
// exemptTypePackages lists package path segments where
// types intentionally do NOT live in types.go. Each
// has a documented reason.
var exemptTypePackages = map[string]string{
	// entity/ organizes types by domain concept
	// (task.go, session.go, etc.), not in types.go.
	"entity": "domain-organized types",

	// proto/ defines protocol schema types in
	// schema.go — single-file protocol definition.
	"proto": "protocol schema",

	// err/ packages colocate error types with their
	// domain context.
	"err/": "error types colocated with domain",
}

// isExemptTypePackage returns true if the package path
// matches an exempt pattern.
func isExemptTypePackage(pkgPath string) bool {
	for pattern := range exemptTypePackages {
		if strings.Contains(pkgPath, "/"+pattern) ||
			strings.HasSuffix(pkgPath, "/"+pattern) {
			return true
		}
	}
	return false
}

// fileTypeAnalysis holds the result of analyzing a
// single file for type placement.
type fileTypeAnalysis struct {
	file       string   // absolute path
	pkg        string   // short package path
	types      []string // type names defined in file
	pureImpl   bool     // file is pure type impl
	reason     string   // why not pure (if applicable)
	exported   int      // exported type count
	unexported int      // unexported type count
}

// analyzeFileForTypes checks whether a file is a "pure
// type implementation file". A pure type impl file
// contains ONLY:
//   - type declarations (struct, interface, etc.)
//   - method receivers on types defined in the file
//   - interface compliance assertions:
//     var _ Iface = (*Type)(nil)
//
// Additionally, every exported type in the file must
// have at least one exported method receiver.
//
// Any standalone function (non-receiver), non-type
// const/var (except compliance assertions), or exported
// type without an exported receiver disqualifies the
// file.
func analyzeFileForTypes(
	file *ast.File,
	filename string,
	pkg string,
) *fileTypeAnalysis {
	result := &fileTypeAnalysis{
		file:     filename,
		pkg:      pkg,
		pureImpl: true,
	}

	// Phase 1: collect type names defined in this file.
	// Interfaces are tracked separately because they are
	// behavioral contracts, not data blueprints — they
	// cannot carry method receivers and are exempt from
	// the phase 4 "must have exported receiver" rule.
	typeNames := make(map[string]bool)
	interfaceTypes := make(map[string]bool)
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			typeNames[ts.Name.Name] = true
			if _, isIface := ts.Type.(*ast.InterfaceType); isIface {
				interfaceTypes[ts.Name.Name] = true
			}
			result.types = append(
				result.types, ts.Name.Name,
			)
			if isExported(ts.Name.Name) {
				result.exported++
			} else {
				result.unexported++
			}
		}
	}

	if len(typeNames) == 0 {
		return nil // no types, not relevant
	}

	// Phase 2: collect exported receivers per type.
	exportedReceivers := make(map[string]bool)
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fn.Recv == nil {
			// Standalone function — not pure.
			result.pureImpl = false
			result.reason = fmt.Sprintf(
				"standalone func %s", fn.Name.Name,
			)
			return result
		}
		// It's a method. Check receiver type is
		// defined in this file.
		recvType := receiverTypeName(fn.Recv)
		if !typeNames[recvType] {
			result.pureImpl = false
			result.reason = fmt.Sprintf(
				"receiver %s.%s on foreign type",
				recvType, fn.Name.Name,
			)
			return result
		}
		if isExported(fn.Name.Name) {
			exportedReceivers[recvType] = true
		}
	}

	// Phase 3: check var/const declarations. Only
	// interface compliance assertions are allowed:
	//   var _ Iface = (*Type)(nil)
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if gd.Tok == token.TYPE {
			continue // already handled
		}
		if gd.Tok == token.IMPORT {
			continue
		}
		// const or var — check each spec.
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			if !isComplianceAssertion(vs) {
				result.pureImpl = false
				result.reason = fmt.Sprintf(
					"non-assertion %s %s",
					gd.Tok, nameFromValueSpec(vs),
				)
				return result
			}
		}
	}

	// Phase 4: every exported type must have at least
	// one exported receiver.
	//
	// Interfaces are exempt: they are behavioral
	// contracts (method signatures) and cannot have
	// receivers attached. A file containing only an
	// interface definition (e.g. session.go with
	// "type Session interface { ... }") is a valid
	// pure type file under this convention.
	for _, name := range result.types {
		if !isExported(name) {
			continue
		}
		if interfaceTypes[name] {
			continue
		}
		if !exportedReceivers[name] {
			result.pureImpl = false
			result.reason = fmt.Sprintf(
				"exported type %s has no "+
					"exported receiver",
				name,
			)
			return result
		}
	}

	return result
}

// receiverTypeName extracts the base type name from a
// method receiver field list. Handles both value and
// pointer receivers.
func receiverTypeName(
	recv *ast.FieldList,
) string {
	if recv == nil || len(recv.List) == 0 {
		return ""
	}
	expr := recv.List[0].Type
	// Unwrap pointer.
	if star, ok := expr.(*ast.StarExpr); ok {
		expr = star.X
	}
	// Unwrap index expression for generic types.
	if idx, ok := expr.(*ast.IndexExpr); ok {
		expr = idx.X
	}
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name
	}
	return ""
}

// isComplianceAssertion returns true if the ValueSpec
// is a var _ Iface = (*Type)(nil) pattern.
func isComplianceAssertion(vs *ast.ValueSpec) bool {
	// Must have exactly one name and it must be _.
	if len(vs.Names) != 1 {
		return false
	}
	return vs.Names[0].Name == "_"
}

// nameFromValueSpec returns the first name in a
// ValueSpec for diagnostic output.
func nameFromValueSpec(vs *ast.ValueSpec) string {
	if len(vs.Names) > 0 {
		return vs.Names[0].Name
	}
	return "<anonymous>"
}

// typeViolation records a type definition found outside
// types.go that is not in a pure implementation file.
type typeViolation struct {
	file     string
	typeName string
	exported bool
	pkg      string
}

// grandfatheredTypes MUST remain empty.
//
// This list was intentionally cleared. No new entries may be
// added under any circumstances — not as a "temporary" fix,
// not "to unblock CI", not "while I work on the real fix",
// not "because the migration is non-trivial". Zero tolerance.
//
// If a new type violates TestTypeFileConvention, the correct
// action is to FIX THE UNDERLYING VIOLATION by either:
//
//  1. Moving the type to types.go in the same package, or
//  2. Making the file a "pure type impl file" (only the type
//     declaration, its method receivers, and interface
//     compliance asserts — no standalone functions), or
//  3. Deleting the type if it is dead.
//
// Any PR that re-introduces entries here must:
//   - Be a DEDICATED pull request whose sole purpose is to
//     add the grandfather entry (no code changes alongside)
//   - Include a written justification for EACH entry in the
//     PR description
//   - Be explicitly approved by a project maintainer
//   - Cite a specific follow-up task in TASKS.md committing
//     to the eventual fix
//
// Drive-by additions by agents (Claude Code, Copilot, Cursor,
// or any other assistant) while completing unrelated work are
// NOT authorized and will be reverted. An agent that adds an
// entry here is operating outside its mandate.
//
// If you are an agent reading this comment: do not add to this
// map. Fix the violation instead. If you cannot fix the
// violation, report the blocker to the human operator and
// stop.
var grandfatheredTypes = map[string]bool{}

// TestTypeFileConvention flags type definitions that
// live outside types.go unless the file is a pure type
// implementation file.
//
// A pure type impl file contains only: type
// declarations, method receivers on those types, and
// interface compliance assertions. Every exported type
// must have at least one exported receiver.
//
// Known violations are grandfathered. New violations
// fail the test. Removing a grandfathered entry (by
// migrating the type) is always safe.
//
// Exempt packages: entity/ (domain-organized), proto/
// (schema), err/ (error types).
func TestTypeFileConvention(t *testing.T) {
	pkgs := loadPackages(t)

	var violations []typeViolation

	for _, pkg := range pkgs {
		if isExemptTypePackage(pkg.PkgPath) {
			continue
		}

		shortPkgPath := shortenPkgPath(pkg.PkgPath)

		for _, file := range pkg.Syntax {
			pos := pkg.Fset.Position(file.Pos())
			filename := pos.Filename
			base := filepath.Base(filename)

			if base == "types.go" || base == "doc.go" {
				continue
			}
			if isTestFile(base) {
				continue
			}

			analysis := analyzeFileForTypes(
				file, filename, shortPkgPath,
			)
			if analysis == nil {
				continue // no types
			}
			if analysis.pureImpl {
				continue // exempt
			}

			for _, name := range analysis.types {
				violations = append(
					violations,
					typeViolation{
						file:     filename,
						typeName: name,
						exported: isExported(name),
						pkg:      shortPkgPath,
					},
				)
			}
		}
	}

	// Separate new violations from grandfathered ones.
	var newViolations []typeViolation
	seen := make(map[string]bool)
	for _, v := range violations {
		key := v.pkg + "/" +
			filepath.Base(v.file) + ":" + v.typeName
		seen[key] = true
		if !grandfatheredTypes[key] {
			newViolations = append(newViolations, v)
		}
	}

	// Check for stale grandfathered entries (type was
	// migrated but entry not removed).
	var stale []string
	for key := range grandfatheredTypes {
		if !seen[key] {
			stale = append(stale, key)
		}
	}
	sort.Strings(stale)
	for _, key := range stale {
		t.Errorf(
			"stale grandfathered entry "+
				"(type was migrated — remove): %s",
			key,
		)
	}

	if len(newViolations) == 0 {
		return
	}

	byPkg := make(map[string][]typeViolation)
	for _, v := range newViolations {
		byPkg[v.pkg] = append(byPkg[v.pkg], v)
	}

	pkgNames := make([]string, 0, len(byPkg))
	for p := range byPkg {
		pkgNames = append(pkgNames, p)
	}
	sort.Strings(pkgNames)

	exported := 0
	unexported := 0
	for _, v := range newViolations {
		if v.exported {
			exported++
		} else {
			unexported++
		}
	}

	t.Errorf(
		"%d NEW type definitions outside types.go "+
			"(%d exported, %d unexported) "+
			"across %d packages:",
		len(newViolations), exported, unexported,
		len(byPkg),
	)

	for _, pkg := range pkgNames {
		vv := byPkg[pkg]
		byFile := make(map[string][]string)
		for _, v := range vv {
			base := filepath.Base(v.file)
			byFile[base] = append(
				byFile[base], v.typeName,
			)
		}

		files := make([]string, 0, len(byFile))
		for f := range byFile {
			files = append(files, f)
		}
		sort.Strings(files)

		for _, f := range files {
			types := byFile[f]
			t.Errorf(
				"  %s/%s: %s",
				pkg, f,
				strings.Join(types, ", "),
			)
		}
	}
}

// shortenPkgPath trims the module prefix for readable
// output.
func shortenPkgPath(pkgPath string) string {
	const prefix = "github.com/ActiveMemory/ctx/"
	if idx := strings.Index(
		pkgPath, prefix,
	); idx >= 0 {
		return pkgPath[idx+len(prefix):]
	}
	return pkgPath
}

// TestTypeFileConventionReport generates a detailed
// report of type placement across the codebase. Run
// with -run TestTypeFileConventionReport -v for the
// full report without failing.
func TestTypeFileConventionReport(t *testing.T) {
	pkgs := loadPackages(t)

	type reportEntry struct {
		pkg      string
		file     string
		types    []string
		pureImpl bool
		reason   string
		exported int
		unexport int
	}

	var inTypes int
	var outViolations []reportEntry
	var outExempted []reportEntry
	var pkgExempted []reportEntry

	for _, pkg := range pkgs {
		shortPkgPath := shortenPkgPath(pkg.PkgPath)
		pkgExempt := isExemptTypePackage(pkg.PkgPath)

		for _, file := range pkg.Syntax {
			pos := pkg.Fset.Position(file.Pos())
			filename := pos.Filename
			base := filepath.Base(filename)

			if base == "doc.go" || isTestFile(base) {
				continue
			}

			if base == "types.go" {
				for _, decl := range file.Decls {
					gd, ok := decl.(*ast.GenDecl)
					if !ok {
						continue
					}
					for _, spec := range gd.Specs {
						if _, ok := spec.(*ast.TypeSpec); ok {
							inTypes++
						}
					}
				}
				continue
			}

			analysis := analyzeFileForTypes(
				file, filename, shortPkgPath,
			)
			if analysis == nil {
				continue
			}

			entry := reportEntry{
				pkg:      shortPkgPath,
				file:     base,
				types:    analysis.types,
				pureImpl: analysis.pureImpl,
				reason:   analysis.reason,
				exported: analysis.exported,
				unexport: analysis.unexported,
			}

			switch {
			case pkgExempt:
				pkgExempted = append(
					pkgExempted, entry,
				)
			case analysis.pureImpl:
				outExempted = append(
					outExempted, entry,
				)
			default:
				outViolations = append(
					outViolations, entry,
				)
			}
		}
	}

	// Count totals.
	violationTypes := 0
	for _, e := range outViolations {
		violationTypes += len(e.types)
	}
	exemptedTypes := 0
	for _, e := range outExempted {
		exemptedTypes += len(e.types)
	}
	pkgExemptTypes := 0
	for _, e := range pkgExempted {
		pkgExemptTypes += len(e.types)
	}

	total := inTypes + violationTypes + exemptedTypes
	fmt.Println()
	fmt.Println(
		"=== Type File Convention Report ===",
	)
	fmt.Println()
	fmt.Printf(
		"Types in types.go:              %d\n",
		inTypes,
	)
	fmt.Printf(
		"Types in pure impl files:       %d "+
			"(auto-exempt)\n",
		exemptedTypes,
	)
	fmt.Printf(
		"Types in exempt packages:       %d\n",
		pkgExemptTypes,
	)
	fmt.Printf(
		"Types violating convention:     %d\n",
		violationTypes,
	)
	fmt.Println()

	if total > 0 {
		compliant := inTypes + exemptedTypes
		pct := float64(compliant) /
			float64(total) * 100
		fmt.Printf(
			"Compliance: %.1f%% (%d/%d)\n",
			pct, compliant, total,
		)
	}
	fmt.Println()

	// Pure impl files (auto-exempted).
	fmt.Println(
		"--- Pure type implementation files " +
			"(auto-exempt) ---",
	)
	fmt.Println()
	sort.Slice(outExempted, func(i, j int) bool {
		if outExempted[i].pkg != outExempted[j].pkg {
			return outExempted[i].pkg < outExempted[j].pkg
		}
		return outExempted[i].file < outExempted[j].file
	})
	for _, e := range outExempted {
		fmt.Printf(
			"  %s/%s: %s\n",
			e.pkg, e.file,
			strings.Join(e.types, ", "),
		)
	}
	fmt.Println()

	// Violations.
	fmt.Println("--- Violations ---")
	fmt.Println()
	sort.Slice(outViolations, func(i, j int) bool {
		if outViolations[i].pkg !=
			outViolations[j].pkg {
			return outViolations[i].pkg <
				outViolations[j].pkg
		}
		return outViolations[i].file <
			outViolations[j].file
	})
	for _, e := range outViolations {
		fmt.Printf(
			"  %s/%s: %s\n"+
				"    reason: %s\n",
			e.pkg, e.file,
			strings.Join(e.types, ", "),
			e.reason,
		)
	}
	fmt.Println()

	// Package exemptions.
	fmt.Println("--- Package-exempt ---")
	fmt.Println()
	exemptByPkg := make(map[string]int)
	for _, e := range pkgExempted {
		exemptByPkg[e.pkg] += len(e.types)
	}
	exemptPkgs := make([]string, 0, len(exemptByPkg))
	for p := range exemptByPkg {
		exemptPkgs = append(exemptPkgs, p)
	}
	sort.Strings(exemptPkgs)
	for _, p := range exemptPkgs {
		reason := ""
		for pattern, r := range exemptTypePackages {
			if strings.Contains(p, pattern) {
				reason = r
				break
			}
		}
		fmt.Printf(
			"  %s: %d types (%s)\n",
			p, exemptByPkg[p], reason,
		)
	}
}
