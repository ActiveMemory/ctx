//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package fallback returns method-not-found errors
// for unrecognized MCP method names.
//
// # Handler
//
// DispatchErr builds a JSON-RPC error response with
// the standard method-not-found error code. It
// includes the unrecognized method name in the error
// message so the client can diagnose the issue.
//
//	resp := fallback.DispatchErr(req)
//
// # Role in Dispatch
//
// The main dispatch package calls DispatchErr as the
// default case in its method switch. Any method that
// does not match a known handler (initialize, ping,
// resources/*, tools/*, prompts/*) falls through to
// this package.
package fallback
