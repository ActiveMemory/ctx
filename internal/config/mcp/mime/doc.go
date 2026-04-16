//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package mime defines MIME type and content-type
// identifiers used in MCP resource and tool responses.
//
// Every resource the ctx MCP server exposes carries a
// MIME type that tells the AI client how to interpret the
// payload. Similarly, tool results declare a content type
// so the client can render text, markdown, or other
// formats correctly.
//
// # Key Constants
//
//   - [Markdown] ("text/markdown") -- the MIME type
//     assigned to context-file resources
//     (TASKS.md, DECISIONS.md, etc.). Clients that
//     understand markdown can render headings, lists,
//     and code fences.
//   - [ContentTypeText] ("text") -- the content type
//     value used in tool result objects. The MCP
//     specification defines "text" and "image" as the
//     two content type discriminators; ctx uses "text"
//     exclusively because all tool output is textual.
//
// # Why These Are Centralized
//
// Resource registration, resource read handlers, and
// tool result builders all emit these values. A typo in
// a MIME string would cause clients to fall back to
// plain-text rendering silently. Typed constants
// eliminate that risk.
package mime
