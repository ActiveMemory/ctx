//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package mcp defines the typed error constructors
// for the MCP (Model Context Protocol) server. These
// errors fire when MCP tool calls are missing
// required parameters, when the context directory
// cannot be read during search, or when an
// unrecognized session event type is received.
//
// # Domain
//
// Errors fall into two categories:
//
//   - **Validation** -- a required field is missing
//     from a tool call payload. Constructors:
//     [TypeContentRequired], [QueryRequired],
//     [UnknownEventType].
//   - **Search IO** -- the context directory could
//     not be read during a search operation.
//     Constructor: [SearchRead].
//
// # Wrapping Strategy
//
// [SearchRead] wraps its cause with fmt.Errorf %w
// so callers can inspect the underlying read error.
// Validation constructors return plain errors.
// All user-facing text is resolved through
// [internal/assets/read/desc].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package mcp
