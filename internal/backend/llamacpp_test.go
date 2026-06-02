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
	"testing"
	"time"

	errBackend "github.com/ActiveMemory/ctx/internal/err/backend"
)

func TestNewLlamaCpp_DefaultName(t *testing.T) {
	srv := fakeOpenAIServer(t, testModel)
	defer srv.Close()
	b, err := newLlamaCpp(Config{Endpoint: srv.URL})
	if err != nil {
		t.Fatalf("ctor: %v", err)
	}
	if b.Name() != "llamacpp" {
		t.Errorf("Name = %q, want llamacpp", b.Name())
	}
}

func TestLlamaCpp_Ping_HappyDelegatesToOpenAICompat(t *testing.T) {
	srv := fakeOpenAIServer(t, testModel)
	defer srv.Close()
	b, _ := newLlamaCpp(Config{Endpoint: srv.URL})
	if err := b.Ping(context.Background()); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}

func TestLlamaCpp_Ping_NonDialErrorReturnsImmediately(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "broken", http.StatusInternalServerError)
	}))
	defer srv.Close()
	b, _ := newLlamaCpp(Config{Endpoint: srv.URL})
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
