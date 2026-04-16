//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package bootstrap provides terminal output for the
// bootstrap command (ctx system bootstrap).
//
// # Exported Functions
//
// [Dir] prints just the context directory path for
// quiet/machine-readable mode. This is used when agents
// need only the path without surrounding text.
//
// [Text] prints the full human-readable bootstrap
// output: a title banner, the context directory path,
// a wrapped file list, numbered constitution rules,
// numbered next-step suggestions, and an optional
// warning block.
//
// [JSON] prints the bootstrap output as a structured
// JSON object containing the context directory, file
// list, rules, next steps, and optional warnings. If
// JSON encoding fails, a structured error object is
// printed instead.
//
// [CommunityFooter] prints the community link footer
// shown at the bottom of help output.
//
// # Message Categories
//
//   - Info: directory path, file list, rules, next steps
//   - Warning: optional initialization warnings
//   - Error: JSON encoding failures (to stderr)
//
// # Usage
//
//	if quiet {
//	    bootstrap.Dir(cmd, contextDir)
//	} else if jsonMode {
//	    bootstrap.JSON(cmd, dir, files, rules, steps, warn)
//	} else {
//	    bootstrap.Text(cmd, dir, fileList, rules, steps, warn)
//	}
package bootstrap
