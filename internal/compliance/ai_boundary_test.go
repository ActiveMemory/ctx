//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compliance

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDeterministicCoreDoesNotDependOnAIBackends(t *testing.T) {
	root := projectRoot(t)
	guardedPrefixes := []string{
		"internal/cli/agent/",
		"internal/cli/status/",
		"internal/cli/hook/",
	}

	for _, path := range nonTestGoFiles(t, root) {
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			t.Fatalf("filepath.Rel: %v", relErr)
		}
		rel = filepath.ToSlash(rel)
		if !hasGuardedPrefix(rel, guardedPrefixes) {
			continue
		}

		assertNoBackendImport(t, path, rel)
		assertNoAIInvocation(t, path, rel)
	}
}

func hasGuardedPrefix(path string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func assertNoBackendImport(t *testing.T, path string, rel string) {
	t.Helper()
	fset := token.NewFileSet()
	file, parseErr := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
	if parseErr != nil {
		t.Fatalf("parser.ParseFile(%s): %v", rel, parseErr)
	}
	for _, imported := range file.Imports {
		importPath := strings.Trim(imported.Path.Value, "\"")
		if importPath == "github.com/ActiveMemory/ctx/internal/backend" {
			t.Fatalf("%s imports internal/backend", rel)
		}
	}
}

func assertNoAIInvocation(t *testing.T, path string, rel string) {
	t.Helper()
	data, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		t.Fatalf("os.ReadFile(%s): %v", rel, readErr)
	}
	if strings.Contains(string(data), "ctx ai") {
		t.Fatalf("%s invokes or documents ctx ai from deterministic core", rel)
	}
}
