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
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	cfgBackend "github.com/ActiveMemory/ctx/internal/config/backend"
	errBackend "github.com/ActiveMemory/ctx/internal/err/backend"
)

func TestOpenAICompatiblePingSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != cfgBackend.ModelsPath {
			t.Fatalf("path = %s, want %s", r.URL.Path, cfgBackend.ModelsPath)
		}
		_, writeErr := w.Write([]byte(`{"data":[{"id":"first"}]}`))
		if writeErr != nil {
			t.Fatalf("Write() error = %v", writeErr)
		}
	}))
	defer server.Close()
	backend := testOpenAIBackend(t, Config{Endpoint: server.URL})

	pingErr := backend.Ping(context.Background())
	if pingErr != nil {
		t.Fatalf("Ping() error = %v", pingErr)
	}
}

func TestOpenAICompatibleCompleteSuccess(t *testing.T) {
	t.Setenv("CTX_TEST_API_KEY", "token")
	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != cfgBackend.ChatCompletionsPath {
			t.Fatalf(
				"path = %s, want %s",
				r.URL.Path,
				cfgBackend.ChatCompletionsPath,
			)
		}
		if got := r.Header.Get(cfgBackend.HeaderAuthorization); got == "" {
			t.Fatalf("authorization header missing")
		}
		var request chatRequest
		decodeErr := json.NewDecoder(r.Body).Decode(&request)
		if decodeErr != nil {
			t.Fatalf("Decode() error = %v", decodeErr)
		}
		if request.Model != "configured-model" {
			t.Fatalf("model = %q", request.Model)
		}
		if request.ResponseFormat == nil || request.ResponseFormat.JSONSchema == nil {
			t.Fatalf("response_format missing")
		}
		if request.ResponseFormat.Type != cfgBackend.ResponseFormatJSONSchema {
			t.Fatalf("response_format type = %q", request.ResponseFormat.Type)
		}
		if len(request.ResponseFormat.JSONSchema.Schema) == 0 {
			t.Fatalf("schema missing")
		}
		_, writeErr := w.Write([]byte(`{"model":"configured-model","choices":[{"message":{"content":"ok"}}]}`))
		if writeErr != nil {
			t.Fatalf("Write() error = %v", writeErr)
		}
	}))
	defer server.Close()
	backend := testOpenAIBackend(t, Config{ //nolint:gosec // G101: test fixture, value is an env var name, not a credential
		Endpoint:     server.URL,
		APIKeyEnv:    "CTX_TEST_API_KEY",
		DefaultModel: "configured-model",
	})

	response, completeErr := backend.Complete(
		context.Background(),
		Request{Prompt: "hello", Schema: Schema{Name: "proposal", Schema: json.RawMessage(`{"type":"object"}`)}},
	)
	if completeErr != nil {
		t.Fatalf("Complete() error = %v", completeErr)
	}
	if response.Model != "configured-model" {
		t.Fatalf("Model = %q", response.Model)
	}
	if response.Text != "ok" {
		t.Fatalf("Text = %q", response.Text)
	}
}

func TestOpenAICompatibleUpstream4xxBody(t *testing.T) {
	backend := testStatusBackend(t, http.StatusUnauthorized, "nope")

	_, completeErr := backend.Complete(context.Background(), Request{})
	assertUpstreamBody(t, completeErr, "nope")
}

func TestOpenAICompatibleUpstream5xxBody(t *testing.T) {
	backend := testStatusBackend(t, http.StatusBadGateway, "bad gateway")

	_, completeErr := backend.Complete(context.Background(), Request{})
	assertUpstreamBody(t, completeErr, "bad gateway")
}

func TestOpenAICompatibleUnreachableServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(
		http.ResponseWriter,
		*http.Request,
	) {
	}))
	endpoint := server.URL
	server.Close()
	backend := testOpenAIBackend(t, Config{Endpoint: endpoint})

	pingErr := backend.Ping(context.Background())
	var unreachable errBackend.Unreachable
	if !errors.As(pingErr, &unreachable) {
		t.Fatalf("Ping() error = %T, want Unreachable", pingErr)
	}
}

func TestOpenAICompatibleInvalidEndpointScheme(t *testing.T) {
	backend := testOpenAIBackend(t, Config{Endpoint: "file:///tmp/socket"})

	pingErr := backend.Ping(context.Background())
	var invalid errBackend.InvalidEndpoint
	if !errors.As(pingErr, &invalid) {
		t.Fatalf("Ping() error = %T, want InvalidEndpoint", pingErr)
	}
}

func TestOpenAICompatibleTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		time.Sleep(20 * time.Millisecond)
		_, writeErr := w.Write([]byte(`{"data":[]}`))
		if writeErr != nil {
			t.Fatalf("Write() error = %v", writeErr)
		}
	}))
	defer server.Close()
	backend := testOpenAIBackend(t, Config{
		Endpoint: server.URL,
		Timeout:  "1ms",
	})

	pingErr := backend.Ping(context.Background())
	var unreachable errBackend.Unreachable
	if !errors.As(pingErr, &unreachable) {
		t.Fatalf("Ping() error = %T, want Unreachable", pingErr)
	}
}

func TestOpenAICompatiblePingJoinsEndpointBasePath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if r.URL.Path != "/proxy/openai"+cfgBackend.ModelsPath {
			t.Fatalf(
				"path = %s, want %s",
				r.URL.Path,
				"/proxy/openai"+cfgBackend.ModelsPath,
			)
		}
		_, writeErr := w.Write([]byte(`{"data":[{"id":"first"}]}`))
		if writeErr != nil {
			t.Fatalf("Write() error = %v", writeErr)
		}
	}))
	defer server.Close()
	backend := testOpenAIBackend(t, Config{Endpoint: server.URL + "/proxy/openai"})

	pingErr := backend.Ping(context.Background())
	if pingErr != nil {
		t.Fatalf("Ping() error = %v", pingErr)
	}
}

func TestOpenAICompatibleCompleteJoinsEndpointBasePath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if r.URL.Path != "/proxy/openai"+cfgBackend.ChatCompletionsPath {
			t.Fatalf(
				"path = %s, want %s",
				r.URL.Path,
				"/proxy/openai"+cfgBackend.ChatCompletionsPath,
			)
		}
		_, writeErr := w.Write([]byte(
			`{"model":"configured-model","choices":[{"message":{"content":"ok"}}]}`,
		))
		if writeErr != nil {
			t.Fatalf("Write() error = %v", writeErr)
		}
	}))
	defer server.Close()
	backend := testOpenAIBackend(t, Config{
		Endpoint:     server.URL + "/proxy/openai/",
		DefaultModel: "configured-model",
	})

	response, completeErr := backend.Complete(
		context.Background(),
		Request{Prompt: "hello"},
	)
	if completeErr != nil {
		t.Fatalf("Complete() error = %v", completeErr)
	}
	if response.Text != "ok" {
		t.Fatalf("Text = %q", response.Text)
	}
}

func TestOpenAICompatibleRequestModelOverride(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		var request chatRequest
		decodeErr := json.NewDecoder(r.Body).Decode(&request)
		if decodeErr != nil {
			t.Fatalf("Decode() error = %v", decodeErr)
		}
		if request.Model != "override-model" {
			t.Fatalf("model = %q", request.Model)
		}
		_, writeErr := w.Write([]byte(
			`{"model":"override-model","choices":[{"message":{"content":"ok"}}]}`,
		))
		if writeErr != nil {
			t.Fatalf("Write() error = %v", writeErr)
		}
	}))
	defer server.Close()
	backend := testOpenAIBackend(t, Config{
		Endpoint:     server.URL,
		DefaultModel: "configured-model",
	})

	response, completeErr := backend.Complete(
		context.Background(),
		Request{Model: "override-model"},
	)
	if completeErr != nil {
		t.Fatalf("Complete() error = %v", completeErr)
	}
	if response.Model != "override-model" {
		t.Fatalf("Model = %q", response.Model)
	}
}

func TestVLLMFactoryNameBehavior(t *testing.T) {
	backend, factoryErr := vllmFactory(Config{})
	if factoryErr != nil {
		t.Fatalf("vllmFactory() error = %v", factoryErr)
	}
	if backend.Name() != cfgBackend.NameVLLM {
		t.Fatalf("Name() = %q", backend.Name())
	}
}

func testOpenAIBackend(t *testing.T, config Config) Backend {
	t.Helper()
	backend, factoryErr := openAICompatibleFactory(config)
	if factoryErr != nil {
		t.Fatalf("openAICompatibleFactory() error = %v", factoryErr)
	}
	return backend
}

func testStatusBackend(t *testing.T, status int, body string) Backend {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.WriteHeader(status)
		_, writeErr := w.Write([]byte(body))
		if writeErr != nil {
			t.Fatalf("Write() error = %v", writeErr)
		}
	}))
	t.Cleanup(server.Close)
	return testOpenAIBackend(t, Config{Endpoint: server.URL})
}

func assertUpstreamBody(t *testing.T, targetErr error, body string) {
	t.Helper()
	var upstream errBackend.Upstream
	if !errors.As(targetErr, &upstream) {
		t.Fatalf("error = %T, want Upstream", targetErr)
	}
	if upstream.Body != body {
		t.Fatalf("Body = %q", upstream.Body)
	}
}
