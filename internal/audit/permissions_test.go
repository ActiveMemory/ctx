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

// knownPerms lists octal permission literals that must use
// config/fs constants instead.
var knownPerms = map[string]bool{
	"0644": true,
	"0755": true,
	"0750": true,
	"0600": true,
	"0700": true,
}

// TestNoRawPermissions ensures octal file permission literals only
// appear in internal/config/fs/ and internal/io/. All other packages
// must use fs.PermFile, fs.PermExec, fs.PermSecret, etc.
//
// Test files are exempt.
//
// See specs/ast-audit-tests.md for rationale.
func TestNoRawPermissions(t *testing.T) {
	pkgs := loadPackages(t)
	var violations []string

	for _, pkg := range pkgs {
		if strings.Contains(pkg.PkgPath, "internal/config/fs") {
			continue
		}

		for _, file := range pkg.Syntax {
			fpath := pkg.Fset.Position(file.Pos()).Filename
			if isTestFile(fpath) {
				continue
			}

			ast.Inspect(file, func(n ast.Node) bool {
				lit, ok := n.(*ast.BasicLit)
				if !ok || lit.Kind != token.INT {
					return true
				}

				if knownPerms[lit.Value] {
					violations = append(violations,
						posString(pkg.Fset, lit.Pos())+
							": raw permission "+lit.Value+" — use config/fs constants",
					)
				}

				return true
			})
		}
	}

	for _, v := range violations {
		t.Error(v)
	}
}
