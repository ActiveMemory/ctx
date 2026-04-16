//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package mcp is the root namespace for Model Context
// Protocol configuration constants used by the ctx MCP
// server.
//
// The MCP server exposes context files (TASKS.md,
// DECISIONS.md, LEARNINGS.md, etc.) as resources, tools,
// and prompts to AI coding agents over JSON-RPC 2.0 on
// stdio. Every string literal that appears in an MCP
// message -- method names, tool names, field keys, MIME
// types, notification subjects, resource identifiers,
// prompt names, and error codes -- is defined as a typed
// constant in a domain-specific sub-package below this
// one.
//
// # Why Centralize MCP Constants
//
// The MCP server wire format is a contract between ctx
// and every AI client that connects to it. Scattering
// string literals across handler code invites silent
// breakage when one site is renamed but another is not.
// Centralizing them here means:
//
//   - A typo fails to compile instead of producing a
//     runtime mismatch.
//   - Refactoring a tool name or method path is a
//     single-constant change verified by the compiler.
//   - The full surface area of the protocol is
//     discoverable via godoc on this package tree.
//
// # Sub-Package Layout
//
// Each sub-package owns one domain of the protocol:
//
//   - [cfg]        -- server tuning: buffer sizes,
//     default limits, word-overlap thresholds.
//   - [event]      -- session lifecycle markers
//     ("start", "end").
//   - [field]      -- JSON property key names for tool
//     input schemas.
//   - [governance] -- timing thresholds for drift and
//     persist nudge hooks.
//   - [method]     -- JSON-RPC method strings
//     ("tools/call", "resources/list", etc.).
//   - [mime]       -- MIME type and content-type
//     identifiers.
//   - [notify]     -- notification method strings sent
//     to clients on resource changes.
//   - [prompt]     -- prompt registration names that
//     mirror ctx CLI skills.
//   - [resource]   -- resource name constants used as
//     URI path segments.
//   - [schema]     -- JSON Schema type identifiers and
//     JSON-RPC error codes.
//   - [server]     -- server identity, URI prefix, and
//     launch arguments.
//   - [tool]       -- tool registration names that map
//     to ctx CLI subcommands.
package mcp
