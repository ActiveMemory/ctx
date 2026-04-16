//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package loop provides terminal output for the loop
// script generation command (ctx loop).
//
// # Exported Functions
//
// [InfoGenerated] reports successful loop script
// generation with full details: the output file path,
// heading text, selected AI tool, prompt file path,
// maximum iteration count (or "unlimited"), and the
// completion signal string. The output is formatted
// as a multi-line block so the user can review all
// parameters at a glance.
//
// # Message Categories
//
//   - Info: script generation confirmation with all
//     configuration parameters
//
// # Usage
//
//	loop.InfoGenerated(cmd,
//	    outputFile, heading, tool,
//	    promptFile, maxIter, completionMsg)
package loop
