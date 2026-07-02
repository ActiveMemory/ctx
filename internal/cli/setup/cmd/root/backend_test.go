//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSetupBackendNoArgsAccepted(t *testing.T) {
	out, runErr := executeSetup(t,
		"--backend", "vllm",
		"--endpoint", "http://localhost:8000",
	)
	if runErr != nil {
		t.Fatalf("setup --backend failed: %v", runErr)
	}
	if !strings.Contains(out, "backends:") {
		t.Fatalf("output = %q", out)
	}
}

func TestSetupBackendOnlyDoesNotWriteOpenCodeConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("OPENCODE_HOME", home)
	_, runErr := executeSetup(t,
		"--backend", "openai",
		"--endpoint", "https://example.com/v1",
		"--write",
	)
	if runErr != nil {
		t.Fatalf("setup --backend failed: %v", runErr)
	}
	if _, err := os.Stat(filepath.Join(home, "opencode.json")); !os.IsNotExist(err) {
		t.Fatalf("unexpected opencode.json write")
	}
}

func TestSetupNoArgsNoBackendRejected(t *testing.T) {
	_, runErr := executeSetup(t)
	if runErr == nil {
		t.Fatalf("setup without args should fail")
	}
}

func TestSetupExistingToolStillWorks(t *testing.T) {
	_, runErr := executeSetup(t, "aider")
	if runErr != nil {
		t.Fatalf("setup aider failed: %v", runErr)
	}
}

func TestSetupBackendModeWinsOverToolArg(t *testing.T) {
	_, runErr := executeSetup(t,
		"aider",
		"--backend", "vllm",
		"--endpoint", "http://localhost:8000",
	)
	if runErr != nil {
		t.Fatalf("setup --backend with tool arg failed: %v", runErr)
	}
}

func TestSetupBackendUnsupportedRejected(t *testing.T) {
	_, runErr := executeSetup(t, "--backend", "unknown")
	if runErr == nil {
		t.Fatalf("unsupported backend should fail")
	}
}

func TestSetupBackendMissingEndpointRejected(t *testing.T) {
	_, runErr := executeSetup(t, "--backend", "openai-compatible")
	if runErr == nil {
		t.Fatalf("missing endpoint should fail")
	}
}

func TestSetupBackendDryRunPrintsYaml(t *testing.T) {
	out, runErr := executeSetup(t,
		"--backend", "vllm",
		"--endpoint", "http://localhost:8000",
	)
	if runErr != nil {
		t.Fatalf("setup --backend failed: %v", runErr)
	}
	for _, want := range []string{"backends:", "vllm:", "endpoint:"} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q: %q", want, out)
		}
	}
}

func TestSetupBackendWriteCreatesCtxRC(t *testing.T) {
	tmpDir := chdirTemp(t)
	_, runErr := executeSetup(t,
		"--backend", "vllm",
		"--endpoint", "http://localhost:8000",
		"--write",
	)
	if runErr != nil {
		t.Fatalf("setup --backend --write failed: %v", runErr)
	}
	data, readErr := os.ReadFile(filepath.Join(tmpDir, ".ctxrc"))
	if readErr != nil {
		t.Fatalf("ReadFile() error = %v", readErr)
	}
	if !strings.Contains(string(data), "vllm:") {
		t.Fatalf(".ctxrc = %q", string(data))
	}
}

func TestSetupBackendWritePreservesUnrelatedFields(t *testing.T) {
	tmpDir := chdirTemp(t)
	writeErr := os.WriteFile(
		filepath.Join(tmpDir, ".ctxrc"),
		[]byte("token_budget: 123\n"),
		0o644,
	)
	if writeErr != nil {
		t.Fatalf("WriteFile() error = %v", writeErr)
	}
	_, runErr := executeSetup(t,
		"--backend", "vllm",
		"--endpoint", "http://localhost:8000",
		"--write",
	)
	if runErr != nil {
		t.Fatalf("setup --backend --write failed: %v", runErr)
	}
	data, readErr := os.ReadFile(filepath.Join(tmpDir, ".ctxrc"))
	if readErr != nil {
		t.Fatalf("ReadFile() error = %v", readErr)
	}
	if !strings.Contains(string(data), "token_budget: 123") {
		t.Fatalf(".ctxrc = %q", string(data))
	}
}

func TestSetupBackendEnvConflictWarning(t *testing.T) {
	t.Setenv("CTX_TEST_BACKEND_KEY", "set")
	out, runErr := executeSetup(t,
		"--backend", "openai",
		"--api-key-env", "CTX_TEST_BACKEND_KEY",
	)
	if runErr != nil {
		t.Fatalf("setup --backend failed: %v", runErr)
	}
	if !strings.Contains(out, "warning:") {
		t.Fatalf("output = %q", out)
	}
}

func TestSetupBackendOpenCodeWritesProviderConfig(t *testing.T) {
	project := chdirTemp(t)
	home := t.TempDir()
	t.Setenv("OPENCODE_HOME", home)
	_, runErr := executeSetup(t,
		"opencode",
		"--backend", "openai",
		"--endpoint", "https://example.com/v1",
		"--write",
	)
	if runErr != nil {
		t.Fatalf("setup opencode --backend failed: %v", runErr)
	}
	data, err := os.ReadFile(filepath.Join(home, "opencode.json"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if !strings.Contains(string(data), `"provider"`) ||
		!strings.Contains(string(data), `"baseURL": "https://example.com/v1"`) ||
		!strings.Contains(string(data), `"mcp"`) {
		t.Fatalf("opencode.json = %s", string(data))
	}
	if _, err := os.Stat(filepath.Join(project, ".ctxrc")); err != nil {
		t.Fatalf(".ctxrc not written: %v", err)
	}
}

func TestSetupBackendOpenCodeUsesDefaultEndpoint(t *testing.T) {
	project := chdirTemp(t)
	home := t.TempDir()
	t.Setenv("OPENCODE_HOME", home)
	_, runErr := executeSetup(t,
		"opencode",
		"--backend", "openai",
		"--write",
	)
	if runErr != nil {
		t.Fatalf("setup opencode --backend failed: %v", runErr)
	}
	data, err := os.ReadFile(filepath.Join(home, "opencode.json"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if !strings.Contains(string(data), `"baseURL": "https://api.openai.com"`) {
		t.Fatalf("opencode.json = %s", string(data))
	}
	if _, err := os.Stat(filepath.Join(project, ".ctxrc")); err != nil {
		t.Fatalf(".ctxrc not written: %v", err)
	}
}

func TestSetupBackendMalformedCtxRCRejectedAndUntouched(t *testing.T) {
	tmpDir := chdirTemp(t)
	path := filepath.Join(tmpDir, ".ctxrc")
	original := []byte("backends: [broken\n")
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	_, runErr := executeSetup(t, "--backend", "vllm", "--endpoint", "http://localhost:8000", "--write")
	if runErr == nil {
		t.Fatalf("malformed .ctxrc should fail")
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(after) != string(original) {
		t.Fatalf(".ctxrc changed: %q", string(after))
	}
}

func executeSetup(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := new(bytes.Buffer)
	cmd := Cmd()
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	runErr := cmd.Execute()
	return buf.String(), runErr
}

func chdirTemp(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	origDir, wdErr := os.Getwd()
	if wdErr != nil {
		t.Fatalf("Getwd() error = %v", wdErr)
	}
	if chdirErr := os.Chdir(tmpDir); chdirErr != nil {
		t.Fatalf("Chdir() error = %v", chdirErr)
	}
	t.Cleanup(func() {
		if chdirErr := os.Chdir(origDir); chdirErr != nil {
			t.Fatalf("Chdir() cleanup error = %v", chdirErr)
		}
	})
	return tmpDir
}
