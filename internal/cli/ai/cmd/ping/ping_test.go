//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ping_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
	aiPing "github.com/ActiveMemory/ctx/internal/cli/ai/cmd/ping"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	errBackend "github.com/ActiveMemory/ctx/internal/err/backend"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}

func newModelsServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/models", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"object": "list",
			"data":   []map[string]any{{"id": "test-model"}},
		})
	})
	return httptest.NewServer(mux)
}

func declareProject(t *testing.T, ctxrc string) {
	t.Helper()
	tempDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tempDir, dir.Context), 0o700); err != nil {
		t.Fatalf("mkdir .context: %v", err)
		return
	}
	if ctxrc != "" {
		if err := os.WriteFile(
			filepath.Join(tempDir, ".ctxrc"), []byte(ctxrc), 0o600,
		); err != nil {
			t.Fatalf("write .ctxrc: %v", err)
			return
		}
	}
	t.Chdir(tempDir)
	rc.Reset()
	t.Cleanup(rc.Reset)
}

func TestPing_HappyPath(t *testing.T) {
	srv := newModelsServer(t)
	defer srv.Close()
	declareProject(t, `default_backend: openai-compatible
backends:
  - name: openai-compatible
    endpoint: `+srv.URL+`
`)
	cmd := aiPing.Cmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetContext(context.Background())
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
		return
	}
	if !strings.Contains(out.String(), "reachable") {
		t.Errorf("expected 'reachable' in output; got %q", out.String())
	}
}

func TestPing_NoBackendsConfigured(t *testing.T) {
	declareProject(t, "")
	cmd := aiPing.Cmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetContext(context.Background())
	err := cmd.Execute()
	if !errors.Is(err, errBackend.ErrNoBackends) {
		t.Fatalf("got %v, want ErrNoBackends", err)
	}
}

func TestPing_UnknownBackendNamed(t *testing.T) {
	srv := newModelsServer(t)
	defer srv.Close()
	declareProject(t, `backends:
  - name: openai-compatible
    endpoint: `+srv.URL+`
`)
	cmd := aiPing.Cmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"--backend", "anthropic"})
	cmd.SetContext(context.Background())
	err := cmd.Execute()
	if !errors.Is(err, errBackend.ErrBackendNotFound) {
		t.Fatalf("got %v, want ErrBackendNotFound", err)
	}
}

func TestPing_UnreachableSurfacesError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	url := srv.URL
	srv.Close()
	declareProject(t, `backends:
  - name: openai-compatible
    endpoint: `+url+`
`)
	cmd := aiPing.Cmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetContext(context.Background())
	err := cmd.Execute()
	if !errors.Is(err, errBackend.ErrUnreachable) {
		t.Fatalf("got %v, want ErrUnreachable", err)
	}
}
