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

type fakeBackend struct {
	name string
}

func (backend fakeBackend) Name() string {
	return backend.name
}

func (backend fakeBackend) Ping(context.Context) error {
	return nil
}

func (backend fakeBackend) Complete(
	context.Context,
	Request,
) (Response, error) {
	return Response{Model: backend.name}, nil
}

func TestRegistryResolvesSingleBackendAsDefault(t *testing.T) {
	registry := &Registry{}
	registerTestBackend(t, registry, "vllm")

	resolved, resolveErr := registry.Default()
	if resolveErr != nil {
		t.Fatalf("Default() error = %v", resolveErr)
	}
	if resolved.Name() != "vllm" {
		t.Fatalf("Default().Name() = %q", resolved.Name())
	}
}

func TestRegistryMultipleBackendsWithoutDefaultFails(t *testing.T) {
	registry := &Registry{}
	registerTestBackend(t, registry, "vllm")
	registerTestBackend(t, registry, "openai")

	_, resolveErr := registry.Default()
	var multiple errBackend.MultipleBackends
	if !errors.As(resolveErr, &multiple) {
		t.Fatalf("Default() error = %T, want MultipleBackends", resolveErr)
	}
}

func TestRegistryExplicitDefault(t *testing.T) {
	registry := &Registry{}
	registerTestBackend(t, registry, "vllm")
	registerTestBackend(t, registry, "openai")
	registry.SetDefault("openai")

	resolved, resolveErr := registry.Default()
	if resolveErr != nil {
		t.Fatalf("Default() error = %v", resolveErr)
	}
	if resolved.Name() != "openai" {
		t.Fatalf("Default().Name() = %q", resolved.Name())
	}
}

func TestRegistryMissingDefault(t *testing.T) {
	registry := &Registry{}
	registerTestBackend(t, registry, "vllm")
	registry.SetDefault("openai")

	_, resolveErr := registry.Default()
	var missing errBackend.MissingBackend
	if !errors.As(resolveErr, &missing) {
		t.Fatalf("Default() error = %T, want MissingBackend", resolveErr)
	}
}

func TestRegistryEmptyDefault(t *testing.T) {
	registry := &Registry{}

	_, resolveErr := registry.Default()
	var empty errBackend.NoBackendConfigured
	if !errors.As(resolveErr, &empty) {
		t.Fatalf("Default() error = %T, want NoBackendConfigured", resolveErr)
	}
}

func TestRegistryMissingBackend(t *testing.T) {
	registry := &Registry{}

	_, resolveErr := registry.Resolve("missing")
	var missing errBackend.MissingBackend
	if !errors.As(resolveErr, &missing) {
		t.Fatalf("Resolve() error = %T, want MissingBackend", resolveErr)
	}
}

func TestRegistryDuplicateRegistration(t *testing.T) {
	registry := &Registry{}
	registerTestBackend(t, registry, "vllm")

	registerErr := registry.Register("vllm", Config{}, testFactory("vllm"))
	var duplicate errBackend.DuplicateRegistration
	if !errors.As(registerErr, &duplicate) {
		t.Fatalf("Register() error = %T, want DuplicateRegistration", registerErr)
	}
}

func TestRegistryFactoryErrors(t *testing.T) {
	factoryErr := errors.New("boom")
	registry := &Registry{}
	registerErr := registry.Register("vllm", Config{}, func(Config) (Backend, error) {
		return nil, factoryErr
	})
	if registerErr != nil {
		t.Fatalf("Register() error = %v", registerErr)
	}

	_, resolveErr := registry.Resolve("vllm")
	var wrapped errBackend.Factory
	if !errors.As(resolveErr, &wrapped) {
		t.Fatalf("Resolve() error = %T, want Factory", resolveErr)
	}
	if !errors.Is(resolveErr, factoryErr) {
		t.Fatalf("Resolve() does not wrap factory error")
	}
}

func registerTestBackend(t *testing.T, registry *Registry, name string) {
	t.Helper()
	registerErr := registry.Register(name, Config{Name: name}, testFactory(name))
	if registerErr != nil {
		t.Fatalf("Register() error = %v", registerErr)
	}
}

func testFactory(name string) Factory {
	return func(Config) (Backend, error) {
		return fakeBackend{name: name}, nil
	}
}
