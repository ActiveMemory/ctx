//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package mcp provides shared helpers for deploying MCP
// server configuration files across AI tool integrations.
//
// Each tool (Cursor, Kiro, Cline) has a unique JSON
// structure for its mcp.json file, but the deployment
// workflow is identical: check if file exists, create
// directory, marshal config, write file, print
// confirmation.
//
// The [Deploy] function encapsulates this shared workflow.
// Tool-specific packages build their config struct and
// pass it here for writing.
package mcp
