//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package setup defines display names and deploy-path
// constants for the ctx setup command, which configures
// third-party AI editor integrations.
//
// ctx can inject MCP server configuration and steering
// files into Kiro, Cursor, and Cline. Each editor
// stores configuration in a different directory tree.
// This package maps each tool to its expected paths so
// the setup logic stays generic.
//
// # Supported Tools
//
//   - **Kiro**: config dir .kiro/, MCP config at
//     [MCPConfigPathKiro], steering at
//     [SteeringDeployPathKiro].
//   - **Cursor**: config dir .cursor/, MCP config at
//     [MCPConfigPathCursor], steering at
//     [SteeringPathCursor].
//   - **Cline**: MCP config at [MCPConfigPathCline],
//     steering at [SteeringPathCline].
//
// # Display Names
//
// [DisplayKiro], [DisplayCursor], [DisplayCline]
// provide user-facing labels for menus and status
// output.
//
// # Why Centralized
//
// The setup command, the steering sync engine, and
// the init scaffolder all need to agree on where
// each tool's files live. Centralizing paths here
// prevents one subsystem from writing to a path
// another does not read.
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
package setup
