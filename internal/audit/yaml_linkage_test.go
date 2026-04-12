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
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestDescKeyYAMLLinkage verifies bidirectional linkage
// between DescKey constants and YAML asset files.
//
// Direction 1: every YAML key must have a corresponding
// DescKey constant in config/embed/{text,cmd,flag}/.
//
// Direction 2: every DescKey string value must exist as
// a key in the YAML files.
//
// This catches orphan YAML entries (dead text) and
// orphan DescKeys (constants pointing to nothing).
//
// See specs/ast-audit-tests.md for rationale.
func TestDescKeyYAMLLinkage(t *testing.T) {
	// Phase 1: collect all YAML keys.
	yamlKeys := collectYAMLKeys(t)

	// Phase 2: collect all DescKey string values from
	// Go constants.
	descKeys := collectDescKeyValues(t)

	// Phase 3: diff.
	var orphanYAML []string
	for key := range yamlKeys {
		if !descKeys[key] {
			orphanYAML = append(orphanYAML, key)
		}
	}
	sort.Strings(orphanYAML)

	var orphanDesc []string
	for key := range descKeys {
		if !yamlKeys[key] {
			orphanDesc = append(orphanDesc, key)
		}
	}
	sort.Strings(orphanDesc)

	if len(orphanYAML) > 0 {
		t.Errorf(
			"%d YAML keys with no DescKey:",
			len(orphanYAML),
		)
		for _, k := range orphanYAML {
			t.Errorf("  orphan YAML: %s", k)
		}
	}

	if len(orphanDesc) > 0 {
		t.Errorf(
			"%d DescKeys with no YAML entry:",
			len(orphanDesc),
		)
		for _, k := range orphanDesc {
			t.Errorf("  orphan DescKey: %s", k)
		}
	}
}

// collectYAMLKeys parses all YAML files under
// internal/assets/commands/ and returns a set of
// top-level keys.
func collectYAMLKeys(t *testing.T) map[string]bool {
	t.Helper()
	keys := make(map[string]bool)
	root := filepath.Join("..", "assets", "commands")

	walkErr := filepath.WalkDir(
		root,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			ext := filepath.Ext(d.Name())
			if ext != ".yaml" && ext != ".yml" {
				return nil
			}
			// examples.yaml has a different structure
			// (example categories, not DescKey entries).
			// examples.yaml and registry.yaml have
			// different structures validated by
			// TestExamplesYAMLLinkage and
			// TestRegistryYAMLLinkage respectively.
			if d.Name() == "examples.yaml" ||
				d.Name() == "registry.yaml" {
				return nil
			}

			data, readErr := os.ReadFile(path) //nolint:gosec
			if readErr != nil {
				return readErr
			}

			var doc map[string]any
			if unmarshalErr := yaml.Unmarshal(
				data, &doc,
			); unmarshalErr != nil {
				// Some YAML files may have nested
				// structure; try flat string map.
				var flat map[string]string
				if e2 := yaml.Unmarshal(
					data, &flat,
				); e2 != nil {
					return unmarshalErr
				}
				for k := range flat {
					keys[k] = true
				}
				return nil
			}

			collectNestedKeys(doc, "", keys)
			return nil
		},
	)
	if walkErr != nil {
		t.Fatalf("walk YAML: %v", walkErr)
	}

	return keys
}

// collectNestedKeys flattens nested YAML maps into
// dot-separated keys. Stops at the level where the
// value has a "short" or "long" field (the text entry).
func collectNestedKeys(
	m map[string]any, prefix string,
	keys map[string]bool,
) {
	for k, v := range m {
		full := k
		if prefix != "" {
			full = prefix + "." + k
		}

		sub, ok := v.(map[string]any)
		if !ok {
			// Leaf value: this key is a text entry.
			keys[full] = true
			continue
		}

		// If sub has "short" or "long", this is the
		// text entry level.
		if _, hasShort := sub["short"]; hasShort {
			keys[full] = true
			continue
		}
		if _, hasLong := sub["long"]; hasLong {
			keys[full] = true
			continue
		}

		// Otherwise recurse deeper.
		collectNestedKeys(sub, full, keys)
	}
}

// collectDescKeyValues extracts string values from all
// DescKey* constants in config/embed/{text,cmd,flag}/.
func collectDescKeyValues(t *testing.T) map[string]bool {
	t.Helper()
	pkgs := loadPackages(t)
	values := make(map[string]bool)

	for _, pkg := range pkgs {
		if !strings.Contains(pkg.PkgPath, "config/embed/text") &&
			!strings.Contains(pkg.PkgPath, "config/embed/cmd") &&
			!strings.Contains(pkg.PkgPath, "config/embed/flag") {
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

					for i, name := range vs.Names {
						if !strings.HasPrefix(
							name.Name, "DescKey",
						) {
							continue
						}

						if i >= len(vs.Values) {
							continue
						}

						lit, ok := vs.Values[i].(*ast.BasicLit)
						if !ok ||
							lit.Kind != token.STRING {
							continue
						}

						s, err := strconv.Unquote(
							lit.Value,
						)
						if err != nil {
							continue
						}

						values[s] = true
					}
				}
			}
		}
	}

	return values
}
