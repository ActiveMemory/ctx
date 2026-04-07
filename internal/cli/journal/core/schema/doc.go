//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package schema provides core logic for journal schema
// validation.
//
// It resolves which directories to scan based on CLI flags,
// runs validation across all JSONL files in those directories,
// and manages the drift report lifecycle in .context/reports/.
// Used by both the standalone check command and the import
// integration, which validates source files after importing
// sessions and prints a summary if drift is found.
package schema
