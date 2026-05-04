//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package run provides the shared add execution logic invoked
// by every noun-first add subcommand (ctx task add, ctx
// decision add, ctx learning add, ctx convention add).
//
// The function is split out from the noun cobra wiring so all
// four noun packages can call into the same validation,
// extraction, formatting, insertion, and trace-recording
// pipeline without duplicating it.
//
// # Public Surface
//
//   - [Run] executes one add operation. Caller passes the
//     noun as args[0] (followed by any positional content
//     tokens) and an entity.AddConfig with the resolved
//     flag values.
//
// # Validation Boundaries
//
// Hard checks (required fields, secret patterns, length
// limits, provenance requirements per .ctxrc) live in
// [internal/entry] so the rules are identical regardless of
// caller (CLI here, MCP ctx_add tool elsewhere).
//
// # Concurrency
//
// Single-process, sequential.
package run
