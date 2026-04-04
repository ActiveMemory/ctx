//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

// ProtocolVersion is the MCP protocol version string.
const ProtocolVersion = "2024-11-05"

// Standard JSON-RPC error codes.
const (
	// ErrCodeParse indicates malformed JSON.
	ErrCodeParse = -32700
	// ErrCodeNotFound indicates method not found.
	ErrCodeNotFound = -32601
	// ErrCodeInvalidArg indicates invalid parameters.
	ErrCodeInvalidArg = -32602
	// ErrCodeInternal indicates an internal error.
	ErrCodeInternal = -32603
)

// JSON Schema type constants.
const (
	// Object is the JSON Schema type for objects.
	Object = "object"
	// String is the JSON Schema type for strings.
	String = "string"
	// Number is the JSON Schema type for numbers.
	Number = "number"
	// Boolean is the JSON Schema type for booleans.
	Boolean = "boolean"
)
