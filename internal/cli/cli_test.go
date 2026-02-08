//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cli

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestBinaryIntegration is an integration test that builds and runs the actual binary.
//
// This test builds the ctx binary and exercises multiple commands to ensure
// they work correctly end-to-end. It verifies that subcommands execute properly
// (not falling through to root help) and produce expected output.
func TestBinaryIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tmpDir, err := os.MkdirTemp("", "cli-binary-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Build the binary
	binaryPath := filepath.Join(tmpDir, "ctx-test-binary")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/ctx") //nolint:gosec // G204: test builds local binary
	buildCmd.Env = append(os.Environ(), "CGO_ENABLED=0")

	// Get the project root (go up from internal/cli)
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("failed to get project root: %v", err)
	}
	buildCmd.Dir = projectRoot

	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build binary: %v\n%s", err, output)
	}

	// Create a test directory
	testDir := filepath.Join(tmpDir, "test-project")
	if err := os.Mkdir(testDir, 0750); err != nil {
		t.Fatalf("failed to create test dir: %v", err)
	}

	// Subtest: ctx init creates expected files
	t.Run("init creates expected files", func(t *testing.T) {
		initCmd := exec.Command(binaryPath, "init") //nolint:gosec // G204: test runs locally-built binary
		initCmd.Dir = testDir
		if output, err := initCmd.CombinedOutput(); err != nil {
			t.Fatalf("ctx init failed: %v\n%s", err, output)
		}

		// Check .context directory exists
		ctxDir := filepath.Join(testDir, ".context")
		if _, err := os.Stat(ctxDir); os.IsNotExist(err) {
			t.Fatal(".context directory was not created")
		}

		// Check required files exist
		requiredFiles := []string{
			"CONSTITUTION.md",
			"TASKS.md",
			"DECISIONS.md",
			"LEARNINGS.md",
			"CONVENTIONS.md",
			"ARCHITECTURE.md",
		}
		for _, name := range requiredFiles {
			path := filepath.Join(ctxDir, name)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("required file %s was not created", name)
			}
		}
	})

	// Subtest: ctx status returns valid status (not just help text)
	t.Run("status returns valid status", func(t *testing.T) {
		statusCmd := exec.Command(binaryPath, "status") //nolint:gosec // G204: test runs locally-built binary
		statusCmd.Dir = testDir
		output, err := statusCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("ctx status failed: %v\n%s", err, output)
		}

		outputStr := string(output)
		// Verify it's actual status output, not help text
		if strings.Contains(outputStr, "Usage:") || strings.Contains(outputStr, "Available Commands:") {
			t.Error("ctx status returned help text instead of status")
		}
		// Check for expected status output markers
		if !strings.Contains(outputStr, "Context Status") && !strings.Contains(outputStr, "Context Directory") {
			t.Errorf("ctx status did not return expected status output, got:\n%s", outputStr)
		}
	})

	// Subtest: ctx add learning modifies LEARNINGS.md
	t.Run("add learning modifies LEARNINGS.md", func(t *testing.T) {
		addCmd := exec.Command(binaryPath, "add", "learning", "Test learning from integration test", //nolint:gosec // G204: test runs locally-built binary
			"--context", "Testing integration",
			"--lesson", "Integration tests catch bugs",
			"--application", "Always run integration tests")
		addCmd.Dir = testDir
		if output, err := addCmd.CombinedOutput(); err != nil {
			t.Fatalf("ctx add learning failed: %v\n%s", err, output)
		}

		// Verify learning was added
		learningsPath := filepath.Join(testDir, ".context", "LEARNINGS.md")
		content, err := os.ReadFile(filepath.Clean(learningsPath))
		if err != nil {
			t.Fatalf("failed to read LEARNINGS.md: %v", err)
		}
		if !strings.Contains(string(content), "Test learning from integration test") {
			t.Error("learning was not added to LEARNINGS.md")
		}
	})

	// Subtest: ctx session save creates session file
	t.Run("session save creates session file", func(t *testing.T) {
		saveCmd := exec.Command(binaryPath, "session", "save") //nolint:gosec // G204: test runs locally-built binary
		saveCmd.Dir = testDir
		if output, err := saveCmd.CombinedOutput(); err != nil {
			t.Fatalf("ctx session save failed: %v\n%s", err, output)
		}

		// Check that sessions directory exists and has at least one file
		sessionsDir := filepath.Join(testDir, ".context", "sessions")
		entries, err := os.ReadDir(sessionsDir)
		if err != nil {
			t.Fatalf("failed to read sessions directory: %v", err)
		}
		if len(entries) == 0 {
			t.Error("no session file was created")
		}
	})

	// Subtest: ctx agent returns context packet
	t.Run("agent returns context packet", func(t *testing.T) {
		agentCmd := exec.Command(binaryPath, "agent") //nolint:gosec // G204: test runs locally-built binary
		agentCmd.Dir = testDir
		output, err := agentCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("ctx agent failed: %v\n%s", err, output)
		}

		outputStr := string(output)
		// Verify it's actual agent output, not help text
		if strings.Contains(outputStr, "Usage:") || strings.Contains(outputStr, "Available Commands:") {
			t.Error("ctx agent returned help text instead of context packet")
		}
		// Check for expected context packet markers
		if !strings.Contains(outputStr, "CONSTITUTION") && !strings.Contains(outputStr, "TASKS") {
			t.Errorf("ctx agent did not return expected context packet, got:\n%s", outputStr)
		}
	})

	// Subtest: ctx drift runs without error
	t.Run("drift runs without error", func(t *testing.T) {
		driftCmd := exec.Command(binaryPath, "drift") //nolint:gosec // G204: test runs locally-built binary
		driftCmd.Dir = testDir
		if output, err := driftCmd.CombinedOutput(); err != nil {
			t.Fatalf("ctx drift failed: %v\n%s", err, output)
		}
	})

	// Subtest: verify all subcommands execute (not falling through to root help)
	t.Run("subcommands execute without falling through to root help", func(t *testing.T) {
		// Commands that should produce output without "Available Commands:"
		// (which would indicate they fell through to root help)
		subcommands := []struct {
			args     []string
			checkFor string // expected output marker
		}{
			{[]string{"status"}, "Context"},
			{[]string{"agent"}, "Context Packet"},
			{[]string{"drift"}, "Drift"},
			{[]string{"load"}, ""},                 // load outputs context, varies by content
			{[]string{"hook", "cursor"}, "Cursor"}, // hook outputs integration instructions
		}

		for _, tc := range subcommands {
			t.Run(strings.Join(tc.args, "_"), func(t *testing.T) {
				cmd := exec.Command(binaryPath, tc.args...) //nolint:gosec // G204: test runs locally-built binary
				cmd.Dir = testDir
				output, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("ctx %s failed: %v\n%s", strings.Join(tc.args, " "), err, output)
				}

				outputStr := string(output)
				// Critical check: should NOT contain root help indicators
				if strings.Contains(outputStr, "Available Commands:") {
					t.Errorf("ctx %s fell through to root help:\n%s", strings.Join(tc.args, " "), outputStr)
				}
				// If we have an expected marker, check for it
				if tc.checkFor != "" && !strings.Contains(outputStr, tc.checkFor) {
					t.Errorf("ctx %s missing expected output %q:\n%s", strings.Join(tc.args, " "), tc.checkFor, outputStr)
				}
			})
		}
	})
}

// TestNoDirectFmtPrint ensures CLI code uses cmd.Print* instead of fmt.Print*.
//
// In Cobra commands, output should go through cmd.OutOrStdout() so that:
// - Tests can capture output
// - --quiet flags work correctly
// - Output can be redirected properly
//
// This test parses all non-test Go files in internal/cli and fails if any
// function that receives a *cobra.Command uses fmt.Print* directly.
func TestNoDirectFmtPrint(t *testing.T) {
	cliDir := "."

	// Forbidden fmt functions that should use cmd.Print* instead
	forbidden := map[string]bool{
		"Print":   true,
		"Println": true,
		"Printf":  true,
	}

	var violations []string

	err := filepath.Walk(cliDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-Go files and test files
		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			t.Errorf("failed to parse %s: %v", path, err)
			return nil
		}

		// Track if this file imports "fmt"
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

		// No fmt import, nothing to check
		if fmtAlias == "" {
			return nil
		}

		// Find functions that have a *cobra.Command parameter
		for _, decl := range node.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Type.Params == nil {
				continue
			}

			// Check if this function has a *cobra.Command parameter
			hasCobraCmd := false
			for _, param := range fn.Type.Params.List {
				if starExpr, ok := param.Type.(*ast.StarExpr); ok {
					if sel, ok := starExpr.X.(*ast.SelectorExpr); ok {
						if ident, ok := sel.X.(*ast.Ident); ok {
							if ident.Name == "cobra" && sel.Sel.Name == "Command" {
								hasCobraCmd = true
								break
							}
						}
					}
				}
			}

			if !hasCobraCmd {
				continue
			}

			// This function has *cobra.Command - check for fmt.Print* calls
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
					violations = append(violations,
						path+":"+
							strings.TrimPrefix(pos.String(), pos.Filename+":")+
							" in "+fn.Name.Name+"()")
				}

				return true
			})
		}

		return nil
	})

	if err != nil {
		t.Fatalf("failed to walk directory: %v", err)
	}

	if len(violations) > 0 {
		t.Errorf("found %d uses of fmt.Print* in functions with *cobra.Command (use cmd.Print* instead):\n  %s",
			len(violations), strings.Join(violations, "\n  "))
	}
}
