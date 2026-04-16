//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package warn defines Printf-style format string
// constants for best-effort warning messages routed
// through the structured logger.
//
// ctx treats many I/O failures as non-fatal: a
// failed file close, a directory walk error, or a
// marshal failure should not crash the CLI. Instead
// these are logged as warnings. This package provides
// the format strings so every call site uses the same
// phrasing and the same argument order (path, error).
//
// # File I/O Formats
//
//   - [Close]: file close failure.
//   - [Write]: file write failure.
//   - [Remove]: file remove failure.
//   - [Mkdir]: directory creation failure.
//   - [Rename]: file rename failure.
//   - [Walk]: directory walk failure.
//   - [Readdir]: directory read failure.
//   - [Getwd]: working directory resolution.
//
// # Serialization Formats
//
//   - [Marshal]: JSON marshal failure.
//   - [JSONEncode]: JSON-safe error for encoding
//     failures (returns valid JSON).
//
// # Specialized Formats
//
//   - [ParseConfig]: config file parse failure
//     with a "using defaults" note.
//   - [CloseResponse]: HTTP response body close.
//   - [CopilotClose]: Copilot CLI file close.
//
// # Context Identifiers
//
//   - [IndexHeader], [IndexSeparator], [IndexRow],
//     [ResponseBody]: labels attached to warnings
//     for traceability in index generation and HTTP
//     operations.
//
// # Why Centralized
//
// Using constants prevents typo drift across 40+
// call sites and makes it easy to grep for all
// warning messages in the codebase.
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
package warn
