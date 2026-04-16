//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package notify defines MCP notification method strings
// that the server sends to connected clients when
// server-side state changes.
//
// Notifications are JSON-RPC 2.0 messages with no "id"
// field -- they are fire-and-forget. The ctx MCP server
// uses them to tell clients that a subscribed resource
// has changed on disk, so the client can re-read the
// resource and update its view.
//
// # Key Constants
//
//   - [ResourcesUpdated]
//     ("notifications/resources/updated") -- emitted
//     when a context file that a client has subscribed
//     to via resources/subscribe is modified. The
//     notification includes the resource URI so the
//     client knows which file to re-fetch.
//
// # How Notifications Flow
//
// The server polls context files at a fixed interval
// (see [server.PollIntervalSec]). When a file's
// modification time changes and at least one client has
// subscribed to the corresponding resource, the server
// writes a ResourcesUpdated notification to stdout.
//
// # Why These Are Centralized
//
// The notification sender and the subscription manager
// must agree on the method string. A constant ensures
// compile-time agreement and makes the notification
// surface area discoverable via godoc.
package notify
