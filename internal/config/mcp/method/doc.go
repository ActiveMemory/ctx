//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package method defines the JSON-RPC 2.0 method strings
// that the ctx MCP server recognizes. Each constant maps
// to a top-level route in the server's request dispatcher.
//
// When an AI client sends a JSON-RPC request, the
// "method" field is matched against these constants to
// select the appropriate handler. The set covers the full
// MCP specification surface that ctx implements.
//
// # Key Constants
//
//   - [Initialize] ("initialize") -- the handshake that
//     negotiates protocol version and capabilities.
//   - [Ping] ("ping") -- a keep-alive probe.
//   - [ResourceList] ("resources/list") -- enumerates
//     available context-file resources.
//   - [ResourceRead] ("resources/read") -- returns the
//     content of a single context resource.
//   - [ResourceSubscribe] / [ResourceUnsubscribe] --
//     manage change-notification subscriptions.
//   - [ToolList] ("tools/list") -- enumerates the tools
//     the server exposes.
//   - [ToolCall] ("tools/call") -- invokes a tool by
//     name with a JSON input object.
//   - [PromptList] ("prompts/list") -- lists available
//     prompts.
//   - [PromptGet] ("prompts/get") -- retrieves a prompt
//     template by name.
//
// # Why These Are Centralized
//
// The dispatcher switch and integration tests both
// reference method strings. A constant ensures that a
// renamed method breaks at compile time rather than
// producing an unmatched route at runtime.
package method
