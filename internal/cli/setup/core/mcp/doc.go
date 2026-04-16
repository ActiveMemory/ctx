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
// directory, marshal config, write file, and print a
// confirmation message.
//
// # Deploy
//
// [Deploy] encapsulates the shared file-writing
// workflow. Tool-specific packages build their config
// struct and pass it here for writing. The function:
//
//  1. Checks whether the target file already exists.
//     If so, prints a skip message and returns nil.
//  2. Creates the parent directory if it does not
//     exist.
//  3. Marshals the config struct as indented JSON.
//  4. Writes the JSON to the target path with a
//     trailing newline.
//
// # SyncSteering
//
// [SyncSteering] synchronizes steering files from the
// context directory to a tool-native format. It checks
// whether a steering directory exists, then delegates
// to the steering package for the actual file sync.
// The report lists files that were written or skipped.
//
// If the steering directory does not exist, the
// function prints a message and returns nil without
// error.
package mcp
