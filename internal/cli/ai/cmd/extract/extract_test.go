//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package extract_test

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
	aiExtract "github.com/ActiveMemory/ctx/internal/cli/ai/cmd/extract"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}

func newExtractServer(t *testing.T, jsonReply string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/chat/completions", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"model": "fake-extractor",
			"choices": []map[string]any{{
				"message": map[string]string{
					"role":    "assistant",
					"content": jsonReply,
				},
				"finish_reason": "stop",
			}},
		})
	})
	return httptest.NewServer(mux)
}

func declareProject(t *testing.T, ctxrc string) string {
	t.Helper()
	tempDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tempDir, dir.Context), 0o700); err != nil {
		t.Fatalf("mkdir .context: %v", err)
		return ""
	}
	if ctxrc != "" {
		if err := os.WriteFile(
			filepath.Join(tempDir, ".ctxrc"), []byte(ctxrc), 0o600,
		); err != nil {
			t.Fatalf("write .ctxrc: %v", err)
			return ""
		}
	}
	t.Chdir(tempDir)
	rc.Reset()
	t.Cleanup(rc.Reset)
	return tempDir
}

func TestExtract_HappyPath_WritesProposal(t *testing.T) {
	srv := newExtractServer(t, `{"decisions":[{"title":"x"}]}`)
	defer srv.Close()
	projectRoot := declareProject(t, `backends:
  - name: openai-compatible
    endpoint: `+srv.URL+`
    default_model: fake-extractor
`)
	cmd := aiExtract.Cmd()
	cmd.SetIn(strings.NewReader("we decided x for reasons y"))
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetContext(context.Background())
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
		return
	}
	if !strings.Contains(out.String(), "proposal written to") {
		t.Errorf("expected confirmation; got %q", out.String())
	}
	proposalsDir := filepath.Join(projectRoot, dir.Context, "proposals")
	entries, err := os.ReadDir(proposalsDir)
	if err != nil {
		t.Fatalf("readdir proposals: %v", err)
		return
	}
	if len(entries) != 1 {
		t.Fatalf("got %d proposal files, want 1", len(entries))
		return
	}
	body, _ := os.ReadFile(filepath.Join(proposalsDir, entries[0].Name()))
	got := string(body)
	for _, want := range []string{
		"# ctx ai extract proposal",
		"backend: `openai-compatible`",
		"model: `fake-extractor`",
		`{"decisions":[{"title":"x"}]}`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("proposal missing %q; got:\n%s", want, got)
		}
	}
}

func TestExtract_EmptyInput(t *testing.T) {
	srv := newExtractServer(t, `{}`)
	defer srv.Close()
	declareProject(t, `backends:
  - name: openai-compatible
    endpoint: `+srv.URL+`
    default_model: fake-extractor
`)
	cmd := aiExtract.Cmd()
	cmd.SetIn(strings.NewReader("   \n  \t  "))
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetContext(context.Background())
	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected error on empty input")
		return
	}
	if !strings.Contains(err.Error(), "empty") {
		t.Errorf("error should mention 'empty'; got %v", err)
	}
}

func TestExtract_NoBackendConfigured(t *testing.T) {
	declareProject(t, "")
	cmd := aiExtract.Cmd()
	cmd.SetIn(strings.NewReader("some text"))
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetContext(context.Background())
	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected ErrNoBackends")
		return
	}
	if !strings.Contains(err.Error(), "no backend") {
		t.Errorf("expected 'no backend' in error; got %v", err)
	}
}

func TestExtract_DoesNotTouchCanonicalFiles(t *testing.T) {
	srv := newExtractServer(t, `{"decisions":[]}`)
	defer srv.Close()
	projectRoot := declareProject(t, `backends:
  - name: openai-compatible
    endpoint: `+srv.URL+`
    default_model: fake-extractor
`)
	ctxDir := filepath.Join(projectRoot, dir.Context)
	pre := []byte("pre-existing tasks content\n")
	if err := os.WriteFile(filepath.Join(ctxDir, "TASKS.md"), pre, 0o600); err != nil {
		t.Fatalf("seed TASKS.md: %v", err)
		return
	}
	cmd := aiExtract.Cmd()
	cmd.SetIn(strings.NewReader("some session text"))
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetContext(context.Background())
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
		return
	}
	post, _ := os.ReadFile(filepath.Join(ctxDir, "TASKS.md"))
	if !bytes.Equal(pre, post) {
		t.Errorf("canonical TASKS.md was mutated by extract;\nbefore: %q\nafter:  %q",
			pre, post)
	}
	_ = errors.New // keep errors import alive for future test additions
}
