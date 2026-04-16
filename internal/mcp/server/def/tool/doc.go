//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package tool defines the **MCP tool catalog** ctx
// advertises in `tools/list` responses: every callable
// MCP tool's schema, parameter definitions, and the
// shared property builders that keep parameter shapes
// consistent across tools.
//
// The package is the *catalog declaration*; the
// dispatch is in [internal/mcp/server/route/tool] and
// the actual logic is in [internal/mcp/handler].
//
// # Public Surface
//
//   - **[Defs]** — the slice of tool definitions
//     advertised in `tools/list`. Each definition
//     carries a name, description, JSON-schema
//     parameters, and an "annotations" map for
//     UI hints.
//   - **[MergeProps](base, extra)** — composes two
//     property maps so a tool can layer its
//     specific arguments on top of the shared
//     entry-attribute boilerplate.
//   - **[EntryAttrProps]** — the canonical property
//     map shared by `ctx_add` variants (priority,
//     branch, commit, session-id, etc.) so the
//     four entry-add tools have an identical
//     argument shape.
//
// # Why a Definitions Package
//
// MCP clients consume `tools/list` once at session
// start and cache the schemas. Centralizing the
// declarations makes the surface stable across
// versions: dispatch and handler refactors do not
// change what the client sees.
//
// # Concurrency
//
// All exports are immutable. Safe for concurrent
// reads.
package tool
