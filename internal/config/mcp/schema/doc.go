//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package schema defines JSON Schema type identifiers and
// JSON-RPC 2.0 error codes used throughout the ctx MCP
// server.
//
// Two groups of constants live here:
//
// # Protocol Version
//
//   - [ProtocolVersion] ("2024-11-05") -- the MCP
//     protocol version negotiated during the
//     initialize handshake.
//
// # JSON-RPC Error Codes
//
// Standard error codes returned in JSON-RPC error
// response objects:
//
//   - [ErrCodeParse] (-32700)      -- malformed JSON
//     in the request body.
//   - [ErrCodeNotFound] (-32601)   -- the requested
//     method does not exist.
//   - [ErrCodeInvalidArg] (-32602) -- the parameters
//     object is invalid or missing required fields.
//   - [ErrCodeInternal] (-32603)   -- an internal
//     server error occurred.
//
// # JSON Schema Type Constants
//
// Type identifiers used when declaring tool input
// schemas:
//
//   - [Object]  -- "object", the root type of every
//     tool input schema.
//   - [String]  -- "string", for text fields like
//     content, query, section.
//   - [Number]  -- "number", for numeric fields like
//     limit.
//   - [Boolean] -- "boolean", for toggle fields like
//     archive.
//
// # Why These Are Centralized
//
// Error codes appear in the dispatcher, in individual
// handlers, and in tests. Schema types appear in every
// tool registration. Centralizing them avoids magic
// numbers and makes the set of recognized types
// explicit.
package schema
