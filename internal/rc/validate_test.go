//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	"strings"
	"testing"
)

func TestValidate_ValidConfig(t *testing.T) {
	data := []byte("token_budget: 4000\nauto_archive: false\n")
	warnings, err := Validate(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %v", warnings)
	}
}

func TestValidate_UnknownTopLevelField(t *testing.T) {
	data := []byte("scratchpad_encypt: true\n")
	warnings, err := Validate(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(warnings) == 0 {
		t.Fatal("expected warning for unknown field")
	}
	if !strings.Contains(warnings[0], "scratchpad_encypt") {
		t.Errorf("warning should mention field name, got: %s", warnings[0])
	}
}

func TestValidate_UnknownNestedField(t *testing.T) {
	data := []byte("notify:\n  evnts:\n    - loop\n")
	warnings, err := Validate(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(warnings) == 0 {
		t.Fatal("expected warning for unknown nested field")
	}
	if !strings.Contains(warnings[0], "evnts") {
		t.Errorf("warning should mention field name, got: %s", warnings[0])
	}
}

func TestValidate_MultipleUnknowns(t *testing.T) {
	data := []byte("tokan_budget: 4000\nauto_archve: true\n")
	warnings, err := Validate(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(warnings) < 2 {
		t.Errorf("expected at least 2 warnings, got %d: %v", len(warnings), warnings)
	}
}

func TestValidate_MalformedYAML(t *testing.T) {
	data := []byte(":\n  :\n  - [invalid yaml")
	_, err := Validate(data)
	if err == nil {
		t.Fatal("expected error for malformed YAML")
	}
}

func TestValidate_EmptyFile(t *testing.T) {
	data := []byte("")
	warnings, err := Validate(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(warnings) != 0 {
		t.Errorf("expected no warnings for empty file, got %v", warnings)
	}
}

func TestValidate_FullValidConfig(t *testing.T) {
	data := []byte(`token_budget: 8000
priority_order:
  - TASKS.md
  - DECISIONS.md
auto_archive: true
archive_after_days: 7
scratchpad_encrypt: true
entry_count_learnings: 30
entry_count_decisions: 20
convention_line_count: 200
injection_token_warn: 15000
context_window: 200000
event_log: false
key_rotation_days: 90
key_path: /tmp/key
notify:
  events:
    - loop
    - nudge
  key_rotation_days: 90
`)
	warnings, err := Validate(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(warnings) != 0 {
		t.Errorf("expected no warnings for full valid config, got %v", warnings)
	}
}

func TestValidate_BackendsWellFormed(t *testing.T) {
	data := []byte(`backends:
  default: vllm
  vllm:
    type: openai-compatible
    endpoint: http://localhost:8000
    api_key_env: ""
    timeout: 30s
    default_model: openai/gpt-oss-120b
`)
	warnings, err := Validate(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %v", warnings)
	}
}

func TestValidate_BackendsMissingEndpoint(t *testing.T) {
	data := []byte(`backends:
  vllm:
    timeout: 30s
`)
	_, err := Validate(data)
	if err == nil {
		t.Fatal("expected error for missing endpoint")
	}
	if !strings.Contains(err.Error(), "vllm") {
		t.Errorf("error should mention backend name, got: %v", err)
	}
}

func TestValidate_BackendsMalformedShape(t *testing.T) {
	data := []byte(`backends:
  - vllm
`)
	_, err := Validate(data)
	if err == nil {
		t.Fatal("expected error for malformed backends shape")
	}
}

func TestValidate_BackendsDefaultMissing(t *testing.T) {
	data := []byte(`backends:
  default: openai
  vllm:
    endpoint: http://localhost:8000
`)
	_, err := Validate(data)
	if err == nil {
		t.Fatal("expected error for missing default backend")
	}
	if !strings.Contains(err.Error(), "openai") {
		t.Errorf("error should mention missing default, got: %v", err)
	}
}

func TestValidate_BackendsMultipleWithoutDefault(t *testing.T) {
	data := []byte(`backends:
  vllm:
    endpoint: http://localhost:8000
  openai:
    endpoint: https://api.openai.com
    api_key_env: OPENAI_API_KEY
`)
	warnings, err := Validate(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %v", warnings)
	}
}

func TestValidate_BackendsUnknownNestedField(t *testing.T) {
	data := []byte(`backends:
  vllm:
    endpoint: http://localhost:8000
    api_token: nope
`)
	warnings, err := Validate(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(warnings) == 0 {
		t.Fatal("expected warning for unknown backend field")
	}
	if !strings.Contains(warnings[0], "api_token") {
		t.Errorf("warning should mention field name, got: %s", warnings[0])
	}
}

func TestValidate_BackendsEmptyOK(t *testing.T) {
	data := []byte("backends: {}\n")
	warnings, err := Validate(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %v", warnings)
	}
}
