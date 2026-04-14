//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package prompt is the **MCP `prompts/*` dispatcher** —
// the layer that takes a JSON-RPC request, validates
// the parameters, and routes to the right handler in
// [internal/mcp/handler] (or, for static catalog
// queries, returns the cached catalog directly).
//
// MCP exposes two prompt RPCs:
//
//   - **`prompts/list`** — return the catalog of
//     server-curated prompts the client may invoke.
//   - **`prompts/get`** — render one prompt by name
//     with its argument values filled in.
//
// # Public Surface
//
//   - **[DispatchList](req, deps)** — handles
//     `prompts/list`. Returns the static catalog
//     of ctx-curated prompts (session-start
//     ceremony, decision-add wizard, etc.).
//   - **[DispatchGet](req, deps)** — handles
//     `prompts/get`. Validates the requested
//     prompt name + arguments, calls into
//     [internal/mcp/handler] for the rendering.
//
// # Concurrency
//
// Each request runs in the read goroutine of
// [internal/mcp/server]; concurrent requests
// against the same `MCPDeps` are sequential by
// MCP design.
//
// # Related Packages
//
//   - [internal/mcp/server]           — owner of the
//     dispatch loop.
//   - [internal/mcp/handler]          — domain
//     logic for prompt rendering.
//   - [internal/mcp/server/def/tool]  — sister
//     package for tool catalog declarations.
//   - [internal/mcp/proto]            — wire types.
package prompt
