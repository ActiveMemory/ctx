//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

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
	"strconv"
	"strings"
	"testing"
)

// deterministicCorePrefixes lists Go package-path prefixes
// that make up ctx's deterministic core: the surfaces
// users hit on every invocation of ctx and that the spec
// Invariant 2 ("zero runtime deps for core functionality")
// requires to keep working without an AI backend.
//
// Adding a prefix here ratchets the boundary; never
// remove one.
var deterministicCorePrefixes = []string{
	"github.com/ActiveMemory/ctx/internal/cli/agent",
	"github.com/ActiveMemory/ctx/internal/cli/status",
	"github.com/ActiveMemory/ctx/internal/cli/load",
	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/",
}

// forbiddenAIPrefixes lists Go package-path prefixes that
// pull in the AI-backend layer (HTTP client, factories,
// CLI surface). Deterministic-core packages must not
// import any of these.
//
// The setup-time backend writer
// (internal/cli/setup/core/backend) is intentionally NOT
// here: it mutates .ctxrc and does not import
// internal/backend at runtime. Same for the rc accessor
// (rc.Backends), which returns plain config data and
// makes no HTTP call.
var forbiddenAIPrefixes = []string{
	"github.com/ActiveMemory/ctx/internal/backend",
	"github.com/ActiveMemory/ctx/internal/cli/ai",
	"github.com/ActiveMemory/ctx/internal/write/ai",
}

// TestDeterministicCoreNoAIImports enforces Invariant 2
// of `specs/ctx-ai-backend.md`: the deterministic ctx
// core (`ctx agent`, `ctx status`, `ctx load`, and every
// `ctx system` ceremony hook) must not import the
// AI-backend layer. Without this gate the additive /
// optional discipline is honour-system only.
//
// Test files are exempt; only non-test source imports
// are checked.
func TestDeterministicCoreNoAIImports(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		if !inDeterministicCore(pkg.PkgPath) {
			continue
		}
		for _, file := range pkg.Syntax {
			fpath := pkg.Fset.Position(file.Pos()).Filename
			if isTestFile(fpath) {
				continue
			}
			for _, imp := range file.Imports {
				path, _ := strconv.Unquote(imp.Path.Value)
				if forbidden := matchForbidden(path); forbidden != "" {
					violations = append(violations,
						posString(pkg.Fset, imp.Pos())+
							": "+pkg.PkgPath+" imports "+path+
							" — deterministic core may not depend on the AI layer ("+forbidden+")",
					)
				}
			}
		}
	}

	for _, v := range violations {
		t.Error(v)
	}
}

// inDeterministicCore reports whether pkgPath belongs to
// a package in the deterministic-core boundary.
//
// Parameters:
//   - pkgPath: full import path to check.
//
// Returns:
//   - bool: true if pkgPath matches any prefix in
//     [deterministicCorePrefixes].
func inDeterministicCore(pkgPath string) bool {
	for _, prefix := range deterministicCorePrefixes {
		if pkgPath == strings.TrimSuffix(prefix, "/") ||
			strings.HasPrefix(pkgPath, prefix) {
			return true
		}
	}
	return false
}

// matchForbidden returns the matching forbidden prefix
// for path, or "" when path is allowed.
//
// Parameters:
//   - path: a non-test source-file import path.
//
// Returns:
//   - string: the offending prefix, or "" when allowed.
func matchForbidden(path string) string {
	for _, prefix := range forbiddenAIPrefixes {
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return prefix
		}
	}
	return ""
}
