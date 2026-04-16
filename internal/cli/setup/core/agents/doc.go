//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package agents deploys AGENTS.md for universal agent
// instructions during project setup.
//
// AGENTS.md provides baseline instructions that any AI
// coding agent can follow, regardless of vendor. This
// ensures consistent behavior across tools like Claude
// Code, Cursor, Cline, and Kiro.
//
// # Deploy Algorithm
//
// [Deploy] generates or merges AGENTS.md in the project
// root using a three-way decision:
//
//  1. If the file exists and contains ctx marker
//     comments (ctx:agents), the file is skipped to
//     avoid duplicating content.
//  2. If the file exists without ctx markers, the
//     template content is appended to preserve any
//     existing non-ctx instructions.
//  3. If the file does not exist, it is created from
//     the embedded template.
//
// The template content is loaded from embedded assets
// via the agent asset reader. Marker detection uses
// string matching against the configured marker
// constant.
//
// # Data Flow
//
// The setup core orchestrator calls [Deploy] during
// ctx init. Deploy reads the embedded AGENTS.md
// template, checks for existing content, and writes
// the result. Status messages are emitted through
// the write/setup package.
package agents
