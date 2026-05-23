//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package backend abstracts an OpenAI-compatible HTTP
// inference surface so the rest of ctx can dispatch
// schema-constrained completions through a configured
// remote without knowing which vendor sits behind it.
//
// # Contract floor
//
// Every backend speaks `/v1/chat/completions` and
// `/v1/models` over HTTP. The [Backend] interface exposes
// three operations: [Backend.Name] (the registered type
// label), [Backend.Ping] (a reachability probe on
// `/v1/models`), and [Backend.Complete] (a single
// non-streaming chat completion). vLLM is the canonical
// local implementation; OpenAI, Anthropic, Ollama, and
// LM Studio are thin wrappers over the same wire shape.
//
// # Registry
//
// [Registry] decouples backend *type* (e.g., "vllm") from
// per-project *configuration*. A type is bound to a
// [Factory] via [Registry.Register]; a project's
// `.ctxrc` `[backends.<name>]` table becomes a [Config]
// via [Registry.Configure]; [Registry.Resolve] and
// [Registry.Default] return a constructed [Backend].
// Resolution failure modes carry typed sentinels from
// [github.com/ActiveMemory/ctx/internal/err/backend] so
// callers can `errors.Is` them and surface
// user-actionable recovery hints.
//
// # Determinism boundary
//
// This package is intentionally invisible to ctx's
// deterministic ceremony surface: `ctx agent`,
// `ctx status`, and the canonical hook commands must
// not import it. A dedicated audit test enforces the
// rule structurally; the rationale is Invariant 2 of
// `specs/ctx-ai-backend.md` (zero runtime deps for
// core functionality).
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/err/backend]
//     supplies the typed sentinel errors and wrapping
//     constructors for resolution failures.
//   - [github.com/ActiveMemory/ctx/internal/rc] parses
//     the `[backends]` table that feeds [Config] values
//     into [Registry.Configure].
//   - `internal/cli/ai/` (forthcoming) is the primary
//     caller; it dispatches user-issued AI verbs through
//     a [Backend] resolved at command time.
package backend
