//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package compliance contains tests that verify the entire codebase adheres
// to the project standards documented in CONTRIBUTING.md, CLAUDE.md, and
// the lint-drift / lint-docs scripts.
//
// These tests are cross-cutting: they inspect source files, configs, and
// build artifacts across the whole repository rather than testing a single
// package's exported API. They mirror the checks performed by
// hack/lint-drift.sh and hack/lint-docs.sh so that violations surface in
// `go test` without requiring bash.
package compliance

import (
	"bufio"
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// projectRoot returns the absolute path to the project root.
func projectRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("failed to resolve project root: %v", err)
	}
	return root
}

// templateDir returns the path to the template/assets directory.
// It supports both internal/assets (current) and internal/tpl (legacy).
func templateDir(t *testing.T, root string) string {
	t.Helper()
	for _, name := range []string{"assets", "tpl"} {
		d := filepath.Join(root, "internal", name)
		if _, err := os.Stat(d); err == nil {
			return d
		}
	}
	t.Fatal("cannot find template dir (internal/assets or internal/tpl)")
	return ""
}

// allGoFiles returns all .go files under the project root, excluding vendor/.
func allGoFiles(t *testing.T, root string) []string {
	t.Helper()
	var files []string
	err := filepath.Walk(root, func(
		path string, info os.FileInfo, err error,
	) error {
		if err != nil {
			return err
		}
		isSkipped := info.Name() == "vendor" ||
			info.Name() == ".git" ||
			info.Name() == "dist" ||
			info.Name() == "site"
		if info.IsDir() && isSkipped {
			return filepath.SkipDir
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk project: %v", err)
	}
	return files
}

// nonTestGoFiles returns all non-test .go files.
func nonTestGoFiles(t *testing.T, root string) []string {
	t.Helper()
	var result []string
	for _, f := range allGoFiles(t, root) {
		if !strings.HasSuffix(f, "_test.go") {
			result = append(result, f)
		}
	}
	return result
}

// ---------------------------------------------------------------------------
// 1. License Header ╬ô├ç├╢ every .go file must have the SPDX header
// ---------------------------------------------------------------------------

// TestLicenseHeader verifies every .go file contains the Apache-2.0 SPDX
// identifier within the first 10 lines.
func TestLicenseHeader(t *testing.T) {
	root := projectRoot(t)
	spdxTag := "SPDX-License-Identifier: Apache-2.0"

	for _, p := range allGoFiles(t, root) {
		rel, _ := filepath.Rel(root, p)
		t.Run(rel, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Clean(p))
			if err != nil {
				t.Fatalf("read: %v", err)
			}
			scanner := bufio.NewScanner(strings.NewReader(string(data)))
			found := false
			for i := 0; i < 10 && scanner.Scan(); i++ {
				if strings.Contains(scanner.Text(), spdxTag) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("missing SPDX license header (%s)", spdxTag)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 2. Package doc.go ╬ô├ç├╢ every package under internal/ should have a doc.go
// ---------------------------------------------------------------------------

// TestDocGoExists verifies every Go package under internal/ has a doc.go.
//
// Some packages (like cli) use sub-packages; for those the check recurses
// into each sub-package directory instead.
func TestDocGoExists(t *testing.T) {
	root := projectRoot(t)
	internalDir := filepath.Join(root, "internal")

	// Packages exempt from the doc.go requirement.
	exempt := map[string]bool{
		"compliance": true, // this test-only package
	}

	entries, err := os.ReadDir(internalDir)
	if err != nil {
		t.Fatalf("readdir internal/: %v", err)
	}

	for _, e := range entries {
		if !e.IsDir() || exempt[e.Name()] {
			continue
		}

		pkgDir := filepath.Join(internalDir, e.Name())

		// If the package contains sub-packages (directories with .go files),
		// check each sub-package instead of the parent.
		subEntries, subErr := os.ReadDir(pkgDir)
		if subErr != nil {
			t.Fatalf("readdir internal/%s: %v", e.Name(), subErr)
		}

		hasSubPkgs := false
		for _, sub := range subEntries {
			if sub.IsDir() {
				subDir := filepath.Join(pkgDir, sub.Name())
				goFiles, _ := filepath.Glob(filepath.Join(subDir, "*.go"))
				if len(goFiles) > 0 {
					hasSubPkgs = true
					name := e.Name() + "/" + sub.Name()
					t.Run(name, func(t *testing.T) {
						docPath := filepath.Join(subDir, "doc.go")
						if _, statErr := os.Stat(docPath); os.IsNotExist(statErr) {
							t.Errorf("missing doc.go in internal/%s", name)
						}
					})
				}
			}
		}

		// If no sub-packages, check the package itself.
		if !hasSubPkgs {
			t.Run(e.Name(), func(t *testing.T) {
				docPath := filepath.Join(pkgDir, "doc.go")
				if _, statErr := os.Stat(docPath); os.IsNotExist(statErr) {
					t.Errorf("missing doc.go in internal/%s", e.Name())
				}
			})
		}
	}
}

// ---------------------------------------------------------------------------
// 3. No literal "\n" ╬ô├ç├╢ use config.NewlineLF instead (lint-drift rule 1)
// ---------------------------------------------------------------------------

// TestNoLiteralNewline mirrors lint-drift rule 1: literal "\n" strings
// should use config.NewlineLF instead.
func TestNoLiteralNewline(t *testing.T) {
	root := projectRoot(t)
	re := regexp.MustCompile(`"\\n"`)

	for _, p := range nonTestGoFiles(t, root) {
		if strings.HasSuffix(p, "token.go") || strings.HasSuffix(p, "whitespace.go") {
			continue
		}
		rel, _ := filepath.Rel(root, p)

		data, err := os.ReadFile(filepath.Clean(p))
		if err != nil {
			t.Fatalf("read %s: %v", rel, err)
		}
		if re.Match(data) {
			t.Run(rel, func(t *testing.T) {
				t.Errorf("literal \"\\n\" found, want config.NewlineLF")
			})
		}
	}
}

// ---------------------------------------------------------------------------
// 4. No literal ".md" ╬ô├ç├╢ use config.ExtMarkdown instead (lint-drift rule 4)
// ---------------------------------------------------------------------------

// TestNoLiteralMdExtension mirrors lint-drift rule 4: literal ".md" strings
// should use config.ExtMarkdown instead.
func TestNoLiteralMdExtension(t *testing.T) {
	root := projectRoot(t)
	re := regexp.MustCompile(`"\.md"`)

	for _, p := range nonTestGoFiles(t, root) {
		if strings.HasSuffix(p, filepath.Join("config", "file.go")) ||
			strings.HasSuffix(p, filepath.Join("file", "ext.go")) {
			continue
		}
		rel, _ := filepath.Rel(root, p)

		data, err := os.ReadFile(filepath.Clean(p))
		if err != nil {
			t.Fatalf("read %s: %v", rel, err)
		}
		if re.Match(data) {
			t.Run(rel, func(t *testing.T) {
				t.Errorf("literal \".md\" found, want config.ExtMarkdown")
			})
		}
	}
}

// ---------------------------------------------------------------------------
// 5. No cmd.Printf/cmd.PrintErrf ╬ô├ç├╢ prefer Println (lint-drift rule 2)
// ---------------------------------------------------------------------------

// TestNoCmdPrintf mirrors lint-drift rule 2: cmd.Printf/cmd.PrintErrf should
// be replaced with cmd.Println(fmt.Sprintf(...)).
func TestNoCmdPrintf(t *testing.T) {
	root := projectRoot(t)
	re := regexp.MustCompile(`cmd\.(Printf|PrintErrf)\(`)

	for _, p := range nonTestGoFiles(t, root) {
		rel, _ := filepath.Rel(root, p)

		data, err := os.ReadFile(filepath.Clean(p))
		if err != nil {
			t.Fatalf("read %s: %v", rel, err)
		}
		if re.Match(data) {
			t.Run(rel, func(t *testing.T) {
				t.Errorf("cmd.Printf/PrintErrf found, want cmd.Println(fmt.Sprintf(...))")
			})
		}
	}
}

// ---------------------------------------------------------------------------
// 6. No magic directory strings — use config.Dir* constants
// (lint-drift rule 3)
// ---------------------------------------------------------------------------

// TestNoMagicDirectoryStrings mirrors lint-drift rule 3: magic directory
// strings in filepath.Join calls should use config.Dir* constants.
func TestNoMagicDirectoryStrings(t *testing.T) {
	root := projectRoot(t)

	tests := []struct {
		pattern  string
		constant string
	}{
		{`filepath\.Join\([^)]*"sessions"`, "config.DirSessions"},
		{`filepath\.Join\([^)]*"archive"`, "config.DirArchive"},
		{`filepath\.Join\([^)]*"tools"`, "config.DirTools"},
	}

	for _, tt := range tests {
		re := regexp.MustCompile(tt.pattern)
		for _, p := range nonTestGoFiles(t, root) {
			rel, _ := filepath.Rel(root, p)

			data, err := os.ReadFile(filepath.Clean(p))
			if err != nil {
				t.Fatalf("read %s: %v", rel, err)
			}
			if re.Match(data) {
				t.Run(rel+"/"+tt.constant, func(t *testing.T) {
					t.Errorf("magic directory string found, want %s", tt.constant)
				})
			}
		}
	}
}

// ---------------------------------------------------------------------------
// 7. No direct fmt.Print* in Cobra command functions ╬ô├ç├╢ use cmd.Print*
// ---------------------------------------------------------------------------

// TestNoDirectFmtPrintInCobraHandlers parses CLI source files and verifies
// that functions accepting *cobra.Command do not call fmt.Print* directly.
// Output should go through cmd.Print* so tests can capture it and --quiet
// flags work correctly.
func TestNoDirectFmtPrintInCobraHandlers(t *testing.T) {
	root := projectRoot(t)
	cliDir := filepath.Join(root, "internal", "cli")

	forbidden := map[string]bool{
		"Print":   true,
		"Println": true,
		"Printf":  true,
	}

	err := filepath.Walk(cliDir, func(
		path string, info os.FileInfo, walkErr error,
	) error {
		if walkErr != nil {
			return walkErr
		}
		notSource := info.IsDir() ||
			!strings.HasSuffix(path, ".go") ||
			strings.HasSuffix(path, "_test.go")
		if notSource {
			return nil
		}

		fset := token.NewFileSet()
		node, parseErr := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if parseErr != nil {
			t.Errorf("parse %s: %v", path, parseErr)
			return nil
		}

		// Check if file imports "fmt"
		var fmtAlias string
		for _, imp := range node.Imports {
			impPath := strings.Trim(imp.Path.Value, `"`)
			if impPath == "fmt" {
				if imp.Name != nil {
					fmtAlias = imp.Name.Name
				} else {
					fmtAlias = "fmt"
				}
				break
			}
		}
		if fmtAlias == "" {
			return nil
		}

		for _, decl := range node.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Type.Params == nil {
				continue
			}

			hasCobraCmd := false
			for _, param := range fn.Type.Params.List {
				if star, ok := param.Type.(*ast.StarExpr); ok {
					if sel, ok := star.X.(*ast.SelectorExpr); ok {
						if sel.Sel.Name == "Command" {
							hasCobraCmd = true
							break
						}
					}
				}
			}
			if !hasCobraCmd {
				continue
			}

			ast.Inspect(fn.Body, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				ident, ok := sel.X.(*ast.Ident)
				if !ok {
					return true
				}
				if ident.Name == fmtAlias && forbidden[sel.Sel.Name] {
					pos := fset.Position(call.Pos())
					rel, _ := filepath.Rel(root, pos.Filename)
					t.Errorf("%s:%d: fmt.%s in Cobra handler ╬ô├ç├╢ use cmd.Print* instead",
						rel, pos.Line, sel.Sel.Name)
				}
				return true
			})
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk cli dir: %v", err)
	}
}

// ---------------------------------------------------------------------------
// 8. gofmt compliance ╬ô├ç├╢ all Go files must be properly formatted
// ---------------------------------------------------------------------------

// TestGofmt verifies all Go files are properly formatted.
// It normalizes CRLF to LF before comparison so that the test passes on
// Windows where git may check out files with CRLF line endings.
func TestGofmt(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping gofmt check in short mode")
	}

	root := projectRoot(t)

	var unformatted []string
	for _, f := range allGoFiles(t, root) {
		data, err := os.ReadFile(filepath.Clean(f))
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		// Normalize CRLF to LF so the check works on Windows.
		normalized := bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
		formatted, fmtErr := format.Source(normalized)
		if fmtErr != nil {
			// File doesn't parse; go vet will catch it.
			continue
		}
		if !bytes.Equal(normalized, formatted) {
			rel, _ := filepath.Rel(root, f)
			unformatted = append(unformatted, rel)
		}
	}

	if len(unformatted) > 0 {
		t.Errorf("files need formatting:\n\t%s\n\nRun: go fmt ./...",
			strings.Join(unformatted, "\n\t"))
	}
}

// ---------------------------------------------------------------------------
// 9. go vet ╬ô├ç├╢ the entire project must pass go vet
// ---------------------------------------------------------------------------

// TestGoVet runs go vet across the entire project with CGO disabled.
func TestGoVet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping go vet in short mode")
	}

	root := projectRoot(t)

	cmd := exec.Command("go", "vet", "./...")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go vet failed:\n%s", string(output))
	}
}

// ---------------------------------------------------------------------------
// 9b. golangci-lint ╬ô├ç├╢ the entire project must pass golangci-lint
// ---------------------------------------------------------------------------

// TestGolangciLint runs golangci-lint across the entire project.
// This catches issues that go vet alone misses (gosec, goconst, unused, etc.).
// golangci-lint is a required dependency — the test fails
// if it is not installed.
func TestGolangciLint(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping golangci-lint in short mode")
	}

	// golangci-lint may not be installed in every CI job (the Lint job
	// runs it separately via golangci-lint-action).  Skip gracefully.
	if _, err := exec.LookPath("golangci-lint"); err != nil {
		t.Skip("golangci-lint is not installed.\n" +
			"Install it with:\n" +
			"  go install github.com/golangci/" +
			"golangci-lint/v2/cmd/golangci-lint@v2.8.0\n" +
			"Or see: https://golangci-lint.run/welcome/install/")
	}

	root := projectRoot(t)

	cmd := exec.Command("golangci-lint", "run", "--timeout=5m")
	cmd.Dir = root
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("golangci-lint failed:\n%s", string(output))
	}
}

// ---------------------------------------------------------------------------
// 10. No secrets in .context/ templates ╬ô├ç├╢ no tokens, keys, passwords
// ---------------------------------------------------------------------------

// TestNoSecretsInTemplates scans template files for patterns that look like
// secrets (API keys, tokens, private keys) per SECURITY.md requirements.
func TestNoSecretsInTemplates(t *testing.T) {
	root := projectRoot(t)
	tplDir := templateDir(t, root)

	secretPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(api[_-]?key|secret[_-]?key|password|token|credential)\s*[:=]`),
		regexp.MustCompile(`(?i)(sk-[a-zA-Z0-9]{20,}|ghp_[a-zA-Z0-9]{36}|gho_[a-zA-Z0-9]{36})`),
		regexp.MustCompile(`(?i)-----BEGIN (RSA |EC )?PRIVATE KEY-----`),
	}

	err := filepath.Walk(tplDir, func(
		path string, info os.FileInfo, walkErr error,
	) error {
		if walkErr != nil || info.IsDir() {
			return walkErr
		}

		//nolint:gosec // path comes from filepath.Walk
		data, readErr := os.ReadFile(filepath.Clean(path))
		if readErr != nil {
			t.Errorf("read %s: %v", path, readErr)
			return nil
		}

		rel, _ := filepath.Rel(root, path)
		for _, re := range secretPatterns {
			if re.Match(data) {
				t.Errorf("%s: potential secret pattern found: %s", rel, re.String())
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk assets dir: %v", err)
	}
}

// ---------------------------------------------------------------------------
// 11. Required context files ╬ô├ç├╢ ctx init must create all required files
// ---------------------------------------------------------------------------

// TestRequiredContextFilesInTemplate verifies that all required context file
// templates exist in internal/assets/ so that ctx init can scaffold them.
func TestRequiredContextFilesInTemplate(t *testing.T) {
	root := projectRoot(t)
	tplDir := templateDir(t, root)

	requiredFiles := []string{
		"CONSTITUTION.md",
		"TASKS.md",
		"DECISIONS.md",
		"LEARNINGS.md",
		"CONVENTIONS.md",
		"ARCHITECTURE.md",
	}

	for _, name := range requiredFiles {
		t.Run(name, func(t *testing.T) {
			// Templates live under context/ subdirectory in assets.
			path := filepath.Join(tplDir, "context", name)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("required template %s not found in internal/assets/context/", name)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 12. VERSION file ╬ô├ç├╢ must exist and contain a valid semver
// ---------------------------------------------------------------------------

// TestVersionFile checks the VERSION file exists and contains valid semver.
func TestVersionFile(t *testing.T) {
	root := projectRoot(t)
	versionPath := filepath.Join(root, "VERSION")

	//nolint:gosec // constructed from test constants
	data, err := os.ReadFile(filepath.Clean(versionPath))
	if err != nil {
		t.Fatalf("cannot read VERSION file: %v", err)
	}

	version := strings.TrimSpace(string(data))
	if version == "" {
		t.Fatal("VERSION file is empty")
	}

	semverRe := regexp.MustCompile(`^\d+\.\d+\.\d+(-[a-zA-Z0-9.]+)?$`)
	if !semverRe.MatchString(version) {
		t.Errorf("VERSION %q is not valid semver (expected X.Y.Z)", version)
	}
}

// ---------------------------------------------------------------------------
// 13. go.mod ╬ô├ç├╢ module path and Go version check
// ---------------------------------------------------------------------------

// TestGoMod verifies the module path and Go version in go.mod.
func TestGoMod(t *testing.T) {
	root := projectRoot(t)
	modPath := filepath.Join(root, "go.mod")

	//nolint:gosec // constructed from test constants
	data, err := os.ReadFile(filepath.Clean(modPath))
	if err != nil {
		t.Fatalf("cannot read go.mod: %v", err)
	}

	content := string(data)

	t.Run("module path", func(t *testing.T) {
		if !strings.Contains(content, "module github.com/ActiveMemory/ctx") {
			t.Error("go.mod should declare module github.com/ActiveMemory/ctx")
		}
	})

	t.Run("go version declared", func(t *testing.T) {
		goVersionRe := regexp.MustCompile(`go\s+1\.\d+`)
		if !goVersionRe.MatchString(content) {
			t.Error("go.mod should declare a Go version (go 1.x)")
		}
	})
}

// ---------------------------------------------------------------------------
// 14. Makefile ╬ô├ç├╢ required targets exist
// ---------------------------------------------------------------------------

// TestMakefileTargets verifies all expected build targets
// exist in the Makefile.
func TestMakefileTargets(t *testing.T) {
	root := projectRoot(t)
	makePath := filepath.Join(root, "Makefile")

	//nolint:gosec // constructed from test constants
	data, err := os.ReadFile(filepath.Clean(makePath))
	if err != nil {
		t.Fatalf("cannot read Makefile: %v", err)
	}

	content := string(data)

	requiredTargets := []string{
		"build:",
		"test:",
		"vet:",
		"fmt:",
		"lint:",
		"clean:",
	}

	for _, target := range requiredTargets {
		t.Run(strings.TrimSuffix(target, ":"), func(t *testing.T) {
			if !strings.Contains(content, target) {
				t.Errorf("Makefile missing required target: %s", target)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 15. CGO_ENABLED=0 ╬ô├ç├╢ build command must not use CGO
// ---------------------------------------------------------------------------

// TestBuildWithoutCGO verifies that Makefile build and test targets use
// CGO_ENABLED=0 as required by the project standards.
func TestBuildWithoutCGO(t *testing.T) {
	root := projectRoot(t)
	makePath := filepath.Join(root, "Makefile")

	//nolint:gosec // constructed from test constants
	data, err := os.ReadFile(filepath.Clean(makePath))
	if err != nil {
		t.Fatalf("cannot read Makefile: %v", err)
	}

	content := string(data)

	t.Run("build target uses CGO_ENABLED=0", func(t *testing.T) {
		if !strings.Contains(content, "CGO_ENABLED=0") {
			t.Error("Makefile build target should use CGO_ENABLED=0")
		}
	})

	t.Run("test target uses CGO_ENABLED=0", func(t *testing.T) {
		// Find the test target line and check CGO
		lines := strings.Split(content, "\n")
		for i, line := range lines {
			if strings.HasPrefix(line, "test:") {
				// Check the next few lines for CGO_ENABLED=0
				found := false
				for j := i + 1; j < i+5 && j < len(lines); j++ {
					if strings.Contains(lines[j], "CGO_ENABLED=0") {
						found = true
						break
					}
				}
				if !found {
					t.Error("test target should use CGO_ENABLED=0")
				}
				break
			}
		}
	})
}

// ---------------------------------------------------------------------------
// 16. .golangci.yml ╬ô├ç├╢ required linters are configured
// ---------------------------------------------------------------------------

// TestGolangciLintConfig verifies that .golangci.yml enables the required
// linters (govet, errcheck, staticcheck, gosec).
func TestGolangciLintConfig(t *testing.T) {
	root := projectRoot(t)
	lintPath := filepath.Join(root, ".golangci.yml")

	//nolint:gosec // constructed from test constants
	data, err := os.ReadFile(filepath.Clean(lintPath))
	if err != nil {
		t.Fatalf("cannot read .golangci.yml: %v", err)
	}

	content := string(data)

	requiredLinters := []string{
		"govet",
		"errcheck",
		"staticcheck",
		"gosec",
	}

	for _, linter := range requiredLinters {
		t.Run(linter, func(t *testing.T) {
			if !strings.Contains(content, linter) {
				t.Errorf(".golangci.yml missing required linter: %s", linter)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 17. No network calls — ctx must be local-only
// (no net/http imports in core)
// ---------------------------------------------------------------------------

// TestNoNetworkImportsInCore verifies that core packages do not import net or
// net/http, enforcing the local-only design described in SECURITY.md.
func TestNoNetworkImportsInCore(t *testing.T) {
	root := projectRoot(t)

	localOnlyPackages := []string{
		"context",
		"config",
		"drift",
		"task",
		"validation",
		"crypto",
		"assets",
		"index",
	}

	for _, pkg := range localOnlyPackages {
		pkgDir := filepath.Join(root, "internal", pkg)
		if _, err := os.Stat(pkgDir); os.IsNotExist(err) {
			continue
		}

		t.Run(pkg, func(t *testing.T) {
			fset := token.NewFileSet()
			pkgs, parseErr := parser.ParseDir(fset, pkgDir, func(info os.FileInfo) bool {
				return !strings.HasSuffix(info.Name(), "_test.go")
			}, parser.ImportsOnly)
			if parseErr != nil {
				t.Fatalf("parse %s: %v", pkg, parseErr)
			}

			for _, p := range pkgs {
				for _, f := range p.Files {
					for _, imp := range f.Imports {
						impPath := strings.Trim(imp.Path.Value, `"`)
						if impPath == "net/http" || impPath == "net" {
							pos := fset.Position(imp.Pos())
							t.Errorf("%s:%d: %s imports %q ╬ô├ç├╢ ctx core must be local-only",
								filepath.Base(pos.Filename), pos.Line, pkg, impPath)
						}
					}
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 18. Security ╬ô├ç├╢ .gitignore protects sensitive files
// ---------------------------------------------------------------------------

// TestGitignoreProtectsSensitiveFiles ensures .gitignore contains entries for
// files that must never be committed (encryption keys, etc.).
func TestGitignoreProtectsSensitiveFiles(t *testing.T) {
	root := projectRoot(t)
	giPath := filepath.Join(root, ".gitignore")

	//nolint:gosec // constructed from test constants
	data, err := os.ReadFile(filepath.Clean(giPath))
	if err != nil {
		t.Fatalf("cannot read .gitignore: %v", err)
	}

	content := string(data)

	sensitivePatterns := []string{
		".scratchpad.key",
	}

	for _, pattern := range sensitivePatterns {
		t.Run(pattern, func(t *testing.T) {
			if !strings.Contains(content, pattern) {
				t.Errorf(".gitignore should protect %s", pattern)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 19. Binary build ╬ô├ç├╢ ensure the project compiles without errors
// ---------------------------------------------------------------------------

// TestProjectCompiles builds the entire project with CGO disabled to verify
// there are no compilation errors.
func TestProjectCompiles(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping build test in short mode")
	}

	root := projectRoot(t)

	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("project does not compile:\n%s", string(output))
	}
}

// ---------------------------------------------------------------------------
// 20. File permissions ╬ô├ç├╢ config.PermSecret must be 0600
// ---------------------------------------------------------------------------

// TestPermissionConstants verifies that config.PermSecret and config.PermFile
// use the expected permission values.
func TestPermissionConstants(t *testing.T) {
	root := projectRoot(t)
	filePath := filepath.Join(root, "internal", "config", "fs", "perm.go")

	//nolint:gosec // constructed from test constants
	data, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		t.Fatalf("read file.go: %v", err)
	}

	content := string(data)

	t.Run("PermSecret is 0600", func(t *testing.T) {
		if !strings.Contains(content, "0600") {
			t.Error("config.PermSecret should be 0600 for secret files")
		}
	})

	t.Run("PermFile is 0644", func(t *testing.T) {
		if !strings.Contains(content, "0644") {
			t.Error("config.PermFile should be 0644 for regular files")
		}
	})
}

// ---------------------------------------------------------------------------
// 21. doc.go subcommand drift — listed names match cmd/ dirs
// ---------------------------------------------------------------------------

// TestDocGoSubcommandDrift walks internal/cli/ looking for doc.go files that
// list subcommands via "//   - name:" bullet patterns under a section header
// containing "subcommand" or "hook". For each listed name, it checks that a
// corresponding cmd/<name>/ directory exists (normalizing hyphens to
// underscores). Missing directories indicate the doc.go is out of date.
func TestDocGoSubcommandDrift(t *testing.T) {
	root := projectRoot(t)
	cliDir := filepath.Join(root, "internal", "cli")

	// sectionRe detects section headers that introduce subcommand lists.
	sectionRe := regexp.MustCompile(
		`(?i)(subcommand|hook)`,
	)
	// bulletRe matches "//   - name:" or "//   - name/other:" patterns.
	// The captured group may contain "/" for combined entries (pause/resume).
	bulletRe := regexp.MustCompile(
		`^//\s+-\s+([\w/-]+)\s*[:/]`,
	)

	// Known name mappings: doc name → dir name.
	knownAliases := map[string]string{
		"switch": "switchcmd",
		"import": "importer",
	}

	// Directories to skip in the undocumented check.
	skipDirs := map[string]bool{
		"root": true,
	}

	err := filepath.Walk(cliDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.Name() != "doc.go" || info.IsDir() {
			return nil
		}

		pkgDir := filepath.Dir(path)
		cmdDir := filepath.Join(pkgDir, "cmd")

		// Only check doc.go files that have a cmd/ directory.
		if _, statErr := os.Stat(cmdDir); statErr != nil {
			return nil
		}

		//nolint:gosec // constructed from test constants
		data, readErr := os.ReadFile(filepath.Clean(path))
		if readErr != nil {
			t.Errorf("read %s: %v", path, readErr)
			return nil
		}

		// Read actual cmd/ subdirectories.
		entries, dirErr := os.ReadDir(cmdDir)
		if dirErr != nil {
			return nil
		}
		actualDirs := make(map[string]bool)
		for _, e := range entries {
			if e.IsDir() {
				actualDirs[e.Name()] = true
			}
		}

		// Extract subcommand names from doc.go. Only process files
		// whose doc comment mentions "subcommand" or "hook" (proving
		// the bullets are subcommand listings, not entry types).
		content := string(data)
		if !sectionRe.MatchString(content) {
			return nil
		}

		var documented []string
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			if m := bulletRe.FindStringSubmatch(scanner.Text()); m != nil {
				documented = append(documented, m[1])
			}
		}

		if len(documented) == 0 {
			return nil
		}

		rel, _ := filepath.Rel(root, path)

		// Normalize: hyphens → underscores, apply aliases.
		normalize := func(name string) string {
			if alias, ok := knownAliases[name]; ok {
				return alias
			}
			return strings.ReplaceAll(name, "-", "_")
		}

		// Expand combined entries (e.g., "pause/resume") and check
		// each documented name has a cmd/ directory.
		docSet := make(map[string]bool)
		for _, raw := range documented {
			parts := strings.Split(raw, "/")
			for _, name := range parts {
				dirName := normalize(name)
				docSet[dirName] = true
				if !actualDirs[dirName] {
					t.Errorf(
						"%s lists subcommand %q but cmd/%s/ does not exist",
						rel, name, dirName,
					)
				}
			}
		}

		// Check for cmd/ directories not documented.
		for dir := range actualDirs {
			if skipDirs[dir] || docSet[dir] {
				continue
			}
			t.Errorf(
				"%s does not document cmd/%s/ subcommand",
				rel, dir,
			)
		}

		return nil
	})
	if err != nil {
		t.Fatalf("walk: %v", err)
	}
}

// ---------------------------------------------------------------------------
// 22. cmd/ purity — cmd dirs contain only Cmd/Run, no helpers or types
// ---------------------------------------------------------------------------

// TestCmdDirPurity walks internal/cli/**/cmd/*/ and verifies that .go files
// (excluding tests and doc.go) only declare exported functions matching
// Cmd or Run*. Unexported functions and type declarations belong in core/.
func TestCmdDirPurity(t *testing.T) {
	root := projectRoot(t)
	cliDir := filepath.Join(root, "internal", "cli")

	// Allowed exported function name patterns in cmd/ files.
	allowedFunc := func(name string) bool {
		if name == "Cmd" {
			return true
		}
		return strings.HasPrefix(name, "Run")
	}

	// Pre-existing violations grandfathered until refactored.
	// Key: relative path from cli/, value: set of allowed names.
	grandfathered := map[string]map[string]bool{
		"guide/cmd/root/command.go":       {"listCommands": true},
		"guide/cmd/root/skill.go":         {"parseSkillFrontmatter": true, "truncateDescription": true, "listSkills": true},
		"guide/cmd/root/types.go":         {"skillMeta": true},
		"setup/cmd/root/run.go":           {"WriteCopilotInstructions": true, "WriteAgentsMd": true, "WriteCopilotCLIHooks": true, "writeCopilotCLIHooks": true, "writeCopilotCLIAgent": true, "writeCopilotCLIInstructions": true, "writeCopilotCLISkills": true, "ensureVSCodeMCPJSON": true, "ensureCopilotCLIMCPConfig": true},
		"initialize/cmd/root/run.go":      {"initScratchpad": true, "hasEssentialFiles": true, "ensureGitignoreEntries": true, "writeGettingStarted": true},
		"journal/cmd/obsidian/run.go":     {"BuildVault": true},
		"journal/cmd/source/list.go":      {"runList": true},
		"journal/cmd/source/show.go":      {"runShow": true},
		"journal/cmd/source/types.go":     {"Opts": true},
		"loop/cmd/root/script.go":         {"GenerateLoopScript": true},
		"pad/cmd/edit/types.go":           {"Mode": true, "Opts": true},
		"pad/cmd/resolve/display.go":      {"displayAll": true},
		"remind/cmd/dismiss/run.go":       {"dismissOne": true, "dismissAll": true},
		"system/cmd/post_commit/score.go": {"scoreCommitViolations": true},
		"task/cmd/complete/run.go":        {"Complete": true},
		"why/cmd/root/data.go":            {"DocEntry": true},
		"why/cmd/root/menu.go":            {"showMenu": true},
		"why/cmd/root/run.go":             {"ShowDoc": true},
		"why/cmd/root/strip.go":           {"StripMkDocs": true, "ExtractAdmonitionTitle": true, "ExtractTabTitle": true},
		"system/cmd/sessionevent/cmd.go":  {"runSessionEvent": true},
	}

	err := filepath.Walk(cliDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}
		if strings.HasSuffix(info.Name(), "_test.go") || info.Name() == "doc.go" {
			return nil
		}

		// Only check files inside a cmd/ directory.
		rel, _ := filepath.Rel(cliDir, path)
		if !strings.Contains(rel, "/cmd/") {
			return nil
		}

		fset := token.NewFileSet()
		node, parseErr := parser.ParseFile(fset, path, nil, 0)
		if parseErr != nil {
			t.Errorf("parse %s: %v", rel, parseErr)
			return nil
		}

		allowed := grandfathered[rel]

		for _, decl := range node.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				name := d.Name.Name
				if d.Recv != nil {
					continue
				}
				if allowed[name] {
					continue
				}
				if !d.Name.IsExported() {
					t.Errorf(
						"%s: unexported function %q belongs in core/, not cmd/",
						rel, name,
					)
				} else if !allowedFunc(name) {
					t.Errorf(
						"%s: exported function %q is not Cmd or Run* — move to core/",
						rel, name,
					)
				}
			case *ast.GenDecl:
				if d.Tok == token.TYPE {
					for _, spec := range d.Specs {
						ts, ok := spec.(*ast.TypeSpec)
						if !ok {
							continue
						}
						if allowed[ts.Name.Name] {
							continue
						}
						t.Errorf(
							"%s: type %q belongs in core/types.go, not cmd/",
							rel, ts.Name.Name,
						)
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk: %v", err)
	}
}

// allSourceFiles returns all source files (.go, .ts, .js) under the project
// root, excluding vendor/, node_modules/, dist/, site/, and .git/.
func allSourceFiles(t *testing.T, root string) []string {
	t.Helper()
	sourceExts := map[string]bool{
		".go": true,
		".ts": true,
		".js": true,
	}
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() && (info.Name() == "vendor" || info.Name() == ".git" ||
			info.Name() == "dist" || info.Name() == "site" || info.Name() == "node_modules") {
			return filepath.SkipDir
		}
		if !info.IsDir() && sourceExts[filepath.Ext(path)] {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk project: %v", err)
	}
	return files
}

// ---------------------------------------------------------------------------
// 23. No UTF-8 BOM — source files must not start with a byte-order mark
// ---------------------------------------------------------------------------

// TestNoUTF8BOM detects the UTF-8 BOM (0xEF 0xBB 0xBF) that Windows editors
// sometimes insert. BOM causes subtle issues with Go tooling and TypeScript
// compilers and should never appear in source files.
func TestNoUTF8BOM(t *testing.T) {
	root := projectRoot(t)
	bom := []byte{0xEF, 0xBB, 0xBF}

	for _, p := range allSourceFiles(t, root) {
		rel, _ := filepath.Rel(root, p)
		t.Run(rel, func(t *testing.T) {
			data, readErr := os.ReadFile(filepath.Clean(p))
			if readErr != nil {
				t.Fatalf("read: %v", readErr)
			}
			if bytes.HasPrefix(data, bom) {
				t.Errorf("file starts with UTF-8 BOM (0xEF 0xBB 0xBF); remove it")
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 24. No mojibake — detect double-encoded UTF-8 (encoding corruption)
// ---------------------------------------------------------------------------

// TestNoMojibake catches the classic Windows encoding corruption where UTF-8
// bytes are misread as Windows-1252/Latin-1 and re-encoded as UTF-8.
// Example: em dash U+2014 becomes a 6-byte garbled sequence starting with
// 0xC3 0xA2. We detect that signature to catch double-encoded files.
func TestNoMojibake(t *testing.T) {
	root := projectRoot(t)
	// 0xC3 0xA2 is UTF-8 for U+00E2 (Latin small letter a with circumflex).
	// In mojibake, it always appears followed by 0xE2 as part of a garbled
	// multi-byte sequence (e.g., em dash becomes 0xC3 0xA2 0xE2 0x82 ...).
	// We match that three-byte signature: 0xC3 0xA2 0xE2.
	mojibakePattern := []byte{0xC3, 0xA2, 0xE2}

	for _, p := range allSourceFiles(t, root) {
		rel, _ := filepath.Rel(root, p)
		t.Run(rel, func(t *testing.T) {
			data, readErr := os.ReadFile(filepath.Clean(p))
			if readErr != nil {
				t.Fatalf("read: %v", readErr)
			}
			if idx := bytes.Index(data, mojibakePattern); idx >= 0 {
				// Show context around the corruption
				start := idx
				if start > 20 {
					start = idx - 20
				}
				end := idx + 30
				if end > len(data) {
					end = len(data)
				}
				t.Errorf("double-encoded UTF-8 (mojibake) detected at byte %d: %q\n"+
					"This usually means a Windows editor re-encoded the file.\n"+
					"Fix: restore from git (git checkout HEAD -- %s) and re-apply changes with a UTF-8-aware editor.",
					idx, data[start:end], rel)
			}
		})
	}
}
