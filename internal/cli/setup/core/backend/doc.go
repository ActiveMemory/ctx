//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package backend implements the `.ctxrc` templating used
// by `ctx setup --backend <name>`. Adds or updates a
// single entry in the `backends:` list keyed by `name:`,
// preserving every other top-level key, comment, and
// formatting in the file via yaml.v3's Node API.
//
// # Idempotency
//
// Re-running `ctx setup --backend X` with the same name
// updates the existing entry in place; same with same
// flags produces the same bytes. New names append.
//
// # Preservation
//
// yaml.v3's Node-tree round-trip preserves comments and
// the ordering of unrelated top-level keys. The only
// reformatting happens inside the `backends:` block this
// package mutates, and within the entries it adds.
//
// # Out of scope
//
// Downstream-tool env-var templating (e.g., writing
// `ANTHROPIC_BASE_URL` into a Claude Code settings file
// to route an upstream tool through a configured backend)
// is intentionally not handled here. Per-tool policies
// differ; the spec calls that work "(where applicable)"
// and a follow-up commit will add per-tool templating
// behind explicit opt-in flags.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/rc] reads the
//     `.ctxrc` file Apply writes; the on-wire shape
//     (`backends:` list of objects + `default_backend:`
//     scalar) is defined there.
//   - [github.com/ActiveMemory/ctx/internal/entity]
//     supplies the [entity.BackendConfig] shape Apply
//     accepts.
//   - [github.com/ActiveMemory/ctx/internal/err/setup]
//     owns the typed errors Apply returns on read/parse/
//     marshal/write failure.
package backend
