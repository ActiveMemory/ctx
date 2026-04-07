//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

// Scanner buffer size constants.
const (
	// BufInitSize is the initial scanner buffer size
	// for session parsing (64 KB).
	BufInitSize = 64 * 1024
	// BufMaxSize is the maximum scanner buffer size
	// for session parsing (1 MB).
	BufMaxSize = 1024 * 1024
	// BufMaxSizeSchema is the maximum scanner buffer size
	// for schema validation (4 MB). Larger than BufMaxSize
	// because tool result lines can exceed 1 MB.
	BufMaxSizeSchema = 4 * 1024 * 1024
)
