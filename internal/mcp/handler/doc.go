//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package handler contains domain logic for MCP tool operations.
//
// Functions accept typed Go parameters and return (string, error) pairs.
// The server package handles JSON-RPC protocol translation, argument
// extraction from MCP maps, and response wrapping. This separation
// keeps domain logic testable without protocol coupling.
package handler
