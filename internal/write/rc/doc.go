//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package rc provides terminal output for runtime
// configuration loading warnings.
//
// During startup, ctx loads YAML configuration files
// from the context directory. When a file cannot be
// parsed, the warning must be emitted before any
// cobra command is available, so this package writes
// directly to os.Stderr through the log/warn layer
// rather than through cobra's output stream.
//
// # Output
//
// [ParseWarning] prints a YAML parse warning that
// includes the filename that failed and the parse
// error. It delegates to the structured warning
// system in [internal/log/warn] with the config
// warning category from [internal/config/warn].
//
// This runs at config-load time, before any cobra
// command exists, so it bypasses the usual
// *cobra.Command output pattern used by other
// write packages.
package rc
