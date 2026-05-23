//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"syscall"
	"testing"

	errBackend "github.com/ActiveMemory/ctx/internal/err/backend"
)

const (
	testModel       = "test/dummy-1b"
	testUserContent = "ping?"
	testAssistant   = "pong."
)

// fakeOpenAIServer returns an httptest server that responds
// to /v1/models with `{"object":"list","data":[{"id":<model>}]}`
// and to /v1/chat/completions with a canned chat reply that
// echoes the supplied model.
func fakeOpenAIServer(t *testing.T, model string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/models", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"object": "list",
			"data": []map[string]any{
				{"id": model, "object": "model", "owned_by": "vllm"},
			},
		})
	})
	mux.HandleFunc("/v1/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		var body chatRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"model": body.Model,
			"choices": []map[string]any{{
				"message":       map[string]string{"role": "assistant", "content": testAssistant},
				"finish_reason": "stop",
			}},
		})
	})
	return httptest.NewServer(mux)
}

func TestNewOpenAICompat_MissingEndpoint(t *testing.T) {
	_, err := newOpenAICompat(Config{Name: "x"})
	if !errors.Is(err, errBackend.ErrMissingEndpoint) {
		t.Fatalf("got %v, want ErrMissingEndpoint", err)
	}
}

func TestNewOpenAICompat_InvalidEndpoint(t *testing.T) {
	cases := []string{
		"ftp://example.com",
		"://no-scheme",
		"http://",
	}
	for _, ep := range cases {
		_, err := newOpenAICompat(Config{Name: "x", Endpoint: ep})
		if !errors.Is(err, errBackend.ErrInvalidEndpoint) {
			t.Errorf("endpoint %q: got %v, want ErrInvalidEndpoint", ep, err)
		}
	}
}

func TestOpenAICompat_Ping_Happy(t *testing.T) {
	srv := fakeOpenAIServer(t, testModel)
	defer srv.Close()
	b, err := newOpenAICompat(Config{Name: "x", Endpoint: srv.URL})
	if err != nil {
		t.Fatalf("ctor: %v", err)
		return
	}
	if pingErr := b.Ping(context.Background()); pingErr != nil {
		t.Fatalf("Ping: %v", pingErr)
	}
}

func TestOpenAICompat_Ping_Unhealthy(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "engine starting", http.StatusServiceUnavailable)
	}))
	defer srv.Close()
	b, _ := newOpenAICompat(Config{Name: "x", Endpoint: srv.URL})
	err := b.Ping(context.Background())
	if !errors.Is(err, errBackend.ErrUnhealthyStatus) {
		t.Fatalf("got %v, want ErrUnhealthyStatus", err)
	}
}

func TestOpenAICompat_Models_Happy(t *testing.T) {
	srv := fakeOpenAIServer(t, testModel)
	defer srv.Close()
	b, err := newOpenAICompat(Config{Name: "x", Endpoint: srv.URL})
	if err != nil {
		t.Fatalf("ctor: %v", err)
		return
	}
	models, modelsErr := b.Models(context.Background())
	if modelsErr != nil {
		t.Fatalf("Models: %v", modelsErr)
	}
	if len(models) != 1 || models[0] != testModel {
		t.Errorf("Models = %v, want [%q]", models, testModel)
	}
}

func TestOpenAICompat_Models_OrderPreserved(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"object":"list","data":[{"id":"a"},{"id":"b"},{"id":"c"}]}`))
	}))
	defer srv.Close()
	b, _ := newOpenAICompat(Config{Name: "x", Endpoint: srv.URL})
	models, err := b.Models(context.Background())
	if err != nil {
		t.Fatalf("Models: %v", err)
	}
	want := []string{"a", "b", "c"}
	if len(models) != len(want) {
		t.Fatalf("len = %d, want %d", len(models), len(want))
	}
	for i, m := range want {
		if models[i] != m {
			t.Errorf("models[%d] = %q, want %q", i, models[i], m)
		}
	}
}

func TestOpenAICompat_Models_Empty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"object":"list","data":[]}`))
	}))
	defer srv.Close()
	b, _ := newOpenAICompat(Config{Name: "x", Endpoint: srv.URL})
	_, err := b.Models(context.Background())
	if !errors.Is(err, errBackend.ErrEmptyModels) {
		t.Fatalf("got %v, want ErrEmptyModels", err)
	}
}

func TestOpenAICompat_Models_Unhealthy(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "engine starting", http.StatusServiceUnavailable)
	}))
	defer srv.Close()
	b, _ := newOpenAICompat(Config{Name: "x", Endpoint: srv.URL})
	_, err := b.Models(context.Background())
	if !errors.Is(err, errBackend.ErrUnhealthyStatus) {
		t.Fatalf("got %v, want ErrUnhealthyStatus", err)
	}
}

func TestOpenAICompat_Models_BadJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("not json"))
	}))
	defer srv.Close()
	b, _ := newOpenAICompat(Config{Name: "x", Endpoint: srv.URL})
	_, err := b.Models(context.Background())
	if err == nil || !strings.Contains(err.Error(), "parse response") {
		t.Fatalf("got %v, want parse-response error", err)
	}
}

func TestOpenAICompat_Models_Unreachable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	url := srv.URL
	srv.Close()
	b, _ := newOpenAICompat(Config{Name: "x", Endpoint: url})
	_, err := b.Models(context.Background())
	if !errors.Is(err, errBackend.ErrUnreachable) {
		t.Fatalf("got %v, want ErrUnreachable", err)
	}
}

func TestOpenAICompat_Ping_Unreachable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	url := srv.URL
	srv.Close()
	b, _ := newOpenAICompat(Config{Name: "x", Endpoint: url})
	err := b.Ping(context.Background())
	if !errors.Is(err, errBackend.ErrUnreachable) {
		t.Fatalf("got %v, want ErrUnreachable", err)
	}
	if !errors.Is(err, syscall.ECONNREFUSED) {
		t.Fatalf("expected cause to be ECONNREFUSED, got %v", err)
	}
}

func TestOpenAICompat_Complete_Happy(t *testing.T) {
	srv := fakeOpenAIServer(t, testModel)
	defer srv.Close()
	b, _ := newOpenAICompat(Config{
		Name: "x", Endpoint: srv.URL, DefaultModel: testModel,
	})
	resp, err := b.Complete(context.Background(), Request{
		Messages: []Message{{Role: "user", Content: testUserContent}},
	})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if resp.Content != testAssistant {
		t.Errorf("Content = %q, want %q", resp.Content, testAssistant)
	}
	if resp.Model != testModel {
		t.Errorf("Model = %q, want %q", resp.Model, testModel)
	}
	if len(resp.Raw) == 0 {
		t.Error("Raw bytes are empty")
	}
}

func TestOpenAICompat_Complete_MissingModel(t *testing.T) {
	srv := fakeOpenAIServer(t, testModel)
	defer srv.Close()
	b, _ := newOpenAICompat(Config{Name: "x", Endpoint: srv.URL})
	_, err := b.Complete(context.Background(), Request{
		Messages: []Message{{Role: "user", Content: testUserContent}},
	})
	if !errors.Is(err, errBackend.ErrMissingModel) {
		t.Fatalf("got %v, want ErrMissingModel", err)
	}
}

func TestOpenAICompat_Complete_UpstreamStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"invalid api key"}`))
	}))
	defer srv.Close()
	b, _ := newOpenAICompat(Config{
		Name: "x", Endpoint: srv.URL, DefaultModel: testModel,
	})
	_, err := b.Complete(context.Background(), Request{
		Messages: []Message{{Role: "user", Content: testUserContent}},
	})
	if !errors.Is(err, errBackend.ErrUpstreamStatus) {
		t.Fatalf("got %v, want ErrUpstreamStatus", err)
	}
	if !strings.Contains(err.Error(), "invalid api key") {
		t.Errorf("error should include body excerpt; got %q", err.Error())
	}
}

func TestOpenAICompat_Complete_BadJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("not json"))
	}))
	defer srv.Close()
	b, _ := newOpenAICompat(Config{
		Name: "x", Endpoint: srv.URL, DefaultModel: testModel,
	})
	_, err := b.Complete(context.Background(), Request{
		Messages: []Message{{Role: "user", Content: testUserContent}},
	})
	if err == nil {
		t.Fatalf("expected parse error")
	}
}

func TestOpenAICompat_Complete_EmptyChoices(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"model":"x","choices":[]}`))
	}))
	defer srv.Close()
	b, _ := newOpenAICompat(Config{
		Name: "x", Endpoint: srv.URL, DefaultModel: testModel,
	})
	_, err := b.Complete(context.Background(), Request{
		Messages: []Message{{Role: "user", Content: testUserContent}},
	})
	if err == nil || !strings.Contains(err.Error(), "no completion choices") {
		t.Fatalf("got %v, want empty-choices error", err)
	}
}

func TestOpenAICompat_Auth_BearerHeader(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"object":"list","data":[]}`)
	}))
	defer srv.Close()
	t.Setenv("FAKE_AI_KEY", "sk-test-123")
	b, _ := newOpenAICompat(Config{
		Name: "x", Endpoint: srv.URL, APIKeyEnv: "FAKE_AI_KEY",
	})
	if err := b.Ping(context.Background()); err != nil {
		t.Fatalf("Ping: %v", err)
	}
	if gotAuth != "Bearer sk-test-123" {
		t.Errorf("Authorization = %q, want %q", gotAuth, "Bearer sk-test-123")
	}
}

func TestOpenAICompat_Complete_RespectsResponseFormat(t *testing.T) {
	var gotPayload chatRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&gotPayload)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"model":"x","choices":[{"message":{"role":"assistant","content":"{}"}}]}`))
	}))
	defer srv.Close()
	b, _ := newOpenAICompat(Config{
		Name: "x", Endpoint: srv.URL, DefaultModel: testModel,
	})
	_, err := b.Complete(context.Background(), Request{
		Messages: []Message{{Role: "user", Content: testUserContent}},
		ResponseFormat: &ResponseFormat{
			Type:       "json_schema",
			JSONSchema: map[string]any{"name": "x"},
		},
		Temperature: -1, // unset
	})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if gotPayload.ResponseFormat == nil ||
		gotPayload.ResponseFormat.Type != "json_schema" {
		t.Errorf("ResponseFormat not propagated: %+v", gotPayload.ResponseFormat)
	}
	if gotPayload.Temperature != nil {
		t.Errorf("Temperature should be omitted when negative; got %v",
			*gotPayload.Temperature)
	}
}
