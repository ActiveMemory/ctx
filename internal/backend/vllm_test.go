//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	errBackend "github.com/ActiveMemory/ctx/internal/err/backend"
)

func TestNewVLLM_DefaultName(t *testing.T) {
	srv := fakeOpenAIServer(t, testModel)
	defer srv.Close()
	b, err := newVLLM(Config{Endpoint: srv.URL})
	if err != nil {
		t.Fatalf("ctor: %v", err)
		return
	}
	if b.Name() != "vllm" {
		t.Errorf("Name = %q, want vllm", b.Name())
	}
}

func TestVLLM_Ping_HappyDelegatesToOpenAICompat(t *testing.T) {
	srv := fakeOpenAIServer(t, testModel)
	defer srv.Close()
	b, _ := newVLLM(Config{Endpoint: srv.URL})
	if err := b.Ping(context.Background()); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}

func TestVLLM_Ping_NonDialErrorReturnsImmediately(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "broken", http.StatusInternalServerError)
	}))
	defer srv.Close()
	b, _ := newVLLM(Config{Endpoint: srv.URL})
	// Tight retry window — if vllm.Ping were to retry on
	// 500 (it must not), this test would block on the
	// interval. Setting window=10ms keeps the test fast
	// even if the retry guard regresses.
	b.coldStartWindow = 10 * time.Millisecond
	b.coldStartInterval = 5 * time.Millisecond
	start := time.Now()
	err := b.Ping(context.Background())
	elapsed := time.Since(start)
	if !errors.Is(err, errBackend.ErrUnhealthyStatus) {
		t.Fatalf("got %v, want ErrUnhealthyStatus", err)
	}
	if elapsed > 200*time.Millisecond {
		t.Errorf("Ping retried on 500; elapsed %v should be near-zero", elapsed)
	}
}

func TestColdStartRetry_SucceedsAfterRefused(t *testing.T) {
	var attempts atomic.Int32
	ping := func(_ context.Context) error {
		n := attempts.Add(1)
		if n < 3 {
			return errBackend.Unreachable("vllm", "http://x", syscall.ECONNREFUSED)
		}
		return nil
	}
	err := coldStartRetry(
		context.Background(), ping,
		200*time.Millisecond, 10*time.Millisecond,
	)
	if err != nil {
		t.Fatalf("retry: %v", err)
	}
	if got := attempts.Load(); got != 3 {
		t.Errorf("attempts = %d, want 3", got)
	}
}

func TestColdStartRetry_NonDialErrorReturnsImmediately(t *testing.T) {
	var attempts atomic.Int32
	ping := func(_ context.Context) error {
		attempts.Add(1)
		return errBackend.UnhealthyStatus("vllm", 500, "broken")
	}
	err := coldStartRetry(
		context.Background(), ping,
		200*time.Millisecond, 10*time.Millisecond,
	)
	if !errors.Is(err, errBackend.ErrUnhealthyStatus) {
		t.Fatalf("got %v, want ErrUnhealthyStatus", err)
	}
	if got := attempts.Load(); got != 1 {
		t.Errorf("attempts = %d, want 1 (no retry on non-dial)", got)
	}
}

func TestColdStartRetry_WindowExpires(t *testing.T) {
	var attempts atomic.Int32
	ping := func(_ context.Context) error {
		attempts.Add(1)
		return errBackend.Unreachable("vllm", "http://x", syscall.ECONNREFUSED)
	}
	start := time.Now()
	err := coldStartRetry(
		context.Background(), ping,
		30*time.Millisecond, 10*time.Millisecond,
	)
	elapsed := time.Since(start)
	if !errors.Is(err, errBackend.ErrUnreachable) {
		t.Fatalf("got %v, want ErrUnreachable after window", err)
	}
	if elapsed < 30*time.Millisecond {
		t.Errorf("retry gave up too early: %v", elapsed)
	}
	if got := attempts.Load(); got < 2 {
		t.Errorf("attempts = %d, expected at least 2", got)
	}
}

func TestColdStartRetry_ContextCancelled(t *testing.T) {
	ping := func(_ context.Context) error {
		return errBackend.Unreachable("vllm", "http://x", syscall.ECONNREFUSED)
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()
	err := coldStartRetry(
		ctx, ping,
		1*time.Second, 5*time.Millisecond,
	)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v, want context.Canceled", err)
	}
}
