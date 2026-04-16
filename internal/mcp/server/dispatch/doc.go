//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package dispatch routes incoming MCP requests to
// domain-specific handlers based on the JSON-RPC
// method name.
//
// # Routing
//
// Do is the single entry point for request dispatch.
// It receives a parsed JSON-RPC request and delegates
// to the appropriate handler:
//
//   - initialize -> initialize.Dispatch
//   - ping       -> ping.Dispatch
//   - resources/list       -> resource.DispatchList
//   - resources/read       -> resource.DispatchRead
//   - resources/subscribe  -> resource.DispatchSubscribe
//   - resources/unsubscribe -> resource.DispatchUnsubscribe
//   - tools/list  -> tool.DispatchList
//   - tools/call  -> tool.DispatchCall
//   - prompts/list -> prompt.DispatchList
//   - prompts/get  -> prompt.DispatchGet
//
// Unrecognized methods fall through to the fallback
// handler, which returns a method-not-found error.
//
// # Dependencies
//
// Do accepts an entity.MCPDeps struct that carries
// runtime dependencies (context directory, token
// budget, session info) needed by domain handlers.
// It also takes the pre-built resource list and a
// poller for resource subscriptions.
//
// # Usage
//
//	resp := dispatch.Do(
//	    version, deps, resList, poller, req,
//	)
package dispatch
