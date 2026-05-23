//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	"context"
	"errors"
	"testing"

	errBackend "github.com/ActiveMemory/ctx/internal/err/backend"
)

// fakeBackend is a no-op Backend used to exercise the
// registry without standing up an HTTP server.
type fakeBackend struct{ name string }

func (b *fakeBackend) Name() string               { return b.name }
func (b *fakeBackend) Ping(context.Context) error { return nil }
func (b *fakeBackend) Models(context.Context) ([]string, error) {
	return []string{b.name + "/dummy"}, nil
}
func (b *fakeBackend) Complete(context.Context, Request) (Response, error) {
	return Response{}, nil
}

func newFake(cfg Config) (Backend, error) {
	return &fakeBackend{name: cfg.Name}, nil
}

func TestRegister_AndResolve(t *testing.T) {
	r := New()
	if err := r.Register("vllm", newFake); err != nil {
		t.Fatalf("register: %v", err)
	}
	r.Configure(Config{Name: "vllm", Endpoint: "http://x"})
	b, err := r.Resolve("vllm")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if b.Name() != "vllm" {
		t.Fatalf("Name() = %q, want vllm", b.Name())
	}
}

func TestRegister_Duplicate(t *testing.T) {
	r := New()
	if err := r.Register("vllm", newFake); err != nil {
		t.Fatalf("first register: %v", err)
	}
	err := r.Register("vllm", newFake)
	if !errors.Is(err, errBackend.ErrDuplicateRegistration) {
		t.Fatalf("second register: got %v, want ErrDuplicateRegistration", err)
	}
}

func TestResolve_Unknown(t *testing.T) {
	r := New()
	if err := r.Register("vllm", newFake); err != nil {
		t.Fatalf("register: %v", err)
	}
	r.Configure(Config{Name: "vllm"})
	_, err := r.Resolve("openai")
	if !errors.Is(err, errBackend.ErrBackendNotFound) {
		t.Fatalf("resolve unknown: got %v, want ErrBackendNotFound", err)
	}
}

func TestResolve_NoBackends(t *testing.T) {
	r := New()
	_, err := r.Resolve("vllm")
	if !errors.Is(err, errBackend.ErrNoBackends) {
		t.Fatalf("resolve empty: got %v, want ErrNoBackends", err)
	}
}

func TestDefault_NoneConfigured(t *testing.T) {
	r := New()
	_, err := r.Default()
	if !errors.Is(err, errBackend.ErrNoBackends) {
		t.Fatalf("default empty: got %v, want ErrNoBackends", err)
	}
}

func TestDefault_SingleImplicit(t *testing.T) {
	r := New()
	if err := r.Register("vllm", newFake); err != nil {
		t.Fatalf("register: %v", err)
	}
	r.Configure(Config{Name: "vllm"})
	b, err := r.Default()
	if err != nil {
		t.Fatalf("default: %v", err)
	}
	if b.Name() != "vllm" {
		t.Fatalf("Name() = %q, want vllm", b.Name())
	}
}

func TestDefault_Explicit(t *testing.T) {
	r := New()
	if err := r.Register("vllm", newFake); err != nil {
		t.Fatalf("register vllm: %v", err)
	}
	if err := r.Register("openai", newFake); err != nil {
		t.Fatalf("register openai: %v", err)
	}
	r.Configure(Config{Name: "vllm"})
	r.Configure(Config{Name: "openai"})
	r.SetDefault("openai")
	b, err := r.Default()
	if err != nil {
		t.Fatalf("default: %v", err)
	}
	if b.Name() != "openai" {
		t.Fatalf("Name() = %q, want openai", b.Name())
	}
}

func TestDefault_Ambiguous(t *testing.T) {
	r := New()
	if err := r.Register("vllm", newFake); err != nil {
		t.Fatalf("register vllm: %v", err)
	}
	if err := r.Register("openai", newFake); err != nil {
		t.Fatalf("register openai: %v", err)
	}
	r.Configure(Config{Name: "vllm"})
	r.Configure(Config{Name: "openai"})
	_, err := r.Default()
	if !errors.Is(err, errBackend.ErrAmbiguousDefault) {
		t.Fatalf("default: got %v, want ErrAmbiguousDefault", err)
	}
}
