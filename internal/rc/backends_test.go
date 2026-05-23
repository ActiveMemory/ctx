//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	"testing"
	"time"
)

func TestBackends_Empty(t *testing.T) {
	declareContext(t, "")
	if got := Backends(); got != nil {
		t.Fatalf("Backends() = %v, want nil for empty .ctxrc", got)
	}
	if got := DefaultBackend(); got != "" {
		t.Fatalf("DefaultBackend() = %q, want empty", got)
	}
}

func TestBackends_Single(t *testing.T) {
	declareContext(t, `default_backend: vllm
backends:
  - name: vllm
    endpoint: http://localhost:8000
    timeout: 30s
    default_model: Qwen/Qwen2.5-1.5B-Instruct
`)
	got := Backends()
	if len(got) != 1 {
		t.Fatalf("Backends() len = %d, want 1", len(got))
	}
	b := got[0]
	if b.Name != "vllm" {
		t.Errorf("Name = %q, want vllm", b.Name)
	}
	if b.Endpoint != "http://localhost:8000" {
		t.Errorf("Endpoint = %q", b.Endpoint)
	}
	if b.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want 30s", b.Timeout)
	}
	if b.DefaultModel != "Qwen/Qwen2.5-1.5B-Instruct" {
		t.Errorf("DefaultModel = %q", b.DefaultModel)
	}
	if DefaultBackend() != "vllm" {
		t.Errorf("DefaultBackend() = %q, want vllm", DefaultBackend())
	}
}

func TestBackends_MultipleOrderPreserved(t *testing.T) {
	declareContext(t, `backends:
  - name: vllm
    endpoint: http://localhost:8000
  - name: openai
    endpoint: https://api.openai.com
    api_key_env: OPENAI_API_KEY
    timeout: 60s
`)
	got := Backends()
	if len(got) != 2 {
		t.Fatalf("Backends() len = %d, want 2", len(got))
	}
	if got[0].Name != "vllm" || got[1].Name != "openai" {
		t.Errorf("order = [%s, %s], want [vllm, openai]",
			got[0].Name, got[1].Name)
	}
	if got[1].APIKeyEnv != "OPENAI_API_KEY" {
		t.Errorf("openai.APIKeyEnv = %q", got[1].APIKeyEnv)
	}
}

func TestBackends_MalformedTimeoutYieldsZero(t *testing.T) {
	declareContext(t, `backends:
  - name: vllm
    endpoint: http://localhost:8000
    timeout: not-a-duration
`)
	got := Backends()
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].Timeout != 0 {
		t.Errorf("Timeout = %v, want 0 (parse failure)", got[0].Timeout)
	}
}

func TestBackends_EmptyTimeoutYieldsZero(t *testing.T) {
	declareContext(t, `backends:
  - name: vllm
    endpoint: http://localhost:8000
`)
	got := Backends()
	if got[0].Timeout != 0 {
		t.Errorf("Timeout = %v, want 0 when omitted", got[0].Timeout)
	}
}
