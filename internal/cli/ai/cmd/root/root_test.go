//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/rc"
)

func TestNoBackendConfiguredFails(t *testing.T) {
	chdirProject(t)

	_, runErr := executeAI(t, "ping")
	if runErr == nil {
		t.Fatalf("ping without backend should fail")
	}
}

func TestMultipleBackendsWithoutDefaultFails(t *testing.T) {
	project := chdirProject(t)
	writeCtxRC(t, project, multipleBackendsYAML())

	_, runErr := executeAI(t, "ping")
	if runErr == nil {
		t.Fatalf("ping with ambiguous backends should fail")
	}
}

func TestPingSucceeds(t *testing.T) {
	project := chdirProject(t)
	server := modelsServer(t)
	t.Cleanup(server.Close)
	writeCtxRC(t, project, singleBackendYAML(server.URL))

	out, runErr := executeAI(t, "ping")
	if runErr != nil {
		t.Fatalf("ping failed: %v", runErr)
	}
	for _, want := range []string{"vllm", server.URL, "model-a"} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q: %q", want, out)
		}
	}
}

func TestProposeWritesArtifactOnly(t *testing.T) {
	project := chdirProject(t)
	server := completionServer(t)
	t.Cleanup(server.Close)
	writeCtxRC(t, project, singleBackendYAML(server.URL))
	contextFile := filepath.Join(project, ".context", "DECISIONS.md")
	before, readErr := os.ReadFile(contextFile)
	if readErr != nil {
		t.Fatalf("ReadFile() error = %v", readErr)
	}
	input := filepath.Join(project, "input.txt")
	if writeErr := os.WriteFile(input, []byte("source"), 0o644); writeErr != nil {
		t.Fatalf("WriteFile() error = %v", writeErr)
	}

	_, runErr := executeAI(
		t,
		"propose",
		input,
		"--emit",
		"decisions,learnings,tasks,open-questions",
	)
	if runErr != nil {
		t.Fatalf("propose failed: %v", runErr)
	}
	after, readErr := os.ReadFile(contextFile)
	if readErr != nil {
		t.Fatalf("ReadFile() error = %v", readErr)
	}
	if string(before) != string(after) {
		t.Fatalf("context file changed")
	}
	entries, readDirErr := os.ReadDir(
		filepath.Join(project, ".context", "proposals", "ai"),
	)
	if readDirErr != nil {
		t.Fatalf("ReadDir() error = %v", readDirErr)
	}
	if len(entries) != 1 {
		t.Fatalf("artifact count = %d", len(entries))
	}
	data, readErr := os.ReadFile(filepath.Join(project, ".context", "proposals", "ai", entries[0].Name()))
	if readErr != nil {
		t.Fatalf("ReadFile() error = %v", readErr)
	}
	var artifact map[string]any
	if err := json.Unmarshal(data, &artifact); err != nil {
		t.Fatalf("artifact json invalid: %v", err)
	}
	if artifact["response"] == nil {
		t.Fatalf("artifact missing response")
	}
}

func TestProposeRejectsEmptyEmit(t *testing.T) {
	project := chdirProject(t)
	server := completionServer(t)
	t.Cleanup(server.Close)
	writeCtxRC(t, project, singleBackendYAML(server.URL))
	input := filepath.Join(project, "input.txt")
	if writeErr := os.WriteFile(input, []byte("source"), 0o644); writeErr != nil {
		t.Fatalf("WriteFile() error = %v", writeErr)
	}
	_, runErr := executeAI(t, "propose", input, "--emit", "")
	if runErr == nil {
		t.Fatalf("empty emit should fail")
	}
}

func TestCommandRegistrationPathResolves(t *testing.T) {
	cmd := Cmd()
	cmd.SetArgs([]string{"ping", "--backend", "vllm"})
	leaf, _, findErr := cmd.Find([]string{"ping"})
	if findErr != nil {
		t.Fatalf("Find() error = %v", findErr)
	}
	if leaf == nil {
		t.Fatalf("ping command not registered")
	}
}

func executeAI(t *testing.T, args ...string) (string, error) {
	t.Helper()
	rc.Reset()
	t.Cleanup(rc.Reset)
	buf := new(bytes.Buffer)
	cmd := Cmd()
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	runErr := cmd.Execute()
	return buf.String(), runErr
}

func chdirProject(t *testing.T) string {
	t.Helper()
	project := t.TempDir()
	contextDir := filepath.Join(project, ".context")
	if mkdirErr := os.MkdirAll(contextDir, 0o755); mkdirErr != nil {
		t.Fatalf("MkdirAll() error = %v", mkdirErr)
	}
	if writeErr := os.WriteFile(
		filepath.Join(contextDir, "DECISIONS.md"),
		[]byte("decisions"),
		0o644,
	); writeErr != nil {
		t.Fatalf("WriteFile() error = %v", writeErr)
	}
	origDir, wdErr := os.Getwd()
	if wdErr != nil {
		t.Fatalf("Getwd() error = %v", wdErr)
	}
	if chdirErr := os.Chdir(project); chdirErr != nil {
		t.Fatalf("Chdir() error = %v", chdirErr)
	}
	t.Cleanup(func() {
		if chdirErr := os.Chdir(origDir); chdirErr != nil {
			t.Fatalf("Chdir() cleanup error = %v", chdirErr)
		}
	})
	return project
}

func writeCtxRC(t *testing.T, project string, content string) {
	t.Helper()
	if writeErr := os.WriteFile(
		filepath.Join(project, ".ctxrc"),
		[]byte(content),
		0o644,
	); writeErr != nil {
		t.Fatalf("WriteFile() error = %v", writeErr)
	}
}

func modelsServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		_, writeErr := w.Write([]byte(`{"data":[{"id":"model-a"}]}`))
		if writeErr != nil {
			t.Fatalf("Write() error = %v", writeErr)
		}
	}))
}

func completionServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		_, writeErr := w.Write([]byte(completionResponseJSON()))
		if writeErr != nil {
			t.Fatalf("Write() error = %v", writeErr)
		}
	}))
}

func singleBackendYAML(endpoint string) string {
	return "backends:\n  default: vllm\n  vllm:\n    endpoint: " +
		endpoint + "\n"
}

func multipleBackendsYAML() string {
	return "backends:\n  vllm:\n    endpoint: http://127.0.0.1\n" +
		"  openai:\n    endpoint: http://127.0.0.2\n"
}

func completionResponseJSON() string {
	return `{"model":"m","choices":[{"message":{"content":"{\"rows\":[{\"emit\":\"decisions\",\"text\":\"ok\"}],\"metadata\":{\"backend\":\"vllm\",\"model\":\"m\",\"input\":\"` + "input.txt" + `\",\"status\":\"proposed\"}}"}}]}`
}
