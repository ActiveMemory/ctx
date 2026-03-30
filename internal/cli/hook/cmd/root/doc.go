//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the ctx hook command for generating
// AI tool integration configs.
//
// [Cmd] builds the cobra.Command with --write flag. [Run] generates
// hook configurations for Claude Code, Cursor, Copilot, and others.
// [WriteCopilotInstructions] deploys the embedded copilot template.
package root
