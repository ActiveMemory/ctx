//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestFlagYAMLMatchesConstants verifies that the flag name
// segment of every flags.yaml key matches an actual flag
// constant value defined in config/flag/.
//
// For example, flags.yaml key "journal.schema.check.dir"
// has flag segment "dir", which must appear as a constant
// value in config/flag/ (e.g., Dir = "dir").
//
// This catches content drift where a flag is renamed in
// code but the YAML key still references the old name.
func TestFlagYAMLMatchesConstants(t *testing.T) {
	flagValues := collectFlagConstValues(t)
	yamlKeys := loadFlagsYAML(t)

	var violations []string
	for key := range yamlKeys {
		parts := strings.Split(key, ".")
		flagName := parts[len(parts)-1]
		if !flagValues[flagName] {
			violations = append(violations,
				key+": flag name "+flagName+
					" not found in config/flag/ constants")
		}
	}

	for _, v := range violations {
		t.Error(v)
	}
}

// TestCommandYAMLMatchesDescKeys verifies that every
// commands.yaml key has a corresponding DescKey constant
// whose string value matches the YAML key.
//
// This is already partially covered by TestDescKeyYAMLLinkage
// but that test covers all YAML files. This test focuses
// specifically on commands.yaml and verifies the key format
// matches the dot-notation convention.
func TestCommandYAMLMatchesDescKeys(t *testing.T) {
	path := filepath.Join(
		"..", "assets", "commands", "commands.yaml",
	)
	data, readErr := os.ReadFile(path) //nolint:gosec
	if readErr != nil {
		t.Fatalf("read commands.yaml: %v", readErr)
	}

	var doc map[string]any
	if unmarshalErr := yaml.Unmarshal(
		data, &doc,
	); unmarshalErr != nil {
		t.Fatalf("parse commands.yaml: %v", unmarshalErr)
	}

	descKeys := collectDescKeys(t)

	var violations []string
	for key := range doc {
		if !descKeys[key] {
			violations = append(violations,
				"commands.yaml key "+key+
					" has no matching DescKey constant")
		}
	}

	for _, v := range violations {
		t.Error(v)
	}
}

// collectFlagConstValues extracts all string constant values
// from config/flag/.
//
// Returns:
//   - map[string]bool: set of flag name values
func collectFlagConstValues(t *testing.T) map[string]bool {
	t.Helper()
	pkgs := loadPackages(t)
	values := make(map[string]bool)

	for _, pkg := range pkgs {
		if !strings.HasSuffix(
			pkg.PkgPath, "config/flag",
		) {
			continue
		}

		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				gd, ok := decl.(*ast.GenDecl)
				if !ok || gd.Tok != token.CONST {
					continue
				}
				for _, spec := range gd.Specs {
					vs, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}
					for _, val := range vs.Values {
						lit, ok := val.(*ast.BasicLit)
						if !ok ||
							lit.Kind != token.STRING {
							continue
						}
						s := strings.Trim(
							lit.Value, `"`,
						)
						values[s] = true
					}
				}
			}
		}
	}

	return values
}

// loadFlagsYAML reads and parses the flags.yaml file.
//
// Returns:
//   - map[string]any: YAML keys from flags.yaml
func loadFlagsYAML(t *testing.T) map[string]any {
	t.Helper()
	path := filepath.Join(
		"..", "assets", "commands", "flags.yaml",
	)
	data, readErr := os.ReadFile(path) //nolint:gosec
	if readErr != nil {
		t.Fatalf("read flags.yaml: %v", readErr)
	}

	var doc map[string]any
	if unmarshalErr := yaml.Unmarshal(
		data, &doc,
	); unmarshalErr != nil {
		t.Fatalf("parse flags.yaml: %v", unmarshalErr)
	}

	return doc
}
