//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package node

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuilder_SinglePackage(t *testing.T) {
	orig, getErr := os.Getwd()
	if getErr != nil {
		t.Fatal(getErr)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	tmp := t.TempDir()
	if chdirErr := os.Chdir(tmp); chdirErr != nil {
		t.Fatal(chdirErr)
	}

	pkg := `{
		"name": "my-app",
		"dependencies": {
			"express": "^4.18.0",
			"lodash": "^4.17.0"
		},
		"devDependencies": {
			"jest": "^29.0.0"
		}
	}`
	pkgPath := filepath.Join(tmp, "package.json")
	if writeErr := os.WriteFile(
		pkgPath, []byte(pkg), 0o644,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	b := &Builder{}

	internal, buildErr := b.Build(false)
	if buildErr != nil {
		t.Fatalf("Build(false) failed: %v", buildErr)
	}
	if len(internal) != 0 {
		t.Errorf(
			"Build(false) for single package should be empty, got %v",
			internal,
		)
	}

	external, buildErr := b.Build(true)
	if buildErr != nil {
		t.Fatalf("Build(true) failed: %v", buildErr)
	}
	deps, ok := external["my-app"]
	if !ok {
		t.Fatal("Build(true) missing 'my-app' key")
	}
	if len(deps) != 3 {
		t.Errorf(
			"expected 3 deps, got %d: %v",
			len(deps), deps,
		)
	}
}

func TestBuilder_Workspaces(t *testing.T) {
	orig, getErr := os.Getwd()
	if getErr != nil {
		t.Fatal(getErr)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	tmp := t.TempDir()
	if chdirErr := os.Chdir(tmp); chdirErr != nil {
		t.Fatal(chdirErr)
	}

	root := `{
		"name": "monorepo",
		"workspaces": ["packages/*"]
	}`
	rootPath := filepath.Join(tmp, "package.json")
	if writeErr := os.WriteFile(
		rootPath, []byte(root), 0o644,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	pkgsDir := filepath.Join(tmp, "packages")

	pkgADir := filepath.Join(pkgsDir, "pkg-a")
	if mkErr := os.MkdirAll(pkgADir, 0o755); mkErr != nil {
		t.Fatal(mkErr)
	}
	pkgAJSON := `{"name":"@mono/pkg-a",` +
		`"dependencies":{"lodash":"^4.0.0"}}`
	if writeErr := os.WriteFile(
		filepath.Join(pkgADir, "package.json"),
		[]byte(pkgAJSON), 0o644,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	pkgBDir := filepath.Join(pkgsDir, "pkg-b")
	if mkErr := os.MkdirAll(pkgBDir, 0o755); mkErr != nil {
		t.Fatal(mkErr)
	}
	pkgBJSON := `{"name":"@mono/pkg-b",` +
		`"dependencies":{"@mono/pkg-a":"*",` +
		`"express":"^4.0.0"}}`
	if writeErr := os.WriteFile(
		filepath.Join(pkgBDir, "package.json"),
		[]byte(pkgBJSON), 0o644,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	b := &Builder{}

	internal, buildErr := b.Build(false)
	if buildErr != nil {
		t.Fatalf("Build(false) failed: %v", buildErr)
	}
	deps, ok := internal["@mono/pkg-b"]
	if !ok {
		t.Fatalf(
			"Build(false) missing '@mono/pkg-b' key, got %v",
			internal,
		)
	}
	if len(deps) != 1 || deps[0] != "@mono/pkg-a" {
		t.Errorf("expected [@mono/pkg-a], got %v", deps)
	}
	if _, ok := internal["@mono/pkg-a"]; ok {
		t.Error("@mono/pkg-a should not have internal deps")
	}
}

func TestBuilder_WorkspacesObject(t *testing.T) {
	orig, getErr := os.Getwd()
	if getErr != nil {
		t.Fatal(getErr)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	tmp := t.TempDir()
	if chdirErr := os.Chdir(tmp); chdirErr != nil {
		t.Fatal(chdirErr)
	}

	root := `{
		"name": "monorepo",
		"workspaces": {"packages": ["libs/*"]}
	}`
	rootPath := filepath.Join(tmp, "package.json")
	if writeErr := os.WriteFile(
		rootPath, []byte(root), 0o644,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	libDir := filepath.Join(tmp, "libs", "core")
	if mkErr := os.MkdirAll(libDir, 0o755); mkErr != nil {
		t.Fatal(mkErr)
	}
	coreJSON := `{"name":"@mono/core",` +
		`"dependencies":{"react":"^18.0.0"}}`
	if writeErr := os.WriteFile(
		filepath.Join(libDir, "package.json"),
		[]byte(coreJSON), 0o644,
	); writeErr != nil {
		t.Fatal(writeErr)
	}

	b := &Builder{}
	graph, buildErr := b.Build(true)
	if buildErr != nil {
		t.Fatalf("Build(true) failed: %v", buildErr)
	}
	if _, ok := graph["@mono/core"]; !ok {
		t.Errorf(
			"expected @mono/core in graph, got %v", graph,
		)
	}
}
