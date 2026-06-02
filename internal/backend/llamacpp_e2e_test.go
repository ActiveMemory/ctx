//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

//go:build e2e

package backend

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestLlamaCpp_E2E_PingRealServer(t *testing.T) {
	endpoint := os.Getenv("LLAMACPP_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8080"
	}
	b, err := newLlamaCpp(Config{Endpoint: endpoint})
	if err != nil {
		t.Fatalf("ctor: %v", err)
	}
	if b.Name() != "llamacpp" {
		t.Errorf("Name = %q, want llamacpp", b.Name())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := b.Ping(ctx); err != nil {
		t.Fatalf("Ping real server: %v", err)
	}
	t.Logf("Ping OK: %s", endpoint)
}

func TestLlamaCpp_E2E_CompleteRealServer(t *testing.T) {
	endpoint := os.Getenv("LLAMACPP_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8080"
	}
	b, err := newLlamaCpp(Config{
		Endpoint:     endpoint,
		DefaultModel: "test",
	})
	if err != nil {
		t.Fatalf("ctor: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := b.Complete(ctx, Request{
		Messages: []Message{
			{Role: "user", Content: "Say hello in one word."},
		},
		MaxTokens:   32,
		Temperature: 0.0,
	})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	t.Logf("Model: %s", resp.Model)
	t.Logf("Content: %q", resp.Content)
	t.Logf("Raw: %s", string(resp.Raw))
	if resp.Content == "" && len(resp.Raw) == 0 {
		t.Fatal("Complete returned empty content and empty raw")
	}
}
