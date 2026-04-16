//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package err provides shared error and warning output
// helpers used across all CLI commands.
//
// # Exported Functions
//
// [With] prints a formatted error message to the
// command's stderr stream with an "Error: " prefix.
// This is the standard way to report fatal errors
// in ctx commands before returning an error code.
//
// [WarnFile] prints a non-fatal file operation warning
// to stdout with the file path and underlying error.
// This is used when a file operation fails but the
// command can continue (e.g. a missing optional file).
//
// # Nil Safety
//
// Both functions treat a nil *cobra.Command as a no-op,
// making them safe to call from code paths where a
// command may not be available.
//
// # Message Categories
//
//   - Error: prefixed error messages to stderr
//   - Warning: file-specific warnings to stdout
//
// # Usage
//
//	if err != nil {
//	    writeerr.With(cmd, err)
//	    return err
//	}
//	writeerr.WarnFile(cmd, path, err)
package err
