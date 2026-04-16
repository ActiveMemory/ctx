//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package event defines the session lifecycle event type
// values used by the MCP server to mark the boundaries of
// an AI coding session.
//
// When an AI agent connects to the ctx MCP server it
// signals session start and session end via the
// ctx_session_event tool. The value passed in the event
// field is one of the constants defined here.
//
// # Key Constants
//
//   - [Start] ("start"): sent when the agent begins a
//     new session. The MCP server uses this to
//     initialize per-session state such as the tool
//     call counter and the drift check timer.
//   - [End] ("end"): sent when the agent is about to
//     exit. The server flushes pending writes,
//     records a journal entry, and tears down session
//     state.
//
// # Why These Are Centralized
//
// The event values appear in tool input validation, hook
// dispatch, and journal recording. A constant ensures
// that the handler's switch statement and the client's
// request always agree on the exact string.
package event
