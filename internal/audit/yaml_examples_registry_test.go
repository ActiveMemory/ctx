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
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestExamplesYAMLLinkage verifies that every key in
// examples.yaml matches a known command desc key constant
// from config/embed/cmd or is an add-subtype key
// (add.<entry-type> or add.default).
//
// Prevents orphan example entries that would never be
// looked up.
func TestExamplesYAMLLinkage(t *testing.T) {
	path := filepath.Join(
		"..", "assets", "commands", "examples.yaml",
	)
	data, readErr := os.ReadFile(path) //nolint:gosec
	if readErr != nil {
		t.Fatalf("read examples.yaml: %v", readErr)
	}

	var doc map[string]any
	if unmarshalErr := yaml.Unmarshal(
		data, &doc,
	); unmarshalErr != nil {
		t.Fatalf("parse examples.yaml: %v", unmarshalErr)
	}

	// Collect command desc key values and add-subtype
	// keys from config/embed/cmd and config/entry.
	descKeys := collectDescKeys(t)
	entryTypes := collectEntryTypes(t)

	// Build allowed set: desc keys + add.<entry-type>
	// + add.default
	allowed := make(map[string]bool, len(descKeys)+len(entryTypes)+1)
	for k := range descKeys {
		allowed[k] = true
	}
	for k := range entryTypes {
		allowed["add."+k] = true
	}
	allowed["add.default"] = true

	for key := range doc {
		if !allowed[key] {
			t.Errorf(
				"examples.yaml key %q has no "+
					"matching desc key or "+
					"add-subtype constant",
				key,
			)
		}
	}
}

// collectDescKeys returns the set of command description
// key string values from config/embed/cmd DescKey constants.
func collectDescKeys(t *testing.T) map[string]bool {
	t.Helper()
	pkgs := loadPackages(t)
	keys := make(map[string]bool)

	for _, pkg := range pkgs {
		if !strings.HasSuffix(
			pkg.PkgPath, "config/embed/cmd",
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

					// Only collect DescKey constants.
					for i, name := range vs.Names {
						if !strings.HasPrefix(
							name.Name, "DescKey",
						) {
							continue
						}
						if i >= len(vs.Values) {
							continue
						}
						lit, ok :=
							vs.Values[i].(*ast.BasicLit)
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

						keys[s] = true
					}
				}
			}
		}
	}

	return keys
}

// collectEntryTypes returns the set of entry type
// string values from config/entry constants.
func collectEntryTypes(t *testing.T) map[string]bool {
	t.Helper()
	pkgs := loadPackages(t)
	types := make(map[string]bool)

	for _, pkg := range pkgs {
		if !strings.HasSuffix(
			pkg.PkgPath, "config/entry",
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

					for i, val := range vs.Values {
						lit, ok :=
							val.(*ast.BasicLit)
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

						_ = i
						types[s] = true
					}
				}
			}
		}
	}

	return types
}

// registryEntry matches the YAML structure of
// registry.yaml entries.
type registryEntry struct {
	Hook    string `yaml:"hook"`
	Variant string `yaml:"variant"`
}

// TestRegistryYAMLLinkage verifies that every
// hook/variant pair in registry.yaml:
//
//  1. Has a hook name matching a config/hook constant
//  2. Has a corresponding .txt template file at
//     internal/assets/hooks/messages/<hook>/<variant>.txt
//
// Prevents orphan registry entries pointing to
// non-existent hooks or missing template files.
func TestRegistryYAMLLinkage(t *testing.T) {
	path := filepath.Join(
		"..", "assets", "hooks",
		"messages", "registry.yaml",
	)
	data, readErr := os.ReadFile(path) //nolint:gosec
	if readErr != nil {
		t.Fatalf("read registry.yaml: %v", readErr)
	}

	var entries []registryEntry
	if unmarshalErr := yaml.Unmarshal(
		data, &entries,
	); unmarshalErr != nil {
		t.Fatalf("parse registry.yaml: %v", unmarshalErr)
	}

	// Collect hook name constants from config/hook.
	hookNames := collectHookNames(t)

	// Template root.
	tplRoot := filepath.Join(
		"..", "assets", "hooks", "messages",
	)

	for _, e := range entries {
		// Check hook name exists as a constant.
		if !hookNames[e.Hook] {
			t.Errorf(
				"registry.yaml hook %q has no "+
					"matching config/hook constant",
				e.Hook,
			)
		}

		// Check template file exists.
		tplPath := filepath.Join(
			tplRoot, e.Hook, e.Variant+".txt",
		)
		if _, statErr := os.Stat(tplPath); os.IsNotExist(statErr) {
			t.Errorf(
				"registry.yaml %s/%s: "+
					"template file missing: %s",
				e.Hook, e.Variant, tplPath,
			)
		}
	}
}

// collectHookNames returns the set of hook name
// string values from config/hook constants.
func collectHookNames(t *testing.T) map[string]bool {
	t.Helper()
	pkgs := loadPackages(t)
	names := make(map[string]bool)

	for _, pkg := range pkgs {
		if !strings.HasSuffix(
			pkg.PkgPath, "config/hook",
		) {
			continue
		}

		for _, file := range pkg.Syntax {
			fpath := pkg.Fset.Position(file.Pos()).Filename
			// Only check hook.go, not notify.go etc.
			if !strings.HasSuffix(fpath, "hook.go") {
				continue
			}

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
						lit, ok :=
							val.(*ast.BasicLit)
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

						names[s] = true
					}
				}
			}
		}
	}

	return names
}
