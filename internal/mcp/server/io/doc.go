//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package io provides thread-safe JSON-RPC message
// writing over an io.Writer.
//
// # Writer
//
// Writer wraps an io.Writer with a mutex so that
// concurrent goroutines can safely emit JSON-RPC
// messages without interleaving. Each call to
// WriteJSON marshals the value, appends a newline,
// and writes the result as a single atomic operation.
//
//	w := io.NewWriter(os.Stdout)
//	err := w.WriteJSON(response)
//
// # Thread Safety
//
// The internal mutex serializes all writes. This is
// required because the MCP server may send responses
// and notifications from different goroutines (for
// example, resource change notifications alongside
// tool call responses).
//
// # Message Format
//
// Each message is marshaled to compact JSON and
// terminated with a single LF newline. This matches
// the JSON-RPC 2.0 line-delimited transport used by
// MCP over stdin/stdout.
package io
