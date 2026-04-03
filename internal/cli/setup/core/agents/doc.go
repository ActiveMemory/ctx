//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package agents deploys AGENTS.md for universal agent instructions.
//
// [Deploy] generates or merges AGENTS.md in the project root,
// preserving existing non-ctx content via marker detection. The file
// provides baseline instructions that any AI coding agent can follow,
// regardless of vendor, ensuring consistent behavior across tools.
//
// Key exports: [Deploy].
// Called by the setup core orchestrator during ctx init.
package agents
