//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package resource handles MCP resource requests
// including list, read, subscribe, and unsubscribe
// operations.
//
// # Dispatchers
//
// DispatchList returns the pre-built resource list
// that the catalog package constructed at startup.
//
// DispatchRead loads context from disk and returns
// the requested resource content. It handles two
// kinds of resources:
//
//   - Individual file resources: looked up via
//     catalog.FileForURI and returned as-is.
//   - Agent packet: assembled from all context files
//     in read order, respecting the token budget.
//     Files that exceed the budget are listed as
//     "Also noted" summaries instead.
//
// DispatchSubscribe and DispatchUnsubscribe parse
// the subscription params and delegate to a callback
// function with the validated URI.
//
// # Agent Packet Assembly
//
// The readAgentPacket function assembles context files
// in priority order. Each file is formatted as a
// labeled section. When the cumulative token count
// exceeds the budget, remaining files are omitted
// and listed as summaries.
//
// # Subscription Handling
//
// The applySubscription helper validates params for
// both subscribe and unsubscribe, then calls the
// provided callback with the extracted URI.
package resource
