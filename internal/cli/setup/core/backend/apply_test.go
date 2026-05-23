//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ActiveMemory/ctx/internal/entity"
)

const (
	testEndpointVLLM     = "http://localhost:8000"
	testEndpointOpenAI   = "https://api.openai.com"
	testEndpointOverride = "http://updated:9000"
)

func writeCtxrc(t *testing.T, dir, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, ".ctxrc"), []byte(content), 0o600); err != nil {
		t.Fatalf("write .ctxrc: %v", err)
		return
	}
}

func readCtxrc(t *testing.T, dir string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, ".ctxrc"))
	if err != nil {
		t.Fatalf("read .ctxrc: %v", err)
		return ""
	}
	return string(data)
}

func TestApply_GreenfieldCreate(t *testing.T) {
	dir := t.TempDir()
	res, err := Apply(dir, entity.BackendConfig{
		Name:         "vllm",
		Endpoint:     testEndpointVLLM,
		Timeout:      30 * time.Second,
		DefaultModel: "Qwen/Qwen2.5-1.5B-Instruct",
	})
	if err != nil {
		t.Fatalf("Apply: %v", err)
		return
	}
	if !res.Created {
		t.Errorf("Created = false, expected true (no prior .ctxrc)")
	}
	if res.Updated {
		t.Errorf("Updated = true, expected false (no prior entry)")
	}
	body := readCtxrc(t, dir)
	if !strings.Contains(body, "backends:") {
		t.Errorf("expected backends key in output, got:\n%s", body)
	}
	if !strings.Contains(body, "name: vllm") {
		t.Errorf("expected name: vllm, got:\n%s", body)
	}
	if !strings.Contains(body, "endpoint: "+testEndpointVLLM) {
		t.Errorf("expected endpoint, got:\n%s", body)
	}
	if !strings.Contains(body, "timeout: 30s") {
		t.Errorf("expected timeout: 30s, got:\n%s", body)
	}
}

func TestApply_AppendsToExistingBackends(t *testing.T) {
	dir := t.TempDir()
	writeCtxrc(t, dir, `token_budget: 12000
backends:
  - name: vllm
    endpoint: http://localhost:8000
`)
	res, err := Apply(dir, entity.BackendConfig{
		Name:      "openai",
		Endpoint:  testEndpointOpenAI,
		APIKeyEnv: "OPENAI_API_KEY",
	})
	if err != nil {
		t.Fatalf("Apply: %v", err)
		return
	}
	if res.Created {
		t.Errorf("Created = true, expected false (.ctxrc preexisted)")
	}
	if res.Updated {
		t.Errorf("Updated = true, expected false (openai is new entry)")
	}
	body := readCtxrc(t, dir)
	if !strings.Contains(body, "token_budget: 12000") {
		t.Errorf("preexisting key dropped:\n%s", body)
	}
	if !strings.Contains(body, "name: vllm") {
		t.Errorf("preexisting backend dropped:\n%s", body)
	}
	if !strings.Contains(body, "name: openai") {
		t.Errorf("new backend missing:\n%s", body)
	}
	if !strings.Contains(body, "api_key_env: OPENAI_API_KEY") {
		t.Errorf("api_key_env not written:\n%s", body)
	}
}

func TestApply_IdempotentUpdate(t *testing.T) {
	dir := t.TempDir()
	writeCtxrc(t, dir, `backends:
  - name: vllm
    endpoint: http://localhost:8000
    timeout: 30s
`)
	res, err := Apply(dir, entity.BackendConfig{
		Name:     "vllm",
		Endpoint: testEndpointOverride,
		Timeout:  60 * time.Second,
	})
	if err != nil {
		t.Fatalf("Apply: %v", err)
		return
	}
	if !res.Updated {
		t.Errorf("Updated = false, expected true (existing vllm replaced)")
	}
	body := readCtxrc(t, dir)
	if strings.Contains(body, "localhost:8000") {
		t.Errorf("old endpoint not replaced:\n%s", body)
	}
	if !strings.Contains(body, testEndpointOverride) {
		t.Errorf("new endpoint not written:\n%s", body)
	}
	if !strings.Contains(body, "timeout: 1m0s") {
		t.Errorf("new timeout not written:\n%s", body)
	}
	// Should not introduce a duplicate name: vllm.
	if strings.Count(body, "name: vllm") != 1 {
		t.Errorf("duplicate entry:\n%s", body)
	}
}

func TestApply_NoTimeoutOmitsField(t *testing.T) {
	dir := t.TempDir()
	if _, err := Apply(dir, entity.BackendConfig{
		Name:     "vllm",
		Endpoint: testEndpointVLLM,
		// Timeout intentionally zero.
	}); err != nil {
		t.Fatalf("Apply: %v", err)
		return
	}
	body := readCtxrc(t, dir)
	if strings.Contains(body, "timeout:") {
		t.Errorf("timeout field should be omitted when zero, got:\n%s", body)
	}
}

func TestApply_PreservesOtherKeys(t *testing.T) {
	dir := t.TempDir()
	writeCtxrc(t, dir, `token_budget: 12000
tool: claude
default_backend: vllm
backends: []
`)
	if _, err := Apply(dir, entity.BackendConfig{
		Name: "vllm", Endpoint: testEndpointVLLM,
	}); err != nil {
		t.Fatalf("Apply: %v", err)
		return
	}
	body := readCtxrc(t, dir)
	for _, key := range []string{
		"token_budget: 12000",
		"tool: claude",
		"default_backend: vllm",
	} {
		if !strings.Contains(body, key) {
			t.Errorf("preexisting key %q dropped:\n%s", key, body)
		}
	}
}

func TestApply_EmptyButExistingFile(t *testing.T) {
	dir := t.TempDir()
	writeCtxrc(t, dir, "")
	res, err := Apply(dir, entity.BackendConfig{
		Name: "vllm", Endpoint: testEndpointVLLM,
	})
	if err != nil {
		t.Fatalf("Apply: %v", err)
		return
	}
	if res.Created {
		t.Errorf("Created = true, expected false (file existed even if empty)")
	}
	body := readCtxrc(t, dir)
	if !strings.Contains(body, "name: vllm") {
		t.Errorf("entry not written:\n%s", body)
	}
}
