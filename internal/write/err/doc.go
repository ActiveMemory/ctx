//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package err provides shared error output helpers for CLI commands.
//
// [With] prints a formatted error message to the command's error
// output. [WarnFile] prints a file-specific warning with the path
// and underlying error. Both route through cobra's error writer
// to respect output redirection.
//
// Example:
//
//	if err != nil {
//	    write.With(cmd, err)
//	    return
//	}
package err
