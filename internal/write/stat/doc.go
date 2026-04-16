//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package stat provides terminal output for the stats
// command (ctx stats).
//
// The stats command displays context usage metrics
// in a tabular format. This package handles both
// cobra-routed output and direct writer output for
// streaming scenarios.
//
// # Table Output
//
// [Table] prints pre-formatted stats lines through
// cobra's output stream. It accepts a slice of
// strings containing the header, separator, and
// data rows, and delegates to [line.All] for
// sequential printing. A nil *cobra.Command is
// treated as a no-op.
//
// # Streaming Output
//
// [StreamLine] writes a single formatted stats line
// to an arbitrary io.Writer. This is used when
// stats are emitted outside a cobra command context,
// such as piped output or background reporting.
package stat
